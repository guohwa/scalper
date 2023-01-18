package customers

import (
	"context"
	"scalper/models"

	"github.com/uncle-gua/log"
	"go.mongodb.org/mongo-driver/bson"
)

var data = make([]models.Customer, 0)

func init() {
	cursor, err := models.CustomerCollection.Find(
		context.TODO(),
		bson.M{},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(context.TODO(), &data); err != nil {
		log.Fatal(err)
	}
}

func Data() []models.Customer {
	return data
}

func Add(customer models.Customer) {
	data = append(data, customer)
}

func Set(customer models.Customer) {
	for i := 0; i < len(data); i++ {
		if data[i].ID == customer.ID {
			data[i].Name = customer.Name
			data[i].ApiKey = customer.ApiKey
			data[i].ApiSecret = customer.ApiSecret
			data[i].Capital = customer.Capital
			data[i].Status = customer.Status
		}
	}
}

func SetPosition(customer models.Customer, position string) {
	for i := 0; i < len(data); i++ {
		if data[i].ID == customer.ID {
			data[i].Position = position
		}
	}

	go func() {
		if _, err := models.CustomerCollection.UpdateByID(
			context.TODO(),
			customer.ID,
			bson.M{
				"$set": bson.M{
					"position": "",
				},
			},
		); err != nil {
			log.Error(err)
		}
	}()
}

func DisableByUser(user models.User) {
	for i := 0; i < len(data); i++ {
		if data[i].UserID == user.ID {
			data[i].Status = "Disable"
		}
	}
}

func EnableByUser(user models.User) {
	for i := 0; i < len(data); i++ {
		if data[i].UserID == user.ID {
			data[i].Status = "Enable"
		}
	}
}
