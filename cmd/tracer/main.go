package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getTracemeValue() int {
	resp, err := http.Get("http://traceme.westus2.azurecontainer.io:3333")
	if err != nil {
		log.Println("error getting data from traceme: ", err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response from traceme: ", err)
	}
	val, err := strconv.Atoi(string(b))
	if err != nil {
		log.Println("error parsing int from traceme: ", err)
	}
	return val
}

// recordTracemeMetrics is a function that will run in a goroutine and will
// periodically scrape the traceme app and set the value of the gauge
func recordTracemeMetrics(g prometheus.Gauge, interval int) {
	go func() {
		for {
			log.Println("scraping traceme value")
			v := getTracemeValue()
			g.Set(float64(v))
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}()
}

func main() {
	// Create a new gauge.
	// A gauge is a metric that represents a single numerical value that can arbitrarily go up and down.
	tracemeGuage := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "traceme_value",
		Help: "Value of the traceme app",
	})
	recordTracemeMetrics(tracemeGuage, 5)

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.Handler())
	log.Println("starting server on port 80")
	http.ListenAndServe(":80", nil)
}
