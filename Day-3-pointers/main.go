package main

import "fmt"

func main() {
	x := 10
	p := &x
	fmt.Println(*p) // Output: 10

	a := 10
	b := &a
	*b = 50
	fmt.Println(a) // Output: 50

	n := 10
	update(n)
	fmt.Println(n) // Output: 10

	m := 10
	update2(&m)
	fmt.Println(m) // Output: 100

	y := 10
	change(&y)
	fmt.Println(y) // Output: 10

	d := 10
	e := &d
	f := &e
	**f = 300
	fmt.Println(d) // Output: 300

	x1 := 10
	y1 := 20
	swap(&x1, &y1)
	fmt.Println("x1:", x1)
	fmt.Println("y1:", y1)

	p1 := Point{X: 1, Y: 2}
	fmt.Println("Before:", p1)
	updatePoint(&p1)
	fmt.Println("After:", p1)

	num := 10
	increment(&num)
	fmt.Println("Incremented:", num)

	printAddresses()
}

func update(x int) {
	x = 100
}

func update2(x *int) {
	*x = 100
}

func change(p *int) {
	p = new(int)
	*p = 100
}

// Swap two numbers using pointers.
func swap(p, q *int) {
	*p, *q = *q, *p
}

// Update a struct through a pointer.
type Point struct {
	X, Y int
}

func updatePoint(p *Point) {
	p.X = 100
	p.Y = 200
}

// Create a function that increments a number using *int.
func increment(p *int) {
	*p++
}

// Print addresses using &
func printAddresses() {
	x := 10
	p := &x
	fmt.Printf("Address of x: %p\n", p)
}

// Why use pointers?
// To modify variables in the caller's scope, to avoid copying large data structures, and to enable shared access to data.

// What is a pointer?
// A variable that holds the memory address of another variable.

// How to declare a pointer?
// var p *int

// How to get the address of a variable?
// Use the & operator, e.g., p := &x

// How to dereference a pointer?
// Use the * operator, e.g., *p to access the value at the address.

// What is nil pointer?
// A pointer that does not point to any valid memory address, often used to indicate the absence of a value.

// What is a pointer to a struct?
// A pointer that holds the address of a struct variable, allowing you to modify the struct's fields through the pointer.

// How to create a pointer to a struct?
// Use the & operator, e.g., p := &Point{X: 1, Y: 2}

// What is a pointer to a pointer?
// A pointer that holds the address of another pointer, allowing for multiple levels of indirection.

// How to declare a pointer to a pointer?
// var pp **int

// What is the zero value of a pointer?
// nil, which indicates that the pointer does not point to any valid memory address.
