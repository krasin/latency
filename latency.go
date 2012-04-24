package latency

import (
	"time"
)

type Reporter func()

type Tracker interface {
	Track() Reporter
}

type ident int

type startReq struct {
	resp chan ident
}

type LatencyReport map[time.Duration]int

type tracker struct {
	count   ident
	started map[ident]time.Time
	lat     LatencyReport
	stop    chan ident
	start   chan startReq
	reports chan LatencyReport
}

func (t *tracker) Track() Reporter {
	req := startReq{resp: make(chan ident)}
	t.start <- req
	id := <-req.resp
	return func() { t.stop <- id }
}

func (t *tracker) run() {
	ticker := time.NewTicker(30 * time.Second)
	var req startReq
	var id ident
	for {
		select {
		case req = <-t.start:
			id = t.count
			t.count++
			t.started[id] = time.Now().UTC()
			req.resp <- id
		case id = <-t.stop:
			lat := time.Now().UTC().Sub(t.started[id])
			delete(t.started, id)
			t.lat[lat]++
		case <-ticker.C:
			report := t.lat
			t.lat = make(map[time.Duration]int)
			t.reports <- report
		}
	}
}

func NewTracker(reports chan LatencyReport) Tracker {
	t := &tracker{
		started: make(map[ident]time.Time),
		lat:     make(map[time.Duration]int),
		stop:    make(chan ident),
		start:   make(chan startReq),
		reports: reports,
	}
	go t.run()
	return t
}
