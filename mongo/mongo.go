package mongo

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type PbTechItem struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Date  []string           `bson:"date"`
	Price []int              `bson:"price"`
}

func ProcessData(itemName string, price int) bool {
	item, _ := checkIfItemExists(itemName)
	if item != nil {
		return false
	}
	insertPbTechItem(itemName, price)
	return true
}

func checkIfItemExists(itemName string) (*PbTechItem, error) {
	coll := client.Database("pbtech_item").Collection("keyboards")
	filter := bson.D{{Key: "name", Value: itemName}}
	var result PbTechItem
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}
	if err != nil {
		panic(err)
	}
	return &result, nil
}

func insertPbTechItem(itemName string, price int) {
	coll := client.Database("pbtech_item").Collection("keyboards")
	newItem := PbTechItem{Name: itemName, Date: []string{getDateString()}, Price: []int{price}}
	_, err := coll.InsertOne(context.TODO(), newItem)
	if err != nil {
		panic(err)
	}
}

func ConnectToMongo() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
}

func getDateString() string {
	nzLocation, err := time.LoadLocation("Pacific/Auckland")
	if err != nil {
		fmt.Println("Error loading location:", err)
		panic(err)
	}

	nzTime := time.Now().In(nzLocation)

	day := nzTime.Day()
	month := nzTime.Month()
	year := nzTime.Year()

	return fmt.Sprintf("%02d-%02d-%d", day, month, year)
}
