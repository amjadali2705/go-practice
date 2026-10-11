package main

import (
	"fmt"

	"example.com/day12/calculator"
)

func main() {
	res := calculator.Add(10, 20)
	fmt.Println("Result:", res)

	res1 := calculator.Subtract(20, 10)
	fmt.Println("Result:", res1)
}



