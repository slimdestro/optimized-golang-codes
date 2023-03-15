/**
* total cost required to make an array sorted
* @slimdestro
*/

package main

import "fmt"

func totalCost(arr []int) int {
	// Initialize cost
	cost := 0

	// Iterate over the array
	for i := 0; i < len(arr)-1; i++ {
		// If the current element is greater than the next element
		if arr[i] > arr[i+1] {
			// Calculate the difference
			diff := arr[i] - arr[i+1]
			// Add the difference to the cost
			cost += diff
			// Swap the elements
			arr[i], arr[i+1] = arr[i+1], arr[i]
		}
	}

	// Return the cost
	return cost
}

func main() {
	arr := []int{4, 3, 2, 1}
	fmt.Println("Total cost of making array sorted:", totalCost(arr))
}