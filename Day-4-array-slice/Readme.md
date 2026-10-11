# Go Day 4 — Arrays and Slices

Arrays and slices are fundamental Go data structures. Understanding their memory behavior, length, capacity, `append()`, and shared backing arrays is especially important for writing efficient backend code.

## 1. Arrays in Go

An array stores a fixed number of elements of the same type.

Example:

```go
var arr [3]int

fmt.Println(arr)
```

Output:

```text
[0 0 0]
```

The array has three elements. Since the element type is `int`, each element is initialized to its zero value, `0`.

You can initialize an array directly:

```go
arr := [3]int{10, 20, 30}

fmt.Println(arr)
```

Output:

```text
[10 20 30]
```

You can also let Go count the elements:

```go
arr := [...]int{10, 20, 30}
```

Here, `...` tells Go to determine the array length from the initializer. The array length is still fixed after declaration.

### Important properties of arrays

- An array has a fixed length.
- All elements have the same type.
- The array's length is part of its type.
- Assigning an array copies its elements.

For example:

```go
a := [3]int{10, 20, 30}
b := a

b[0] = 100

fmt.Println(a)
fmt.Println(b)
```

Output:

```text
[10 20 30]
[100 20 30]
```

Changing `b` does not change `a`, because the array assignment copies the elements.

## 2. Array Length Is Part of Its Type

These are different types:

```go
var a [3]int
var b [4]int
```

You cannot directly assign one to the other:

```go
// a = b // Compile error
```

Even though both contain integers, their lengths differ.

Remember:

```text
[3]int != [4]int
```

This fixed-length behavior is one reason slices are more commonly used for collections in everyday Go code.

---

## 3. What Is a Slice?

A slice is a flexible view into an underlying array.

Example:

```go
s := []int{10, 20, 30}

fmt.Println(s)
```

Output:

```text
[10 20 30]
```

Unlike an array, a slice does not have a fixed length as part of its type.

Compare:

```go
arr := [3]int{10, 20, 30}
s := []int{10, 20, 30}
```

Their types are:

```text
arr → [3]int
s   → []int
```

A slice is not simply an array with a flexible length. It is a small descriptor that refers to elements in an underlying array.

Conceptually, a slice contains:

```text
Slice
├── Pointer → underlying array
├── Length
└── Capacity
```

The exact runtime implementation is an internal detail, but this mental model explains most slice behavior.

---

## 4. Length vs Capacity

These are two of the most important slice concepts.

- **Length (`len`)**: The number of elements currently accessible through the slice.
- **Capacity (`cap`)**: The number of elements available from the slice's starting position to the end of its underlying array.

Example:

```go
s := make([]int, 3, 5)

fmt.Println(len(s))
fmt.Println(cap(s))
fmt.Println(s)
```

Output:

```text
3
5
[0 0 0]
```

Why?

```go
make([]int, 3, 5)
```

means:

```text
Element type: int
Length:       3
Capacity:     5
```

The first three elements are accessible through `s`, and their values are initially zero.

Conceptually:

```text
Underlying array:
[0 0 0 _ _]
 └─────┘
  Length = 3

Capacity = 5
```

The extra two elements exist in the backing array but are not directly accessible through ordinary indexing of this slice until the slice is extended.

### Important

You cannot access an element using an index greater than or equal to the slice length, even if that index is below the capacity.

```go
s := make([]int, 3, 5)

// s[3] = 10 // Panic: index out of range
```

You can extend the slice using `append()` or, in appropriate cases, a reslicing operation.

---

## 5. Creating Slices

### Using a slice literal

```go
s := []int{10, 20, 30}
```

### Using `make`

```go
s := make([]int, 3)
```

This creates a slice with length `3` and capacity at least `3`.

You can specify both length and capacity:

```go
s := make([]int, 3, 5)
```

### Declaring a nil slice

```go
var s []int
```

This creates a nil slice.

Its length and capacity are both zero:

```go
fmt.Println(s == nil) // true
fmt.Println(len(s))   // 0
fmt.Println(cap(s))   // 0
```

### Empty but non-nil slice

```go
s := []int{}
```

This slice is non-nil but has length and capacity zero.

```go
fmt.Println(s == nil) // false
fmt.Println(len(s))   // 0
```

A nil slice and an empty non-nil slice are both valid to range over and append to.

**Important:** In Go, the zero value of a slice is `nil`.

---

## 6. The `append()` Function

`append()` adds elements to a slice and returns the resulting slice.

