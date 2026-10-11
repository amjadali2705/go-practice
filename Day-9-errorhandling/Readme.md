# Go Day 9 — Error Handling

## 1. What Is an Error in Go?

In Go, an error is an interface defined by the standard library:

```go
type error interface {
    Error() string
}
```

Any type that implements the `Error() string` method can satisfy the `error` interface.

Go uses explicit error handling rather than exceptions for most expected failures.

Example:

```go
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }

    return a / b, nil
}
```

Usage:

```go
result, err := divide(10, 0)

if err != nil {
    fmt.Println("Error:", err)
    return
}

fmt.Println(result)
```

**Key rule:** A `nil` error usually means success; a non-`nil` error indicates failure.

---

## 2. `errors.New()`

Use `errors.New()` to create an error with a fixed message.

```go
import "errors"

err := errors.New("user not found")
```

Example:

```go
func getUser(id int) error {
    if id <= 0 {
        return errors.New("invalid user ID")
    }

    return nil
}
```

Use it when the error message is known in advance.

---

## 3. `fmt.Errorf()`

Use `fmt.Errorf()` to create a formatted error message.

```go
err := fmt.Errorf("user %d not found", 101)
```

Output:

```text
user 101 not found
```

Unlike `fmt.Println()`, `fmt.Errorf()` creates and returns an error; it does not print the message.

It can also wrap an existing error when `%w` is used.

---

## 4. Error Wrapping — `%w` vs `%v`

This is one of the most important concepts in Go error handling.

### Using `%w`

```go
return fmt.Errorf("get user failed: %w", err)
```

`%w` wraps the original error and preserves the error chain.

### Using `%v`

```go
return fmt.Errorf("get user failed: %v", err)
```

`%v` formats the underlying error's text into a new error message, but does not preserve the original error as a wrapped error.

### Example

```go
var ErrNotFound = errors.New("not found")

func repository() error {
    return ErrNotFound
}

func service() error {
    err := repository()
    return fmt.Errorf("service failed: %w", err)
}
```

Usage:

```go
err := service()

fmt.Println(err)
fmt.Println(errors.Is(err, ErrNotFound))
```

Output:

```text
service failed: not found
true
```

The original error is still detectable because `%w` preserved it.

**Remember:**

```text
%w → wrap error + preserve error chain
%v → format error text only
```

---

## 5. `errors.Is()`

`errors.Is()` checks whether an error matches a target error, including through its wrapping chain.

Example:

```go
var ErrUserNotFound = errors.New("user not found")

func getUser() error {
    return fmt.Errorf("repository failed: %w", ErrUserNotFound)
}
```

Test:

```go
err := getUser()

fmt.Println(errors.Is(err, ErrUserNotFound))
```

Output:

```text
true
```

This works even though the original error was wrapped.

### Avoid comparing error strings

Bad:

```go
if err.Error() == "user not found" {
    // handle error
}
```

Good:

```go
if errors.Is(err, ErrUserNotFound) {
    // handle error
}
```

String comparisons are fragile because messages can change. `errors.Is()` is designed for identifying errors and their wrapped forms.

---

## 6. Sentinel Errors

A **sentinel error** is a predefined error value used to represent a recognizable failure condition.

Example:

```go
var (
    ErrUserNotFound = errors.New("user not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrUserExists   = errors.New("user already exists")
)
```

Use it like this:

```go
if errors.Is(err, ErrUserNotFound) {
    // handle not found
}
```

Sentinel errors are often declared at package level, but the key concept is that callers can identify the predefined error value.

Use clear names, commonly beginning with `Err`.

**Terminology note:** `ErrUserNotFound` is a predefined error value, not a Go predeclared identifier such as `int`, `string`, or `nil`.

---

## 7. `errors.As()`

`errors.As()` searches an error chain for an error that can be assigned to a target type.

Example:

```go
type ValidationError struct {
    Field string
    Msg   string
}

func (e *ValidationError) Error() string {
    return e.Field + ": " + e.Msg
}
```

Create and wrap the error:

```go
original := &ValidationError{
    Field: "email",
    Msg:   "invalid email",
}

err := fmt.Errorf("request validation failed: %w", original)
```

Detect its type:

```go
var validationErr *ValidationError

if errors.As(err, &validationErr) {
    fmt.Println(validationErr.Field)
    fmt.Println(validationErr.Msg)
}
```

Output:

```text
email
invalid email
```

`errors.As()` gives access to the matched error value, including its fields.

