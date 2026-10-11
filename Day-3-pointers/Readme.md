# Go Day 3 — Pointers

Pointers are an important Go concept because they allow you to work with the memory location of a value. They are especially useful when modifying structs, passing large values to functions, sharing state, and working with linked data structures.

## 1. What Is a Pointer?

A pointer is a value that holds the memory address of another value.

Example:

```go
package main

import "fmt"

func main() {
    x := 10
    p := &x

    fmt.Println(x)
    fmt.Println(p)
}
```

Possible output:

```text
10
0xc000012080
```

The memory address shown is only an example; it varies between executions.

Here:

- `x` stores the value `10`.
- `&x` means the address of `x`.
- `p` stores that address.

Conceptually:

```text
Variable x
┌──────────────┐
│ Value: 10    │
│ Address: A   │
└──────────────┘
       ↑
       │
Variable p
┌──────────────┐
│ Value: A     │
└──────────────┘
```

The pointer `p` points to `x`.

## 2. The Address-of Operator `&`

The `&` operator obtains the address of a variable.

```go
x := 100
p := &x
```

This means:

```text
x  → stores 100
&x → address of x
p  → stores the address of x
```

You can print the value and the address separately:

```go
fmt.Println(x)
fmt.Println(&x)
fmt.Println(p)
```

The last two expressions represent the same address.

Remember:

```go
&x
```

means **the address of `x`**.

## 3. The Dereference Operator `*`

The `*` operator can access or modify the value at a pointer's address.

Example:

```go
x := 10
p := &x

fmt.Println(*p)
```

Output:

```text
10
```

Here:

```go
*p
```

means the value stored at the address held by `p`.

You can also update the original value:

```go
x := 10
p := &x

*p = 20

fmt.Println(x)
```

Output:

```text
20
```

Why?

Because `p` points to `x`. Assigning to `*p` changes the value stored at that address.

**Golden rule:**

```text
&x  → address of x
*p  → value at the address stored in p
*p = value → modify the pointed-to value
```

## 4. Pointer Declaration

Suppose:

```go
var x int = 10
var p *int = &x
```

Here:

```text
x → int
p → *int
```

The type `*int` means "pointer to an int."

Other examples:

```go
var p1 *int
var p2 *string
var p3 *bool
```

These are pointer variables, each capable of holding an address of the corresponding type or `nil`.

A pointer to an `int` cannot directly point to a `string`.

```go
x := 10
p := &x
```

The inferred type of `p` is `*int`.

## 5. Changing a Value Through a Pointer

This is one of the most important pointer concepts.

```go
package main

import "fmt"

func main() {
    x := 10
    p := &x

    *p = 20

    fmt.Println(x)
    fmt.Println(*p)
}
```

Output:

```text
20
20
```

Both print `20` because `p` points to `x`.

You have changed the original variable rather than creating a separate integer.

Conceptually:

```text
Before:
x = 10
p ──→ x

After *p = 20:
x = 20
p ──→ x
```

## 6. Copying a Value Is Different from Copying a Pointer

Consider:

```go
x := 10
p := &x

y := *p
y = 50

fmt.Println(x)
fmt.Println(y)
```

Output:

```text
10
50
```

Why?

```go
y := *p
```

copies the integer value into `y`. It does not make `y` a pointer to `x`.

The values are now separate:

```text
x = 10
y = 10
```

After:

```go
y = 50
```

the result is:

```text
x = 10
y = 50
```

### Compare these two declarations

```go
y := *p
```

`y` receives a copy of the pointed-to value.

```go
q := p
```

`q` receives a copy of the pointer value, so both `p` and `q` point to the same original variable.

Example:

```go
x := 10

p := &x
q := p

*q = 100

fmt.Println(x)
fmt.Println(*p)
fmt.Println(*q)
```

Output:

```text
100
100
100
```

The pointer value was copied, but both pointers still point to `x`.

## 7. Reassigning a Pointer

Consider:

```go
x := 10
p := &x

p = new(int)

*p = 50

fmt.Println(x)
fmt.Println(*p)
```

Output:

```text
10
50
```

This is a common interview question.

Initially:

```go
p := &x
```

means `p` points to `x`.

Then:

```go
p = new(int)
```

makes `p` point to a newly allocated integer initialized to zero.

Now `p` no longer points to `x`.

The assignment:

```go
*p = 50
```

changes the new integer, not `x`.

Conceptually:

```text
Initially:

p ──→ x = 10

After p = new(int):

x = 10

p ──→ new integer = 0

After *p = 50:

x = 10

p ──→ new integer = 50
```

**Important:** Changing the value through a pointer is different from changing which address the pointer stores.

## 8. Nil Pointers

A pointer's zero value is `nil`.

Example:

```go
var p *int

fmt.Println(p == nil)
```

Output:

```text
true
```

A nil pointer does not point to a valid value.

