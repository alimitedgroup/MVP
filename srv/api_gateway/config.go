package main

import (
	"log"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type Config struct {
	// Configuration for the API
	API          APIConfig           `yaml:"api"`
	BrokerConfig broker.BrokerConfig `yaml:"broker"`
}

func loadConfig() fx.Option {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.AddConfigPath(".")

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
