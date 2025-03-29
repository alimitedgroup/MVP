package portout

import "github.com/alimitedgroup/MVP/srv/notification/types"

type IStockEventPublisher interface {
	PublishStockAlert(alert types.StockAlertEvent) error
}
