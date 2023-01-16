package service

import (
	"context"
	"scalper/config"
	"scalper/models"
	"scalper/utils"

	"github.com/uncle-gua/gobinance/futures"
	"github.com/uncle-gua/log"
)

var position = &Position{
	customers: make([]models.Customer, 0),
	hold:      "NONE",
	entry:     0.0,
	peak:      -1,
	reach:     false,
}

type Position struct {
	customers []models.Customer `bson:"-"`
	hold      string            `bson:"hold"`
	entry     float64           `bson:"entry"`
	peak      float64           `bson:"-"`
	reach     bool              `bson:"-"`
}

func (cache *Position) Save() error {
	return nil
}

func (cache *Position) Add(customer models.Customer) {
	cache.customers = append(cache.customers, customer)
}

func (cache *Position) Remove(customer models.Customer) {
	customers := cache.customers[:0]
	for _, c := range cache.customers {
		if c.ID != customer.ID {
			customers = append(customers, c)
		}
	}
	cache.customers = customers
}

func (cache *Position) Open(positionSide string, price float64) {
	if positionSide == "LONG" && cache.hold == "SHORT" {
		cache.Close("SHORT", price)
	}
	if positionSide == "SHORT" && cache.hold == "LONG" {
		cache.Close("LONG", price)
	}

	for _, customer := range cache.customers {
		go func(customer models.Customer) {
			side := futures.SideTypeBuy
			if positionSide == "SHORT" {
				side = futures.SideTypeSell
			}

			quantity := utils.FormatQuantity(customer.Capital / price)
			_, err := futures.NewClient(customer.ApiKey, customer.ApiSecret).
				NewCreateOrderService().
				Symbol(config.Param.Symbol.Name).
				Side(side).
				PositionSide(futures.PositionSideType(positionSide)).
				Type(futures.OrderTypeMarket).
				Quantity(quantity).
				Do(context.Background())
			if err != nil {
				log.Error(err)
				return
			}
		}(customer)
	}
}

func (cache *Position) Close(positionSide string, price float64) {
	for _, customer := range cache.customers {
		go func(customer models.Customer) {
			side := futures.SideTypeBuy
			if positionSide == "LONG" {
				side = futures.SideTypeSell
			}

			quantity := utils.Abs(customer.Position)
			_, err := futures.NewClient(customer.ApiKey, customer.ApiSecret).
				NewCreateOrderService().
				Symbol(config.Param.Symbol.Name).
				Side(side).
				PositionSide(futures.PositionSideType(positionSide)).
				Type(futures.OrderTypeMarket).
				Quantity(quantity).
				Do(context.Background())
			if err != nil {
				log.Error(err)
			}
		}(customer)
	}
}
