package main

import (
	"scalper/position"
	"time"

	"github.com/uncle-gua/log"
)

func main() {
	// web.Start()
	position.Close("LONG", 1735.00)
	log.Info(position.Entry)
	log.Info(position.Hold)
	time.Sleep(10 * time.Second)
}
