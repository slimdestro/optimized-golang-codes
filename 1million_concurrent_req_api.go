package main

import (
	"fmt"
	"net/http"
	"sync"
)

// Create a struct to store the data
type Data struct {
	Name string
	Age  int
}

// Create a global map to store the data
var dataMap map[string]Data

// Create a global mutex to protect the data
var mutex sync.Mutex

// Create a handler function to handle the requests
func handler(w http.ResponseWriter, r *http.Request) {
	// Lock the mutex
	mutex.Lock()
	// Get the data from the map
	data := dataMap[r.URL.Path]
	// Unlock the mutex
	mutex.Unlock()
	// Write the response
	fmt.Fprintf(w, "Name: %s, Age: %d", data.Name, data.Age)
}

func main() {
	// Initialize the data map
	dataMap = make(map[string]Data)
	// Add some data to the map
	dataMap["/user1"] = Data{Name: "John", Age: 25}
	dataMap["/user2"] = Data{Name: "Jane", Age: 30}
	// Create a server
	server := http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(handler),
	}
	// Set the max concurrent users to 1 million
	server.SetMaxConcurrentUsers(1000000)
	// Start the server
	server.ListenAndServe()
}