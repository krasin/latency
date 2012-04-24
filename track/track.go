package track

import (
	"github.com/krasin/latency"
	"log"
	"time"
)

var reports = make(chan latency.LatencyReport)
var ticker = time.NewTicker(30 * time.Second).C

var defaultTracker = latency.NewTracker(reports, ticker)

func logReports() {
	for report := range reports {
		log.Printf("Latency report: %v", report)
	}
}

func init() {
	go logReports()
}

func Track() latency.Reporter {
	return defaultTracker.Track()
}
