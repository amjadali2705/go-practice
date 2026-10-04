package main

import (
	"testing"
)

// func TestAdd(t *testing.T) {
// 	got := Add(2, 3)

// 	if got != 5 {
// 		t.Errorf("got %d, want %d", got, 5)
// 	}
// }

// func TestExample(t *testing.T) {
// 	t.Error("something went wrong") // This will mark the test as failed but continue execution

// 	fmt.Println("hello") // This will be printed even though the test failed
// }

// func TestExample1(t *testing.T) {
// 	t.Fatal("something went wrong") // This will mark the test as failed and stop execution

// 	fmt.Println("hello") // This will not be printed because the test failed and execution stopped
// }

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(10, 20)
	}
}

// # Go Day 11 — Testing

// Go has a built-in testing framework through the standard library package:

// ```go
// import "testing"
// ```

// Testing is especially important for backend development because it helps verify business logic, API behavior, database interactions, concurrency, and performance.

// ---

// # 1. Test File Naming

// Go recognizes test files ending with:

// ```text
// _test.go
// ```

// Example:

// ```text
// project/
// ├── calculator.go
// └── calculator_test.go
// ```

// A normal test function starts with:

// ```go
// func TestXxx(t *testing.T)
// ```

// A benchmark starts with:

// ```go
// func BenchmarkXxx(b *testing.B)
// ```

// ---

// # 2. Basic Unit Test

// Production code:

// ```go
// func Add(a, b int) int {
// 	return a + b
// }
// ```

// Test:

// ```go
// func TestAdd(t *testing.T) {
// 	got := Add(2, 3)

// 	if got != 5 {
// 		t.Errorf("got %d, want %d", got, 5)
// 	}
// }
// ```

// Run:

// ```bash
// go test
// ```

// The basic pattern is:

// ```text
// Arrange
//    ↓
// Act
//    ↓
// Assert
// ```

// Example:

// ```text
// Arrange → input = 2, 3
// Act     → Add(2, 3)
// Assert  → result == 5
// ```

// ---

// # 3. `testing.T`

// A test receives:

// ```go
// *t testing.T
// ```

// `t` provides methods for controlling and reporting the test.

// Common methods:

// ```text
// t.Error()
// t.Errorf()
// t.Fatal()
// t.Fatalf()
// t.Log()
// t.Run()
// t.Helper()
// t.Parallel()
// ```

// ---

// # 4. `t.Error` vs `t.Fatal`

// ## `t.Error`

// Marks the test as failed but continues execution.

// ```go
// func TestExample(t *testing.T) {
// 	t.Error("something went wrong")

// 	fmt.Println("hello")
// }
// ```

// `hello` is still executed.

// Mental model:

// ```text
// t.Error()
//     ↓
// test fails
//     ↓
// execution continues
// ```

// ## `t.Fatal`

// Marks the test as failed and stops the current test function.

// ```go
// func TestExample(t *testing.T) {
// 	t.Fatal("something went wrong")

// 	fmt.Println("hello")
// }
// ```

// `hello` will not execute.

// Mental model:

// ```text
// t.Fatal()
//     ↓
// test fails
//     ↓
// current test stops
// ```

// ---

// # 5. Table-Driven Tests

// A very common Go testing pattern is the table-driven test.

// Instead of writing multiple test functions:

// ```text
// TestAddPositive
// TestAddZero
// TestAddNegative
// ```

// put the cases into a slice:

// ```go
// func TestAdd(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		a    int
// 		b    int
// 		want int
// 	}{
// 		{
// 			name: "positive",
// 			a:    2,
// 			b:    3,
// 			want: 5,
// 		},
// 		{
// 			name: "zero",
// 			a:    0,
// 			b:    5,
// 			want: 5,
// 		},
// 		{
// 			name: "negative",
// 			a:    -2,
// 			b:    -3,
// 			want: -5,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := Add(tt.a, tt.b)

// 			if got != tt.want {
// 				t.Errorf("got %d, want %d", got, tt.want)
// 			}
// 		})
// 	}
// }
// ```

// Benefits:

// ```text
// Less duplicate code
// Easy to add new cases
// Clear test data
// Readable failures
// ```

// ---

// # 6. `t.Run`

