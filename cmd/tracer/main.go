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

func recordMetrics() {
	go func() {
		for {
            v := getTracemeValue()
            trace
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	tracemeGuage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "traceme_value",
		Help: "Value of the traceme app",
	})
)

func main() {
	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

