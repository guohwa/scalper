package service

import (
	"context"
	"scalper/config"
	"scalper/customers"
	"scalper/log"
	"scalper/models"
	"scalper/utils"
	"time"

	"github.com/uncle-gua/gobinance/futures"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var position = new(Position)

type Position struct {
	Hold  string
	Entry float64
	Peak  float64
	Reach bool
}

func (position *Position) Default() *Position {
	return &Position{
		Hold:  "NONE",
		Entry: 0.0,
		Peak:  -1,
		Reach: false,
	}
}

func (position *Position) Load() error {
	if err := models.PositionCollection.FindOne(
		context.TODO(),
		bson.M{},
	).Decode(position); err != nil {
		if err != mongo.ErrNoDocuments {
			position = position.Default()
		}
		if _, err := models.PositionCollection.InsertOne(
			context.TODO(),
			position,
		); err != nil {
			return err
		}
	}

	log.Infof("hold: %s, entry: %.2f", position.Hold, position.Entry)
	return nil
}

func (position *Position) Save() error {
	filter := bson.M{}
	if _, err := models.PositionCollection.UpdateOne(
		context.TODO(),
		filter,
		bson.M{
			"$set": bson.M{
				"hold":  position.Hold,
				"entry": position.Entry,
			},
		},
	); err != nil {
		return err
	}

	return nil
}

func (position *Position) Open(positionSide string, price float64) {
	defer func() {
		position.Hold = positionSide
		position.Entry = price

		go func() {
			if err := position.Save(); err != nil {
				log.Error(err)
			}
		}()
	}()

	for i, customer := range customers.Data() {
		if customer.Status == "Disable" {
			continue
		}
		go func(i int, customer models.Customer, hold string) {
			service := futures.NewClient(customer.ApiKey, customer.ApiSecret).NewCreateOrderService()

			if positionSide == "LONG" && hold == "SHORT" {
				log.Infof("Reverse customer: %s, hold: %s, quantity: %s", customer.Name, hold, customer.Position)
				side := futures.SideTypeBuy
				_, err := service.Symbol(config.Param.Symbol.Name).
					Side(side).
					PositionSide(futures.PositionSideType("SHORT")).
					Type(futures.OrderTypeMarket).
					Quantity(utils.Abs(customer.Position)).
					Do(context.Background())
				if err != nil {
					log.Error(err)
				}
			}

			if positionSide == "SHORT" && hold == "LONG" {
				log.Infof("Reverse customer: %s, hold: %s, quantity: %s", customer.Name, hold, customer.Position)
				side := futures.SideTypeSell
				_, err := service.Symbol(config.Param.Symbol.Name).
					Side(side).
					PositionSide(futures.PositionSideType("LONG")).
					Type(futures.OrderTypeMarket).
					Quantity(utils.Abs(customer.Position)).
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
			log.Infof("Open customer: %s, positionSide: %s, quantity: %s", customer.Name, positionSide, quantity)
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

			customers.SetPosition(customer,
				func(positionSide, quantity string) string {
					if positionSide == "SHORT" {
						return "-" + quantity
					}
					return quantity
				}(positionSide, quantity),
			)

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
			if _, err := models.OrderCollection.InsertOne(
				context.TODO(),
				order,
			); err != nil {
				log.Error(err)
			}
		}(i, customer, position.Hold)
	}
}

func (position *Position) Close(positionSide string, price float64) {
	position.Entry = 0.0
	position.Hold = "NONE"

	go func() {
		if err := position.Save(); err != nil {
			log.Error(err)
		}
	}()

	for i, customer := range customers.Data() {
		go func(i int, customer models.Customer) {
			side := futures.SideTypeBuy
			if positionSide == "LONG" {
				side = futures.SideTypeSell
			}

			log.Infof("Close customer: %s, positionSide: %s, quantity: %s", customer.Name, positionSide, customer.Position)
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

			customers.SetPosition(customer, "")

			if err := models.OrderCollection.FindOneAndUpdate(
				context.TODO(),
				bson.M{
					"customerId": customer.ID,
					"closeTime":  0,
				},
				bson.M{
					"$set": bson.M{
						"closePrice": price,
						"closeTime":  time.Now().UnixMilli(),
					},
				},
				options.FindOneAndUpdate().SetSort(bson.M{
					"entryTime": -1,
				}),
			).Err(); err != nil {
				log.Error(err)
			}
		}(i, customer)
	}
}