Example:

```go
s := []int{10, 20}

s = append(s, 30)

fmt.Println(s)
```

Output:

```text
[10 20 30]
```

Notice:

```go
s = append(s, 30)
```

We assign the returned slice back to `s`.

This is important because `append()` may return a slice with a different length and, if needed, a different underlying array.

You can append multiple elements:

```go
s := []int{10, 20}

s = append(s, 30, 40, 50)

fmt.Println(s)
```

Output:

```text
[10 20 30 40 50]
```

You can append the elements of another slice using `...`:

```go
a := []int{10, 20}
b := []int{30, 40}

a = append(a, b...)

fmt.Println(a)
```

Output:

```text
[10 20 30 40]
```

---

## 7. What Happens When `append()` Runs Out of Capacity?

Suppose:

```go
s := make([]int, 2, 2)

s[0] = 10
s[1] = 20

s = append(s, 30)
```

Initially:

```text
Length:   2
Capacity: 2
Values:   [10 20]
```

There is no spare capacity, so `append()` must provide storage that can hold the additional element. It does this by allocating a new underlying array when necessary and copying the existing elements.

After the append:

```text
Values: [10 20 30]
Length: 3
Capacity: at least 3
```

The exact new capacity depends on the Go implementation and allocation size. Do not assume that capacity always doubles.

**Golden rule:**

```text
Enough capacity
    → append may reuse the underlying array

Not enough capacity
    → append allocates a new underlying array when necessary
```

Always use the returned slice because its length, capacity, or underlying array may have changed.

---

## 8. Slice Assignment Does Not Copy All Elements

This is a very important concept.

Example:

```go
s1 := []int{10, 20, 30}
s2 := s1

s2[0] = 100

fmt.Println(s1)
fmt.Println(s2)
```

Output:

```text
[100 20 30]
[100 20 30]
```

Why?

The slice assignment copies the slice descriptor, not all the underlying elements.

Conceptually:

```text
s1 ──┐
     ├──→ Underlying array: [100 20 30]
s2 ──┘
```

Both slices refer to the same underlying array.

Therefore, changing an element through one slice can be visible through the other.

This is different from ordinary array assignment, which copies the array's elements.

---

## 9. Slices and Shared Backing Arrays

Consider this example:

```go
s1 := []int{10, 20, 30, 40}
s2 := s1[:2]

s2[0] = 100

fmt.Println(s1)
fmt.Println(s2)
fmt.Println(len(s2), cap(s2))
```

Output:

```text
[100 20 30 40]
[100 20]
2 4
```

Why?

```go
s2 := s1[:2]
```

creates a slice that starts at the same position as `s1`, but its length is only two.

Its capacity remains four because four elements are available from its starting position to the end of the underlying array.

Conceptually:

```text
Underlying array:
[10 20 30 40]
 └──────────┘ s1: len 4, cap 4
 └──────┘     s2: len 2, cap 4
```

After changing `s2[0]`, both slices see the updated element because they share the same array.

---

## 10. A Tricky `append()` Example

Consider:

```go
s1 := []int{10, 20, 30, 40}
s2 := s1[:2]

s2 = append(s2, 100)

fmt.Println(s1)
fmt.Println(s2)
```

Output:

```text
[10 20 100 40]
[10 20 100]
```

Why did `s1[2]` change to `100`?

Initially:

```text
s1 = [10 20 30 40]
s2 = [10 20]
```

The slices share the same backing array, and `s2` has capacity four.

Appending `100` does not need a new array. It reuses the existing one and writes `100` into the third position.

Therefore, `s1` observes the change too.

### Key lesson

`append()` doesn't always allocate a new array. If there is enough capacity, it can reuse the existing underlying array.

---

## 11. Why Might Two Slices Show Different Results?

Consider:

```go
s1 := make([]int, 2, 2)
s1[0] = 10
s1[1] = 20

s2 := s1

s2 = append(s2, 30)

s2[0] = 100

fmt.Println(s1)
fmt.Println(s2)
```

Output:

```text
[10 20]
[100 20 30]
```

Why?

Initially, `s1` and `s2` share an underlying array:

```text
[10 20]
```

But the capacity is only two. Appending `30` requires additional storage, so a new underlying array is allocated.

After `append()`:

```text
s1 ──→ Original array: [10 20]

s2 ──→ New array:      [10 20 30]
```

When `s2[0]` changes to `100`, the original array is unaffected.

This is one of the most important slice interview questions.