// `t.Run()` creates a named subtest.

// Example:

// ```go
// t.Run("positive", func(t *testing.T) {
// 	...
// })
// ```

// The structure becomes:

// ```text
// TestAdd
// ├── positive
// ├── zero
// └── negative
// ```

// You can run a particular subtest using an appropriate test name pattern.

// ---

// # 7. Testing Errors

// When your application returns errors, test them properly.

// Suppose:

// ```go
// var ErrNotFound = errors.New("not found")
// ```

// and:

// ```go
// func GetUser() error {
// 	return fmt.Errorf("database: %w", ErrNotFound)
// }
// ```

// Test:

// ```go
// func TestGetUser(t *testing.T) {
// 	err := GetUser()

// 	if !errors.Is(err, ErrNotFound) {
// 		t.Errorf("expected ErrNotFound, got %v", err)
// 	}
// }
// ```

// Use:

// ```go
// errors.Is()
// ```

// rather than comparing error strings.

// Avoid:

// ```go
// err.Error() == "not found"
// ```

// because error strings are not a stable error identity.

// ---

// # 8. Testing Custom Error Types

// For a custom error:

// ```go
// type ValidationError struct {
// 	Field string
// }

// func (e *ValidationError) Error() string {
// 	return "validation failed"
// }
// ```

// Use `errors.As()` in the test:

// ```go
// var validationErr *ValidationError

// if !errors.As(err, &validationErr) {
// 	t.Fatal("expected ValidationError")
// }
// ```

// This lets you verify:

// ```text
// error type
// +
// error-specific data
// ```

// This connects directly to Day 9 error handling.

// ---

// # 9. Test Helpers

// Suppose you repeatedly use an assertion helper:

// ```go
// func assertEqual(t *testing.T, got, want int) {
// 	if got != want {
// 		t.Errorf("got %d, want %d", got, want)
// 	}
// }
// ```

// Use:

// ```go
// func assertEqual(t *testing.T, got, want int) {
// 	t.Helper()

// 	if got != want {
// 		t.Errorf("got %d, want %d", got, want)
// 	}
// }
// ```

// `t.Helper()` tells the testing framework that this function is a helper, which improves failure reporting.

// ---

// # 10. `t.Parallel`

// Tests can be marked as safe to run in parallel:

// ```go
// func TestA(t *testing.T) {
// 	t.Parallel()

// 	...
// }
// ```

// This can make a test suite faster.

// But parallel tests require care with:

// ```text
// Shared state
// Global variables
// Files
// Database records
// Environment variables
// External resources
// ```

// Don't blindly add `t.Parallel()` to every test.

// ---

// # 11. Testing HTTP Handlers

// Go provides:

// ```go
// net/http/httptest
// ```

// This lets you test HTTP handlers without starting a real HTTP server.

// Example:

// ```go
// req := httptest.NewRequest(
// 	http.MethodGet,
// 	"/users/1",
// 	nil,
// )

// rec := httptest.NewRecorder()

// handler.ServeHTTP(rec, req)
// ```

// You can check:

// ```go
// rec.Code
// rec.Body
// rec.Header()
// ```

// Example:

// ```go
// if rec.Code != http.StatusOK {
// 	t.Errorf("expected 200, got %d", rec.Code)
// }
// ```

// ---

// # 12. Testing Gin Handlers

// For Gin:

// ```go
// gin.SetMode(gin.TestMode)

// router := gin.New()

// router.GET("/hello", func(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "hello",
// 	})
// })
// ```

// Then:

// ```go
// req := httptest.NewRequest(
// 	http.MethodGet,
// 	"/hello",
// 	nil,
// )

// rec := httptest.NewRecorder()

// router.ServeHTTP(rec, req)
// ```

// Check:

// ```go
// if rec.Code != http.StatusOK {
// 	t.Errorf("expected 200, got %d", rec.Code)
// }
// ```

// This is very useful for testing your Go/Gin APIs.

// ---

// # 13. Dependency Injection for Testing

// Suppose a service directly creates a database connection:

// ```go
// func CreateUser() error {
// 	db := connectDatabase()
// 	...
// }
// ```

// Testing this is difficult because the test depends on a real database.

// Instead, define an interface:

