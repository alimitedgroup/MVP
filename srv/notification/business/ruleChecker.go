package business

import (
	"context"
	"log"
	"time"

	"go.uber.org/fx"
)

type RuleChecker struct {
	ticker  *time.Ticker
	service Business
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewRuleChecker(lc fx.Lifecycle, service Business) *RuleChecker {
	ctx, cancel := context.WithCancel(context.Background())

	rc := &RuleChecker{
		ticker:  time.NewTicker(5 * time.Second),
		service: service,
		ctx:     ctx,
		cancel:  cancel,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("RuleChecker avviato.")
			go rc.run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("RuleChecker fermato.")
			rc.stop()
			return nil
		},
	})

	return rc
}

func (rc *RuleChecker) run() {
	for {
		select {
		case <-rc.ctx.Done():
			return
		case <-rc.ticker.C:
			rc.checkAllRules()
		}
	}
}

func (rc *RuleChecker) checkAllRules() {
	log.Println("[RuleChecker] Controllo periodico delle regole avviato.")

	rules, err := rc.service.ListQueryRules() // recupera tutte le regole dal repository in memoria
	if err != nil {
		log.Println(err)
	}

	if len(rules) == 0 {
		log.Println("[RuleChecker] Nessuna regola trovata.")
		return
	}

	// Per ogni regola, interroga Influx e confronta la quantity con la threshold
	for _, rule := range rules {
		log.Printf("[RuleChecker] Controllo regola: %+v", rule)

		// Esempio: se rule è un AddQueryRuleCmd con metodi GetGoodID, GetOperator e GetThreshold
		goodID := rule.GoodId
		operator := rule.Operator
		threshold := rule.Threshold

		// Invoca il metodo del service che interroga Influx
		resp := rc.service.GetCurrentQuantityByGoodID(goodID)
		if err := resp.GetOperationResult(); err != nil {
			log.Printf("[RuleChecker] Errore nel recupero quantity per goodID=%s: %v", goodID, err)
			continue
		}

		currentQuantity := resp.CurrentQuantity // supponendo esista un metodo GetValue() o simile

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
			log.Printf("[RuleChecker] Operatore non valido: %s", operator)
			continue
		}

		if condTrue {
			log.Printf("[ALERT] good_id=%s quantity=%d %s %d",
				goodID, currentQuantity, operator, threshold)
			// INVIO DELLA NOTIFICA

		} else {
			log.Printf("Nessun alert: good_id=%s (quantity=%d, threshold=%d %s)",
				goodID, currentQuantity, threshold, operator)
		}
	}
}

func (rc *RuleChecker) stop() {
	rc.cancel()
	rc.ticker.Stop()
}
