package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	count := 1

	c := colly.NewCollector()
	c.OnHTML(".titleline", func(e *colly.HTMLElement) {
		fmt.Printf("%4d: %s\n", count, e.Text)
		count++
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("error:", e, r.Request.URL, string(r.Body))
	})
	c.Visit("https://news.ycombinator.com/news")
}
