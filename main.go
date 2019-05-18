package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()
	url := "https://www.thisamericanlife.org/archive"

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		matched, err := regexp.MatchString(`\A\/\d+\/\S+\z`, link)
		if err != nil {
			log.Fatal(err)
		}
		if matched {
			fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(url)
}
