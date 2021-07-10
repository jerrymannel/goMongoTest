package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Test random comment so that VSCode shuts up
type Execution struct {
	ID      string `bson:"_id"`
	GroupID string `bson:"groupId,omitempty"`
	DealID  string `bson:"dealId,omitempty"`
}

func main() {
	url := "mongodb://localhost:27017"

	if os.Getenv("URL") != "" {
		fmt.Println("Using the URL :: ", os.Getenv("URL"))
		url = os.Getenv("URL")
	} else {
		fmt.Println("Using the default URL :: ", url)
	}

	clientOptions := options.Client().ApplyURI(url)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("xcro").Collection("execution.executions")
	collection2 := client.Database("xcro").Collection("execution.ecomm")

	// Pass these options to the Find method
	findOptions := options.Find()

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.Background(), bson.D{{"dealId", "PT-DEAL-1625934489122-1"}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	var elem Execution
	for cur.Next(context.Background()) {
		// create a value into which the single document can be decoded
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println("\t", elem)
		insertResult, err := collection2.InsertOne(context.Background(), elem)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted from exec to ecomm: ", insertResult.InsertedID)
	}

	pipeline := bson.D{
		{"$match", bson.D{{"dealId", "PT-DEAL-1625934489122-1"}}},
	}

	aggregateOptions := options.Aggregate()

	aggregateCursor, err := collection2.Aggregate(context.Background(), mongo.Pipeline{pipeline}, aggregateOptions)
	for aggregateCursor.Next(context.Background()) {
		err := aggregateCursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		result, err := collection2.UpdateOne(context.Background(), bson.D{{"_id", elem.ID}}, bson.D{{"$set", bson.D{{"groupId", elem.DealID + " ## " + elem.ID}}}})
		if err != nil {
			log.Fatal(err)
		}
		if result.MatchedCount != 0 {
			fmt.Println("Matched and updated an existing document", elem.ID)
		}
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.Background())
}
