package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Test random comment so that VSCode shuts up
type Test struct {
	String string
	Number int
	Bool   bool
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

	// Check the connection
	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("test").Collection("testCollection")

	// INSERT TEST
	insertData := Test{"String", 10, false}
	insertResult, err := collection.InsertOne(context.Background(), insertData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	// FIND TEST
	fmt.Println("Fetching all doccuments")
	// Pass these options to the Find method
	findOptions := options.Find()

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.Background(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.Background()) {
		// create a value into which the single document can be decoded
		var elem Test
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\t", elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.Background())

	// DELETE TEST
	deleteResult, err := collection.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the users collection\n", deleteResult.DeletedCount)

	// DISCONNECT TEST
	err = client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

	fmt.Print("Hold the console: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	fmt.Println(input.Text())
}
