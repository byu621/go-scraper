package main

import (
	"fmt"
	"net/http"

	"github.com/byu621/go-scraper/mongo"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)





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



func main() {
	mongo.PingMongo()
	mongo.GetData()

	router := gin.Default()
	router.GET("/keyb", getKeyboards)

	router.Run("localhost:8080")
}
