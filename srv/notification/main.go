package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/nats-io/nats.go"
)

type StockUpdate struct {
	WarehouseID string `json:"warehouse_id"`
	GoodID      string `json:"good_id"`
	Quantity    int    `json:"quantity"`
}

func main() {
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatalf("Errore connessione a NATS: %v", err)
	}
	defer nc.Close()

	
	token := "my-token"            
	url := "http://influxdb:8086"   
	org := "my-org"
	bucket := "stockdb"


	client := influxdb2.NewClientWithOptions(url, token,
		influxdb2.DefaultOptions().
			SetBatchSize(10).
			SetFlushInterval(1000), 
	)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(org, bucket)

	_, err = nc.Subscribe("stock.update.>", func(msg *nats.Msg) {
		var su StockUpdate
		if err := json.Unmarshal(msg.Data, &su); err != nil {
			log.Printf("Errore: %v\n", err)
			return
		}

		log.Printf("Ricevuto update da nats: %+v\n", su)

		p := write.NewPoint(
			"stock_measurement",
			map[string]string{
				"warehouse_id": su.WarehouseID,
				"good_id":      su.GoodID,
			},
			map[string]interface{}{
				"quantity": su.Quantity,
			},
			time.Now(),
		)

		err := writeAPI.WritePoint(context.Background(), p)
		if err != nil {
			log.Printf("Errore scrittura InfluxDB: %v\n", err)
		} else {
			log.Printf("Dato scritto su InfluxDB: %+v\n", su)
		}
	})
	if err != nil {
		log.Fatalf("Errore : %v", err)
	}

	log.Println("in ascolto su stock.update.> ...")
	select {}
}
