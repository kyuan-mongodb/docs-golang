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

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your `MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// begin insert docs
	coll := client.Database("tea").Collection("ratings")

	docs := []interface{}{
		bson.D{{"name", "Chai"}, {"rating", 10}},
		bson.D{{"name", "Earl Grey"}, {"rating", 5}},
		bson.D{{"name", "Chai"}, {"rating", 7}},
		bson.D{{"name", "Earl Grey"}, {"rating", 9}},
	}
	// end insert docs
	result, err := coll.InsertMany(context.TODO(), docs)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d documents inserted with IDs:\n", len(result.InsertedIDs))
	for _, id := range result.InsertedIDs {
		fmt.Printf("\t%s\n", id)
	}

	// begin find docs
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"rating", bson.D{{"$gt", 5}}}},
				bson.D{{"rating", bson.D{{"$lt", 10}}}},
			}},
	}

	findCursor, findErr := coll.Find(context.TODO(), filter)
	if findErr != nil {
		panic(findErr)
	}

	var findResults []bson.D
	if findErr = findCursor.All(context.TODO(), &findResults); findErr != nil {
		panic(findErr)
	}
	for _, result := range findResults {
		fmt.Println(result)
	}
	// end find docs

	// begin aggregate docs
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$name"},
			{"average", bson.D{
				{"$avg", "$rating"},
			}},
		}}}

	aggCursor, aggErr := coll.Aggregate(context.TODO(), mongo.Pipeline{groupStage})
	if aggErr != nil {
		panic(aggErr)
	}

	var aggResults []bson.M
	if aggErr = aggCursor.All(context.TODO(), &aggResults); aggErr != nil {
		panic(aggErr)
	}
	for _, result := range aggResults {
		fmt.Printf("%v has an average rating of %v \n", result["_id"], result["average"])
	}
	// end aggregate docs
}