// ```go
// type UserRepository interface {
// 	CreateUser(*User) error
// }
// ```

// Then inject it:

// ```go
// type UserService struct {
// 	repo UserRepository
// }
// ```

// Production:

// ```text
// UserService
//     ↓
// Real UserRepository
//     ↓
// Database
// ```

// Testing:

// ```text
// UserService
//     ↓
// Fake/Mock Repository
// ```

// This is called **dependency injection**.

// ---

// # 14. Fake vs Mock

// A **fake** is usually a simplified working implementation.

// Example:

// ```text
// FakeRepository
// → stores data in memory
// ```

// A **mock** generally focuses on expected interactions, such as:

// ```text
// CreateUser()
// must be called once
// ```

// Terminology can vary between teams and libraries, but the main idea is:

// > Replace external dependencies with controlled test doubles.

// ---

// # 15. Unit Test

// A unit test tests a small piece of logic in isolation.

// Example:

// ```text
// Service
//    ↓
// Fake Repository
// ```

// No real database is required.

// Characteristics:

// ```text
// Fast
// Isolated
// Repeatable
// Easy to debug
// ```

// ---

// # 16. Integration Test

// An integration test checks that multiple components work together.

// Example:

// ```text
// Repository
//     ↓
// MariaDB
// ```

// or:

// ```text
// Handler
//    ↓
// Service
//    ↓
// Repository
//    ↓
// MariaDB
// ```

// The key difference is that integration testing involves real component interactions.

// ---

// # 17. End-to-End Test

// An E2E test exercises the complete system from the external user's perspective.

// For example:

// ```text
// HTTP Request
//     ↓
// Gin Handler
//     ↓
// Service
//     ↓
// Repository
//     ↓
// Database
//     ↓
// HTTP Response
// ```

// Difference:

// ```text
// Unit
// → isolated component

// Integration
// → multiple components working together

// E2E
// → whole system / user flow
// ```

// Integration testing is therefore **not automatically the same as E2E testing**.

// ---

// # 18. Happy Path vs Edge Cases

// Don't only test successful input.

// For a backend endpoint, test things such as:

// ```text
// Valid request
// Invalid request
// Missing fields
// Empty values
// Boundary values
// Not found
// Duplicate data
// Unauthorized request
// Expired token
// Database failure
// Timeout
// ```

// Good tests cover both:

// ```text
// Expected behavior
// +
// Failure behavior
// ```

// ---

// # 19. Boundary Testing

// Suppose:

// ```go
// func IsValidAge(age int) bool {
// 	return age >= 18
// }
// ```

// Don't test only:

// ```text
// 20 → true
// ```

// Also test:

// ```text
// 17 → false
// 18 → true
// 19 → true
// 0  → false
// ```

// The boundary is:

// ```text
// 18
// ```

// Boundary conditions are common sources of bugs.

// ---

// # 20. Testing Concurrent Code

// Concurrent code needs synchronization in the **test itself**.

// This is problematic:

// ```go
// func TestCounter(t *testing.T) {
// 	var counter int

// 	for i := 0; i < 100; i++ {
// 		go func() {
// 			counter++
// 		}()
// 	}
// }
// ```

// Problems:

// ```text
// 1. Data race
// 2. No synchronization
// 3. Test may finish before goroutines finish
// ```

// A `WaitGroup` can solve completion:

// ```go
// var wg sync.WaitGroup

// for i := 0; i < 100; i++ {
// 	wg.Add(1)

// 	go func() {
// 		defer wg.Done()
// 		// work
// 	}()
// }

// wg.Wait()
// ```

// But `WaitGroup` does **not** protect `counter`.

// You still need:

// ```text
// Mutex
// or
// Atomic
// ```

// This connects directly to Day 8.

// ---

// # 21. Race Detector

// Run:

// ```bash
// go test -race ./...
// ```

// The race detector can detect many data races.

// It is useful for concurrent code.

// But:

// ```text
// -race
// does NOT guarantee:
//     no deadlocks
//     no goroutine leaks
//     no logical bugs
//     no incorrect results
// ```

// It specifically targets data races.

// ---

// # 22. Benchmarks

// Benchmarks measure performance.

// Example:

