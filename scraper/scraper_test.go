package scraper

import (
	"testing"
)

// func TestScraper(t *testing.T) {
// 	// want := "Hello, world."
// 	fmt.Println("starting testing scraper")

// 	got := Scrape("https://youtube.com")
// 	if got == nil {
// 		t.Errorf("Scrape = %q, want %q", got, "not nil")
// 	}
// 	// log.Println("Data")
// 	// fmt.Println(got)

// }

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
