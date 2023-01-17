package service

import (
	"context"
	"scalper/config"
	"scalper/models"
	"scalper/utils"
	"time"

	"github.com/uncle-gua/gobinance/futures"
	"github.com/uncle-gua/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var position = Load()

type Position struct {
	customers []models.Customer `bson:"-"`
	Hold      string            `bson:"hold"`
	Entry     float64           `bson:"entry"`
	peak      float64           `bson:"-"`
	reach     bool              `bson:"-"`
}

func Default() *Position {
	return &Position{
		customers: make([]models.Customer, 0),
		Hold:      "NONE",
		Entry:     0.0,
		peak:      -1,
		reach:     false,
	}
}

func Load() *Position {
	filter := bson.M{}

	var position = new(Position)
	if err := models.PositionCollection.FindOne(
		context.TODO(),
		filter,
	).Decode(&position); err != nil {
		if err != mongo.ErrNoDocuments {
			log.Fatal(err)
		}
		position = Default()
		go func() {
			if _, err := models.PositionCollection.InsertOne(
				context.TODO(),
				position,
			); err != nil {
				log.Error(err)
			}
		}()
	}

	return position
}

func (position *Position) Save() error {
	filter := bson.M{}

	if _, err := models.PositionCollection.UpdateOne(
		context.TODO(),
		filter,
		position,
	); err != nil {
		return err
	}

	return nil
}

func (position *Position) Add(customer models.Customer) {
	position.customers = append(position.customers, customer)
}

func (position *Position) Remove(customer models.Customer) {
	customers := position.customers[:0]
	for _, c := range position.customers {
		if c.ID != customer.ID {
			customers = append(customers, c)
		}
	}
	position.customers = customers
}

func (position *Position) Open(positionSide string, price float64) {
	position.Entry = price
	position.Hold = positionSide
	go func() {
		if err := position.Save(); err != nil {
			log.Error(err)
		}
	}()

	for i, customer := range position.customers {
		go func(i int, customer models.Customer) {
			service := futures.NewClient(customer.ApiKey, customer.ApiSecret).NewCreateOrderService()

			if positionSide == "LONG" && position.Hold == "SHORT" {
				side := futures.SideTypeBuy
				_, err := service.Symbol(config.Param.Symbol.Name).
					Side(side).
					PositionSide(futures.PositionSideType("SHORT")).
					Type(futures.OrderTypeMarket).
					Quantity(customer.Position).
					Do(context.Background())
				if err != nil {
					log.Error(err)
				}
			}

			if positionSide == "SHORT" && position.Hold == "LONG" {
				side := futures.SideTypeSell
				_, err := service.Symbol(config.Param.Symbol.Name).
					Side(side).
					PositionSide(futures.PositionSideType("LONG")).
					Type(futures.OrderTypeMarket).
					Quantity(customer.Position).
					Do(context.Background())
				if err != nil {
					log.Error(err)
				}
			}

			side := futures.SideTypeBuy
			if positionSide == "SHORT" {
				side = futures.SideTypeSell
			}

			quantity := utils.FormatQuantity(customer.Capital / price)
			_, err := service.
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

			if positionSide == "SHORT" {
				quantity = "-" + quantity
			}
			position.customers[i].Position = quantity
			if _, err = models.CustomerCollection.UpdateByID(
				context.TODO(),
				customer.ID,
				bson.M{
					"position": quantity,
				},
			); err != nil {
				log.Error(err)
			}

			order := models.Order{
				ID:           primitive.NewObjectID(),
				CustomerID:   customer.ID,
				Symbol:       config.Param.Symbol.Name,
				Type:         "MARKET",
				Side:         string(side),
				PositionSide: positionSide,
				Quantity:     quantity,
				EntryPrice:   price,
				EntryTime:    time.Now().UnixMilli(),
				ClosePrice:   0.0,
				CloseTime:    0,
			}
			if _, err = models.OrderCollection.InsertOne(
				context.TODO(),
				order,
			); err != nil {
				log.Error(err)
			}
		}(i, customer)
	}
}

func (position *Position) Close(positionSide string, price float64) {
	position.Entry = 0.0
	position.Hold = "NONE"
	position.reach = false
	position.peak = -1

	go func() {
		if err := position.Save(); err != nil {
			log.Error(err)
		}
	}()

	for i, customer := range position.customers {
		go func(i int, customer models.Customer) {
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

			position.customers[i].Position = "0"
			if _, err = models.CustomerCollection.UpdateByID(
				context.TODO(),
				customer.ID,
				bson.M{
					"position": "0",
				},
			); err != nil {
				log.Error(err)
			}

			if err = models.OrderCollection.FindOneAndUpdate(
				context.TODO(),
				bson.M{
					"customerId": customer.ID,
					"closeTime":  0,
				},
				bson.M{
					"closePrice": price,
					"coseTime":   time.Now().UnixMilli(),
				},
			).Err(); err != nil {
				log.Error(err)
			}
		}(i, customer)
	}
}
