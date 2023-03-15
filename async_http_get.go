package main

import (
	"fmt"
	"net/http"
)

func main() {
	urls := []string{
		"http://site.com/1",
		"http://site.com/2",
		"http://site.com/3",
	}

	// Create a channel to receive the responses
	responses := make(chan *http.Response)

	// Make the requests asynchronously
	for _, url := range urls {
		go func(url string) {
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
				return
			}
			responses <- resp
		}(url)
	}

	// Read the responses from the channel
	for range urls {
		resp := <-responses
		fmt.Println(resp.Status)
	}
}