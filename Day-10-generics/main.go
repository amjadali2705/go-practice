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

// # — Generics

// ## 1. What are Generics?

// Generics allow us to write code that works with **multiple types** while maintaining compile-time type safety.

// Instead of writing the same function for `int`, `float64`, `string`, etc., we can write one generic function.

// Example:

// ```go
// func Identity[T any](value T) T {
// 	return value
// }
// ```

// Usage:

// ```go
// x := Identity(10)        // T = int
// y := Identity("hello")  // T = string
// ```

// Generics are mainly useful when:

// ```text
// Same algorithm
//      +
// Different types
//      =
// Generic code
// ```

// ---

// # 2. Type Parameters

// In:

// ```go
// func Identity[T any](value T) T {
// 	return value
// }
// ```

// `T` is a **type parameter**.

// It represents a type that will be determined when the function is used.

// Example:

// ```go
// Identity(100)
// ```

// Go infers:

// ```text
// T = int
// ```

// So conceptually the function becomes:

// ```text
// Identity[int]
// ```

// ---

// # 3. Type Parameter vs Value Parameter

// Example:

// ```go
// func Add[T int | float64](a, b T) T {
// 	return a + b
// }
// ```

// Here:

// ```text
// T   → type parameter
// a,b → value parameters
// ```

// `T` represents a type.

// `a` and `b` represent values.

// ---

// # 4. `any`

// `any` is an alias for:

// ```go
// interface{}
// ```

// Example:

// ```go
// func Identity[T any](value T) T {
// 	return value
// }
// ```

// This means:

// > `T` may be any type.

// But `any` doesn't guarantee that a particular operation is supported.

// For example, this is invalid:

// ```go
// func Equal[T any](a, b T) bool {
// 	return a == b
// }
// ```

// because `any` does not guarantee that every possible `T` can be compared.

// ---

// # 5. Constraints

// A **constraint** specifies which types are allowed for a type parameter and, in some cases, which operations can be performed on that type.

// Example:

// ```go
// func Max[T int | float64](a, b T) T {
// 	if a > b {
// 		return a
// 	}

// 	return b
// }
// ```

// Here:

// ```text
// T = int
// or
// T = float64
// ```

// are allowed.

// A type such as:

// ```go
// string
// ```

// does not satisfy this particular constraint.

// ---

// # 6. Union Constraints

// This:

// ```go
// [T int | float64]
// ```

// means:

// ```text
// T can be int OR float64
// ```

// Example:

// ```go
// func Sum[T int | float64](a, b T) T {
// 	return a + b
// }
// ```

// Usage:

// ```go
// Sum(10, 20)
// Sum(1.5, 2.5)
// ```

// ---

// # 7. Custom Constraints

// Constraints can be defined as interfaces.

// Example:

// ```go
// type Number interface {
// 	int | int64 | float64
// }
// ```

// Then:

// ```go
// func Add[T Number](a, b T) T {
// 	return a + b
// }
// ```

// Now `T` must satisfy `Number`.

// ---

// # 8. The `~` Operator

// `~` is very important in generics.

// Consider:

// ```go
// type UserID int
// ```

// `UserID` is a distinct defined type, even though its underlying type is `int`.

// ```text
// Defined type   → UserID
// Underlying type → int
// ```

// If the constraint is:

// ```go
// type Number interface {
// 	int
// }
// ```

// `UserID` does not satisfy it.

// But:

// ```go
// type Number interface {
// 	~int
// }
// ```

// allows:

// ```text
// int
// UserID
// Other defined types whose underlying type is int
// ```

// So `~int` means:

// > `int` and types whose underlying type is `int`.

// Example:

// ```go
// type UserID int

// type Number interface {
// 	~int | ~float64
// }

// func Add[T Number](a, b T) T {
// 	return a + b
// }
// ```

// This works:

// ```go
// a := UserID(10)
// b := UserID(20)

// result := Add(a, b)
// ```

// The result is:

// ```text
// 30
// ```

// and its type is `UserID`, because `T` was inferred as `UserID`.

// ---

// # 9. Type Inference

// Go can often determine the type parameter automatically.

// Example:

// ```go
// func Identity[T any](value T) T {
// 	return value
// }
// ```

// Instead of:

// ```go
// Identity[int](10)
// ```

// we can simply write:

// ```go
// Identity(10)
// ```

// Go infers:

// ```text
// T = int
// ```

// Similarly:

// ```go
// Identity("hello")
// ```

// gives:

// ```text
// T = string
// ```

// Type inference makes generic code easier to read.

// ---

// # 10. Explicit Type Arguments

// You can also specify the type yourself:

// ```go
// result := Identity[int](100)
// ```

// Both are valid:

// ```go
// Identity(100)
// Identity[int](100)
// ```

// Usually, inference is preferred when it is obvious.

// ---

// # 11. `comparable`

// `comparable` is a built-in constraint.

// ```go
// [T comparable]
// ```

// means that `T` supports:

// ```go
// ==
// !=
// ```

// Example:

// ```go
// func Equal[T comparable](a, b T) bool {
// 	return a == b
// }
// ```

