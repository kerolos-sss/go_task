package scraper

import (
	"fmt"
	"log"
	"testing"
)

func TestScraper(t *testing.T) {
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

func TestIsAccessible(t *testing.T) {
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
