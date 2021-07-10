package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// MongoMetaData random comment so that VSCode shuts up

func main() {
	url := "mongodb://localhost:27017"

	clientOptions := options.Client().ApplyURI(url)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	type MongoMetaData struct {
		Version string `json:"version"`
		Set     string `json:"set"`
	}

	mongoMetaData := MongoMetaData{}
	client.Database("test").RunCommand(nil, bsonx.Doc{{"buildInfo", bsonx.Int32(1)}}).Decode(&mongoMetaData)
	fmt.Println(mongoMetaData.Version)

	client.Database("admin").RunCommand(nil, bsonx.Doc{{"replSetGetStatus", bsonx.Int32(1)}}).Decode(&mongoMetaData)
	fmt.Println(mongoMetaData.Set)

}
