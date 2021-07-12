package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Test random comment so that VSCode shuts up
type Execution struct {
	ID                    string         `json:"_id" bson:"_id"`
	SubInstructionID      string         `json:"subInstructionId,omitempty" bson:"subInstructionId,omitempty"`
	Priority              int            `json:"priority,omitempty" bson:"priority,omitempty"`
	Amount                float64        `json:"amount,omitempty" bson:"amount,omitempty"`
	DealID                string         `json:"dealId" bson:"dealId,omitempty"`
	OriginallyScheduledOn *CustomDateObj `json:"originallyScheduledOn" bson:"originallyScheduledOn"`
}

type CustomDateObj struct {
	Tz     string    `json:"tz,omitempty" bson:"tz,omitempty"`
	TzInfo string    `json:"tzInfo,omitempty" bson:"tzInfo,omitempty"`
	Utc    time.Time `json:"utc,omitempty" bson:"utc,omitempty"`
	Epoch  int64     `json:"epoch,omitempty" bson:"epoch,omitempty"`
}

func main() {
	url := "mongodb://localhost:27017"

	dealId := "PT-DEAL-1626018451387-1"

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
	findOptions.SetProjection(bson.D{{"subInstructionId", 1}, {"dealId", 1}})

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.Background(), bson.D{{"dealId", dealId}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	fmt.Println("Inserting into another collection")
	var elem Execution
	for cur.Next(context.Background()) {
		// create a value into which the single document can be decoded
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(elem.OriginallyScheduledOn)
		collection2.InsertOne(context.Background(), elem)
		fmt.Print("#")
	}
	fmt.Println()

	pipeline := bson.D{
		{"$match", bson.D{{"dealId", dealId}}},
	}

	aggregateOptions := options.Aggregate()

	fmt.Println("Updating documents in new collection")
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
			fmt.Print("#")
		}
	}
	fmt.Println()

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.Background())
}
