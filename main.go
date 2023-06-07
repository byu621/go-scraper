package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/byu621/go-scraper/mongo"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

func convertToMoney(price int) string {
	dollar := float64(price) / 100.0
	dollarStr := strconv.FormatFloat(dollar, 'f', -1, 64)
	return fmt.Sprintf("$%s", dollarStr)
}

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
			isDbUpdated, prefix := mongo.ProcessData(name, priceInt)

			if isDbUpdated {
				line := fmt.Sprintf("%s: %4d: %8s %s", prefix, count, price, name)
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

func getKeyboardsPretty(c *gin.Context) {
	var lines []string
	lines = append(lines, fmt.Sprintf("Number of items: %d", mongo.GetPbTechItemsCount()))
	lines = append(lines, fmt.Sprintf("Number of items with more than one price: %d", mongo.GetPbTechItemsCountWithMoreThanOnePrice()))
	items := mongo.GetPbTechItemsWithMoreThanOnePrice()
	for _, result := range items {
		lines = append(lines, result.Name)
		for j, price := range result.Price {
			lines = append(lines, fmt.Sprintf("%d: %s %s", j, result.Date[j], convertToMoney(price)))
		}
	}
	c.IndentedJSON(http.StatusOK, lines)
}

func main() {
	mongo.ConnectToMongo()

	router := gin.Default()
	router.GET("/keyb", getKeyboards)
	router.GET("/keybpretty", getKeyboardsPretty)

	var port = envPortOr("8080")

	if os.Getenv("ENVIRONMENT") == "DEVELOPMENT" {
		router.Run(fmt.Sprintf("localhost%s", port))
	} else {
		router.Run(port)
	}
}

func envPortOr(port string) string {
	// If `PORT` variable in environment exists, return it
	if envPort := os.Getenv("PORT"); envPort != "" {
	  return ":" + envPort
	}
	// Otherwise, return the value of `port` variable from function argument
	return ":" + port
  }