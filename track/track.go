package track

import (
	"encoding/json"
	"fmt"
	"github.com/krasin/latency"
	"log"
	"time"
)

var reports = make(chan latency.LatencyReport)
var span = 30 * time.Second
var ticker = time.NewTicker(span).C

var defaultTracker = latency.NewTracker(reports, ticker)

type LatencyReport struct {
	QPS       float64
	LatencyMs map[string]int
}

func logReports() {
	for report := range reports {
		lr := LatencyReport{
			QPS:       float64(100*len(report)/int(span.Seconds())) / 100,
			LatencyMs: make(map[string]int),
		}
		for lat, count := range report {
			lr.LatencyMs[fmt.Sprintf("%d", lat.Nanoseconds()/1000/1000)] = count
		}
		data, err := json.Marshal(lr)
		if err != nil {
			// TODO: do not emit incorrect json in case if error contains `"`.
			log.Printf(`Latency report: { "error": "%v" }`, err)
			continue
		}
		log.Printf("Latency report: %v", string(data))
	}
}

func init() {
	go logReports()
}

func Track() latency.Reporter {
	return defaultTracker.Track()
}
