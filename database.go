package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func connect() *mongo.Client {
	// Create client
	dbPassword := os.Getenv("DB_PASSWORD")
	url := "mongodb+srv://admin:" + dbPassword + "@cluster0.wfckv.mongodb.net/" + databaseName + "?retryWrites=true&w=majority"
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Println(err)
	}

	// Create connect
	err = client.Connect(context.TODO())
	if err != nil {
		log.Println(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

func disconnect(client *mongo.Client) {
	err := client.Disconnect(context.TODO())

	if err != nil {
		log.Println(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
