# Go Day 2 — Struct Memory Layout, Alignment & Padding

Understanding struct memory layout is important for writing efficient Go programs. It helps explain memory usage, field offsets, alignment, and why changing the order of struct fields can reduce memory consumption.

## 1. What Is a Struct?

A struct is a collection of fields grouped into one type.

Example:

```go id="m08zg2"
type User struct {
    ID   int
    Age  int16
    Name string
}
```

You can create a value of this type:

```go id="8msn38"
u := User{
    ID:   101,
    Age:  24,
    Name: "Amjad",
}
```

Each field has a type and occupies storage according to that type's representation and alignment requirements.

Go may insert padding bytes between fields or at the end of a struct to satisfy alignment requirements.

---

## 2. How Much Memory Do Data Types Use?

On a typical 64-bit Go implementation, these are the common sizes:

| Type | Typical size |
|---|---:|
| `bool` | 1 byte |
| `byte` (`uint8`) | 1 byte |
| `int8` | 1 byte |
| `int16` | 2 bytes |
| `int32` | 4 bytes |
| `int64` | 8 bytes |
| `float32` | 4 bytes |
| `float64` | 8 bytes |
| `int` | 8 bytes |
| `uint` | 8 bytes |
| `rune` (`int32`) | 4 bytes |

**Important:** The sizes of `int`, `uint`, and pointer-related types depend on the target architecture. The table describes common sizes on a 64-bit platform, not a universal rule for every Go target.

A `rune` is an alias for `int32` and is commonly used to represent a Unicode code point.

A `byte` is an alias for `uint8`.

### Example

```go id="4d68gl"
var a int16
var b int32
var c int64

fmt.Println(unsafe.Sizeof(a)) // 2
fmt.Println(unsafe.Sizeof(b)) // 4
fmt.Println(unsafe.Sizeof(c)) // 8
```

Import:

```go id="psw0h1"
import "unsafe"
```

These sizes are the same on common 32-bit and 64-bit Go targets for the explicitly sized integer types shown.

---

## 3. What Is Memory Alignment?

**Memory alignment** means placing a value at an address that satisfies its type's alignment requirement.

For example, on a typical 64-bit Go implementation, an `int64` commonly has 8-byte alignment.

Its starting address should be a multiple of 8.

Conceptually, valid aligned addresses include:

```text id="ctn4ab"
0
8
16
24
32
```

If a field's next available position isn't suitably aligned, the compiler inserts padding bytes before that field.

The exact alignment requirements can depend on the target architecture and Go implementation.

### Why is alignment needed?

Aligned memory access can be more efficient for a processor and is required by some architectures or operations.

Go handles the necessary padding in ordinary struct layouts automatically.

---

## 4. What Is Padding?

**Padding** consists of bytes inserted into a struct layout to satisfy alignment requirements.

Consider:

```go id="5wk18c"
type Example struct {
    A byte
    B int32
}
```

On a typical 64-bit Go implementation:

- `A` occupies 1 byte at offset 0.
- `B` requires 4-byte alignment.
- The next available offset after `A` is 1.
- Three padding bytes are inserted so that `B` begins at offset 4.

Conceptually:

```text id="jv5boq"
Offset 0       1       2       3       4       5       6       7
      +-------+-------+-------+-------+-------+-------+-------+-------+
      | A     | Padding: 3 bytes        | B: int32 (4 bytes)          |
      +-------+-------+-------+-------+-------+-------+-------+-------+
```

Memory layout:

```text id="d6jmbt"
A       → 1 byte
Padding → 3 bytes
B       → 4 bytes

Total   → 8 bytes
```

Without the padding, `B` would begin at offset 1, which does not satisfy its 4-byte alignment requirement.

---

## 5. Struct Alignment

The alignment requirement of a struct is generally the maximum alignment requirement among its fields on common Go implementations.

For example:

```go id="m1b7hz"
type Example struct {
    A byte
    B int64
}
```

