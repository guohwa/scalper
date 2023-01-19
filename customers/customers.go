package customers

import (
	"context"

	"scalper/log"
	"scalper/models"

	"go.mongodb.org/mongo-driver/bson"
)

var customers = make([]models.Customer, 0)

func init() {
	cursor, err := models.CustomerCollection.Find(
		context.TODO(),
		bson.M{},
	)
	if err != nil {
		log.Fatal(err)
	}

	var items []models.Customer
	if err = cursor.All(context.TODO(), &items); err != nil {
		log.Fatal(err)
	}
	for _, customer := range items {
		Add(customer)
	}
}

func Data() []models.Customer {
	return customers
}

func Add(customer models.Customer) {
	customers = append(customers, customer)
	log.Tracef("customer name: %s, position: %s", customer.Name, customer.Position)
}

func Set(customer models.Customer) {
	for i := 0; i < len(customers); i++ {
		if customers[i].ID == customer.ID {
			customers[i].Name = customer.Name
			customers[i].ApiKey = customer.ApiKey
			customers[i].ApiSecret = customer.ApiSecret
			customers[i].Capital = customer.Capital
			customers[i].Status = customer.Status
		}
	}
}

func SetPosition(customer models.Customer, position string) {
	for i := 0; i < len(customers); i++ {
		if customers[i].ID == customer.ID {
			customers[i].Position = position
		}
	}

	go func() {
		if _, err := models.CustomerCollection.UpdateByID(
			context.TODO(),
			customer.ID,
			bson.M{
				"$set": bson.M{
					"position": position,
				},
			},
		); err != nil {
			log.Error(err)
		}
	}()
}
