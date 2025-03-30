package business

import (
	"context"
	"log"
	"time"

	"github.com/alimitedgroup/MVP/srv/notification/portin"
	"github.com/alimitedgroup/MVP/srv/notification/portout"
	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/google/uuid"

	"go.uber.org/fx"
)

type RuleChecker struct {
	rulePort    portin.QueryRules
	queryPort   portout.RuleQueryRepository
	publishPort portout.StockEventPublisher
	// stop è un canale su cui verranno mandati al massimo due messaggi.
	// Per la logica che ci sta dietro, fare riferimento al commento all'interno di NewRuleChecker.
	stop chan bool
}

func NewRuleChecker(lc fx.Lifecycle, rules portin.QueryRules, queries portout.RuleQueryRepository, publish portout.StockEventPublisher) *RuleChecker {
	rc := &RuleChecker{
		rulePort:    rules,
		queryPort:   queries,
		publishPort: publish,
		stop:        make(chan bool, 1),
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
	ticker := time.NewTicker(5 * time.Second)
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
	log.Println("[RuleChecker] Controllo periodico delle regole avviato.")

	rules, err := rc.rulePort.ListQueryRules() // recupera tutte le regole dal repository in memoria
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
		resp := rc.queryPort.GetCurrentQuantityByGoodID(goodID)
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
				log.Printf("[RuleChecker] Errore nell'invio della notifica: %v", err)
			}

		} else {
			log.Printf("Nessun alert: good_id=%s (quantity=%d, threshold=%d %s)",
				goodID, currentQuantity, threshold, operator)
		}
	}
}