**Remember:** Shared backing arrays explain why modifications are sometimes visible between slices and sometimes aren't. Check capacity and whether `append()` needed new storage.

---

## 12. Slicing Syntax

General syntax:

```go
slice[low:high]
```

The low index is included, but the high index is excluded.

Example:

```go
s := []int{10, 20, 30, 40, 50}

fmt.Println(s[1:4])
```

Output:

```text
[20 30 40]
```

The selected indexes are `1`, `2`, and `3`.

Other examples:

```go
s[:3]  // First three elements
s[2:]  // From index 2 to the end
s[:]   // A slice over all currently accessible elements
```

A simple rule:

```text
low is included
high is excluded
```

The two-index slice expression requires valid bounds, and the high index must not exceed the original slice's capacity.

---

## 13. Three-Index Slicing

Go also supports:

```go
slice[low:high:max]
```

The third index limits the capacity of the resulting slice.

Example:

```go
s1 := []int{10, 20, 30, 40, 50}
s2 := s1[1:3:3]

fmt.Println(s2)
fmt.Println(len(s2))
fmt.Println(cap(s2))
```

Output:

```text
[20 30]
2
2
```

Calculation:

```text
Length   = high - low = 3 - 1 = 2
Capacity = max - low  = 3 - 1 = 2
```

This limits the resulting slice's capacity to two, so appending to it cannot reuse the original array beyond that capacity.

Example:

```go
s2 = append(s2, 100)

fmt.Println(s1)
fmt.Println(s2)
```

Output:

```text
[10 20 30 40 50]
[20 30 100]
```

A new underlying array is needed because `s2` has capacity two.

Three-index slicing is useful when you want to prevent an appended slice from overwriting later elements in the original backing array.

---

## 14. The `copy()` Function

If you need an independent copy of slice elements, use `copy()`.

Example:

```go
s1 := []int{10, 20, 30}
s2 := make([]int, len(s1))

copy(s2, s1)

s2[0] = 100

fmt.Println(s1)
fmt.Println(s2)
```

Output:

```text
[10 20 30]
[100 20 30]
```

Now changing `s2` does not affect `s1`, because each slice has its own backing array in this example.

### Important behavior

```go
copied := copy(dst, src)
```

The function copies up to:

```text
min(len(dst), len(src))
```

elements and returns the number copied.

Example:

```go
src := []int{1, 2, 3, 4}
dst := make([]int, 2)

n := copy(dst, src)

fmt.Println(dst)
fmt.Println(n)
```

Output:

```text
[1 2]
2
```

**Remember:** `copy()` copies elements. It does not automatically create a new destination slice; you must provide the destination.

For slices of pointers, maps, or other reference-containing values, copying the elements is not the same as deeply cloning the objects they refer to.

---

## 15. Passing Slices to Functions

A slice argument is passed by value, meaning the slice descriptor is copied. The copied descriptor still refers to the same underlying array.

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

The function can modify the shared element.

### What happens if the function appends?

```go
func add(s []int) {
    s = append(s, 100)
}
```

Calling it does not update the caller's slice length because the new slice descriptor is assigned only to the local variable `s`.

Example:

```go
func add(s []int) {
    s = append(s, 100)
}

func main() {
    s := []int{10, 20}

    add(s)

    fmt.Println(s)
}
```

Output:

```text
[10 20]
```

Even if the local `append()` reuses the backing array, the caller's slice length is still two.

### Correct approach: return the slice

```go
func add(s []int) []int {
    s = append(s, 100)
    return s
}

func main() {
    s := []int{10, 20}

    s = add(s)

    fmt.Println(s)
}
```

Output:

```text
[10 20 100]
```

**Golden rule:**

When a function needs to return a changed slice length or a potentially different underlying array, return the slice and assign it back to the caller's variable.

---

## 16. The `range` Loop

You can iterate over a slice with `range`.

```go
s := []int{10, 20, 30}

for index, value := range s {
    fmt.Println(index, value)
}
```

Output:

```text
0 10
1 20
2 30
```

If you don't need the index:

```go
for _, value := range s {
    fmt.Println(value)
}
```

If you only need the indexes:

```go
for index := range s {
    fmt.Println(index)
}
```

### Important: `value` is a copy

Consider:

```go
s := []int{10, 20, 30}

for _, value := range s {
    value *= 2
}

fmt.Println(s)
```

Output:

```text
[10 20 30]
```

The `value` variable contains a copy of each element. Changing it does not change the corresponding slice element.

