package latency

import (
	"time"
)

type Reporter func()

type Tracker interface {
	Track() Reporter
	Stop()
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
	start   chan startReq
	finish  chan ident
	stop    chan bool
	reports chan LatencyReport
	ticker  <-chan time.Time
}

func (t *tracker) Track() Reporter {
	req := startReq{resp: make(chan ident)}
	t.start <- req
	id := <-req.resp
	return func() { t.finish <- id }
}

func (t *tracker) Stop() {
	t.stop <- true
}

func roundLat(d, unit time.Duration) time.Duration {
	return ((d + unit/2) / unit) * unit
}

func (t *tracker) run() {
	var req startReq
	var id ident
	for {
		select {
		case req = <-t.start:
			id = t.count
			t.count++
			t.started[id] = time.Now().UTC()
			req.resp <- id
		case id = <-t.finish:
			lat := roundLat(time.Now().UTC().Sub(t.started[id]), 10*time.Millisecond)
			delete(t.started, id)
			t.lat[lat]++
		case <-t.ticker:
			report := t.lat
			t.lat = make(map[time.Duration]int)
			t.reports <- report
		case <-t.stop:
			return
		}
	}
}

func NewTracker(reports chan LatencyReport, ticker <-chan time.Time) Tracker {
	t := &tracker{
		started: make(map[ident]time.Time),
		lat:     make(map[time.Duration]int),
		start:   make(chan startReq),
		finish:  make(chan ident),
		stop:    make(chan bool),
		reports: reports,
		ticker:  ticker,
	}
	go t.run()
	return t
}
