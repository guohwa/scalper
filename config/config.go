package config

import (
	"scalper/log"
	"scalper/models"

	"github.com/sirupsen/logrus"
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
		if err != mongo.ErrNoDocuments {
			log.Fatal(err)
		}

		App.Default()
		if err := App.Save(); err != nil {
			log.Fatal(err)
		}
	}
	level, err := logrus.ParseLevel(App.Level)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)

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
