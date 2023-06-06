package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
			name := productLinks[0]
			price := e.ChildText(".ginc .full-price")

			priceNoDollarSign := price[1:]
			priceNoDecimal := strings.ReplaceAll(priceNoDollarSign, ".", "")
			priceInt, _ := strconv.Atoi(priceNoDecimal)
			isDbUpdated := mongo.ProcessData(name, priceInt)

			if isDbUpdated {
				line := fmt.Sprintf("%4d: %8s %s", count, price, name)
				lines = append(lines, line)
				count++
			}
		}
	})

	co.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	co.OnError(func(r *colly.Response, e error) {
		fmt.Println("error:", e, r.Request.URL, string(r.Body))
	})

	co.Visit("https://www.pbtech.co.nz/category/peripherals/keyboards/gaming-keyboards")
	if len(lines) > 0 {
		c.IndentedJSON(http.StatusOK, lines)
	} else {
		c.IndentedJSON(http.StatusOK, "Status: No Updates :(")
	}
}

func main() {
	mongo.ConnectToMongo()

	router := gin.Default()
	router.GET("/keyb", getKeyboards)

	router.Run("localhost:8080")
}
