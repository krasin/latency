package latency

import (
	"testing"
	"time"
)

func TestOneRequest(t *testing.T) {
	reports := make(chan LatencyReport)
	ticker := make(chan time.Time)

	lat := NewTracker(reports, ticker)
	defer lat.Stop()
	track := lat.Track()
	track()
	ticker <- time.Now().UTC()
	report := <-reports
	if len(report) != 1 {
		t.Fatalf("Expected: 1 element in the report, has: %d elements", len(report))
	}
	ticker <- time.Now().UTC()
	report = <-reports
	if len(report) != 0 {
		t.Fatalf("Expected: 1 element in the report, has: %d elements", len(report))
	}
}
