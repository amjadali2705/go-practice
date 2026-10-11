# Go Day 6 — Functions, Closures, Defer & Methods

Functions and methods are fundamental to Go. Understanding them deeply helps you write reusable code, manage resources safely, understand variable capture, and design maintainable backend services.

## 1. Functions in Go

A function is a reusable block of code that performs a task.

Basic syntax:

```go
func functionName(parameters) returnType {
    // function body
}
```

Example:

```go
func Add(a, b int) int {
    return a + b
}

func main() {
    result := Add(10, 20)
    fmt.Println(result)
}
```

Output:

```text
30
```

Here:
- `a` and `b` are parameters.
- `int` after the parentheses is the return type.
- `return a + b` returns the result.

## 2. Multiple Return Values

Go functions can return multiple values.

Example:

```go
func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }

    return a / b, nil
}
```

Usage:

```go
result, err := Divide(10, 2)

if err != nil {
    fmt.Println(err)
    return
}

fmt.Println(result)
```

Output:

```text
5
```

This pattern is extremely common in Go because functions often return a result along with an error.

**Remember:**

```text
value, err := function()
```

Always handle the error appropriately rather than assuming the operation succeeded.

## 3. Named Return Values

Go allows you to name return values in the function signature.

Example:

```go
func Add(a, b int) (result int) {
    result = a + b
    return
}
```

Here, `result` is a named return variable.

A bare `return` returns the current value of the named return variable.

You can also write:

```go
func Add(a, b int) (result int) {
    result = a + b
    return result
}
```

Both versions are valid.

Named return values can be useful when they improve clarity, but avoid using them when they make the code harder to understand.

## 4. Variadic Functions

A variadic function accepts zero or more arguments of the same declared type.

Example:

```go
func Sum(numbers ...int) int {
    total := 0

    for _, number := range numbers {
        total += number
    }

    return total
}
```

Usage:

```go
fmt.Println(Sum())
fmt.Println(Sum(10))
fmt.Println(Sum(10, 20, 30))
```

Output:

```text
0
10
60
```

The parameter:

```go
numbers ...int
```

is received inside the function as a slice of `int`.

You can pass an existing slice using `...`:

```go
nums := []int{10, 20, 30}
result := Sum(nums...)
```

Important: a variadic parameter must be the final parameter in the function's parameter list.

## 5. Functions as First-Class Values

In Go, functions can be assigned to variables, passed as arguments, and returned from other functions.

Example:

```go
func Add(a, b int) int {
    return a + b
}

func main() {
    operation := Add

    fmt.Println(operation(10, 20))
}
```

Output:

```text
30
```

Here, `operation` holds the function value.

Notice the difference:

```go
operation := Add
```

assigns the function itself.

```go
operation := Add(10, 20)
```

calls the function and assigns its result.

Functions as values are useful for callbacks, middleware, higher-order functions, and configurable behavior.

## 6. Anonymous Functions

An anonymous function is a function without a declared name.

Example:

```go
func() {
    fmt.Println("Hello")
}()
```

The final `()` calls the function immediately.

You can also assign an anonymous function to a variable:

```go
greet := func(name string) {
    fmt.Println("Hello", name)
}

greet("Go")
```

Output:

```text
Hello Go
```

Anonymous functions are frequently used for short callbacks, goroutines, and closures.

## 7. Closures

A **closure** is a function value that can access variables from its surrounding scope.

Example:

```go
func counter() func() int {
    count := 0

    return func() int {
        count++
        return count
    }
}
```

Usage:

```go
c := counter()

fmt.Println(c())
fmt.Println(c())
fmt.Println(c())
```

Output:

```text
1
2
3
```

Why?

The returned function accesses the `count` variable from the surrounding function. That variable remains available while the closure references it.

Conceptually:

```text
counter()
   |
   ├── count = 0
   |
   └── returned closure
          |
          └── accesses count
```

Each call to `counter()` creates a separate `count` variable for its returned closure.

Example:

```go
c1 := counter()
c2 := counter()

fmt.Println(c1()) // 1
fmt.Println(c1()) // 2
fmt.Println(c2()) // 1
```

