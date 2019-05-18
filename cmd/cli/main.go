package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/gocolly/colly"
	"github.com/jonaskay/talsongs/episodes"
	"github.com/jonaskay/talsongs/link"
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

	c.Visit(url)
}

func main() {
	var paths episodes.Episodes
	baseUrl := "https://www.thisamericanlife.org"

	fmt.Println("Fetching episode links...")
	i := 1
	archive := fmt.Sprintf("%s/archive", baseUrl)
	url := archive
	lastPage := false
	for !lastPage {
		lastPage = true
		c := colly.NewCollector()

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			indexLink, err := link.Index(e)
			if err != nil {
				log.Fatal(err)
			}
			if indexLink != "" {
				lastPage = false
			}

			episodeLink, err := link.Episode(e)
			if err != nil {
				log.Fatal(err)
			}
			if episodeLink != "" {
				paths = append(paths, episodeLink)
			}
		})

		c.OnResponse(func(r *colly.Response) {
			url = fmt.Sprintf("%s?page=%d", archive, i)
		})

		c.Visit(url)

		i++
	}
	fmt.Println("Finished fetching episode links!")

	fmt.Println("Fetching songs...")
	paths = paths.Unique()
	for i := 0; i < len(paths); i++ {
		url := baseUrl + paths[i]
		c := colly.NewCollector()

		c.OnHTML("div.field-name-field-song", func(e *colly.HTMLElement) {
			e.ForEach("div", func(i int, e *colly.HTMLElement) {
				if e.Name == "div" {
					e.ForEach("div", func(j int, e *colly.HTMLElement) {
						var url string
						e.ForEach("a[href]", func(k int, e *colly.HTMLElement) {
							url = e.Attr("href")
						})

						fmt.Printf("%s %s\n", e.Text, url)
					})
				}
			})
		})

		c.Visit(url)
	}
	fmt.Println("Finished fetchings songs!")

}