Dereferencing it causes a runtime panic:

```go
var p *int

fmt.Println(*p) // Panic
```

Before dereferencing a pointer that may be nil, check it:

```go
if p != nil {
    fmt.Println(*p)
}
```

This is especially important when functions return pointers or when struct fields contain pointers.

## 9. The `new()` Function

The built-in `new()` function allocates zero-initialized storage for a value and returns a pointer to it.

Example:

```go
p := new(int)

fmt.Println(*p)
```

Output:

```text
0
```

Because the zero value of `int` is `0`.

You can assign through the pointer:

```go
p := new(int)
*p = 100

fmt.Println(*p)
```

Output:

```text
100
```

Compare:

```go
x := 10
p := &x
```

with:

```go
p := new(int)
```

The first obtains the address of an existing variable. The second gives you a pointer to newly allocated zero-initialized storage.

Go manages memory automatically through garbage collection when allocated objects are no longer reachable, subject to the language's memory and lifetime rules.

## 10. Pointers and Functions

Go passes arguments **by value**, including pointer arguments.

This means a function receives a copy of the pointer value, but that copied pointer can still point to the same underlying variable.

Example:

```go
func update(p *int) {
    *p = 100
}

func main() {
    x := 10

    update(&x)

    fmt.Println(x)
}
```

Output:

```text
100
```

Why?

The function receives a copy of the address of `x`. Dereferencing that pointer allows it to modify the original variable.

### What if the function only changes its local parameter?

```go
func update(p *int) {
    p = new(int)
    *p = 100
}

func main() {
    x := 10

    update(&x)

    fmt.Println(x)
}
```

Output:

```text
10
```

Why?

The function reassigns its local copy of the pointer. The caller's variable `x` is unchanged, because the new pointer exists only in the local parameter.

This distinction is critical:

```text
*p = value
→ changes the pointed-to value

p = new(int)
→ changes the local pointer variable
```

## 11. Modifying Structs with Pointers

Pointers are often used with structs.

Example:

```go
type User struct {
    Name string
    Age  int
}

func updateAge(u *User) {
    u.Age = 25
}

func main() {
    user := User{
        Name: "Amjad",
        Age:  24,
    }

    updateAge(&user)

    fmt.Println(user.Age)
}
```

Output:

```text
25
```

The pointer allows the function to modify the original struct.

Go automatically handles dereferencing in selector expressions, so:

```go
u.Age = 25
```

is effectively equivalent to:

```go
(*u).Age = 25
```

When `u` is a pointer to a struct, you normally do not need to write `(*u).Age` explicitly.

## 12. Value Parameters vs Pointer Parameters

Consider these two functions.

### Value parameter

```go
func update(x int) {
    x = 100
}

func main() {
    x := 10

    update(x)

    fmt.Println(x)
}
```

Output:

```text
10
```

The function receives a copy of the integer.

### Pointer parameter

```go
func update(x *int) {
    *x = 100
}

func main() {
    x := 10

    update(&x)

    fmt.Println(x)
}
```

Output:

```text
100
```

The function receives a copy of the pointer, which still points to the original integer.

| Feature | Value parameter | Pointer parameter |
|---|---|---|
| Parameter type | `int` | `*int` |
| Receives | Copy of integer | Copy of pointer |
| Can modify caller's integer through parameter? | No | Yes |
| Can reassign caller's variable to a new address? | Not applicable | Not by reassigning only the local pointer |

## 13. Go Is Still Pass-by-Value

A frequent misconception is that Go passes pointers by reference.

Go passes **all arguments by value**.

For a pointer parameter, it copies the pointer value (the address). The copied pointer can then access the same pointed-to object.

Example:

```go
func modify(p *int) {
    *p = 200
}

func main() {
    x := 10
    modify(&x)

    fmt.Println(x)
}
```

Output:

```text
200
```

The pointer itself was passed by value, but its copied address points to the original integer.

This is why the function can modify the value without passing the integer itself by pointer-to-pointer.

## 14. Pointers to Pointers

A pointer can itself be pointed to by another pointer.

Example:

```go
x := 10

p := &x
q := &p

fmt.Println(x)
fmt.Println(*p)
fmt.Println(**q)
```

Output:

```text
10
10
10
```

Types:

```text
x → int
p → *int
q → **int
```

Here:

```go
*p
```

accesses `x`, while:

```go
**q
```

also accesses `x`.

Pointers to pointers are less common in everyday Go application code than ordinary pointers, but can be useful when an operation must update a pointer variable itself.

## 15. Pointers and Slices

Slices already have reference-like behavior because a slice value refers to an underlying array.

Example:

```go
func update(s []int) {
    s[0] = 100
}

func main() {
    s := []int{10, 20, 30}

    update(s)

    fmt.Println(s)
}
```

Output:

```text
[100 20 30]
```

