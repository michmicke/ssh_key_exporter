package main

import (
	"log"
	"net/http"
	"time"

	"github.com/michmicke/ssh_key_exporter/internal/config"
	"github.com/michmicke/ssh_key_exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	reg := prometheus.NewRegistry()
	m := metrics.NewMetrics(reg)

	log.Printf("config %s", config.GetConfig())

	go func() {
		for {
			c := make(chan bool)
			go m.WatchAuthorizedKeys(c)
			time.Sleep(config.GetConfig().PollingInterval)
			close(c)
			log.Print("Refreshing...")
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
