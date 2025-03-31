package business

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/alimitedgroup/MVP/srv/notification/portin"
	"github.com/alimitedgroup/MVP/srv/notification/portout"
	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const DEFAULT_WAIT_SECONDS = 45

type RuleChecker struct {
	*zap.Logger

	rulePort    portin.QueryRules
	queryPort   portout.RuleQueryRepository
	publishPort portout.StockEventPublisher
	// stop è un canale su cui verranno mandati al massimo due messaggi.
	// Per la logica che ci sta dietro, fare riferimento al commento all'interno di NewRuleChecker.
	stop chan bool
}

func NewRuleChecker(lc fx.Lifecycle, logger *zap.Logger, rules portin.QueryRules, queries portout.RuleQueryRepository, publish portout.StockEventPublisher) *RuleChecker {
	rc := &RuleChecker{
		rulePort:    rules,
		queryPort:   queries,
		publishPort: publish,
		stop:        make(chan bool, 1),
		Logger:      logger.Named("rule-checker"),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go rc.run()
			return nil
		},
		// Questo blocco si occupa di fermare la goroutine che gestisce il controllo delle regole.
		// La logica è semplice:
		// 1) Quando l'applicazione si sta fermando, mandiamo un messaggio sul canale stop
		// 2) Attendiamo finché non arriva una risposta, sempre nel canale stop,
		//    o finché non andiamo in timeout (ossia arriva un messaggio sul canale ctx.Done())
		// 3) Terminiamo
		OnStop: func(ctx context.Context) error {
			rc.stop <- true
			select {
			case <-ctx.Done():
			case <-rc.stop:
			}
			return nil
		},
	})

	return rc
}

func (rc *RuleChecker) run() {
	var wait_seconds int
	if env_value, exist := os.LookupEnv("RULE_CHECKER_TIMER"); exist {
		var err error
		wait_seconds, err = strconv.Atoi(env_value)
		if err != nil {
			rc.Error("Errore nella conversione di RULE_CHECKER_TIMER: %v. Uso valore di default 45.", zap.Error(err))
			wait_seconds = DEFAULT_WAIT_SECONDS
		}
	} else {
		wait_seconds = DEFAULT_WAIT_SECONDS
	}
	ticker := time.NewTicker(time.Duration(wait_seconds) * time.Second)
	for {
		select {
		case <-rc.stop:
			rc.stop <- true
			return
		case <-ticker.C:
			rc.checkAllRules()
		}
	}
}

func (rc *RuleChecker) checkAllRules() {
	rc.Debug("Controllo periodico delle regole avviato")

	rules, err := rc.rulePort.ListQueryRules() // recupera tutte le regole dal repository in memoria
	if err != nil {
		rc.Error("Error while listing rules", zap.Error(err))
	}

	if len(rules) == 0 {
		rc.Debug("Nessuna regola trovata")
		return
	}

	// Per ogni regola, interroga Influx e confronta la quantity con la threshold
	for _, rule := range rules {
		rc.Debug("Controllo regola", zap.Any("rule", rule))

		// Esempio: se rule è un AddQueryRuleCmd con metodi GetGoodID, GetOperator e GetThreshold
		goodID := rule.GoodId
		operator := rule.Operator
		threshold := rule.Threshold

		// Invoca il metodo del service che interroga Influx
		resp := rc.queryPort.GetCurrentQuantityByGoodID(goodID)
		if err := resp.GetOperationResult(); err != nil {
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

		if condTrue {
			err := rc.publishPort.PublishStockAlert(types.StockAlertEvent{
				Id:              uuid.NewString(),
				Status:          "Pending",
				GoodID:          goodID,
				CurrentQuantity: currentQuantity,
				Operator:        operator,
				Threshold:       threshold,
				Timestamp:       time.Now().UnixMilli(),
			})
			if err != nil {
				rc.Error("Errore nell'invio della notifica", zap.String("goodId", goodID), zap.Error(err))
			}
		}
	}
}
