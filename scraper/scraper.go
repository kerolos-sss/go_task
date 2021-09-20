package scraper

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/google/uuid"
)

func Scrape(url string) ([]string, map[string]int) {

	mainCollector := colly.NewCollector()
	aCollections := make(map[string]bool)

	mainCollector.OnHTML("body", func(bodyElement *colly.HTMLElement) {
		fmt.Println("got a body")
	})
	mainCollector.OnHTML("a", func(aElement *colly.HTMLElement) {
		aCollections[aElement.Attr("href")] = true
		fmt.Println(aCollections)
	})
	tagsToCount := []string{"h1", "h2", "h3", "h4", "h5", "h6"}
	// I chose to count using a dict with UUID to avoid count issues with concurrency
	// I thought of regular int inside the dict but I don't know how it would behave in racing
	// I thought of slices they did not seem a good choice at all
	tagsToCountDicts := make(map[string]map[string]bool)
	for _, tag := range tagsToCount {
		tagsToCountDicts[tag] = make(map[string]bool)
		mainCollector.OnHTML("a", func(aElement *colly.HTMLElement) {
			tagsToCountDicts[tag][uuid.NewString()] = true
		})
	}

	fmt.Println("starting to collect from: ")
	fmt.Println(url)
	mainCollector.Visit(url)
	// fmt.Println(aCollections)
	tagsCount := make(map[string]int)
	for tag, dict := range tagsToCountDicts {
		tagsCount[tag] = len(dict)
	}
	urls := make([]string, len(aCollections))
	index := 0
	for key, _ := range aCollections {
		urls[index] = key
		index++
	}

	return urls, tagsCount

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