`c1` and `c2` maintain separate state.

Closures are useful for callbacks, local state, middleware, and function factories.

## 8. A Subtle Closure Example

Consider:

```go
func main() {
    x := 10

    f := func() {
        fmt.Println(x)
    }

    x = 20

    f()
}
```

Output:

```text
20
```

The closure accesses the variable `x`, not an automatic snapshot of its initial value.

When the closure executes, `x` is `20`.

This distinction becomes important when functions capture variables that change over time.

## 9. `defer`

The `defer` statement schedules a function call to execute when the surrounding function returns.

Example:

```go
func main() {
    defer fmt.Println("Second")
    fmt.Println("First")
}
```

Output:

```text
First
Second
```

The deferred call runs after the regular statement.

Common uses include:
- Closing files.
- Releasing resources.
- Unlocking mutexes.
- Cleaning up temporary state.

Example:

```go
func readFile() error {
    file, err := os.Open("data.txt")
    if err != nil {
        return err
    }
    defer file.Close()

    // Read the file.
    return nil
}
```

After `readFile()` returns, `file.Close()` is called.

## 10. Multiple `defer` Calls Run in LIFO Order

LIFO means **Last In, First Out**.

Example:

```go
func main() {
    defer fmt.Println("First")
    defer fmt.Println("Second")
    defer fmt.Println("Third")
}
```

Output:

```text
Third
Second
First
```

The last deferred call registered executes first.

This is useful when releasing nested resources in the reverse order in which they were acquired.

## 11. Deferred Arguments Are Evaluated Immediately

This is a common Go interview question.

Example:

```go
func main() {
    x := 10

    defer fmt.Println(x)

    x = 20
}
```

Output:

```text
10
```

Why?

For this call:

```go
defer fmt.Println(x)
```

the argument `x` is evaluated when the `defer` statement executes. The deferred call then runs later with that saved argument.

### Compare with a deferred closure

```go
func main() {
    x := 10

    defer func() {
        fmt.Println(x)
    }()

    x = 20
}
```

Output:

```text
20
```

The deferred closure reads `x` when the closure runs.

**Golden rule:**

```text
defer function(x)
→ arguments are evaluated when defer is registered

defer func() { use(x) }()
→ closure reads the captured variable when it executes
```

## 12. `defer` and Return Values

The return process is important when `defer` is involved.

For an ordinary return expression, Go evaluates the expression, assigns the resulting return value, runs deferred calls, and then returns.

Example:

```go
func example() int {
    x := 10

    defer func() {
        x = 20
    }()

    return x
}
```

Output from calling the function:

```text
10
```

Why?

The expression `return x` evaluates to `10` before the deferred function changes the local variable.

The deferred function changes `x`, but the already evaluated return value remains `10`.

### Named return values are different

```go
func example() (x int) {
    x = 10

    defer func() {
        x = 20
    }()

    return
}
```

This returns:

```text
20
```

Because `x` is the named return variable. The deferred function modifies that variable before the function finishes returning.

This distinction is useful to understand, though named return values should be used thoughtfully.

## 13. Methods

A method is a function associated with a receiver type.

Example:

```go
type User struct {
    Name string
}

func (u User) Greet() {
    fmt.Println("Hello", u.Name)
}
```

Usage:

```go
u := User{Name: "Amjad"}
u.Greet()
```

Output:

```text
Hello Amjad
```

The receiver:

```go
(u User)
```

associates `Greet` with the `User` type.

A method is called using:

```go
u.Greet()
```

whereas a standalone function is called by its function name.

## 14. Functions vs Methods

Function:

```go
func Greet(name string) {
    fmt.Println("Hello", name)
}
```

Call:

```go
Greet("Amjad")
```

Method:

```go
type User struct {
    Name string
}

func (u User) Greet() {
    fmt.Println("Hello", u.Name)
}
```

Call:

```go
u := User{Name: "Amjad"}
u.Greet()
```

A method has a receiver and can operate on the receiver value.

Use methods when the behavior naturally belongs to a type.

## 15. Value Receivers

A value receiver receives a copy of the receiver value.

