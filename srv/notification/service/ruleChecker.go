package service

import (
	"context"
	"log"
	"time"

	"go.uber.org/fx"
)


type RuleChecker struct {
	ticker  *time.Ticker
	service IService
	ctx     context.Context
	cancel  context.CancelFunc
}


func NewRuleChecker(lc fx.Lifecycle, service IService) *RuleChecker {
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

	rules := rc.service.GetAllQueryRules()

	if len(rules) == 0 {
		log.Println("[RuleChecker] Nessuna regola trovata.")
		return
	}

	for _, rule := range rules {
		log.Printf("[RuleChecker] Controllo regola: %+v", rule)
		//logica verifica...
	}
}

func (rc *RuleChecker) stop() {
	rc.cancel()
	rc.ticker.Stop()
}
