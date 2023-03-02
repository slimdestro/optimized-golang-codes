package main

import "fmt"

// Function to find the nth Fibonacci number using dynamic programming
func Fibonacci(n int) int {
	// Declare an array to store Fibonacci numbers
	f := make([]int, n+1, n+2)
	if n < 2 {
		f = f[0:2]
	}
	f[0] = 0
	f[1] = 1
	for i := 2; i <= n; i++ {
		f[i] = f[i-1] + f[i-2]
	}
	return f[n]
}

func main() {
	fmt.Println(Fibonacci(9))
}