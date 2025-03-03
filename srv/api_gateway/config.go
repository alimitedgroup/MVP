package main

import (
	"log"
	"strings"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type Config struct {
	// Configuration for the API
	API          APIConfig           `mapstructure:"api"`
	BrokerConfig broker.BrokerConfig `mapstructure:"broker"`
}

func loadConfig() fx.Option {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// BUG: need to fix the default values or else viper doesn't the env vars inside the config variable if the config file is not present
	viper.SetDefault("api.host", "")
	viper.SetDefault("api.port", "8080")
	viper.SetDefault("broker.url", nats.DefaultURL)

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
		fx.Supply(&config.API),
	)
}
