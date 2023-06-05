package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func getKeyboards(c *gin.Context) {
	count := 1
	var lines []string

	co := colly.NewCollector()

	co.OnHTML(".row", func(e *colly.HTMLElement) {
		productLinks := e.ChildAttrs(".js-product-link", "title")
		if len(productLinks) == 1 {
			price := e.ChildText(".ginc .full-price")
			line := fmt.Sprintf("%4d: %8s %s", count, price, productLinks[0])
			fmt.Println(line)
			lines = append(lines, line)
			count++
		}
	})

	co.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	co.OnError(func(r *colly.Response, e error) {
		fmt.Println("error:", e, r.Request.URL, string(r.Body))
	})

	co.Visit("https://www.pbtech.co.nz/category/peripherals/keyboards/gaming-keyboards")
	c.IndentedJSON(http.StatusOK, lines)
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
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
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/keyb", getKeyboards)

	router.Run("localhost:8080")

}
