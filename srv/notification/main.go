package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"
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

type QueryRule struct {
    GoodID    string `json:"good_id"`   
    Operator  string `json:"operator"`  // "<", ">", "<=", ">="
    Threshold int    `json:"threshold"` 
}

var queryRules []*QueryRule

const checkIntervalSeconds = 60

var (
    nc *nats.Conn 
)


func main() {

    natsURL := "nats://nats:4222" 
    influxURL := "http://influxdb:8086"

    influxToken := os.Getenv("INFLUXDB_TOKEN")
    if influxToken == "" {
        influxToken = "my-token"
    }

    influxOrg := os.Getenv("INFLUXDB_ORG")
    if influxOrg == "" {
        influxOrg = "my-org"
    }
    influxBucket := "stockdb"

    log.Printf("Avvio Notification con parametri:\n  NATS_URL=%s\n  InfluxURL=%s\n  InfluxToken=%s\n  InfluxOrg=%s\n  InfluxBucket=%s\n",
        natsURL, influxURL, influxToken, influxOrg, influxBucket,
    )

    var err error
    nc, err = nats.Connect(natsURL)
    if err != nil {
        log.Fatalf("Errore connessione a NATS: %v", err)
    }
    defer nc.Close()

    client := influxdb2.NewClientWithOptions(
        influxURL,
        influxToken,
        influxdb2.DefaultOptions().
            SetBatchSize(10).
            SetFlushInterval(1000),
    )
    defer client.Close()

    writeAPI := client.WriteAPIBlocking(influxOrg, influxBucket)

    
    _, err = nc.Subscribe("stock.update.>", func(msg *nats.Msg) {
        var su StockUpdate
        if e := json.Unmarshal(msg.Data, &su); e != nil {
            log.Printf("Errore decoding JSON: %v\n", e)
            return
        }

        log.Printf("Ricevuto update da nats: %+v\n", su)
        if len(su.Goods) == 0 {
            log.Println("Nessun good trovato nel messaggio.")
            return
        }

        good := su.Goods[0]

        p := influxdb2.NewPoint(
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
        log.Fatalf("Errore subscribe (stock.update): %v", err)
    }

    _, err = nc.Subscribe("notification.query.add", func(msg *nats.Msg) {
        var rule QueryRule
        if e := json.Unmarshal(msg.Data, &rule); e != nil {
            log.Printf("Errore parsing regola: %v\n", e)
            return
        }
        queryRules = append(queryRules, &rule)
        log.Printf("Aggiunta nuova regola: %+v\n", rule)
    })
    if err != nil {
        log.Fatalf("Errore subscribe (notification.query.add): %v", err)
    }
		log.Printf("DEBUG: Ora corrente dal container: %s", time.Now().Format(time.RFC3339))

   
    ticker := time.NewTicker(time.Duration(checkIntervalSeconds) * time.Second)
		go func() {
				for range ticker.C {
						// Log di debug
						log.Printf("DEBUG: ticker scattato, regole in memoria = %d", len(queryRules))
						checkAllRules(influxURL, influxToken, influxOrg)
				}
		}()
    log.Println("Servizio notification avviato. In ascolto su stock.update.> e notification.query.add ...")
    select {}
}

// =============================================
//             Utility
// =============================================

func checkAllRules(influxURL, influxToken, influxOrg string) {
    if len(queryRules) == 0 {
        return
    }
    client := influxdb2.NewClient(influxURL, influxToken)
    defer client.Close()
    queryAPI := client.QueryAPI(influxOrg)

    for _, rule := range queryRules {
        checkRule(queryAPI, rule)
    }
}

func checkRule(queryAPI api.QueryAPI, rule *QueryRule) {
		fluxQuery := fmt.Sprintf(`
				from(bucket: "stockdb")
					|> range(start: -7d)
					|> filter(fn: (r) => r["_measurement"] == "stock_measurement")
					|> filter(fn: (r) => r["good_id"] == "%s")
					|> filter(fn: (r) => r["_field"] == "quantity")
					|> yield(name: "mean")
		`, rule.GoodID)


    result, err := queryAPI.Query(context.Background(), fluxQuery)
    if err != nil {
        log.Printf("[checkRule] Errore query Influx per good_id=%s: %v\n", rule.GoodID, err)
        return
    }

		found := false

    for result.Next() {
			  found = true
        val, ok := result.Record().Value().(float64)
        if !ok {
            continue
        }
        currentQuantity := int(val)

        condTrue := false
        switch rule.Operator {
        case "<":
            condTrue = currentQuantity < rule.Threshold
        case ">":
            condTrue = currentQuantity > rule.Threshold
        case "<=":
            condTrue = currentQuantity <= rule.Threshold
        case ">=":
            condTrue = currentQuantity >= rule.Threshold
        default:
            log.Printf("[checkRule] Operatore non valido: %s\n", rule.Operator)
            continue
        }

				

        if condTrue {
            // se vera
            log.Printf("[ALERT] good_id=%s quantity=%d %s %d",
                rule.GoodID, currentQuantity, rule.Operator, rule.Threshold)

            //logica di notifica
        } else {
            log.Printf("Nessun alert: good_id=%s (quantity=%d, threshold=%d %s)",
                rule.GoodID, currentQuantity, rule.Threshold, rule.Operator)
        }
    }

		if !found {
				log.Printf("Nessun record trovato in Influx per good_id=%s", rule.GoodID)
		}

    if result.Err() != nil {
        log.Printf("[checkRule] Errore result: %v\n", result.Err())
    }
}
