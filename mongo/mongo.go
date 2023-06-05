package mongo

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Restaurant struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string
	RestaurantId string `bson:"restaurant_id"`
	Cuisine      string
	Address      interface{}
	Borough      string
	Grades       []interface{}
}

func PingMongo() {
	executeMongo(func(client *mongo.Client) {
		fmt.Println("func")
	})
}

func GetData() {
	executeMongo(func(client *mongo.Client) {
		coll := client.Database("sample_restaurants").Collection("restaurants")
		filter := bson.D{{"name", "Bagels N Buns"}}
		var result Restaurant
		err := coll.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// This error means your query did not match any documents.
				return
			}
			panic(err)
		}

		fmt.Println(result.Name)
	})
}

func executeMongo(f func(*mongo.Client)) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	f(client)

	if err = client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
	// fmt.Println("Successfully disconnected")
}