To modify the original elements, use the indexes:

```go
for i := range s {
    s[i] *= 2
}

fmt.Println(s)
```

Output:

```text
[20 40 60]
```

---

## 17. Slice Comparisons

You cannot compare two slices directly using `==`.

This is invalid:

```go
a := []int{1, 2}
b := []int{1, 2}

// fmt.Println(a == b) // Compile error
```

Slices can only be compared directly with `nil`.

```go
var s []int

fmt.Println(s == nil) // true
```

To compare slice contents, use an element-by-element loop or a suitable comparison helper from the standard library.

For example, for comparable element types:

```go
import "slices"

a := []int{1, 2}
b := []int{1, 2}

fmt.Println(slices.Equal(a, b))
```

Output:

```text
true
```

This checks equality of elements, not whether the slices share an underlying array.

---

## 18. Common Slice Pitfalls

### Mistake 1: Forgetting to assign the result of `append()`

Avoid:

```go
append(s, 10)
```

Prefer:

```go
s = append(s, 10)
```

### Mistake 2: Assuming slice assignment makes a full copy

```go
s2 := s1
```

Copies the descriptor, not the elements.

### Mistake 3: Assuming `append()` always allocates

It can reuse the backing array when there is sufficient capacity.

### Mistake 4: Assuming `append()` never allocates

When capacity is insufficient, new storage is allocated as necessary.

### Mistake 5: Modifying a `range` value

The range value is a copy of an element, so changing it doesn't modify the original slice element.

### Mistake 6: Appending inside a function without returning the slice

If the caller needs the new length or new backing array, return the updated slice.

### Mistake 7: Assuming `copy()` is always a deep copy

`copy()` copies elements. It does not recursively clone objects referenced by those elements.

---

## 19. Day 4 Interview Questions

**1. What is the difference between an array and a slice?**

An array has a fixed length that is part of its type. A slice is a descriptor providing a flexible-length view into an underlying array.

**2. What is the difference between length and capacity?**

Length is the number of accessible elements. Capacity is the number of elements available from the slice's starting position to the end of its backing array.

**3. What happens when a slice is assigned to another variable?**

The slice descriptor is copied, but the two slices can share the underlying array.

**4. What does `append()` return?**

The updated slice, which may have a different length or may refer to a different underlying array.

**5. When does `append()` need new storage?**

When the existing backing array cannot hold the additional elements.

**6. Why should we assign the result of `append()`?**

Because the returned slice reflects the updated length and potentially a new backing array.

**7. What is a nil slice?**

A slice whose zero value is `nil`, with length and capacity both zero.

**8. What is the difference between `s[:2]` and `s[1:3]`?**

The first selects indexes zero and one. The second selects indexes one and two.

**9. What does three-index slicing control?**

It specifies the resulting slice's maximum capacity.

**10. What does `copy()` do?**

It copies elements from a source slice into a destination slice, up to the smaller slice length.

**11. Why doesn't appending inside a function necessarily change the caller's slice length?**

The slice descriptor is passed by value. The function needs to return the updated slice if the caller should see its new length.

**12. Why doesn't modifying a `range` variable modify a slice element?**

Because the range variable contains a copy of the element.

**13. Can slices be compared using `==`?**

Only against `nil`, not directly against another slice.

**14. What happens when two slices share a backing array?**

Element writes through one slice can be observed through the other wherever their accessible indexes overlap.

---

## 20. Day 4 Golden Mental Model

```text
ARRAY
→ fixed length
→ assignment copies elements

SLICE
→ pointer + length + capacity
→ assignment copies descriptor
→ backing array may be shared

len(s)
→ accessible element count

cap(s)
→ available capacity from the slice's start

append(s, values...)
→ returns updated slice
→ may reuse or replace backing array

copy(dst, src)
→ copies elements into destination

range
→ index + copy of element
```

### Most important example to remember

```go
s1 := []int{10, 20, 30}
s2 := s1

s2[0] = 100
```

Both slices now see the update because they share the backing array.

But:

```go
s := make([]int, 2, 2)
s = append(s, 30)
```

may allocate a new backing array because capacity is exhausted.

And when appending inside a function:

```go
func add(s []int) []int {
    return append(s, 100)
}
```

the caller should use the returned slice:

```go
s = add(s)
```

**Final takeaway:** The key to mastering Go slices is to understand that the slice descriptor and its backing array are different things. Once you understand length, capacity, shared arrays, and `append()`, most slice behavior becomes predictable.