Even though the slice was passed by value, its copied descriptor refers to the same underlying array.

However, if you append inside a function, the caller does not automatically receive the updated slice descriptor:

```go
func add(s []int) {
    s = append(s, 100)
}
```

If the caller needs the new length or potential new backing array, return the updated slice:

```go
func add(s []int) []int {
    return append(s, 100)
}
```

Call it like this:

```go
s = add(s)
```

You generally do not need a pointer to a slice just to modify its existing elements.

## 16. Pointers and Maps

Maps also have reference-like behavior.

Example:

```go
func update(m map[string]int) {
    m["A"] = 100
}

func main() {
    m := map[string]int{
        "A": 10,
    }

    update(m)

    fmt.Println(m["A"])
}
```

Output:

```text
100
```

The function receives a copy of the map value, but the copy refers to the same underlying map.

However, reassigning the local map parameter to a new map doesn't replace the caller's map.

**Remember:** Pointers are not required simply because you want to modify the contents of a map or slice. Their values already provide access to shared underlying data.

## 17. Go Does Not Support Pointer Arithmetic

In some programming languages, you can perform arithmetic on pointers.

Go does not support ordinary pointer arithmetic.

For example, you cannot normally write:

```go
// p = p + 1 // Invalid Go
```

Instead, use indexes and slices to navigate through collections:

```go
s := []int{10, 20, 30}

fmt.Println(s[1])
```

Output:

```text
20
```

This is one way Go reduces the risk of unsafe memory manipulation.

The `unsafe` package provides special low-level operations, but it should not be used as a substitute for ordinary pointer and slice operations.

## 18. Common Pointer Mistakes

### Mistake 1: Dereferencing a nil pointer

```go
var p *int
fmt.Println(*p) // Panic
```

Check the pointer before dereferencing it if it may be nil.

### Mistake 2: Confusing `p` and `*p`

```go
p = new(int)
```

changes the pointer variable.

```go
*p = 10
```

changes the value pointed to.

### Mistake 3: Expecting a copied value to modify the original

```go
y := *p
y = 50
```

does not modify the original value.

### Mistake 4: Assuming local pointer reassignment changes the caller's pointer

```go
func update(p *int) {
    p = new(int)
}
```

Only the local parameter is reassigned. To change the caller's pointer variable itself, the function needs another mechanism, such as returning a pointer or accepting a pointer to a pointer.

### Mistake 5: Taking addresses unnecessarily

Use pointers when you need shared identity, mutation, or a suitable API contract. Do not use pointers automatically for every variable.

### Mistake 6: Assuming pointers are always faster

Pointers can prevent copying large values, but they may also introduce indirection and affect allocation behavior. Choose pointer usage based on semantics first, then benchmark performance if necessary.

---

## 19. Day 3 Interview Questions

**1. What is a pointer?**

A value that holds the memory address of another value.

**2. What does `&x` mean?**

It obtains the address of `x`.

**3. What does `*p` mean?**

It accesses the value at the address held by `p`.

**4. What is the zero value of a pointer?**

`nil`.

**5. What happens when you dereference a nil pointer?**

The program panics.

**6. What does `new(int)` return?**

A `*int` pointer to zero-initialized integer storage.

**7. Does Go pass pointers by reference?**

No. Go passes pointer values by value. The copied pointer can still refer to the same original data.

**8. Why can a function modify an integer using a pointer parameter?**

Because the copied pointer still points to the caller's original integer.

**9. Why doesn't reassigning a local pointer change the caller's pointer?**

Because the parameter is a copy of the pointer value. Reassignment affects only that copy.

**10. What is the difference between `p = new(int)` and `*p = 10`?**

The former changes which address the pointer stores; the latter changes the value at the pointed-to address.

**11. Do slices and maps always need pointers to modify their contents?**

No. Their copied values can still access the same underlying data.

**12. Does Go support pointer arithmetic?**

Not through ordinary Go pointer operations.

---

## 20. Day 3 Golden Mental Model

```text
x := 10
p := &x
```

Conceptually:

```text
x = 10
p ──→ x
```

Now:

```go
*p = 20
```

means:

```text
x = 20
```

But:

```go
y := *p
y = 50
```

means:

```text
x = 20
y = 50
```

And:

```go
p = new(int)
```

means `p` points to different storage. It does not modify the original `x`.

### Final rules to remember

- `&x` gets the address of `x`.
- `*p` accesses the value at a pointer's address.
- `*p = value` can modify the pointed-to value.
- A pointer's zero value is `nil`.
- Dereferencing a nil pointer panics.
- Go passes pointers by value.
- Reassigning a local pointer doesn't reassign the caller's pointer.
- Slice and map values can access shared underlying data without explicit pointers.
- Go does not support ordinary pointer arithmetic.

The most important distinction in Day 3 is this:

**Changing the value at an address is not the same as changing which address a pointer holds.**