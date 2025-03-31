package adapterout

import (
	"github.com/alimitedgroup/MVP/srv/notification/config"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func NewInfluxClient(cfg config.NotificationConfig) influxdb2.Client {
	return influxdb2.NewClient(cfg.InfluxUrl, cfg.InfluxToken)
}