Example:

```go
type User struct {
    Age int
}

func (u User) IncreaseAge() {
    u.Age++
}
```

Usage:

```go
u := User{Age: 24}
u.IncreaseAge()

fmt.Println(u.Age)
```

Output:

```text
24
```

Why?

`IncreaseAge()` modifies the copy, not the original `User`.

This assumes the receiver contains ordinary value data. Copying a struct does not deep-copy everything it refers to: fields such as pointers, slices, maps, or other reference-like values can still refer to shared underlying data.

## 16. Pointer Receivers

A pointer receiver receives a pointer to the original value.

Example:

```go
type User struct {
    Age int
}

func (u *User) IncreaseAge() {
    u.Age++
}
```

Usage:

```go
u := User{Age: 24}
u.IncreaseAge()

fmt.Println(u.Age)
```

Output:

```text
25
```

Why?

The pointer receiver allows the method to modify the original struct.

The receiver:

```go
(u *User)
```

is a pointer to `User`.

### Value vs pointer receiver

```text
Value receiver:   (u User)
→ receives a copy of the struct

Pointer receiver: (u *User)
→ receives a pointer to the struct
```

Use a pointer receiver when you need to modify the original value or want to avoid copying a larger struct.

Consistency also matters: types commonly use pointer receivers for most or all of their methods when mutation or shared identity is important.

## 17. Automatic Addressing for Method Calls

Go can automatically take the address of an addressable variable when calling a pointer-receiver method.

Example:

```go
type User struct {
    Age int
}

func (u *User) IncreaseAge() {
    u.Age++
}

func main() {
    u := User{Age: 24}
    u.IncreaseAge()
}
```

This works even though `u` is a value rather than a pointer. Go effectively calls the method using `&u`.

However, this convenience does not mean `User` and `*User` have the same method sets. That distinction becomes important when satisfying interfaces.

## 18. Method Sets

The method set determines which methods a type has for interface satisfaction and other type-system rules.

Consider:

```go
type Speaker interface {
    Speak()
}

type Dog struct{}

func (d *Dog) Speak() {
    fmt.Println("Woof")
}
```

Here, `Speak()` has a pointer receiver.

Therefore:

```go
var s Speaker = &Dog{} // Valid
```

But:

```go
var s Speaker = Dog{} // Compile error
```

Why?

The method set of `*Dog` includes methods declared with receivers `Dog` and `*Dog`. The method set of `Dog` includes methods declared with a `Dog` receiver, but not methods that require a `*Dog` receiver.

### Remember

```text
Value receiver method
→ available in method sets of T and *T

Pointer receiver method
→ available in method set of *T
```

A compiler may automatically take the address for a normal method call on an addressable value, but interface satisfaction follows the method-set rules.

## 19. When Should You Use Value or Pointer Receivers?

Use a pointer receiver when:
- The method must modify the original value.
- Copying the receiver would be expensive.
- Consistency with the type's other methods calls for a pointer receiver.

Use a value receiver when:
- The value is small and naturally treated as immutable.
- A copy is appropriate.
- You do not need to change the original receiver.

Example:

```go
type Counter struct {
    Value int
}

func (c *Counter) Increment() {
    c.Value++
}

func (c Counter) Current() int {
    return c.Value
}
```

A pointer receiver is appropriate for `Increment()` because it changes the receiver. A value receiver can be appropriate for `Current()` because it only reads the value.

For consistency, a type whose methods operate on shared state often uses pointer receivers throughout.

## 20. Closures in Loops

Loop variables captured by closures have an important language-version detail.

In Go 1.22 and later, when a variable is declared by the loop itself, each iteration gets its own variable.

Example with modern Go semantics:

```go
func main() {
    for i := 0; i < 3; i++ {
        go func() {
            fmt.Println(i)
        }()
    }

    time.Sleep(time.Second)
}
```

This commonly demonstrates that the goroutines capture per-iteration variables, so the values are `0`, `1`, and `2`, although their print order is not guaranteed.

**Important:** `time.Sleep` here is only a demonstration technique and is not reliable synchronization. Use a `sync.WaitGroup` to wait for goroutines in real code.

