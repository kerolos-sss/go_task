package scraper

import (
	"fmt"
	"log"
	"testing"
)

func sestScraper(t *testing.T) {
	// want := "Hello, world."
	fmt.Println("starting testing scraper")

	urls, headings := Scrape("https://youtube.com")
	if urls == nil || headings == nil {
		t.Errorf("Scrape = %q, want %q", urls, "not nil")
	}

	log.Println("Data")
	log.Println(urls)
	log.Println(headings)

}

func sestIsAccessible(t *testing.T) {
	url1 := "https://google.com/dsgdfgds"
	accessible := DetectIsAccessible(url1)
	if accessible {
		t.Errorf("wrong accessability detection url: %q , should not be accessable", url1)
	}

	url2 := "https://google.com"
	accessible2 := DetectIsAccessible(url2)
	if accessible2 == false {
		t.Errorf("wrong accessability detection url: %q , should be accessable", url2)
	}

}

func TestGetPageDetails(t *testing.T) {
	url1 := "https://google.com/dsgdfgds"
	accessible, tags := GetPageDetails(url1)
	if len(accessible) > 0 || tags["h1"]+tags["h2"]+tags["h3"]+tags["h4"]+tags["h5"]+tags["h6"] > 0 {
		t.Errorf("wrong accessability detection url: %q , should not be accessable", url1)
		fmt.Println("NON ACCESSIBLE")
		fmt.Println(accessible)
		fmt.Println(tags)
	} else {
		fmt.Println("NON ACCESSIBLE")
		fmt.Println(accessible)
		fmt.Println(tags)
	}

	url2 := "https://accounts.google.com/AccountChooser?continue=https%3A%2F%2Fdocs.google.com%2Fforms%2Fcreate%3Fusp%3Dforms_alc&followup=https%3A%2F%2Fdocs.google.com%2Fforms%2Fcreate%3Fusp%3Dforms_alc&service=wise&ltmpl=forms"
	accessible2, tags2 := GetPageDetails(url2)
	if len(accessible2) > 0 || len(tags2) > 0 {
		fmt.Println("ACCESSIBLE")
		fmt.Println(accessible2)
		fmt.Println(tags2)

	} else {
		t.Errorf("wrong accessability detection url: %q , should not be accessable", url1)
		fmt.Println("ACCESSIBLE")
		fmt.Println(accessible2)
		fmt.Println(tags2)
	}

}
