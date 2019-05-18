package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/gocolly/colly"
)

func visitIndexPage(url string) {
	c := colly.NewCollector()

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

func indexPageLink(e *colly.HTMLElement) (match bool, err error) {
	link := e.Attr("href")
	return regexp.MatchString(`\A\/archive\?page=\d+\z`, link)
}

func episodeLink(e *colly.HTMLElement) (match bool, err error) {
	link := e.Attr("href")

	m, err := regexp.MatchString(`\A\/\d+\/\S+\z`, link)
	if m {
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
	}
	return m, err
}

func main() {
	i := 1
	baseUrl := "https://www.thisamericanlife.org/archive"
	url := baseUrl
	lastPage := false
	for !lastPage {
		lastPage = true
		c := colly.NewCollector()

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
		})

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			indexPage, err := indexPageLink(e)
			if err != nil {
				log.Fatal(err)
			}
			if indexPage {
				lastPage = false
			}

			episodeLink(e)

		})

		c.OnResponse(func(r *colly.Response) {
			url = fmt.Sprintf("%s?page=%d", baseUrl, i)
		})

		c.Visit(url)

		i++
	}
}
