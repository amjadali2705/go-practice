# Go Day 7 — Interfaces

Interfaces are one of the most important concepts in Go. They help you write flexible, maintainable code by allowing different types to provide the same behavior.

For backend development, interfaces are especially useful for repository abstractions, dependency injection, testing, and designing loosely coupled services.

## 1. What Is an Interface?

An **interface** defines a set of method requirements that a type must satisfy to implement it.

Example:

```go
type Speaker interface {
    Speak()
}
```

This interface says:

> Any type with a `Speak()` method matching this signature can satisfy `Speaker`.

Now create two types:

```go
type Dog struct{}
type Cat struct{}

func (d Dog) Speak() {
    fmt.Println("Woof")
}

func (c Cat) Speak() {
    fmt.Println("Meow")
}
```

Both `Dog` and `Cat` satisfy the `Speaker` interface.

We can write:

```go
func MakeSpeak(s Speaker) {
    s.Speak()
}
```

Usage:

```go
MakeSpeak(Dog{})
MakeSpeak(Cat{})
```

Output:

```text
Woof
Meow
```

The function doesn't need to know whether it receives a `Dog` or a `Cat`. It only needs a value that satisfies `Speaker`.

---

## 2. Go Interfaces Are Implemented Implicitly

In some languages, a type must explicitly declare that it implements an interface. Go doesn't require that declaration.

Example:

```go
type Writer interface {
    Write()
}

type File struct{}

func (f File) Write() {
    fmt.Println("Writing file")
}
```

`File` satisfies `Writer` automatically because it has the required method.

You do not write something like:

```go
// Go does not require an explicit declaration like this:
// type File implements Writer
```

**Remember:**

```text
Required methods match
        ↓
Type satisfies interface
        ↓
No explicit declaration required
```

This is called implicit interface implementation.

---

## 3. Why Do We Need Interfaces?

Suppose you have two payment methods:

```go
type CardPayment struct{}
type UpiPayment struct{}

func (c CardPayment) Pay() {
    fmt.Println("Card payment")
}

func (u UpiPayment) Pay() {
    fmt.Println("UPI payment")
}
```

Instead of creating a separate function for each payment type, define an interface:

```go
type PaymentMethod interface {
    Pay()
}
```

Then:

```go
func ProcessPayment(p PaymentMethod) {
    p.Pay()
}
```

Usage:

```go
ProcessPayment(CardPayment{})
ProcessPayment(UpiPayment{})
```

Now `ProcessPayment` can accept different implementations without changing its own code.

This is useful when you want to separate business logic from specific implementations.

---

## 4. Interface Values

An interface value conceptually contains two pieces of information:

```text
Dynamic type + Dynamic value
```

Example:

```go
var s Speaker = Dog{}
```

Conceptually:

```text
Interface value
├── Dynamic type: Dog
└── Dynamic value: Dog{}
```

Another example:

```go
var s Speaker = Cat{}
```

Conceptually:

```text
Interface value
├── Dynamic type: Cat
└── Dynamic value: Cat{}
```

The interface provides a common way to call methods even when the underlying concrete type changes.

This model is useful for understanding interfaces, especially when working with `nil`.

---

## 5. The Nil Interface Trap

This is one of the most important Go interview topics.

Consider:

```go
var x interface{}

fmt.Println(x == nil)
```

Output:

```text
true
```

The interface has neither a dynamic type nor a dynamic value.

Now consider:

```go
var p *int = nil
var x interface{} = p

fmt.Println(x == nil)
```

Output:

```text
false
```

Why?

The interface contains:

```text
Dynamic type:  *int
Dynamic value: nil
```

The interface itself is not nil because it has a dynamic type.

### Remember

```text
Nil interface
├── Dynamic type:  none
└── Dynamic value: none

Interface containing typed nil
├── Dynamic type:  *int
└── Dynamic value: nil
```

This can cause bugs when returning pointers through interfaces.

---

## 6. A Practical Typed-Nil Example

Consider an interface:

```go
type User interface {
    GetName() string
}
```

A concrete type:

```go
type Person struct{}

func (p *Person) GetName() string {
    if p == nil {
        return "nil person"
    }

    return "Amjad"
}
```

Now:

```go
var p *Person = nil
var u User = p

fmt.Println(u == nil)
fmt.Println(u.GetName())
```

Output:

```text
false
nil person
```

The interface isn't nil, so `u == nil` returns `false`. Its dynamic value is a nil pointer, and the method can still be called because `GetName` handles a nil receiver.

If the method dereferenced `p` without checking it, it could panic.

**Practical lesson:** When working with interfaces that may contain pointers, understand both the interface's nil status and the underlying pointer's value.

