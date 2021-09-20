package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"

	"os"
	"project/scraper"
)

func main() {
	fmt.Println("hi I am a func")

	baseUrl := "https://youtube.com"

	scraper.Scrape(baseUrl)

	// avoid declaring root trusted certs in the system
	// BUT SACRIFICE SECURITY
	config := tls.Config{
		InsecureSkipVerify: true,
	}
	transport := http.Transport{
		TLSClientConfig: &config,
	}
	client := http.Client{
		Transport: &transport,
	}

	response, err := client.Get(baseUrl)

	checkError(err)

	fmt.Println(response)
	fmt.Println(response.Body)
	body, err := ioutil.ReadAll(response.Body)

	checkError(err)

	// rootNode, err := htmlquery.LoadDoc(string(body))
	// checkError(err)
	fmt.Println(body)

	// I think using colly will provide an executed rendered page (not just the initial html)
	// I don't think so anymore
	// I will read some docs for that

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
