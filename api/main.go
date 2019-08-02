package main

import (
	"context"
	"log"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	api := &SushiAPI{
		dbClient: createDatabase("mongodb://localhost:27017"),
		router:   mux.NewRouter(),
	}
	api.Start(":5000") //blocks
}

// should this be in it's own package???
func createDatabase(uri string) *mongo.Client {
	// setup client options
	clientOptions := options.Client().ApplyURI(uri)

	// connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
