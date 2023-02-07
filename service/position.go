package service

import (
	"context"
	"scalper/config"
	"scalper/customers"
	"scalper/log"
	"scalper/models"
	"scalper/utils"
	"time"

	"github.com/uncle-gua/gobinance/common"
	"github.com/uncle-gua/gobinance/futures"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var position = new(Position)

type Position struct {
	Hold  string  `bson:"hold"`
	Entry float64 `bson:"entry"`
	Peak  float64 `bson:"-"`
	Reach bool    `bson:"-"`
}

func (position *Position) init() {
	position.Hold = "NONE"
	position.Entry = 0.0
}

func (position *Position) Load() error {
	if err := models.PositionCollection.FindOne(
		context.TODO(),
		bson.M{},
	).Decode(position); err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
		position.init()
		if _, err := models.PositionCollection.InsertOne(
			context.TODO(),
			position,
		); err != nil {
			return err
		}
	}

	position.Peak = -1
	position.Reach = false

	log.Infof("hold: %s, entry: %.2f, peak: %.2f, reach: %v", position.Hold, position.Entry, position.Peak, position.Reach)
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
		position.Peak = -1
		position.Reach = false

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
				go func(amount string) {
					side := futures.SideTypeBuy
					retry := 1
					for {
						_, err := service.Symbol(config.Param.Symbol.Name).
							Side(side).
							PositionSide(futures.PositionSideType("SHORT")).
							Type(futures.OrderTypeMarket).
							Quantity(utils.Abs(amount)).
							Do(context.Background())
						if err == nil {
							return
						}
						log.Error(err)
						if retry > 2 {
							return
						}
						if err, ok := err.(*common.APIError); ok {
							if err.Code != -1001 {
								return
							}
						}
						retry++
						time.Sleep(3 * time.Millisecond)
					}
				}(customer.Position)

				go func() {
					if err := models.OrderCollection.FindOneAndUpdate(
						context.TODO(),
						bson.M{
							"customerId":   customer.ID,
							"positionSide": "SHORT",
							"closeTime":    0,
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
				}()
			}

			if positionSide == "SHORT" && hold == "LONG" {
				log.Infof("Reverse customer: %s, hold: %s, quantity: %s", customer.Name, hold, customer.Position)
				go func(amount string) {
					side := futures.SideTypeSell
					retry := 1
					for {
						_, err := service.Symbol(config.Param.Symbol.Name).
							Side(side).
							PositionSide(futures.PositionSideType("LONG")).
							Type(futures.OrderTypeMarket).
							Quantity(utils.Abs(amount)).
							Do(context.Background())
						if err == nil {
							return
						}
						log.Error(err)
						if retry > 2 {
							return
						}
						if err, ok := err.(*common.APIError); ok {
							if err.Code != -1001 {
								return
							}
						}
						retry++
						time.Sleep(3 * time.Millisecond)
					}
				}(customer.Position)

				go func() {
					if err := models.OrderCollection.FindOneAndUpdate(
						context.TODO(),
						bson.M{
							"customerId":   customer.ID,
							"positionSide": "LONG",
							"closeTime":    0,
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
				}()
			}

			side := futures.SideTypeBuy
			if positionSide == "SHORT" {
				side = futures.SideTypeSell
			}

			quantity := utils.FormatQuantity(customer.Capital / price)
			log.Infof("Open customer: %s, positionSide: %s, quantity: %s", customer.Name, positionSide, quantity)
			retry := 1
			for {
				_, err := service.
					Symbol(config.Param.Symbol.Name).
					Side(side).
					PositionSide(futures.PositionSideType(positionSide)).
					Type(futures.OrderTypeMarket).
					Quantity(quantity).
					Do(context.Background())
				if err == nil {
					break
				}
				log.Error(err)
				if retry > 2 {
					return
				}
				if err, ok := err.(*common.APIError); ok {
					if err.Code != -1001 {
						return
					}
				}
				retry++
				time.Sleep(3 * time.Millisecond)
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
	position.Peak = -1
	position.Reach = false

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
			retry := 1
			for {
				_, err := futures.NewClient(customer.ApiKey, customer.ApiSecret).
					NewCreateOrderService().
					Symbol(config.Param.Symbol.Name).
					Side(side).
					PositionSide(futures.PositionSideType(positionSide)).
					Type(futures.OrderTypeMarket).
					Quantity(quantity).
					Do(context.Background())
				if err == nil {
					break
				}
				log.Error(err)
				if retry > 2 {
					return
				}
				if err, ok := err.(*common.APIError); ok {
					if err.Code != -1001 {
						return
					}
				}
				retry++
				time.Sleep(3 * time.Millisecond)
			}

			customers.SetPosition(customer, "")

			if err := models.OrderCollection.FindOneAndUpdate(
				context.TODO(),
				bson.M{
					"customerId":   customer.ID,
					"positionSide": positionSide,
					"closeTime":    0,
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