On a typical 64-bit platform:

```text id="jp7397"
A       → offset 0, size 1
Padding → 7 bytes
B       → offset 8, size 8

Total   → 16 bytes
```

The struct's alignment is typically 8 bytes because `int64` has 8-byte alignment.

The total struct size is rounded up to a multiple of the struct alignment.

This is one reason a struct can use more memory than the sum of the field sizes.

---

## 6. Trailing Padding

Padding can also appear at the end of a struct.

Consider:

```go id="j5skgc"
type Example struct {
    A int16
    B byte
}
```

On a typical Go implementation:

```text id="t7zopb"
A → offset 0, size 2
B → offset 2, size 1

Field bytes → 3
Trailing padding → 1

Struct size → 4 bytes
```

Why is the last byte needed?

The struct's alignment is typically 2 bytes because of the `int16` field. Its total size must be a multiple of 2, so one trailing padding byte is added.

### Important rule

Padding can exist:

- Between two fields.
- At the end of a struct.

The compiler determines the layout automatically.

---

## 7. Struct Size vs Sum of Field Sizes

Consider:

```go id="q3v4xj"
type User struct {
    A byte
    B int64
    C byte
}
```

On a typical 64-bit Go implementation:

```text id="fopkjj"
A → 1 byte
Padding → 7 bytes
B → 8 bytes
C → 1 byte
Trailing padding → 7 bytes

Total → 24 bytes
```

The field sizes add up to:

```text id="kxwfuk"
1 + 8 + 1 = 10 bytes
```

But the struct occupies 24 bytes in this layout.

The additional 14 bytes are padding.

This is why **the sum of field sizes does not necessarily equal the total struct size**.

---

## 8. Field Ordering Can Reduce Padding

Consider these two structs.

### Struct A

```go id="3l21rw"
type A struct {
    X byte
    Y int64
    Z byte
}
```

Typical 64-bit layout:

```text id="jhhctv"
X       → 1 byte
Padding → 7 bytes
Y       → 8 bytes
Z       → 1 byte
Padding → 7 bytes

Total → 24 bytes
```

### Struct B

```go id="atky7y"
type B struct {
    Y int64
    X byte
    Z byte
}
```

Typical 64-bit layout:

```text id="x4c5p4n"
Y → 8 bytes
X → 1 byte
Z → 1 byte
Padding → 6 bytes

Total → 16 bytes
```

Both structs contain the same field types, but their sizes differ because of field ordering.

```text id="w6d7sf"
Struct A → 24 bytes
Struct B → 16 bytes
```

### Why?

In Struct B, the 8-byte field is placed first. The two 1-byte fields then follow it, and only six trailing bytes are needed to align the total size.

**Important:** Reordering fields can reduce memory consumption, but only change field order when it preserves the intended behavior and doesn't break compatibility requirements such as binary layouts, serialization assumptions, or reflection-based expectations.

---

## 9. A More Complex Padding Example

Consider:

```go id="m3j7ws"
type Example struct {
    A byte
    B int32
    C byte
}
```

On a typical 64-bit Go implementation:

```text id="rlw9cv"
A       → offset 0, size 1
Padding → 3 bytes
B       → offset 4, size 4
C       → offset 8, size 1
Padding → 3 bytes

Total → 12 bytes
```

The fields require 6 bytes in total, but the struct occupies 12 bytes.

The struct's alignment is typically 4 bytes because `int32` requires 4-byte alignment.

Therefore, after the last field, the total size is rounded up to a multiple of 4.

### Calculate it step by step

```text id="jjn9pv"
Start at offset 0.

A: 1 byte
Next offset: 1

B needs 4-byte alignment.
Next aligned offset: 4
Padding before B: 3 bytes

B: 4 bytes
Next offset: 8

C: 1 byte
Next offset: 9

Struct alignment: 4 bytes
Rounded total size: 12 bytes
```

