package config

import (
	"fmt"
	"os"
)

type WarehouseConfig struct {
	ID string
}

func ConfigFromEnv() (*WarehouseConfig, error) {
	var ok bool
	cfg := &WarehouseConfig{}

	cfg.ID, ok = os.LookupEnv("ENV_SERVICE_ID")
	if !ok {
		return nil, fmt.Errorf("ENV_SERVICE_ID environment variable not set")
	}
	return cfg, nil
}
