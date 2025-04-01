package config

import (
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"time"
)

type NotificationConfig struct {
	ServiceId    string
	InfluxUrl    string
	InfluxToken  string
	InfluxOrg    string
	InfluxBucket string
	CheckerTimer time.Duration
}

func requireEnv(name string, logger *zap.Logger) string {
	val, ok := os.LookupEnv(name)
	if !ok {
		logger.Fatal(fmt.Sprintf("%s environment variable not set", name))
	}
	return val
}

func FromEnv(logger *zap.Logger) *NotificationConfig {
	var cfg NotificationConfig

	cfg.ServiceId = requireEnv("ENV_SERVICE_ID", logger)
	cfg.InfluxUrl = requireEnv("INFLUXDB_URL", logger)
	cfg.InfluxToken = requireEnv("INFLUXDB_TOKEN", logger)
	cfg.InfluxOrg = requireEnv("INFLUXDB_ORG", logger)
	cfg.InfluxBucket = requireEnv("INFLUXDB_BUCKET", logger)

	var err error
	timer := requireEnv("RULE_CHECKER_TIMER", logger)
	cfg.CheckerTimer, err = time.ParseDuration(timer)
	if err != nil {
		logger.Fatal("invalid RULE_CHECKER_TIMER", zap.Error(err))
	}

	return &cfg
}

var Module = fx.Provide(FromEnv)
