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

type pAccessable struct {
	count      int
	internal   bool
	accessable bool
}
type pATTERS struct {
	htmlVersion string
	title       string
	hasLogin    bool
}
type Counts struct {
	HtmlVersion string `json:"htmlVersion"`
	Title       string `json:"title"`
	H1          int    `json:"h1"`
	H2          int    `json:"h2"`
	H3          int    `json:"h3"`
	H4          int    `json:"h4"`
	H5          int    `json:"h5"`
	H6          int    `json:"h6"`
	Internal    int    `json:"internal"`
	External    int    `json:"external"`
	Inacessable int    `json:"inacessable"`
	HasLogin    bool   `json:"hasLogin"`
}

/**
* Testing document
* returns (<a> tags Href, h# tags count)
 */
func Scrape(pageURL string) (map[string]int, map[string]int, pATTERS) {

	// I chose to count using a dict with UUID to avoid count issues with concurrency
	// I thought of regular int inside the dict but I don't know how it would behave in racing
	// I thought of slices they did not seem a good choice at all

	mainCollector := colly.NewCollector()
	aCollections := make(map[string]bool)

	separatorUUID := uuid.NewString()
	otherAtters := pATTERS{
		htmlVersion: "1",
		title:       "",
		hasLogin:    false,
	}
	re := regexp.MustCompile(`<[!]{0,1}DOCTYPE.*>`)
	mainCollector.OnScraped(func(r *colly.Response) {
		found := re.Find(r.Body)
		if len(found) > 0 {
			otherAtters.htmlVersion = string(found)
		}
	})

	mainCollector.OnHTML("title", func(titleElement *colly.HTMLElement) {
		otherAtters.title = titleElement.Text
	})

	passwordFields := make(map[string]bool)
	hasNonPasswordNonHiddenFields := false
	mainCollector.OnHTML("form input", func(inputElement *colly.HTMLElement) {
		if inputElement.Attr("type") == "password" {
			passwordFields[uuid.NewString()] = true
		} else if inputElement.Attr("type") != "hidden" {
			hasNonPasswordNonHiddenFields = true
		}
	})
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

	otherAtters.hasLogin = hasNonPasswordNonHiddenFields && 2 > len(passwordFields) && len(passwordFields) > 0
	return urls, tagsCount, otherAtters
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

func GetPageDetails(pageUrl string) (map[string]pAccessable, map[string]int, pATTERS) {

	// parse only base url
	base, err := url.Parse(pageUrl)
	if err != nil {
		log.Fatal(err)
	}

	urls, tagsCount, otherAtters := Scrape(pageUrl)
	processedURLs := make(map[string]int)

	for key, count := range urls {
		// and then use it to parse relative URLs
		u, err := base.Parse(key)
		if err != nil {
			log.Fatal(err)
		}

		processedURLs[u.String()] += count
	}

	accessable := make(map[string]pAccessable)

	var group sync.WaitGroup

	for key, count := range processedURLs {
		// and then use it to parse relative URLs
		group.Add(1)
		// fmt.Println("worker for: " + key)
		go func(key string, count int) {

			// fmt.Println("worker for: " + strconv.Itoa(count) + "   key :" + key)
			u, _ := url.Parse(key)
			accessable[key] = pAccessable{
				count:      count,
				accessable: DetectIsAccessible(key),
				internal:   u.Host == base.Host,
			}
			group.Done()
		}(key, count)
	}

	group.Wait()

	return accessable, tagsCount, otherAtters
}

func GetPageDetailsAndCounts(pageUrl string) Counts {
	fmt.Println("SCRAPPING #####")
	fmt.Println(pageUrl)
	accessable, tagsCount, otherAtters := GetPageDetails(pageUrl)

	internalCount := 0
	externalCount := 0
	inAccessibleCount := 0
	for _, props := range accessable {
		if !props.accessable {
			inAccessibleCount += props.count
		}
		if props.internal {
			internalCount += props.count
		} else {
			externalCount += props.count
		}
	}
	return Counts{
		HtmlVersion: otherAtters.htmlVersion,
		Title:       otherAtters.title,
		H1:          tagsCount["h1"],
		H2:          tagsCount["h2"],
		H3:          tagsCount["h3"],
		H4:          tagsCount["h4"],
		H5:          tagsCount["h5"],
		H6:          tagsCount["h6"],
		Internal:    internalCount,
		External:    externalCount,
		Inacessable: inAccessibleCount,
		HasLogin:    otherAtters.hasLogin,
	}
}
