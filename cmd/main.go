package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"log"
	"os"
	"os/exec"
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

	//Prometheus Metrics using Gauge
	voltage_count := promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "voltage_count",
			Help: "voltage calibarated in powermeter",
		},
	)

	current_count := promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "current_count",
			Help: "current calibarated in powermeter",
		},
	)

	_, err := pm.SwitchPowerMeter()
	if err != nil {
		log.Printf("%v",err)
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	ticker := time.NewTicker(5 * time.Millisecond)
	done := make(chan bool)
	for {
		go getEnergyFromPowerMeter(
			done,
			ticker,
		)
		time.Sleep(20 * time.Second)
	}

}

func getEnergyFromPowerMeter(done chan bool, ticker *time.Ticker) {
	select {
	case <-done:
		return
	case t := <-ticker.C:
		fmt.Println(
			"Tick at",
			t,
		)
		fmt.Println("command started")

		// parse_csv_and_publish(path)
		energyStats := pm.GetEnergyStats()

		//publish
		////Fetch wakeup data
		ptWakeupCount.Set(sysInfo.Wakeups)

		////Fetch cpuUsage data
		ptCpuUsageCount.Set(sysInfo.CpuUsage)

		////Fetch baseLine power
		ptBaselinePowerCount.Set(baseLinePower)

		//Fetch no of tunables
		ptTuCount.Set(float64(tunNum))
	}
}
}
