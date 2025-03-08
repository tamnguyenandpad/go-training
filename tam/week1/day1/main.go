package main

import (
	"fmt"
)

func fib(n int) int {
	fmt.Printf("Finding the Fibonacci number at the position %d ...\n", n)
	defer fmt.Println("Done!")
	if n <= 0 {
		fmt.Println("Invalid input! Please enter a positive integer.")
		return -1
	}

	if n == 1 {
		fmt.Println("\nThe Fibonacci number at position 1 is 0")
		return 0
	}

	n1, n2 := 0, 1
	for i := 2; i < n; i++ {
		n1, n2 = n2, n1+n2
	}
	fmt.Printf("\nThe Fibonacci number at position %d is %d\n", n, n2)
	return n2
}

func main() {
	n := 9 // The position of the Fibonacci number to find
	fib(n)
}