### `errors.Is()` vs `errors.As()`

| Function | Purpose |
|---|---|
| `errors.Is()` | Checks whether an error matches a target error |
| `errors.As()` | Finds an error matching a specified type |

Mental model:

```text
errors.Is() → Does this error match ErrNotFound?

errors.As() → Is there a *ValidationError in the chain?
```

---

## 8. Custom Error Types

Custom error types are useful when an error needs to carry structured information.

```go
type APIError struct {
    Code    int
    Message string
}

func (e *APIError) Error() string {
    return e.Message
}
```

Create an error:

```go
err := &APIError{
    Code:    400,
    Message: "invalid request",
}
```

Read the custom fields:

```go
var apiErr *APIError

if errors.As(err, &apiErr) {
    fmt.Println(apiErr.Code)
}
```

Output:

```text
400
```

Use custom error types when callers need more than a simple error message.

For example, a validation error might include a field name, while a domain error might identify a specific business failure.

---

## 9. What Is an Error Chain?

An error chain is formed when errors wrap other errors.

Example:

```go
dbErr := errors.New("connection refused")

repoErr := fmt.Errorf("repository failed: %w", dbErr)

serviceErr := fmt.Errorf("service failed: %w", repoErr)

handlerErr := fmt.Errorf("request failed: %w", serviceErr)
```

Conceptually:

```text
request failed
       ↓
service failed
       ↓
repository failed
       ↓
connection refused
```

The outer error provides context, while the inner error represents the original cause.

`errors.Is()` and `errors.As()` can inspect the chain.

### Important warning

If an intermediate layer uses `%v` instead of `%w`, it does not preserve the original error as part of the wrapped chain.

```go
return fmt.Errorf("handler failed: %v", err)
```

This can prevent callers from detecting an underlying sentinel error with `errors.Is()`.

Use `%w` when preserving the underlying error is appropriate.

---

## 10. Panic vs Error

Use ordinary errors for expected failures.

Examples:

```text
Invalid input
User not found
Database unavailable
Network timeout
Duplicate record
```

Example:

```go
func getUser(id int) (*User, error) {
    if id <= 0 {
        return nil, errors.New("invalid user ID")
    }

    // Fetch the user...
    return nil, nil
}
```

Use `panic()` for exceptional situations such as violated program invariants where normal execution cannot reasonably continue.

```go
panic("impossible state")
```

Do not use panic as the normal way to report invalid requests or database failures.

**Rule:**

```text
Expected failure → return error
Exceptional invariant/programming failure → panic may be appropriate
```

---

## 11. Understanding `defer`

A deferred function runs when the surrounding function returns or during panic unwinding.

Example:

```go
func main() {
    defer fmt.Println("deferred")

    fmt.Println("main")
}
```

Output:

```text
main
deferred
```

Deferred calls run in last-in, first-out (LIFO) order.

```go
func main() {
    defer fmt.Println("first")
    defer fmt.Println("second")
    defer fmt.Println("third")
}
```

Output:

```text
third
second
first
```

---

## 12. `recover()`

`recover()` can stop a panic when it is called directly from a deferred function executing in the same goroutine.

Example:

```go
func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
        }
    }()

    fmt.Println("before")
    panic("boom")
    fmt.Println("after")
}
```

Output:

```text
before
Recovered: boom
```

`after` is not printed because the function does not resume at the line following `panic()`.

### Important rules

1. `recover()` must be called from a deferred function to recover an active panic.
2. A panic in one goroutine cannot be recovered by a different goroutine.
3. Recovery prevents the panic from continuing to propagate in that goroutine.
4. After recovery, the panicking function returns; execution does not continue from the line after `panic()`.

### Recovering in a goroutine

```go
go func() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
        }
    }()

    panic("boom")
}()
```

The recovery must be inside this goroutine.

Avoid adding `recover()` everywhere without a reason. Recovery is most useful at deliberate boundaries, such as preventing a panic from crashing a server request-handling goroutine.

---

## 13. Converting a Panic into an Error

A deferred function can recover a panic and assign an error to a named return value.

```go
func process() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()

    panic("something failed")
}
```

The named return variable `err` can be modified by the deferred function.

However, this pattern should be used deliberately. It should not replace normal error handling.

---

## 14. Error Handling in a Go Backend

A typical backend architecture looks like this:

```text
HTTP Handler
     ↓
Service
     ↓
Repository / DAO
     ↓
Database
```

Each layer should have a clear responsibility.

### Repository

Works with the database and returns data or errors.