The Go 1.22 loop-variable change applies to loops whose variables are declared by the loop. Variables declared outside the loop and reused across iterations can still be shared by closures.

## 21. Closures and Concurrency

Even with modern loop-variable semantics, be careful when multiple goroutines access the same mutable variable.

Example:

```go
counter := 0

for i := 0; i < 3; i++ {
    go func() {
        counter++
    }()
}
```

All goroutines access the shared `counter`, so there is a data race.

Per-iteration loop variables do not automatically make shared state safe.

Use a mutex or an appropriate atomic operation to protect shared state, and use a `WaitGroup` if completion must be synchronized.

## 22. Common Mistakes

### Mistake 1: Assuming a value receiver modifies the original struct

```go
func (u User) Update() {
    u.Age++
}
```

This modifies the receiver copy, not the original struct's ordinary value field.

### Mistake 2: Assuming `defer` executes immediately

`defer` registers a call for later execution when the surrounding function returns or unwinds.

### Mistake 3: Forgetting that deferred arguments are evaluated immediately

```go
defer fmt.Println(x)
```

saves the current argument value.

### Mistake 4: Confusing closures with snapshots

A closure can access variables from its surrounding scope. Captured variables do not automatically become copies.

### Mistake 5: Ignoring method sets

A pointer-receiver method can mean `*T` satisfies an interface while `T` does not.

### Mistake 6: Assuming `time.Sleep()` is synchronization

It only pauses execution. Use a `WaitGroup`, channel, or another suitable synchronization mechanism when ordering or completion matters.

---

## 23. Day 6 Interview Questions

**1. What is a variadic function?**

A function that accepts zero or more arguments of the same declared type using `...T`.

**2. What are multiple return values used for in Go?**

They allow functions to return multiple results, commonly a value and an error.

**3. What is a closure?**

A function value that can access variables from its surrounding scope.

**4. What does `defer` do?**

It schedules a function call to execute when the surrounding function returns or unwinds.

**5. In what order do multiple deferred calls execute?**

Last In, First Out (LIFO).

**6. When are arguments to a deferred call evaluated?**

When the `defer` statement executes, not when the deferred call finally runs.

**7. What is the difference between a function and a method?**

A method has a receiver associated with a type; a function does not.

**8. What is a value receiver?**

A receiver that receives a copy of the receiver value.

**9. What is a pointer receiver?**

A receiver that receives a pointer to the original value.

**10. Why use pointer receivers?**

To modify the original value, avoid unnecessary copies, or maintain consistent receiver semantics.

**11. What is a method set?**

The set of methods associated with a type for language rules such as interface satisfaction.

**12. Why can `*Dog` satisfy an interface while `Dog` doesn't?**

A required method may have a pointer receiver, making it part of the method set of `*Dog` but not `Dog`.

**13. Does a closure automatically copy captured variables?**

No. A closure can access the variables in its surrounding scope.

**14. What is the difference between these two statements?**

```go
defer fmt.Println(x)
```

```go
defer func() {
    fmt.Println(x)
}()
```

The first evaluates `x` when the defer statement executes. The second closure reads the captured `x` when it executes.

---

## 24. Day 6 Golden Mental Model

```text
Function
→ reusable block of code

Variadic function
→ accepts zero or more arguments

Function value
→ a function can be assigned, passed, or returned

Closure
→ function + access to surrounding variables

defer
→ executes when the function returns or unwinds

Value receiver
→ operates on a receiver copy

Pointer receiver
→ can modify the original receiver

Method set
→ determines method availability for type-system rules
```

### Final takeaway

The most important distinctions to remember are:

```text
Closure
→ captures access to surrounding variables

defer call
→ executes later, but its arguments are evaluated immediately

Value receiver
→ receives a copy

Pointer receiver
→ receives a pointer to the original

Automatic address-taking
→ can make pointer-receiver method calls convenient

Method set
→ still determines interface satisfaction
```

Mastering these concepts will help you write clearer Go code, reason about return values and resource cleanup, and avoid common mistakes when using methods, closures, and interfaces.
