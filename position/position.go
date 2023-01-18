package position

import (
	"context"
	"scalper/config"
	"scalper/customers"
	"scalper/models"
	"scalper/utils"
	"time"

	"github.com/uncle-gua/gobinance/futures"
	"github.com/uncle-gua/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Hold  = "NONE"
	Entry = 0.0
	Peak  = -1.0
	Reach = false
)

type position struct {
	Hold  string  `bson:"hold"`
	Entry float64 `bson:"entry"`
}

func Default() *position {
	return &position{
		Hold:  "NONE",
		Entry: 0.0,
	}
}

func init() {
	var position = new(position)
	if err := models.PositionCollection.FindOne(
		context.TODO(),
		bson.M{},
	).Decode(&position); err != nil {
		if err != mongo.ErrNoDocuments {
			log.Fatal(err)
		}
		position = Default()
		if _, err := models.PositionCollection.InsertOne(
			context.TODO(),
			position,
		); err != nil {
			log.Error(err)
		}
	}
	Hold = position.Hold
	Entry = position.Entry
}

func save() error {
	filter := bson.M{}
	if _, err := models.PositionCollection.UpdateOne(
		context.TODO(),
		filter,
		bson.M{
			"$set": bson.M{
				"hold":  Hold,
				"entry": Entry,
			},
		},
	); err != nil {
		return err
	}

	return nil
}

func Open(positionSide string, price float64) {
	Entry = price
	Hold = positionSide
	go func() {
		if err := save(); err != nil {
			log.Error(err)
		}
	}()

	for i, customer := range customers.Data() {
		if customer.Status == "Disable" {
			continue
		}
		go func(i int, customer models.Customer) {
			service := futures.NewClient(customer.ApiKey, customer.ApiSecret).NewCreateOrderService()

			if positionSide == "LONG" && Hold == "SHORT" {
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

			if positionSide == "SHORT" && Hold == "LONG" {
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

			customers.SetPosition(customer, quantity)

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
		}(i, customer)
	}
}

func Close(positionSide string, price float64) {
	Entry = 0.0
	Hold = "NONE"

	go func() {
		if err := save(); err != nil {
			log.Error(err)
		}
	}()

	for i, customer := range customers.Data() {
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