---

## 7. The `any` Type

In modern Go, `any` is an alias for `interface{}`.

These are equivalent:

```go
var a any
var b interface{}
```

An empty interface has no required methods, so it can hold a value of any type.

Example:

```go
func PrintValue(v any) {
    fmt.Println(v)
}
```

Usage:

```go
PrintValue(100)
PrintValue("hello")
PrintValue(true)
```

This is useful when a function genuinely needs to accept arbitrary values. However, `any` provides no specific method requirements, so it is less restrictive than a custom interface.

---

## 8. Type Assertion

A type assertion lets you retrieve a concrete value from an interface.

Syntax:

```go
value := x.(T)
```

Example:

```go
var x any = "hello"

s := x.(string)

fmt.Println(s)
```

Output:

```text
hello
```

But there is a risk: if the dynamic type doesn't match, a single-value assertion panics.

Example:

```go
var x any = 100

s := x.(string) // Panics
```

### Safe type assertion

Use the two-value form:

```go
value, ok := x.(T)
```

Example:

```go
var x any = 100

s, ok := x.(string)

fmt.Println(s, ok)
```

Output:

```text
 false
```

The zero value for `string` is an empty string, so `s` is `""`, while `ok` is `false`.

**Rule:**

```text
x.(T)       → panics if assertion fails
x.(T), ok   → safe; reports success or failure
```

Use the two-value form when the dynamic type is not guaranteed.

---

## 9. Type Switch

A type switch lets you handle different concrete types held inside an interface.

Example:

```go
func Describe(v any) {
    switch value := v.(type) {
    case int:
        fmt.Println("Integer:", value)

    case string:
        fmt.Println("String:", value)

    case bool:
        fmt.Println("Boolean:", value)

    default:
        fmt.Println("Unknown type")
    }
}
```

Usage:

```go
Describe(100)
Describe("hello")
Describe(true)
Describe(3.14)
```

Output:

```text
Integer: 100
String: hello
Boolean: true
Unknown type
```

Inside each case, `value` has the relevant concrete type.

Use a type switch when different concrete types require different handling.

---

## 10. Method Sets and Pointer Receivers

Interfaces work closely with methods and receivers.

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

The method uses a pointer receiver:

```go
func (d *Dog) Speak()
```

Now:

```go
var s Speaker = &Dog{}
```

is valid.

But:

```go
var s Speaker = Dog{}
```

does not compile.

Why?

A method declared with a pointer receiver belongs to the method set of `*Dog`, not generally the method set of `Dog`.

### Method-set rule

```text
Method with value receiver
    ↓
T and *T have the method

Method with pointer receiver
    ↓
*T has the method
T does not have the method
```

This matters when assigning values to interfaces.

A method call on an addressable variable may allow Go to take its address automatically, but that does not change the method-set rules for interface satisfaction.

---

## 11. Compile-Time Interface Assertion

Sometimes you want the compiler to verify that a type satisfies an interface.

Use:

```go
var _ Speaker = (*Dog)(nil)
```

This means:

> Verify at compile time that `*Dog` satisfies `Speaker`.

If `*Dog` doesn't implement the required methods, the compiler reports an error.

It does not create a useful runtime `Dog` instance; it is a compile-time check.

For example:

```go
type Repository interface {
    GetUser(id int) (*User, error)
}

type MyRepository struct{}

func (r *MyRepository) GetUser(id int) (*User, error) {
    return nil, nil
}

var _ Repository = (*MyRepository)(nil)
```

This can help catch accidental interface implementation changes during refactoring.

---

## 12. Small Interfaces

In Go, small interfaces are often easier to reuse and test.

Example:

```go
type Reader interface {
    Read() error
}
```

Prefer interfaces that describe the behavior required by their consumers rather than creating large interfaces containing every method a type supports.

A useful principle:

> Accept the behavior you need, not a larger concrete abstraction than necessary.

For example, a function that only needs to log a message doesn't need a full logging system interface if a small interface is sufficient.

---

## 13. Interfaces in Backend Development

Suppose you have a user service:

```go
type UserRepository interface {
    GetUser(id int) (*User, error)
    CreateUser(user *User) error
}
```

The service depends on this interface:

```go
type UserService struct {
    repo UserRepository
}
```

Its constructor could be:

```go
func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

Production code can inject a repository backed by MariaDB or another database.

Tests can inject a fake repository.

Conceptually:

```text
                UserService
                     |
             UserRepository
                     |
              +------+------+
              |             |
        Real Repository  Fake Repository
              |             |
           Database       Unit Test
