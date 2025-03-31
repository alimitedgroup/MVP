package portout

import "github.com/alimitedgroup/MVP/srv/notification/types"

type StockEventPublisher interface {
	PublishStockAlert(alert types.StockAlertEvent) error
	RevokeStockAlert(alert types.StockAlertEvent) error
}