```go
return ErrUserNotFound
```

### Service

Applies business logic and may add useful context.

```go
return fmt.Errorf("get user: %w", err)
```

### Handler

Classifies the error and produces the appropriate HTTP response.

```go
if errors.Is(err, ErrUserNotFound) {
    c.JSON(http.StatusNotFound, gin.H{
        "error": "user not found",
    })
    return
}
```

Keep infrastructure-specific details out of the public API response unless there is a deliberate, safe reason to expose them.

---

## 15. Never Leak Sensitive Internal Errors to Clients

Suppose the database returns:

```text
dial tcp 10.20.4.5:3306: connection refused
```

Returning that exact message to an API client could expose infrastructure details.

Instead, return an appropriate public error:

```json
{
    "error": "internal server error"
}
```

Log the technical error internally:

```go
logger.Error("database failure", zap.Error(err))
```

The principle is:

```text
Client → safe, useful response

Server logs → detailed technical context
```

Also avoid logging secrets, credentials, tokens, or sensitive request data.

---

## 16. Common Mistakes

### Mistake 1: Comparing error strings

Avoid:

```go
if err.Error() == "not found" {
}
```

Prefer:

```go
errors.Is(err, ErrNotFound)
```

### Mistake 2: Losing the error chain

Avoid `%v` when callers need to identify the original error.

```go
return fmt.Errorf("service failed: %w", err)
```

### Mistake 3: Ignoring errors

Avoid:

```go
result, _ := divide(10, 0)
```

unless ignoring the error is intentional and safe.

Prefer:

```go
result, err := divide(10, 0)
if err != nil {
    return err
}
```

### Mistake 4: Using panic for normal failures

Invalid input and a missing database record normally should be handled by returning errors.

### Mistake 5: Assuming `recover()` resumes execution

It stops panic propagation when used correctly, but execution does not return to the line immediately following `panic()`.

### Mistake 6: Wrapping every error without adding useful context

Good wrapping should help identify where the failure occurred.

```go
return fmt.Errorf("create user: %w", err)
```

Avoid repeatedly adding meaningless messages at every layer.

---

## 17. Day 9 Interview Questions

**1. What is an error in Go?**

An interface with an `Error() string` method.

**2. What is the difference between `errors.New()` and `fmt.Errorf()`?**

`errors.New()` creates a fixed-message error. `fmt.Errorf()` creates a formatted error and can wrap another error using `%w`.

**3. What is `%w`?**

A formatting verb that wraps an error and preserves its place in the error chain.

**4. What is the difference between `%w` and `%v`?**

`%w` wraps the error; `%v` formats its text without wrapping it.

**5. What does `errors.Is()` do?**

Checks whether an error matches a target, including through its wrapping chain.

**6. What does `errors.As()` do?**

Searches the error chain for an error assignable to a specified type.

**7. What is a sentinel error?**

A predefined error value commonly used to identify a recognizable failure condition.

**8. What is an error chain?**

A sequence of errors connected through wrapping.

**9. When should you use panic?**

For exceptional situations such as violated program invariants, rather than normal, expected failures.

**10. What does `recover()` do?**

Stops an active panic when called correctly from a deferred function in the same goroutine.

**11. Can a panic in one goroutine be recovered by another?**

No.

**12. Why shouldn't raw database errors always be returned to clients?**

They may expose sensitive infrastructure details and are often unsuitable as public API messages.

**13. Why wrap an error with `%w`?**

To add useful context while preserving the underlying error for `errors.Is()` and `errors.As()`.

---

## 18. Day 9 Golden Mental Model

```text
NORMAL FAILURE
      ↓
return error
      ↓
caller checks err != nil
      ↓
add context if useful
      ↓
errors.Is() / errors.As()
      ↓
handle appropriately
```

Error wrapping:

```text
Handler error
      ↓
Service error
      ↓
Repository error
      ↓
Database error
```

Panic recovery:

```text
panic
  ↓
deferred function
  ↓
recover()
  ↓
panic stops propagating
```

### Key takeaways

- `errors.New()` creates a fixed-message error.
- `fmt.Errorf()` creates a formatted error.
- `%w` wraps and preserves the error chain.
- `%v` formats the message without wrapping.
- `errors.Is()` checks an error target.
- `errors.As()` finds a matching error type.
- Sentinel errors represent predefined failure conditions.
- Expected failures should generally return errors.
- `recover()` must be used correctly in the same goroutine.
- Keep public error responses safe and logs informative.
