package main

import (
	//	"fmt"
	"github.com/krasin/latency/track"
	"time"
)

func req(d time.Duration) {
	defer track.Track()()
	time.Sleep(d)
}

func main() {
	for {
		req(1 * time.Second)
		req(1 * time.Second)
		req(10 * time.Second)
		req(2 * time.Second)
	}

}
