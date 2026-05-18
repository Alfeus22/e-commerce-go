package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var UserCollection *mongo.Collection
var ProductCollection *mongo.Collection
var OrderCollection *mongo.Collection

func ConnectDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Berhasil connect")
	UserCollection = client.Database("e-commerce").Collection("users")
	ProductCollection = client.Database("e-commerce").Collection("products")
	OrderCollection = client.Database("e-commerce").Collection("orders")
	return client
}