// ```go
// func BenchmarkAdd(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		Add(10, 20)
// 	}
// }
// ```

// Run:

// ```bash
// go test -bench=BenchmarkAdd
// ```

// Go automatically chooses `b.N`.

// You don't manually set it.

// Conceptually:

// ```text
// b.N
//  ↓
// number of benchmark iterations
// chosen by Go
// ```

// ---

// # 23. Benchmark Output

// You may see something like:

// ```text
// BenchmarkAdd-8    1000000000    0.3 ns/op
// ```

// Important fields:

// ```text
// 1000000000
// → number of iterations

// 0.3 ns/op
// → average time per operation
// ```

// For memory information:

// ```bash
// go test -bench=BenchmarkAdd -benchmem
// ```

// You may see:

// ```text
// 0 B/op
// 0 allocs/op
// ```

// Meaning:

// ```text
// B/op
// → bytes allocated per operation

// allocs/op
// → allocations per operation
// ```

// ---

// # 24. Running All Benchmarks

// Run:

// ```bash
// go test -bench=.
// ```

// Run a specific benchmark:

// ```bash
// go test -bench=BenchmarkAdd
// ```

// Show memory allocations:

// ```bash
// go test -bench=BenchmarkAdd -benchmem
// ```

// For your practice project, `BenchmarkAdd` worked correctly with the specific command:

// ```bash
// go test -bench=BenchmarkAdd
// ```

// ---

// # 25. Fuzz Testing

// Go also supports fuzz testing.

// Example:

// ```go
// func FuzzReverse(f *testing.F) {
// 	f.Add("hello")

// 	f.Fuzz(func(t *testing.T, s string) {
// 		// test property
// 	})
// }
// ```

// Run:

// ```bash
// go test -fuzz=.
// ```

// Fuzzing generates many different inputs automatically.

// Useful for:

// ```text
// Parsers
// String processing
// Input validation
// Encoders/decoders
// Protocol handling
// Security-sensitive input
// ```

// ---

// # 26. Property-Based Thinking

// Instead of testing only known examples, define a property that should always be true.

// For example:

// ```text
// Reverse(Reverse(s)) == s
// ```

// should hold for every valid string `s`.

// Fuzzing can then search for inputs that violate the property.

// ---

// # 27. Test Coverage

// Run:

// ```bash
// go test -cover ./...
// ```

// This tells you how much code was executed by tests.

// You can generate a profile:

// ```bash
// go test -coverprofile=coverage.out ./...
// ```

// Then:

// ```bash
// go tool cover -html=coverage.out
// ```

// ---

// # 28. Coverage Is Not Correctness

// This is a very important point:

// ```text
// 100% coverage ≠ bug-free
// ```

// Coverage tells you:

// > Was this code executed?

// It does not guarantee:

// > Did this code produce the correct result?

// Example:

// ```go
// func TestAdd(t *testing.T) {
// 	Add(2, 3)
// }
// ```

// The function executes, but the test doesn't check whether the result is correct.

// Therefore:

// ```text
// Coverage
// → code execution

// Assertions
// → correctness
// ```

// ---

// # 29. Flaky Tests

// A flaky test sometimes passes and sometimes fails without a relevant code change.

// Common causes:

// ```text
// Race conditions
// Timing assumptions
// time.Sleep()
// External services
// Network calls
// Current time
// Randomness
// Shared global state
// Test ordering
// ```

// Avoid tests like:

// ```go
// time.Sleep(100 * time.Millisecond)
// ```

// just because you hope a goroutine finishes.

// Prefer explicit synchronization:

// ```text
// WaitGroup
// Channel
// Context
// Controlled dependency
// ```

// ---

// # 30. Don't Over-Test Implementation Details

// Tests should generally verify **observable behavior**, not unnecessary implementation details.

// Suppose the implementation changes from:

// ```text
// for loop
// ```

// to:

// ```text
// map
// ```

// The tests shouldn't fail just because the internal implementation changed.

// Prefer:

// ```text
// Input
//   ↓
// Function
//   ↓
// Expected behavior
// ```

// rather than testing exactly how the function internally works.

// ---

// # 31. Backend Testing Architecture

// A useful model for a Go backend:

