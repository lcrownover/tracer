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

// getTracemeValue is a function that will scrape the traceme app and return
// the value as a float64, since that is what the GaugeFunc expects
func getTracemeValue() float64 {
	log.Println("getting traceme value")
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
	return float64(val)
}

// cache is a type that will hold the cached values
// type cache map[string]float64

// startCaching initiates all the goroutines
// and return the cache (which is a reference type, no pointer needed)
// func startCaching() cache {
//     c := make(map[string]float64)
//     go scrapeTracemeValue(c, 5)
//     // add other funcs here to run on a schedule
//     // just make sure they use unique keys
//     return c
// }

// scrapeTracemeValueOnSchedule is a function that will run in a goroutine and will
// set the cached value.
// you can create more of these functions that accept the cache and an interval
// and add them to the startCaching function
// func scrapeTracemeValue(c cache, interval int) {
//     for {
//         log.Println("scraping traceme value")
//         v := getTracemeValue()
//         c["traceme_value"] = v
//         time.Sleep(time.Duration(interval) * time.Second)
//     }
// }

func main() {
	// Create a new gauge.
	// A gauge is a metric that represents a single numerical value that can arbitrarily go up and down.
	tracemeGauge := promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "traceme_value",
			Help: "Value of the traceme app",
		},
		getTracemeValue,
	)
	// recordTracemeMetrics(tracemeGauge, 5)

    // this is an example of cached metrics
    // you would have multiple NewGaugeFunc's that use the same cache
    // and read different keys from the cache
    // cache := startCaching()
    // cachedTracemeGauge := promauto.NewGaugeFunc(
    //     prometheus.GaugeOpts{
    //         Name: "cached_traceme_value",
    //         Help: "Value of the traceme app",
    //     },
    //     func() float64 {
    //         return cache["traceme_value"]
    //     },
    // )

	// set up prometheus registry
	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(tracemeGauge)
	// promRegistry.MustRegister(cachedTracemeGauge) // register the cache gauge

	// set up prometheus handler, disabling all the built-in metrics
	promHandler := promhttp.HandlerFor(
		promRegistry,
		promhttp.HandlerOpts{
			EnableOpenMetrics: false,
		},
	)

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promHandler)
	log.Println("starting server on port 80")
	http.ListenAndServe(":80", nil)
}
