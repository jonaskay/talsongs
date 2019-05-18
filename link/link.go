package link

import (
	"regexp"

	"github.com/gocolly/colly"
)

func match(pattern string, e *colly.HTMLElement) (path string, err error) {
	path = e.Attr("href")

	m, err := regexp.MatchString(`\A\/archive\?page=\d+\z`, path)
	if !m {
		return "", err
	}
	return path, err
}

func Index(e *colly.HTMLElement) (path string, err error) {
	return match(`\A\/archive\?page=\d+\z`, e)
}

func Episode(e *colly.HTMLElement) (path string, err error) {
	return match(`\A\/\d+\/\S+\z`, e)
}
