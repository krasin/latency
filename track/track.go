package track

import (
	"encoding/json"
	"fmt"
	"github.com/krasin/latency"
	"log"
	"time"
)

var instanceId = time.Now().Unix()
var reports = make(chan latency.LatencyReport)
var span = 30 * time.Second
var ticker = time.NewTicker(span).C

var defaultTracker = latency.NewTracker(reports, ticker)

type LatencyReport struct {
	Instance  int64
	QPS       float64
	LatencyMs map[string]int
}

func logReports() {
	for report := range reports {
		lr := LatencyReport{
			Instance:  instanceId,
			LatencyMs: make(map[string]int),
		}
		total := 0
		for lat, count := range report {
			total += count
			lr.LatencyMs[fmt.Sprintf("%d", lat.Nanoseconds()/1000/1000)] = count
		}
		lr.QPS = float64((100*total)/int(span.Seconds())) / 100
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
