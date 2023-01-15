package config

import (
	"context"
	"scalper/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Param *models.Param = new(models.Param)
	App   *models.App   = new(models.App)
)

func Load() error {
	filter := bson.M{}
	if err := models.ParamCollection.FindOne(
		context.Background(),
		filter,
	).Decode(Param); err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
		Param.Default()
	}

	return nil
}

func Save() error {
	if Param.ID.IsZero() {
		Param.ID = primitive.NewObjectID()
		if _, err := models.ParamCollection.InsertOne(
			context.TODO(),
			Param,
		); err != nil {
			return err
		}
	} else {
		if _, err := models.ParamCollection.UpdateByID(
			context.TODO(),
			Param.ID,
			bson.M{"$set": Param},
		); err != nil {
			return err
		}
	}

	return nil
}
