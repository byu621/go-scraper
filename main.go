package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	count := 1

	c := colly.NewCollector()

	c.OnHTML(".row", func(e *colly.HTMLElement) {
		productLinks := e.ChildAttrs(".js-product-link", "title")
		if len(productLinks) == 1 {
			price := e.ChildText(".ginc .full-price")
			fmt.Printf("%4d: %8s %s\n", count, price, productLinks[0])
			count++
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("error:", e, r.Request.URL, string(r.Body))
	})

	c.Visit("https://www.pbtech.co.nz/category/peripherals/keyboards/gaming-keyboards")
}