```

This is a form of dependency injection.

Benefits include:

- Reduced coupling between service and database implementation.
- Easier unit testing.
- Ability to replace implementations without rewriting the service.
- Clearer separation of responsibilities.

A concrete repository can implement the interface implicitly.

---

## 14. Interfaces and Testing

A service that directly creates its database connection is difficult to test in isolation.

Instead, inject a dependency through an interface.

Example:

```go
type UserRepository interface {
    GetUser(id int) (*User, error)
}

type UserService struct {
    repo UserRepository
}

func (s *UserService) GetUser(id int) (*User, error) {
    return s.repo.GetUser(id)
}
```

A test can use a fake:

```go
type FakeUserRepository struct {
    User *User
    Err  error
}

func (f *FakeUserRepository) GetUser(id int) (*User, error) {
    return f.User, f.Err
}
```

Now the test can supply predictable data or errors without connecting to a real database.

This connects Day 7 interfaces directly to Day 11 testing.

---

## 15. Interfaces vs Generics

Interfaces and generics solve different problems.

### Interface

Use when you care about behavior.

```go
type Speaker interface {
    Speak()
}

func MakeSpeak(s Speaker) {
    s.Speak()
}
```

The function accepts any value that satisfies the required behavior.

### Generics

Use when you want reusable code across types while preserving type relationships.

```go
func Identity[T any](value T) T {
    return value
}
```

The function returns the same type it receives.

### Remember

```text
Interface
→ common behavior

Generics
→ common code across types
```

They can complement each other, and neither replaces the other.

---

## 16. Common Interface Mistakes

### Mistake 1: Assuming a typed nil interface is nil

```go
var p *Dog = nil
var s Speaker = p

fmt.Println(s == nil) // false
```

The interface has a dynamic type, so it is not nil.

### Mistake 2: Using an unsafe type assertion

```go
s := value.(string)
```

If the type isn't `string`, it panics.

Prefer:

```go
s, ok := value.(string)
```

when the type isn't guaranteed.

### Mistake 3: Forgetting pointer receiver rules

A type with a pointer-receiver method may not satisfy an interface when a value of the non-pointer type is assigned.

### Mistake 4: Making interfaces unnecessarily large

Small interfaces often make dependencies easier to replace and tests easier to write.

### Mistake 5: Defining interfaces without a real need

Interfaces should represent useful behavior or an abstraction boundary. Don't automatically create an interface for every struct.

---

## 17. Day 7 Interview Questions

**1. What is an interface?**

A type that specifies a set of required methods.

**2. How does a type implement an interface in Go?**

Implicitly, by implementing the required methods with matching signatures.

**3. What is `any`?**

An alias for `interface{}`.

**4. What does an interface value conceptually contain?**

Its dynamic type and dynamic value.

**5. Why can an interface containing a nil pointer be non-nil?**

Because the interface still has a dynamic type, even though its dynamic value is nil.

**6. What is a type assertion?**

A mechanism for retrieving a concrete value from an interface.

**7. How do you prevent a failed type assertion from panicking?**

Use the two-value form:

```go
value, ok := x.(T)
```

**8. What is a type switch?**

A switch that selects a case based on the dynamic type stored in an interface.

**9. Why does `*Dog` satisfy an interface while `Dog` doesn't?**

A required method may use a pointer receiver, making it part of `*Dog`'s method set but not `Dog`'s.

**10. What is a compile-time interface assertion?**

A declaration such as:

```go
var _ Speaker = (*Dog)(nil)
```

that verifies interface satisfaction during compilation.

**11. How do interfaces help testing?**

They allow real dependencies to be replaced with fake or mock implementations.

**12. What is the difference between interfaces and generics?**

Interfaces abstract behavior, while generics support reusable code across types.

---

## 18. Day 7 Golden Mental Model

```text
Interface
   ↓
Defines required methods
   ↓
Concrete type implements them
implicitly
   ↓
Value can be used through interface
```

For interface values:

```text
Interface value
├── Dynamic type
└── Dynamic value
```

For nil:

```text
nil interface
→ no dynamic type or value

interface containing typed nil
→ dynamic type exists, dynamic value is nil
```

For assertions:

```text
x.(T)
→ can panic

x.(T), ok
→ safe type check
```

For architecture:

```text
Interface
   ↓
Dependency Injection
   ↓
Real Implementation / Fake Implementation
   ↓
Production / Testing
```

### Final takeaways

- Go interfaces are implemented implicitly.
- An interface describes behavior through method requirements.
- An interface value conceptually contains a dynamic type and dynamic value.
- A typed nil inside an interface does not make the interface itself nil.
- Use safe type assertions when necessary.
- Understand method sets and pointer receivers.
- Keep interfaces small and useful.
- Use interfaces to decouple services from concrete dependencies.
