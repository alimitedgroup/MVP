package config

// func TestLoadConfig(t *testing.T) {
// 	ctx := t.Context()
// 	config := LoadConfig()

// 	app := fx.New(
// 		config,
// 		fx.Invoke(func(cfg *Config) {
// 			assert.Equal(t, "4", cfg.WarehouseConfig.ID)
// 		}),
// 	)

// 	err := app.Start(ctx)
// 	if err != nil {
// 		t.Errorf("error starting app: %v", err)
// 	}

// 	defer func() {
// 		err := app.Stop(ctx)
// 		if err != nil {
// 			t.Errorf("error stopping app: %v", err)
// 		}
// 	}()
// }
