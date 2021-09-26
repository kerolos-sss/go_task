package scraper

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"

	"net/url"

	"github.com/gocolly/colly/v2"
	"github.com/google/uuid"
)

type Accessable struct {
	count      int
	internal   bool
	accessable bool
}

type Counts struct {
	htmlVersion string
	title       string
	h1          int
	h2          int
	h3          int
	h4          int
	h5          int
	h6          int
	internal    int
	external    int
	inacessable int
	hasLogin    bool
}

/**
* Testing document
* returns (<a> tags Href, h# tags count)
 */
func Scrape(pageURL string) (map[string]int, map[string]int, string) {

	// I chose to count using a dict with UUID to avoid count issues with concurrency
	// I thought of regular int inside the dict but I don't know how it would behave in racing
	// I thought of slices they did not seem a good choice at all

	mainCollector := colly.NewCollector()
	aCollections := make(map[string]bool)

	separatorUUID := uuid.NewString()
	doctype := "1"
	re := regexp.MustCompile(`<[!]{0,1}DOCTYPE.*>`)
	mainCollector.OnScraped(func(r *colly.Response) {
		fmt.Println("got a response")
		fmt.Println(string(r.Body[:100]))
		found := re.Find(r.Body)
		if len(found) > 0 {
			doctype = string(found)
		}

		fmt.Println(doctype)
		fmt.Println(doctype)
		fmt.Println(doctype)
		fmt.Println(doctype)

	})
	// mainCollector.OnHTML("body", func(bodyElement *colly.HTMLElement) {
	// 	fmt.Println("got a body")
	// 	for i, node := range bodyElement.DOM.Parent().Nodes {
	// 		fmt.Println("got a DOCTYPE : index: " + strconv.Itoa(i) + " type: ")
	// 		fmt.Println(node.Data)
	// 		fmt.Println(node.Attr)
	// 		fmt.Println(node.Namespace)
	// 		fmt.Println(*node.Parent)
	// 		fmt.Println("####")
	// 	}
	// })
	// mainCollector.OnHTML("!DOCTYPE", func(element *colly.HTMLElement) {
	// })

	mainCollector.OnHTML("a", func(aElement *colly.HTMLElement) {
		uuidStr := uuid.NewString()
		aCollections[uuidStr+separatorUUID+aElement.Attr("href")] = true
		// fmt.Println(aCollections)
	})
	tagsToCount := []string{"h1", "h2", "h3", "h4", "h5", "h6"}

	tagsToCountDicts := make(map[string]map[string]bool)
	for _, tag := range tagsToCount {
		tagsToCountDicts[tag] = make(map[string]bool)
		mainCollector.OnHTML(tag, func(aElement *colly.HTMLElement) {
			tagsToCountDicts[tag][uuid.NewString()] = true
		})
	}

	// fmt.Println("starting to collect from: ")
	// fmt.Println(pageURL)
	mainCollector.Visit(pageURL)
	// fmt.Println(aCollections)
	tagsCount := make(map[string]int)
	for tag, dict := range tagsToCountDicts {
		tagsCount[tag] = len(dict)
	}
	urls := make(map[string]int)
	for key, _ := range aCollections {
		url := strings.Split(key, separatorUUID)[1]
		urls[url] += 1
	}

	return urls, tagsCount, doctype
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

func GetPageDetails(pageUrl string) (map[string]Accessable, map[string]int, string) {

	// parse only base url
	base, err := url.Parse(pageUrl)
	if err != nil {
		log.Fatal(err)
	}

	urls, tagsCount, doctype := Scrape(pageUrl)
	processedURLs := make(map[string]int)

	for key, count := range urls {
		// and then use it to parse relative URLs
		u, err := base.Parse(key)
		if err != nil {
			log.Fatal(err)
		}

		processedURLs[u.String()] += count
	}

	accessable := make(map[string]Accessable)

	var group sync.WaitGroup

	for key, count := range processedURLs {
		// and then use it to parse relative URLs
		group.Add(1)
		// fmt.Println("worker for: " + key)
		go func(key string, count int) {

			// fmt.Println("worker for: " + strconv.Itoa(count) + "   key :" + key)
			u, _ := url.Parse(key)
			accessable[key] = Accessable{
				count:      count,
				accessable: DetectIsAccessible(key),
				internal:   u.Host == base.Host,
			}
			group.Done()
		}(key, count)
	}

	group.Wait()

	return accessable, tagsCount, doctype
}

// func GetPageDetails(pageUrl string) Counts {
// 	accessable, tagsCount = GetPageDetails(pageUrl)
// 	return Counts{}
// }