// Now:

// ```go
// Equal(10, 10)          // valid
// Equal("go", "go")      // valid
// ```

// But:

// ```text
// slice     → not comparable
// map       → not comparable
// function  → not comparable
// ```

// So:

// ```go
// Equal([]int{1}, []int{1})
// ```

// is not allowed.

// ---

// # 12. `any` vs `comparable`

// Remember:

// ```text
// any
//  ↓
// T can be any type

// comparable
//  ↓
// T must support == and !=
// ```

// Therefore:

// ```go
// func Equal[T any](a, b T) bool
// ```

// does not work.

// But:

// ```go
// func Equal[T comparable](a, b T) bool
// ```

// does.

// ---

// # 13. Generic Functions

// Basic generic function:

// ```go
// func First[T any](items []T) T {
// 	return items[0]
// }
// ```

// Usage:

// ```go
// nums := []int{10, 20, 30}
// x := First(nums)
// ```

// Here:

// ```text
// T = int
// x = int
// ```

// For:

// ```go
// names := []string{"A", "B"}
// ```

// the same function works with:

// ```text
// T = string
// ```

// ---

// # 14. Generic Functions with Multiple Type Parameters

// A function can have multiple type parameters.

// Example:

// ```go
// func Pair[K comparable, V any](key K, value V) {
// 	fmt.Println(key, value)
// }
// ```

// Here:

// ```text
// K → key type
// V → value type
// ```

// Examples:

// ```go
// Pair("user_id", 100)
// Pair(1, "Amjad")
// ```

// ---

// # 15. Generic Types

// Generics can also be used with structs.

// Example:

// ```go
// type Box[T any] struct {
// 	Value T
// }
// ```

// Usage:

// ```go
// intBox := Box[int]{
// 	Value: 100,
// }

// stringBox := Box[string]{
// 	Value: "hello",
// }
// ```

// Therefore:

// ```text
// intBox    → Box[int]
// stringBox → Box[string]
// ```

// and:

// ```text
// intBox.Value    → int
// stringBox.Value → string
// ```

// ---

// # 16. Generic Data Structures

// Generics are useful for data structures.

// Example:

// ```go
// type Stack[T any] struct {
// 	items []T
// }

// func (s *Stack[T]) Push(value T) {
// 	s.items = append(s.items, value)
// }

// func (s *Stack[T]) Pop() T {
// 	last := len(s.items) - 1
// 	value := s.items[last]
// 	s.items = s.items[:last]

// 	return value
// }
// ```

// Now the same stack implementation can be used for:

// ```text
// Stack[int]
// Stack[string]
// Stack[User]
// ```

// without duplicating the implementation.

// ---

// # 17. Generics vs Interfaces

// This is one of the most important concepts.

// ## Generics

// Generics are mainly about:

// > Reusing code across different types.

// Example:

// ```go
// func First[T any](items []T) T
// ```

// The algorithm is the same regardless of whether `T` is:

// ```text
// int
// string
// User
// Transaction
// ```

// ## Interfaces

// Interfaces are mainly about:

// > Defining common behavior.

// Example:

// ```go
// type Speaker interface {
// 	Speak()
// }
// ```

// Any type that implements `Speak()` satisfies the interface.

// ### Mental model

// ```text
// Generics
// → "I want the same code to work with different TYPES."

// Interfaces
// → "I want to work with values that provide the same BEHAVIOR."
// ```

// ---

// # 18. Example: Interface Is Better

// Suppose:

// ```go
// type Speaker interface {
// 	Speak()
// }
// ```

// and:

// ```go
// func MakeSpeak(s Speaker) {
// 	s.Speak()
// }
// ```

// We don't care whether `s` is:

// ```text
// Dog
// Cat
// Human
// ```

// We only care that it has:

// ```go
// Speak()
// ```

// So an interface is the natural abstraction.

// ---

// # 19. Example: Generic Is Better

// Suppose:

// ```go
// func Contains[T comparable](items []T, target T) bool {
// 	for _, item := range items {
// 		if item == target {
// 			return true
// 		}
// 	}

// 	return false
// }
// ```

// The algorithm is the same for:

// ```text
// []int
// []string
// []UserID
// ```

// This is a good generic use case.

// ---

// # 20. `any` vs Generic

// Compare:

// ```go
// func Print(x any) {
// 	fmt.Println(x)
// }
// ```

// with:

// ```go
// func Identity[T any](x T) T {
// 	return x
// }
// ```

// `Print` says:

// > I accept any value.

// `Identity` says:

// > I accept any type, but the input and output are tied to the same type `T`.

// Example:

// ```go
// x := Identity(10)
// ```

// `x` is known to be:

// ```text
// int
// ```

// So:

// ```text
// any
// → arbitrary value

// generic
// → type-safe relationship between values/types
// ```

// ---

// # 21. Generic Repository — Don't Overuse It

// You may see something like:

// ```go
// type Repository[T any] interface {
// 	Get(id int) (T, error)
// 	Save(entity T) error
// }
// ```

// This can be useful when multiple repositories genuinely share the same abstraction.

// But don't automatically convert every repository into a generic repository.

