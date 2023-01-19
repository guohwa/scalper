package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type App struct {
	ID     primitive.ObjectID `bson:"_id"`
	Title  string             `json:"title"`
	Mode   string             `json:"mode"`
	Level  string             `json:"level"`
	Listen string             `json:"listen"`
	Trust  []string           `json:"trust"`
}

func (app *App) Default() {
	app.ID = primitive.NewObjectID()
	app.Title = "Scalper"
	app.Mode = "debug"
	app.Listen = ":8082"
	app.Level = "error"
	app.Trust = []string{
		"127.0.0.1",
	}
}

func (app *App) Load() error {
	filter := bson.M{}

	if err := AppCollection.FindOne(
		context.TODO(),
		filter,
	).Decode(app); err != nil {
		return err
	}

	return nil
}

func (app *App) Save() error {
	if _, err := AppCollection.InsertOne(
		context.TODO(),
		app,
	); err != nil {
		return err
	}

	return nil
}

func (app *App) Update() error {
	if _, err := AppCollection.UpdateByID(
		context.TODO(),
		app.ID,
		bson.M{"$set": app},
	); err != nil {
		return err
	}

	return nil
}
