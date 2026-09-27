package main

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

func repo() error {
	return ErrNotFound
}

func service() error {
	err := repo()
	return fmt.Errorf("service failed: %w", err)
}

var ErrDB = errors.New("database error")

func repo1() error {
	return ErrDB
}

func service1() error {
	return fmt.Errorf("service failed: %w", repo1())
}

func handler() error {
	return fmt.Errorf("handler failed: %v", service1())
}

type MyError struct {
	Code int
}

func (e *MyError) Error() string {
	return "custom error"
}

func main() {
	// err := service()
	// fmt.Println(err)                         // service failed: not found
	// fmt.Println(errors.Is(err, ErrNotFound)) // true

	// err1 := fmt.Errorf("wrapped: %v", ErrNotFound) // wrapped: not found
	// fmt.Println(errors.Is(err1, ErrNotFound))      // false

	// original := &MyError{Code: 500}
	// err2 := fmt.Errorf("operation failed: %w", original)
	// var target *MyError
	// fmt.Println(errors.As(err2, &target)) // true
	// fmt.Println(target.Code)              // 500

	// defer func() {
	// 	fmt.Println("defer") // defer will be executed even if panic occurs
	// }()
	// panic("boom")              // panic will stop the execution of the program
	// fmt.Println("after panic") // this line will not be executed

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("recovered") // recover will catch the panic and allow the program to continue
	// 	}
	// }()
	// fmt.Println("before") // this line will be executed
	// panic("boom")         // panic will stop the execution of the program, but recover will catch it
	// fmt.Println("after")  //this line will not be executed

	err3 := handler()
	fmt.Println(err3)                   // handler failed: service failed: database error
	fmt.Println(errors.Is(err3, ErrDB)) // false, because the error is wrapped with %v, not %w
}
