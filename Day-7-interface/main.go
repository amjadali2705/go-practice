package main

import "fmt"

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	width, height float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.radius * c.radius
}

func (c Circle) Perimeter() float64 {
	return 2 * 3.14 * c.radius
}

type Speaker interface {
	Speak()
}

type Dog struct{}

func (d Dog) Speak() {
	fmt.Println("Woof")
}

type User struct{}

func check(x any) {
	switch v := x.(type) {
	case int:
		fmt.Println("int", v)
	case string:
		fmt.Println("string", v)
	default:
		fmt.Println("unknown")
	}
}

func main() {
	shapes := []Shape{
		Rectangle{width: 5, height: 3},
		Circle{radius: 2},
	}

	for _, shape := range shapes {
		fmt.Printf("Area: %.2f, Perimeter: %.2f\n", shape.Area(), shape.Perimeter())
	}

	var s Speaker = Dog{} // Dog implements the Speaker interface
	s.Speak()             // Output: Woof

	var x interface{}
	fmt.Println(x == nil) // Output: true

	var u *User = nil
	var x1 interface{} = u
	fmt.Println(x1 == nil) // Output: false, because x1 holds a typed nil (*User), which is not equal to nil

	var x3 any = 10
	value := x3.(int)
	fmt.Println(value) // Output: 10

	// var x4 any = "hello"
	// value1 := x4.(int)
	// fmt.Println(value1) // This will cause a panic at runtime because x4 does not hold an int

	var x5 any = "hello"
	value2, ok := x5.(int)
	fmt.Println(value2) // Output: 0, because the assertion failed
	fmt.Println(ok)     // Output: false, indicating that the assertion was not successful

	check(10)   // Output: int 10
	check("Go") // Output: string Go
	check(true) // Output: unknown
}

// What is the difference between a nil interface and an interface containing a typed nil?

// A nil interface is an interface value that has no dynamic type and no value. It is represented as `nil` in Go. When you compare a nil interface to `nil`, it evaluates to true.

// An interface containing a typed nil, on the other hand, has a dynamic type (the type of the value it holds) but the value itself is nil. In this case, when you compare the interface to `nil`, it evaluates to false because the interface has a dynamic type, even though the value is nil.

// INTERFACE
//    │
//    ├── defines behavior
//    │
//    ├── implicit implementation
//    │
//    ├── dynamic type + dynamic value
//    │
//    ├── method sets determine satisfaction
//    │
//    ├── any / interface{}
//    │       │
//    │       ├── type assertion
//    │       └── type switch
//    │
//    └── nil
//          │
//          ├── nil interface
//          │
//          └── typed nil



// Go interfaces are implemented implicitly; there is no explicit implements keyword.

// A type satisfies an interface when its method set contains all methods required by the interface.

// An interface value conceptually contains a dynamic type and a dynamic value.

// A type assertion extracts or checks the concrete dynamic type stored inside an interface.

// A type switch allows branching based on the dynamic type stored in an interface.

// An interface is nil only when both its dynamic type and dynamic value are nil.

// If a method has a pointer receiver, *T can satisfy the interface while T may not.

// Small interfaces are generally easier to implement, test, compose, and reuse.

// Interfaces are most valuable when they provide meaningful abstraction, decoupling, dependency injection, or testability.
