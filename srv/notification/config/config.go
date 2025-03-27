package config

import (
	"log"
	"strings"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type Config struct {
	BrokerConfig broker.BrokerConfig `mapstructure:"broker"`
}

func LoadConfig() fx.Option {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("broker.url", nats.DefaultURL)
	viper.SetDefault("influxdb.url", "http://influxdb:8086")
	viper.SetDefault("influxdb.token", "my-token")
	viper.SetDefault("influxdb.org", "my-org")
	viper.SetDefault("influxdb.bucket", "stockdb")

	viper.SetTypeByDefaultValue(true)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("fatal error config file: %w", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("fatal error unmarshalling config: %w", err)
	}

	return fx.Options(
		fx.Supply(&config),
		fx.Supply(&config.BrokerConfig),
	)
}
