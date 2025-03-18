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
    ID          string `json:"id"`
    WarehouseID string `json:"warehouse_id"`
    Type        string `json:"type"`
    Goods       []struct {
        ID       string `json:"id"`
        Quantity int    `json:"quantity"`
        Delta    int    `json:"delta"`
    } `json:"goods"`
    OrderID    string `json:"order_id"`
    TransferID string `json:"transfer_id"`
    Timestamp  int64  `json:"timestamp"`
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
            log.Printf("Errore decoding JSON: %v\n", err)
            return
        }

        log.Printf("Ricevuto update da nats: %+v\n", su)

        if len(su.Goods) == 0 {
            log.Println("Nessun good trovato nel messaggio.")
            return
        }
        good := su.Goods[0]

        p := write.NewPoint(
            "stock_measurement",
            map[string]string{
                "warehouse_id": su.WarehouseID,
                "good_id":      good.ID,
            },
            map[string]interface{}{
                "quantity": good.Quantity,
                "delta":    good.Delta,
                "type":     su.Type,
            },
            time.Now(),
        )

        if wErr := writeAPI.WritePoint(context.Background(), p); wErr != nil {
            log.Printf("Errore scrittura InfluxDB: %v\n", wErr)
        } else {
            log.Printf("Dato scritto su InfluxDB: WarehouseID=%s GoodID=%s Quantity=%d Delta=%d\n",
                su.WarehouseID, good.ID, good.Quantity, good.Delta)
        }
    })
    if err != nil {
        log.Fatalf("Errore subscribe: %v", err)
    }

    
    log.Println("in ascolto su stock.update.> ...")
    select {}
}
