//simple password generator 

package main

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Create a slice of characters to choose from
	characters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+")

	// Create a variable to store the password
	var password string

	// Generate a random length for the password
	length := rand.Intn(20) + 10

	// Generate the password
	for i := 0; i < length; i++ {
		password += string(characters[rand.Intn(len(characters))])
	}

	// Print the password
	println(password)
}