// For example:

// ```text
// UserRepository
// TransactionRepository
// OrderRepository
// ```

// may have completely different:

// ```text
// queries
// joins
// business rules
// filters
// database operations
// ```

// In that case, forcing them into one generic abstraction can make the code more complicated.

// ### Principle

// > Use generics when there is genuine type-level reuse, not simply because generics exist.

// ---

// # 22. Generic Methods — Important Limitation

// Go allows generic types:

// ```go
// type Box[T any] struct {
// 	Value T
// }
// ```

// But methods cannot introduce their own independent type parameters.

// This is not allowed:

// ```go
// func (b Box[T]) Convert[R any]() R {
// 	// ...
// }
// ```

// Instead, use a generic function:

// ```go
// func Convert[T any, R any](b Box[T]) R {
// 	// ...
// }
// ```

// This is a useful interview/trivia point.

// ---

// # 23. Constraints Can Express Behavior

// Constraints can also require methods.

// Example:

// ```go
// type Stringer interface {
// 	String() string
// }
// ```

// Then:

// ```go
// func PrintAll[T Stringer](items []T) {
// 	for _, item := range items {
// 		fmt.Println(item.String())
// 	}
// }
// ```

// Now `T` must provide:

// ```go
// String() string
// ```

// So constraints can restrict types based on:

// ```text
// type sets
// and/or
// required methods
// ```

// ---

// # 24. Generics Do Not Replace Interfaces

// A common misconception is:

// > "Now that Go has generics, interfaces are no longer needed."

// Incorrect.

// They solve different problems.

// ```text
// Generics
// → type-level reuse

// Interfaces
// → behavioral abstraction
// ```

// Both remain important.

// ---

// # 25. When Should You Use Generics?

// Good use cases:

// ```text
// Reusable algorithms
// Generic data structures
// Collection utilities
// Type-safe helpers
// Code where input/output types have a clear relationship
// ```

// Examples:

// ```go
// Contains[T comparable]
// First[T any]
// Map[T, R]
// Stack[T]
// Box[T]
// ```

// ---

// # 26. When Should You NOT Use Generics?

// Avoid generics when they don't make the code clearer.

// For example:

// ```go
// func GetUser(id int) (*User, error)
// ```

// doesn't automatically become better as:

// ```go
// func Get[T any](id int) (T, error)
// ```

// If the operation is specifically about users, using `User` is usually clearer.

// Don't use generics merely to:

// ```text
// remove every duplicate
// make code look advanced
// avoid defining domain-specific types
// ```

// ---

// # 27. Important Interview Questions

// ### Q1. What are generics?

// Generics allow functions, types, and data structures to work with multiple types while maintaining compile-time type safety.

// ### Q2. What is a type parameter?

// A type parameter is a placeholder representing a type that is supplied or inferred when generic code is instantiated.

// ### Q3. What is a constraint?

// A constraint specifies which types can be used for a type parameter and what operations are valid for that parameter.

// ### Q4. What does `any` mean?

// `any` is an alias for `interface{}` and allows any type.

// ### Q5. What does `comparable` mean?

// It restricts the type parameter to types that support `==` and `!=`.

// ### Q6. What does `~int` mean?

// It allows `int` and defined types whose underlying type is `int`.

// ### Q7. Generics vs interfaces?

// ```text
// Generics → reuse code across types
// Interfaces → abstract common behavior
// ```

// ### Q8. What is type inference?

// Go can infer the generic type parameter from the function arguments, so explicit type arguments are often unnecessary.

// ### Q9. Can generic types have methods?

// Yes.

// Example:

// ```go
// type Box[T any] struct {
// 	Value T
// }
// ```

// Methods can use the type parameter `T`, but a method cannot declare its own additional type parameters.

// ### Q10. Should every repository be generic?

// No. Use generics only when there is genuine reusable type-level abstraction.

// ---

// # 28. Golden Rules

// Remember these:

// ```text
// [T any]
// → T can be any type

// [T comparable]
// → T supports == and !=

// [T int | float64]
// → T can only be int or float64

// [T ~int]
// → T can be int or a defined type whose underlying type is int

// % generics
// → reusable code across types

// Interfaces
// → common behavior
// ```

// Most importantly:

// ```text
// GENERIC
// "I want one algorithm to work with many TYPES."

// INTERFACE
// "I want code to work with many VALUES
// that provide the same BEHAVIOR."
// ```

// ---

// # 29. Day 10 Mental Model

// ```text
//                     GENERICS
//                        |
//           +------------+------------+
//           |                         |
//     Type Parameters             Constraints
//           |                         |
//           T                    any / comparable
//           |                    int | float64
//           |                    ~int
//           |
//       Type Inference
//           |
//    Generic Functions
//    Generic Structs
//    Generic Data Structures
// ```

// For backend development:

// ```text
// Generics
//    ↓
// Reusable + type-safe utilities

// Interfaces
//    ↓
// Abstraction + dependency injection

// Concrete Types
//    ↓
// Domain-specific business logic
// ```

// The goal is not to replace one with another.

// The goal is to know **which abstraction solves which problem**.
