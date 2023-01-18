package config

import (
	"scalper/models"

	"github.com/uncle-gua/log"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	App   *models.App   = new(models.App)
	Param *models.Param = new(models.Param)
)

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
