package business

import (
	"context"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/notification/config"
	"go.uber.org/fx"
	"time"

	"github.com/alimitedgroup/MVP/srv/notification/portin"
	"github.com/alimitedgroup/MVP/srv/notification/portout"
	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RuleCheckerParams struct {
	fx.In
	Lc fx.Lifecycle

	Logger  *zap.Logger
	Brk     *broker.NatsMessageBroker
	Rules   portin.QueryRules
	Queries portout.RuleQueryRepository
	Publish portout.StockEventPublisher
	Cfg     *config.NotificationConfig
}

func NewRuleChecker(
	p RuleCheckerParams,
) *RuleChecker {
	rc := &RuleChecker{
		rulePort:    p.Rules,
		queryPort:   p.Queries,
		publishPort: p.Publish,
		Logger:      p.Logger.Named("rule-checker"),
		cfg:         p.Cfg,
		brk:         p.Brk,
		stop:        make(chan bool, 1),
		stopOk:      make(chan bool, 1),
	}

	p.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go rc.run()
			return nil
		},
		// Questo blocco si occupa di fermare la goroutine che gestisce il controllo delle regole.
		// La logica è semplice:
		// 1) Quando l'applicazione si sta fermando, mandiamo un messaggio sul canale stop
		// 2) Attendiamo finché non arriva una risposta, all'interno del canale stopOk,
		//    o finché non andiamo in timeout (ossia arriva un messaggio sul canale ctx.Done())
		// 3) Terminiamo
		OnStop: func(ctx context.Context) error {
			rc.stop <- true
			select {
			case <-ctx.Done():
			case <-rc.stopOk:
			}
			return nil
		},
	})

	return rc
}

type RuleChecker struct {
	*zap.Logger
	cfg *config.NotificationConfig
	brk *broker.NatsMessageBroker

	rulePort    portin.QueryRules
	queryPort   portout.RuleQueryRepository
	publishPort portout.StockEventPublisher

	// stop e stopOk sono canali su cui verrà mandato al massimo un messaggio
	// Per la logica che ci sta dietro, fare riferimento al commento all'interno di NewRuleChecker.
	stop   chan bool
	stopOk chan bool
}

func (rc *RuleChecker) run() {
	ticker := time.NewTicker(rc.cfg.CheckerTimer)
	for {
		select {
		case <-rc.stop:
			rc.Debug("RuleChecker stopped")
			rc.stopOk <- true
			return
		case <-ticker.C:
			rc.checkAllRules()
		}
	}
}

func (rc *RuleChecker) checkAllRules() {
	rules, err := rc.rulePort.ListQueryRules() // recupera tutte le regole dal repository in memoria
	if err != nil {
		rc.Error("Error while listing rules", zap.Error(err))
	}

	rc.Debug("Controllo periodico delle regole avviato", zap.Int("rulesAmount", len(rules)))
	if len(rules) == 0 {
		return
	}

	// Per ogni regola, interroga Influx e confronta la quantity con la threshold
	for _, rule := range rules {
		// Esempio: se rule è un AddQueryRuleCmd con metodi GetGoodID, GetOperator e GetThreshold
		goodID := rule.GoodId
		operator := rule.Operator
		threshold := rule.Threshold

		// Invoca il metodo del service che interroga Influx
		resp := rc.queryPort.GetCurrentQuantityByGoodID(goodID)
		if resp.Err != nil {
			rc.Error("Errore nel recupero stock", zap.String("goodId", goodID), zap.Error(err))
			continue
		}

		currentQuantity := resp.CurrentQuantity

		// Confronta con l'operatore
		condTrue := false
		switch operator {
		case "<":
			condTrue = currentQuantity < threshold
		case ">":
			condTrue = currentQuantity > threshold
		case "<=":
			condTrue = currentQuantity <= threshold
		case ">=":
			condTrue = currentQuantity >= threshold
		default:
			rc.Error("Operatore non valido", zap.String("operator", operator))
			continue
		}

		alert := types.StockAlertEvent{
			Id:              uuid.NewString(),
			GoodID:          goodID,
			CurrentQuantity: currentQuantity,
			Operator:        operator,
			Threshold:       threshold,
			Timestamp:       time.Now().UnixMilli(),
			RuleId:          rule.RuleId.String(),
		}

		if condTrue {
			alert.Status = types.StockPending
			err = rc.publishPort.PublishStockAlert(alert)
		} else {
			alert.Status = types.StockRevoked
			err = rc.publishPort.RevokeStockAlert(alert)
		}
		if err != nil {
			rc.Error(
				"Errore nell'aggiornamento di stato di una notifica",
				zap.String("goodId", goodID),
				zap.String("ruleId", rule.RuleId.String()),
				zap.String("newStatus", string(alert.Status)),
				zap.Error(err),
			)
		}
	}

	err = rc.brk.Nats.Flush()
	if err != nil {
		rc.Error("Error while flushing messages to broker", zap.Error(err))
	}
}
