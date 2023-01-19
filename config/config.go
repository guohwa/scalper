package config

import (
	"scalper/models"

	"github.com/uncle-gua/log"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Param *models.Param = new(models.Param)
	App   *models.App   = new(models.App)
)

func Default() {
	App.Default()
	Param.Default()
}

func init() {
	if err := App.Load(); err != nil {
		if err == mongo.ErrNoDocuments {
			App.Default()
			if err := App.Save(); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}

		level, err := log.ParseLevel(App.Level)
		if err != nil {
			log.Fatal(err)
		}
		log.SetLevel(level)
	}

	if err := Param.Load(); err != nil {
		if err == mongo.ErrNoDocuments {
			Param.Default()
			if err := Param.Save(); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}
}
