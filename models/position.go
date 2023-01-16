package models

import (
	"github.com/uncle-gua/gobinance/futures"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Position struct {
	CustomerID primitive.ObjectID           `bson:"customerId"`
	Hold       string                       `bson:"hold"`
	Order      *futures.CreateOrderResponse `bson:"-"`
	Risk       *futures.PositionRisk        `bson:"-"`
}
