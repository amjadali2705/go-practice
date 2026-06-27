package main

import "fmt"

type User struct {
	Name string
}

func main() {
	fmt.Println(add(10, 20)) //output - 30

	fmt.Println(add1(1, 2, 3)) //output - 6

	defer fmt.Println(1)
	defer fmt.Println(2)
	fmt.Println(3)

	u := User{"Amjad"}
	u.Change()
	fmt.Println(u.Name) //output - Amjad (because the Change method has a value receiver, it does not modify the original User instance)

	u1 := User{"Amjad"}
	u1.Change1()
	fmt.Println(u1.Name) //output - Ali (because the Change1 method has a pointer receiver, it modifies the original User instance)

	x := 10
	defer fmt.Println(x) //output - 10 (the value of x is evaluated at the time of the defer statement, so it prints 10 even though x is changed later)
	x = 20
	fmt.Println(x) //output - 20

	c := counter()
	fmt.Println(c()) //output - 1
	fmt.Println(c()) //output - 2

	y := 10
	defer func() {
		fmt.Println(y) // output - 20 (the value of y is evaluated at the time of the defer statement, so it prints 20 even though y is changed later)
	}()
	y = 20

	for i := 1; i <= 3; i++ {
		defer fmt.Println(i) //output - 3 2 1 (deferred function calls are executed in LIFO order)
	}

	funcs := []func(){}
	for i := 0; i < 3; i++ {

		funcs = append(funcs, func() {
			fmt.Println(i)
		})
	}
	for _, f := range funcs {
		f()
	}
}

func add(a, b int) int {
	return a + b
}

func add1(nums ...int) int {

	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}

func (u User) Change() {
	u.Name = "Ali"
}

func (u *User) Change1() {
	u.Name = "Ali"
}

func counter() func() int {

	x := 0

	return func() int {
		x++
		return x
	}
}

// Functions are first-class values.
// Go supports multiple return values.
// Errors are commonly returned as the second return value.
// Variadic functions accept any number of arguments using ...
// Anonymous functions can be assigned to variables or invoked immediately.
// Closures capture variables from their surrounding scope.
// defer executes when the surrounding function returns and follows LIFO order.
// Deferred function arguments are evaluated immediately, but deferred closures read captured variables when they execute.
// Methods belong to types.
// Value receivers operate on a copy.
// Pointer receivers operate on the original object.
// Go automatically takes the address when calling pointer receiver methods on addressable values.
