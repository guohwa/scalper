package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	ApiKey    string             `bson:"apiKey"`
	ApiSecret string             `bson:"apiSecret"`
	Role      string             `bson:"role"`
	Service   string             `bson:"service"`
	Status    string             `bson:"status"`
}
