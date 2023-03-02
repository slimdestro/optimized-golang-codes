package main

import (
	"fmt"
)

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Create two channels
	oddChan := make(chan []int)
	evenChan := make(chan []int)

	// Launch two go routines
	go func() {
		oddNums := []int{}
		for _, num := range nums {
			if num%2 != 0 {
				oddNums = append(oddNums, num)
			}
		}
		oddChan <- oddNums
	}()

	go func() {
		evenNums := []int{}
		for _, num := range nums {
			if num%2 == 0 {
				evenNums = append(evenNums, num)
			}
		}
		evenChan <- evenNums
	}()

	// Receive from the channels
	oddNums := <-oddChan
	evenNums := <-evenChan

	fmt.Println("Odd numbers:", oddNums)
	fmt.Println("Even numbers:", evenNums)
}
