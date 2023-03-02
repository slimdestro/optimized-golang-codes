package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Repository struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Stars int `json:"stars"`
	Pulls int `json:"pulls"`
	Tags []string `json:"tags"`
}

func main() {
	// Get the repository data from Docker Hub
	resp, err := http.Get("https://registry.hub.docker.com/v2/repositories/library/ubuntu/")
	if err != nil {
		fmt.Println("Error getting repository data:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Unmarshal the JSON data
	var repo Repository
	err = json.Unmarshal(body, &repo)
	if err != nil {
		fmt.Println("Error unmarshalling JSON data:", err)
		return
	}

	// Print the repository data
	fmt.Println("Name:", repo.Name)
	fmt.Println("Description:", repo.Description)
	fmt.Println("Stars:", repo.Stars)
	fmt.Println("Pulls:", repo.Pulls)
	fmt.Println("Tags:", repo.Tags)
}
      