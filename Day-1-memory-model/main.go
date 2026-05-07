package main

import "fmt"

// 🔹 Escape Analysis (VERY IMPORTANT)

// Go compiler decides:
// 👉 “Stack or Heap?”

// This process is called:

// 👉 Escape Analysis

// How to check if a variable is allocated on the stack or heap?
// 👉 Use the -gcflags="-m" flag when building or running your Go code.
// Example: go run -gcflags="-m" main.go

// 🚫 Heap = Slow
// GC overhead
// More latency
// Memory pressure

// ✅ Stack = Fast
// No GC
// Auto cleanup

// “Not all pointers go to heap”
// 👉 “Escape analysis decides, not developer”

func test() *int {
	a := 5
	return &a
}

func main() {
	// This variable will be allocated on the heap because it escapes the function scope.
	x := test()
	println(*x) // Output: 5

	test1() //stack allocation, no escape

	a := 5
	fmt.Println(a) // Heap

	y := 5 // Stack allocation, no escape
	z := &y
	fmt.Println(*z)
}

func test1() int {
	a := 5
	return a
}

// In the test1 function, the variable 'a' does not escape the function scope, so it will be allocated on the stack.

// Go compiler does:

// Escape analysis
// Inline optimizations
// Stack allocation if safe

// Interview question:

// 🧠 Level 1 (Fundamentals)
// ❓ Q1: What is the difference between stack and heap in Go?

// 👉 Expectation:

// Stack → fast, automatic, function-scoped
// Heap → GC managed, used when data escapes

// ❓ Q2: What is escape analysis in Go?

// 👉 Key point:

// Compiler decides whether variable goes to stack or heap
// Done at compile time using -gcflags="-m"

// ❓ Q3: Does returning a pointer always cause heap allocation?

// 👉 Expected answer:
// ❌ No (trick question)

// ✔ Mostly yes, but:

// Compiler may optimize in some cases
// Decision is made by escape analysis, not rule-based

// Level 2 (Tricky Practical)
// ❓ Q4: What will happen here?
// func test() *int {
//     x := 10
//     return &x
// }

// 👉 Answer:

// x escapes → heap

// ❓ Q5: Does this escape?
// func main() {
//     x := 10
//     fmt.Println(x)
// }

// 👉 Answer:

// Not necessarily (depends on compiler optimization)

// ❓ Q6: What about this?
// func main() {
//     x := 10
//     fmt.Println(&x)
// }

// 👉 Answer:

// Yes → escapes (address passed)

// ⚠️ Level 3 (Where most people fail)
// ❓ Q7: Why do interfaces often cause heap allocation?

// 👉 Expected answer:

// Interface = (type + value)
// May require boxing
// Compiler may move data to heap for safety

// ❓ Q8: Does this escape?
// func main() {
//     x := 10
//     y := x
//     fmt.Println(y)
// }

// 👉 Answer:

// No escape (pure value copy)

// ❓ Q9: Does taking pointer always mean heap?

// 👉 Answer:
// ❌ No

// ✔ Only if:

// Pointer escapes function
// Or compiler cannot guarantee safety

// ❓ Q10: How does escape analysis affect performance?

// 👉 Expected:

// Heap allocation → slower
// GC overhead increases
// More latency in high-scale systems

// ❓ Q11: How do you reduce heap allocations in Go?

// 👉 Expected:

// Avoid unnecessary pointers
// Avoid interface conversions
// Use value types when possible
// Reuse objects (sync.Pool)

// ❓ Q12: How do you check if variable escapes?

// 👉 Answer:

// go build -gcflags="-m"

// 🔥 Level 5 (Expert / System Design Thinking)
// ❓ Q13: Why is Go’s escape analysis important for backend services?

// 👉 Expected:

// High throughput systems need low GC pressure
// Heap allocations increase latency
// Efficient memory = better scalability

// ❓ Q14: What is a hidden escape in Go?

// 👉 Expected:

// Happens when:
// Using interfaces
// Using closures
// Returning references indirectly

// ❓ Q15: Analyze this:
// func get() interface{} {
//     x := 10
//     return x
// }

// 👉 Answer:

// x escapes (stored in interface)

// 🚀 Bonus Rapid-Fire (Very Important)

// Try answering quickly:

// Returning struct vs pointer → which is better? // Depends on size and use case
// Does slice always allocate on heap? // Not necessarily, depends on escape analysis
// Does closure cause escape? // Often yes, but depends on usage
// Does map always live on heap? // Yes, maps are reference types and typically allocated on the heap
// Can compiler keep heap variable on stack? // No, If variable is marked as escaping → must go to heap, Compiler cannot “bring it back”