// ```text
//                     TESTS
//                       |
//          +------------+------------+
//          |            |            |
//         Unit      Integration      API/E2E
//          |            |            |
//       Service      Repository     Handler
//       Utility         ↓             ↓
//          ↓        Database       Service
//    Fake/Mock                      ↓
//                               Repository
//                                   ↓
//                               Database
// ```

// Each level answers a different question.

// ---

// # 32. Common Commands

// Run tests for the current package:

// ```bash
// go test
// ```

// Run all package tests:

// ```bash
// go test ./...
// ```

// Verbose test output:

// ```bash
// go test -v
// ```

// Run a specific test:

// ```bash
// go test -run TestAdd
// ```

// Run benchmarks:

// ```bash
// go test -bench=BenchmarkAdd
// ```

// Benchmark with memory statistics:

// ```bash
// go test -bench=BenchmarkAdd -benchmem
// ```

// Race detection:

// ```bash
// go test -race ./...
// ```

// Coverage:

// ```bash
// go test -cover ./...
// ```

// Fuzzing:

// ```bash
// go test -fuzz=.
// ```

// ---

// # 33. Day 11 Interview Questions

// ### 1. What is `_test.go`?

// A Go source file ending in `_test.go` containing tests, benchmarks, fuzz tests, or examples recognized by the testing tool.

// ### 2. What is the difference between `t.Error()` and `t.Fatal()`?

// ```text
// Error → fail + continue
// Fatal → fail + stop current test
// ```

// ### 3. What is a table-driven test?

// A test pattern where multiple test cases are stored in a table and executed using the same test logic.

// ### 4. Why use `t.Run()`?

// To create named independent subtests and make failures easier to identify.

// ### 5. How do you test wrapped errors?

// Use:

// ```go
// errors.Is()
// errors.As()
// ```

// as appropriate.

// ### 6. Why use dependency injection in tests?

// To replace real external dependencies such as databases or APIs with fakes/mocks.

// ### 7. Unit test vs integration test?

// ```text
// Unit
// → isolated component

// Integration
// → multiple components working together
// ```

// ### 8. Does 100% coverage mean bug-free?

// No. Coverage only measures code execution.

// ### 9. What does `-race` do?

// Detects many data races.

// ### 10. What is `b.N`?

// The number of benchmark iterations chosen by the Go benchmarking framework.

// ### 11. What does `ns/op` mean?

// Average nanoseconds per benchmark operation.

// ### 12. What is fuzz testing?

// Automated testing that generates many inputs to find unexpected behavior or violations of a property.

// ### 13. What makes a test flaky?

// Timing assumptions, races, randomness, external dependencies, shared state, and similar nondeterministic factors.

// ### 14. Why shouldn't every repository automatically be mocked?

// Use test doubles where isolation is useful, but integration tests should also verify that real components work together.

// ---

// # 34. Day 11 Golden Mental Model

// ```text
// UNIT TEST
//     ↓
// Test one piece of logic
//     ↓
// Fast + isolated
// ```

// ```text
// TABLE-DRIVEN TEST
//     ↓
// One test structure
//     ↓
// Many scenarios
//     ↓
// t.Run()
// ```

// ```text
// DEPENDENCY INJECTION
//     ↓
// Interface
//     ↓
// Fake / Mock
//     ↓
// Isolated unit test
// ```

// ```text
// INTEGRATION TEST
//     ↓
// Real component interactions
//     ↓
// Database / external dependency
// ```

// ```text
// BENCHMARK
//     ↓
// b.N iterations
//     ↓
// Measure performance
//     ↓
// ns/op, B/op, allocs/op
// ```

// ```text
// -race
//     ↓
// Find many data races
// ```

// ```text
// COVERAGE
//     ↓
// Code execution measurement
//     ↓
// NOT proof of correctness
// ```

// The most important things to remember from Day 11 are:

// ```text
// Test       → verify behavior
// Table test → test many scenarios cleanly
// t.Error    → fail + continue
// t.Fatal    → fail + stop current test
// Interface  → inject fake/mock dependency
// Unit       → isolated
// Integration→ components together
// -race      → race detection
// Benchmark  → performance
// Fuzz       → discover unexpected inputs
// Coverage   → execution, not correctness
// ```
