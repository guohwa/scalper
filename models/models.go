package models

import (
	"context"
	"time"

	"scalper/utils"

	"github.com/uncle-gua/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database = struct {
	Host string
	Name string
}{
	Host: "mongodb://localhost:27017",
	Name: "scalper",
}

var Client *mongo.Client

var (
	AppCollection      *mongo.Collection
	UserCollection     *mongo.Collection
	OrderCollection    *mongo.Collection
	ParamCollection    *mongo.Collection
	SessionCollection  *mongo.Collection
	CustomerCollection *mongo.Collection
)

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(database.Host))
	if err != nil {
		log.Fatal(err)
	}

	AppCollection = Client.Database(database.Name).Collection("app")
	UserCollection = Client.Database(database.Name).Collection("users")
	OrderCollection = Client.Database(database.Name).Collection("orders")
	ParamCollection = Client.Database(database.Name).Collection("param")
	SessionCollection = Client.Database(database.Name).Collection("sessions")
	CustomerCollection = Client.Database(database.Name).Collection("customers")

	res := UserCollection.FindOne(context.Background(), bson.M{"username": "admin"})
	if err = res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			user := User{
				ID:       primitive.NewObjectID(),
				Username: "admin",
				Password: utils.Encrypt("admin"),
				Role:     "Admin",
				Status:   "Enable",
			}
			result, err := UserCollection.InsertOne(context.Background(), &user)
			if err != nil {
				log.Fatal(err)
			}
			if id, ok := result.InsertedID.(primitive.ObjectID); ok {
				log.Infof("system init, admin id: %s", id.Hex())
			} else {
				log.Fatal("system init failed")
			}
		} else {
			log.Fatal(err)
		}
	}
}
