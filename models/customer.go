package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"userId"`
	Name      string             `bson:"name"`
	ApiKey    string             `bson:"apiKey"`
	ApiSecret string             `bson:"apiSecret"`
	Capital   float64            `bson:"capital"`
	Position  string             `bson:"position"`
	Status    string             `bson:"status"`
}
