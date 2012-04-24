package latency

import (
	"testing"
)

func TestOneRequest(t *testing.T) {
	reports := make(chan LatencyReport)
	lat := NewTracker(reports)
	track := lat.Track()
	track()
	report := <-reports
	if len(report) != 1 {
		t.Fatalf("Expected: 1 element in the report, has: %d elements", len(report))
	}
	report = <-reports
	if len(report) != 0 {
		t.Fatalf("Expected: 1 element in the report, has: %d elements", len(report))
	}
}
