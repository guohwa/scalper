package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID           primitive.ObjectID `bson:"_id"`
	CustomerID   primitive.ObjectID `bson:"customerId"`
	Symbol       string             `bson:"symbol"`
	Type         string             `bson:"type"`
	Side         string             `bson:"side"`
	PositionSide string             `bson:"positionSide"`
	Quantity     string             `bson:"quantity"`
	EntryPrice   string             `bson:"entryPrice"`
	EntryTime    int64              `bson:"entryTime"`
	ClosePrice   string             `bson:"closePrice"`
	CloseTime    int64              `bson:"closeTime"`
}
