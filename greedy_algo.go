package main

import "fmt"

func main() {
	// define the array of items
	items := []int{6, 10, 4, 5, 20}

	// define the capacity of the knapsack
	capacity := 12

	// define the knapsack
	knapsack := []int{}

	// define the current capacity
	currentCapacity := 0

	// loop through the items
	for _, item := range items {
		// check if the item fits in the knapsack
		if currentCapacity+item <= capacity {
			// add the item to the knapsack
			knapsack = append(knapsack, item)
			// update the current capacity
			currentCapacity += item
		}
	}

	// print the knapsack
	fmt.Println(knapsack)
}
