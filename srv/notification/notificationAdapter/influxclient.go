package notificationAdapter

import (
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func NewInfluxClient() influxdb2.Client {
	url := os.Getenv("INFLUXDB_URL")
	token := os.Getenv("INFLUXDB_TOKEN")
	if url == "" {
		url = "http://influxdb:8086"
	}
	if token == "" {
		token = "my-token"
	}

	return influxdb2.NewClient(url, token)
}