This is the process to use when solving struct-layout questions.

---

## 10. Using `unsafe.Sizeof`

The `unsafe.Sizeof()` function reports the size in bytes of a value's type.

Example:

```go id="34ks4z"
type Example struct {
    A byte
    B int64
    C byte
}

func main() {
    e := Example{}

    fmt.Println(unsafe.Sizeof(e))
}
```

On a typical 64-bit Go implementation, the result is:

```text id="1q0s42"
24
```

`unsafe.Sizeof()` helps you inspect memory layouts during experiments.

### Important limitation

It measures the size of the value's type itself. It doesn't count all the memory that might be reachable through pointers, slices, maps, or other reference-containing fields.

For example, the size of a slice value measures its descriptor, not the entire backing array.

---

## 11. Using `unsafe.Offsetof`

The `unsafe.Offsetof()` function reports the byte offset of a struct field from the beginning of the struct.

Example:

```go id="9x8qlu"
type Example struct {
    A byte
    B int64
    C byte
}

func main() {
    var e Example

    fmt.Println(unsafe.Offsetof(e.A))
    fmt.Println(unsafe.Offsetof(e.B))
    fmt.Println(unsafe.Offsetof(e.C))
}
```

On a typical 64-bit Go implementation:

```text id="6uxzw0"
0
8
16
```

Why?

```text id="vgxb0z"
A starts at offset 0
B starts at offset 8
C starts at offset 16
```

The padding between fields changes their offsets.

Import:

```go id="hcf7z7"
import "unsafe"
```

---

## 12. Using `unsafe.Alignof`

The `unsafe.Alignof()` function reports the alignment requirement of a value's type.

Example:

```go id="zlxl6c"
type Example struct {
    A byte
    B int64
}

func main() {
    var e Example

    fmt.Println(unsafe.Alignof(e.A))
    fmt.Println(unsafe.Alignof(e.B))
    fmt.Println(unsafe.Alignof(e))
}
```

On a typical 64-bit Go implementation, the output is commonly:

```text id="2r7nhv"
1
8
8
```

These correspond to the alignment requirements of the fields and the struct.

### Summary

```text id="nca5ul"
unsafe.Sizeof(v)
→ size in bytes

unsafe.Offsetof(v.Field)
→ field offset in bytes

unsafe.Alignof(v)
→ alignment requirement
```

These are useful learning tools, but be cautious about using `unsafe` in production code.

---

## 13. How `int16`, `int32`, and `int64` Affect Layout

These integer types have explicitly defined widths:

```text id="i9x2gu"
int16 → 2 bytes
int32 → 4 bytes
int64 → 8 bytes
```

But field alignment and the resulting struct layout also matter.

Example:

```go id="j69c2g"
type Example struct {
    A int16
    B int32
    C int16
}
```

On a typical 64-bit Go implementation:

```text id="7s2i8l"
A → offset 0, size 2
Padding → 2 bytes
B → offset 4, size 4
C → offset 8, size 2
Trailing padding → 2 bytes

Total → 12 bytes
```

The field sizes alone add up to:

```text id="3jrk91"
2 + 4 + 2 = 8 bytes
```

But the final struct size is 12 bytes due to padding.

---

## 14. Memory Alignment and Struct Optimization

When struct memory usage matters, grouping fields according to alignment can reduce padding.

A common optimization technique is to place fields with stricter alignment requirements before smaller fields.

Example:

```go id="w5ne4r"
type PoorLayout struct {
    A byte
    B int64
    C byte
    D int32
}
```

A reordered version may reduce padding:

```go id="d81187"
type BetterLayout struct {
    B int64
    D int32
    A byte
    C byte
}
```

On typical 64-bit implementations, the reordered layout can require less memory.

However, do not reorder fields automatically just to minimize size. Consider:

- Readability and logical grouping.
- Binary or serialization layout assumptions.
- Compatibility with external memory layouts.
- The number of values actually allocated.
- Measured performance before optimizing.

