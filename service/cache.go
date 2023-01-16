package service

import (
	"errors"
	"scalper/models"

	"github.com/uncle-gua/gobinance/futures"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrNoExists = errors.New("position does not exists")

type Cache struct {
	data map[primitive.ObjectID]models.Position
}

func (cache *Cache) Get(customerID primitive.ObjectID) models.Position {
	if position, ok := cache.data[customerID]; ok {
		return position
	}

	return models.Position{
		CustomerID: customerID,
		Hold:       "NONE",
		Order:      nil,
		Risk:       nil,
	}
}

func (cache *Cache) Set(customerID primitive.ObjectID, position models.Position) {
	cache.data[customerID] = position
}

func (cache Cache) SetHold(customerID primitive.ObjectID, hold string) {
	if position, ok := cache.data[customerID]; ok {
		position.Hold = hold
		cache.Set(customerID, position)
	} else {
		position := models.Position{
			CustomerID: customerID,
			Hold:       hold,
			Order:      nil,
			Risk:       nil,
		}
		cache.Set(customerID, position)
	}
}

func (cache Cache) SetOrder(customerID primitive.ObjectID, order *futures.CreateOrderResponse) {
	if position, ok := cache.data[customerID]; ok {
		position.Order = order
		cache.Set(customerID, position)
	} else {
		position := models.Position{
			CustomerID: customerID,
			Hold:       "NONE",
			Order:      order,
			Risk:       nil,
		}
		cache.Set(customerID, position)
	}
}

func (cache Cache) SetRisk(customerID primitive.ObjectID, risk *futures.PositionRisk) {
	if position, ok := cache.data[customerID]; ok {
		position.Risk = risk
		cache.Set(customerID, position)
	} else {
		position := models.Position{
			CustomerID: customerID,
			Hold:       "NONE",
			Order:      nil,
			Risk:       risk,
		}
		cache.Set(customerID, position)
	}
}
