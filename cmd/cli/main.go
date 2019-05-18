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

func indexPageLink(e *colly.HTMLElement) (link string, err error) {
	path := e.Attr("href")

	m, err := regexp.MatchString(`\A\/archive\?page=\d+\z`, path)
	if !m {
		return "", err
	} else {
		return path, err
	}
}

func episodePageLink(e *colly.HTMLElement) (link string, err error) {
	path := e.Attr("href")

	m, err := regexp.MatchString(`\A\/\d+\/\S+\z`, path)
	if !m {
		return "", err
	}
	fmt.Printf("Link found: %q -> %s\n", e.Text, path)
	return path, err
}

func main() {
	var episodes []string

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
			indexLink, err := indexPageLink(e)
			if err != nil {
				log.Fatal(err)
			}
			if indexLink != "" {
				lastPage = false
			}

			episodeLink, err := episodePageLink(e)
			if err != nil {
				log.Fatal(err)
			}
			if episodeLink != "" {
				episodes = append(episodes, episodeLink)
			}
		})

		c.OnResponse(func(r *colly.Response) {
			url = fmt.Sprintf("%s?page=%d", baseUrl, i)
		})

		c.Visit(url)

		i++
	}
}
