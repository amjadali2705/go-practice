package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")

	// Example usage of the generic function
	printValues(1, 2, 3, 4, 5)
}

// write about generics in go and how to use it with examples
// Generics in Go, introduced in Go 1.18, allow developers to write functions and data structures that can operate on different types without sacrificing type safety. This feature enables code reusability and reduces redundancy by allowing you to define a function or a type that can work with any data type.

/// Basic Example of Generics
func printValues[T any](values ...T) {
	for _, v := range values {
		fmt.Println(v)
	}
}