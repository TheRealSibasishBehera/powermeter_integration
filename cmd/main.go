package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"log"
	//"os"
	//"os/exec"
	"time"
	//"log"
	"net/http"
	pm "powermeter_integration/pkg/powermeter/tasmotta-HLW8032"
)

var (
	address = flag.String(
		"address",
		"0.0.0.0:8881",
		"bind address",
	)
	metricsPath = flag.String(
		"metrics-path",
		"/metrics",
		"metrics path",
	)
)

func main() {

	flag.Parse()

	//register the collector
	err := prometheus.Register(version.NewCollector("powermeter_exporter"))
	if err != nil {
		log.Fatalf(
			"failed to register : %v",
			err,
		)
	}

	if err != nil {
		log.Fatalf(
			"failed to create collector: %v",
			err,
		)
	}
	//prometheus http handler
	go func() {
		http.Handle(
			*metricsPath,
			promhttp.Handler(),
		)
		http.HandleFunc(
			"/",
			func(w http.ResponseWriter, r *http.Request) {
				_, err = w.Write(
					[]byte(`<html>}
	fmt.Println("exporter call over")
}
			<head><title>PowerMeter Exporter</title></head>
			<body>
			<h1>Energy Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`),
				)
				if err != nil {
					log.Fatalf(
						"failed to write response: %v",
						err,
					)
				}
			},
		)

		err = http.ListenAndServe(
			*address,
			nil,
		)
		if err != nil {
			log.Fatalf(
				"failed to bind on %s: %v",
				*address,
				err,
			)
		}
		fmt.Println("exporter call over")
	}()

	_ = pm.GetSwitchStateJSON()
	if err != nil {
		log.Printf(
			"%v",
			err,
		)
	}
	//Prometheus Metrics using Gauge
	voltageCount := promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "voltage_count",
			Help: "voltage calibarated in powermeter",
		},
	)

	currentCount := promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "current_count",
			Help: "current calibarated in powermeter",
		},
	)

	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)
	for {
		go getEnergyFromPowerMeter(
			done,
			ticker,
			voltageCount,
			currentCount,
		)
		time.Sleep(5 * time.Second)
		done <- true
	}
}

func getEnergyFromPowerMeter(done chan bool, ticker *time.Ticker, voltageCount prometheus.Gauge, currentCount prometheus.Gauge) {
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			fmt.Println(
				"Tick at",
				t,
			)
			fmt.Println("command started")

			energyStats := pm.GetEnergyStats()

			voltageCount.Set(float64(energyStats.StatusSNS.Energy.Voltage))
			currentCount.Set(energyStats.StatusSNS.Energy.Current)

		}
	}

}
