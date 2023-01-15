package main

import (
	"scalper/config"

	"github.com/uncle-gua/log"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}

	log.Info(config.Param)

	if err := config.Save(); err != nil {
		log.Fatal(err)
	}
}
