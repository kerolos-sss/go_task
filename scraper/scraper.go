package scraper

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func Scrape(url string) map[string]bool {

	mainCollector := colly.NewCollector()
	aCollections := make(map[string]bool)

	mainCollector.OnHTML("body", func(bodyElement *colly.HTMLElement) {
		fmt.Println("got a body")
	})
	mainCollector.OnHTML("a", func(aElement *colly.HTMLElement) {
		aCollections[aElement.Attr("href")] = true
		fmt.Println(aCollections)
	})

	fmt.Println("starting to collect from: ")
	fmt.Println(url)
	mainCollector.Visit(url)
	// fmt.Println(aCollections)
	return aCollections

}

func DetectIsAccessible(url string) bool {
	collector := colly.NewCollector()
	var err *error
	collector.OnError(func(response *colly.Response, e error) {
		if response.StatusCode >= 400 {
			err = &e
		}
	})

	collector.Visit(url)
	return err == nil
}
