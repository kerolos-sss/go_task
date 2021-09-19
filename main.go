package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	fmt.Println("hi I am a func")

	baseUrl := "http://youtube.com"
	response, err := http.Get(baseUrl)

	checkError(err)

	fmt.Println(response)
	fmt.Println(response.Body)
	body, err := ioutil.ReadAll(response.Body)

	checkError(err)
	fmt.Println(body)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