Use benchmarks if memory layout might materially affect your application's performance.

---

## 15. Structs Containing Pointers, Slices, Maps, or Strings

Be careful when estimating a struct's memory usage.

Consider:

```go id="n7g0ie"
type User struct {
    ID   int64
    Name string
    Tags []string
}
```

`unsafe.Sizeof(user)` counts the struct's immediate representation.

It does not include all the memory occupied by the characters of `Name`, the backing array of `Tags`, or the strings stored in the slice.

Conceptually:

```text id="t3g3bd"
User struct
├── ID
├── string descriptor ──→ string bytes
└── slice descriptor ───→ slice backing array
                              └── string data
```

Thus, there is a distinction between:

```text id="v2rj7m"
Size of the struct value
```

and:

```text id="s1bk5w"
Total memory reachable from the struct
```

This matters when profiling real Go applications.

---

## 16. Common Mistakes

### Mistake 1: Assuming struct size equals the sum of field sizes

Padding can increase the total size.

### Mistake 2: Assuming all types have the same size on every architecture

The sizes of `int`, `uint`, pointers, and their alignment requirements can depend on the target.

### Mistake 3: Ignoring trailing padding

The compiler may add padding after the final field to satisfy the struct's alignment.

### Mistake 4: Assuming field ordering has no effect

Field ordering can change the amount of padding and the total struct size.

### Mistake 5: Thinking `unsafe.Sizeof` reports total reachable memory

It reports the size of the value's type representation, not the total heap data reachable through pointers and reference-containing fields.

### Mistake 6: Using `unsafe` unnecessarily

Prefer safe Go operations in production unless you have a compelling, well-tested low-level requirement.

---

## 17. Day 2 Interview Questions

**1. What is struct padding?**

Padding is unused storage inserted between or after fields to satisfy memory alignment requirements.

**2. Why is padding needed?**

It ensures fields and structs meet their types' alignment requirements.

**3. What is memory alignment?**

A requirement that a value's starting address satisfy a particular boundary.

**4. Why can reordering struct fields reduce memory usage?**

It can reduce padding by placing fields with stricter alignment requirements together.

**5. What does `unsafe.Sizeof()` do?**

It returns the size in bytes of a value's type representation.

**6. What does `unsafe.Offsetof()` do?**

It returns a field's byte offset from the beginning of its struct.

**7. What does `unsafe.Alignof()` do?**

It reports a value type's alignment requirement.

**8. Is `int` always 8 bytes?**

No. Its size depends on the architecture and implementation.

**9. How large is `int16`?**

Two bytes.

**10. Why can a struct containing a `byte` and an `int64` occupy 16 bytes?**

The `int64` commonly requires 8-byte alignment, leading to padding before it and potentially at the end of the struct.

**11. Does `unsafe.Sizeof()` include the memory referenced by a slice?**

No. It measures the slice value itself, not its complete backing array.

**12. Does the Go language guarantee identical struct layouts across all architectures?**

No. Alignment and size details can vary by target and implementation.

---

## 18. Day 2 Golden Mental Model

```text id="3oc9kg"
Struct
   ↓
Fields have sizes and alignment requirements
   ↓
Compiler lays out each field
   ↓
Padding may be inserted between fields
   ↓
Trailing padding may be inserted
   ↓
Final struct size is determined
```

The main rules:

```text id="j2uzg0"
unsafe.Sizeof(v)
→ size of the value's representation

unsafe.Offsetof(v.Field)
→ field offset

unsafe.Alignof(v)
→ alignment requirement

Padding
→ alignment overhead

Field ordering
→ can change padding and struct size
```

### Final takeaway

The key to understanding Go struct memory layout is to calculate each field's offset using its alignment requirement, then account for any padding and round the final size to the struct's required alignment.

These rules help you answer memory-layout interview questions and understand how Go represents structs under the hood.
