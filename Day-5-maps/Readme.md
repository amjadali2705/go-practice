# Go Day 5 — Maps

Maps are a built-in Go data structure used to store key-value pairs. They are useful for lookups, counting, grouping, caching, configuration, and representing JSON-like data.

For backend development, maps are commonly used for request metadata, configuration, lookup tables, and aggregating transaction data.

## 1. What Is a Map?

A map stores values associated with unique keys.

```go
ages := map[string]int{
    "Amjad": 24,
    "Rahul": 25,
}

fmt.Println(ages["Amjad"])
```

Output:

```text
24
```

Here:

- `string` is the key type.
- `int` is the value type.
- `"Amjad"` is a key.
- `24` is its associated value.

The general syntax is:

```go
map[KeyType]ValueType
```

Examples:

```go
map[string]int
map[int]string
map[string]bool
map[int]User
```

Keys must be comparable types. Values can be of any type.

---

## 2. Creating Maps

There are several ways to declare and initialize maps.

### Using `make`

```go
m := make(map[string]int)

m["A"] = 10
m["B"] = 20

fmt.Println(m)
```

`make` creates an initialized map ready for use.

### Using a map literal

```go
m := map[string]int{
    "A": 10,
    "B": 20,
}
```

### Declaring a nil map

```go
var m map[string]int
```

Here, `m` is a nil map.

It has no initialized underlying map storage.

**Important:** A nil map can be read from, but you cannot assign a key-value pair to it.

---

## 3. Nil Map vs Empty Map

These two maps are different.

### Nil map

```go
var m map[string]int

fmt.Println(m == nil)
fmt.Println(len(m))
```

Output:

```text
true
0
```

### Initialized empty map

```go
m := make(map[string]int)

fmt.Println(m == nil)
fmt.Println(len(m))
```

Output:

```text
false
0
```

Both have length zero, but only the first is nil.

An empty map is initialized and ready for writes. A nil map is not.

### Important behavior

```go
var m map[string]int

fmt.Println(m["A"])
```

Output:

```text
0
```

Reading a missing key returns the value type's zero value. For `int`, that is `0`.

But this panics:

```go
var m map[string]int

m["A"] = 10 // Panic: assignment to entry in nil map
```

Use:

```go
m := make(map[string]int)
m["A"] = 10
```

---

## 4. Reading Values from a Map

Example:

```go
m := map[string]int{
    "A": 10,
    "B": 20,
}

fmt.Println(m["A"])
fmt.Println(m["C"])
```

Output:

```text
10
0
```

Because `"C"` doesn't exist, Go returns the zero value for `int`.

This can create a problem: what if the actual stored value is also zero?

Use the comma-ok idiom.

---

## 5. The Comma-OK Idiom

Go allows you to check whether a key exists.

```go
m := map[string]int{
    "A": 0,
}

value, ok := m["A"]

fmt.Println(value, ok)
```

Output:

```text
0 true
```

Now check a missing key:

```go
value, ok := m["B"]

fmt.Println(value, ok)
```

Output:

```text
0 false
```

Here:

- `value` is the stored value or the zero value.
- `ok` is `true` if the key exists and `false` otherwise.

### Common pattern

```go
value, ok := m[key]

if !ok {
    fmt.Println("Key not found")
    return
}

fmt.Println(value)
```

**Remember:** When the zero value is a valid stored value, use comma-ok to distinguish a missing key from a key that exists.

---

## 6. Updating and Inserting Values

Assigning to a map key inserts or replaces a value.

```go
m := map[string]int{
    "A": 10,
}

m["A"] = 50 // Update existing key
m["B"] = 20 // Insert new key

fmt.Println(m)
```

Conceptually:

```text
Before:
A → 10

After:
A → 50
B → 20
```

You don't need separate operations for insertion and updating.

---

## 7. Deleting Map Entries

Use the built-in `delete()` function.

```go
m := map[string]int{
    "A": 10,
    "B": 20,
}

delete(m, "A")

fmt.Println(m)
```

The remaining entry is:

```text
B → 20
```

Deleting a nonexistent key is safe:

```go
delete(m, "C")
```

This does nothing and does not panic.

Deleting from a nil map is also safe.

---

## 8. The `len()` Function

Use `len()` to get the number of entries.

```go
m := map[string]int{
    "A": 10,
    "B": 20,
}

fmt.Println(len(m))
```

Output:

```text
2
```

Remember that `len(m)` counts entries, not the amount of memory the map uses.

For a nil map:

```go
var m map[string]int

fmt.Println(len(m))
```

Output:

```text
0
```

---

## 9. Iterating over a Map

Use `range`.

```go
m := map[string]int{
    "A": 10,
    "B": 20,
    "C": 30,
}

for key, value := range m {
    fmt.Println(key, value)
}
```

The order is **not guaranteed**.

One execution might produce:

```text
A 10
B 20
C 30
```

Another may produce:

```text
C 30
A 10
B 20
```

Both are valid.

### Important rules

- Go maps do not guarantee insertion order.
- Do not write tests that depend on a particular map iteration order.
- If deterministic output is required, collect the keys into a slice, sort them, and iterate over that slice.

Example:

```go
keys := make([]string, 0, len(m))

for key := range m {
    keys = append(keys, key)
}

sort.Strings(keys)

for _, key := range keys {
    fmt.Println(key, m[key])
}
```

Import:

```go
import "sort"
```

---

## 10. Iterating over Keys Only

You can omit the value:

```go
for key := range m {
    fmt.Println(key)
}
```

When only the keys are needed, this is cleaner than declaring an unused value variable.

You can also retrieve a value directly inside the loop if needed.

---

## 11. Valid Map Key Types

Map keys must be comparable.

Common valid key types include:

```go
map[int]string
map[string]int
map[bool]string
map[float64]string
map[UserID]string
```

Structs and arrays can also be keys if their component types are comparable.

These are invalid:

```go
map[[]int]string
map[map[string]int]string
```

Why?

Slices and maps cannot be compared using `==` against another value of the same type, so they cannot be ordinary map key types.

Functions are also not valid map key types.

### Important subtlety

An interface type can be used as a map key type, but inserting an interface value whose dynamic value is not comparable will panic.

For example:

```go
m := make(map[any]string)

m[10] = "number"       // Valid
m["hello"] = "string"  // Valid

// This panics because the dynamic value is a slice.
m[[]int{1, 2}] = "slice"
```

This distinction is useful when working with `any`.

---

## 12. Map Keys Must Be Unique

A map can contain only one entry per key.

```go
m := map[string]int{
    "A": 10,
}

m["A"] = 20

fmt.Println(m["A"])
```

Output:

```text
20
```

The second assignment replaces the first value.

Maps are not designed to store multiple independent values under the same key. If you need multiple values, use a slice as the map value:

```go
m := map[string][]int{
    "A": {10, 20, 30},
}
```

---

## 13. Maps Are Not Directly Comparable

You can compare a map with `nil`:

```go
m := make(map[string]int)

fmt.Println(m == nil)
```

Output:

```text
false
```

But you cannot compare two maps directly:

```go
m1 := map[string]int{"A": 10}
m2 := map[string]int{"A": 10}

// fmt.Println(m1 == m2) // Compile error
```

Go does not support direct equality comparison between maps.

To compare map contents, you must implement the comparison or use a suitable helper from the standard library, such as `maps.Equal` for comparable value types.

Example:

```go
import "maps"

m1 := map[string]int{"A": 10}
m2 := map[string]int{"A": 10}

fmt.Println(maps.Equal(m1, m2))
```

Output:

```text
true
```

Note that `maps.Equal` compares keys and values, not whether the variables refer to the same underlying map.

---

## 14. How Maps Behave When Passed to Functions

This is one of the most important Day 5 concepts.

Consider:

```go
func update(m map[string]int) {
    m["A"] = 100
    m = make(map[string]int)
    m["B"] = 200
}

func main() {
    m := map[string]int{
        "A": 10,
    }

    update(m)

    fmt.Println(m)
}
```

Output:

```text
map[A:100]
```

Why doesn't `"B": 200` appear?

### Step 1: Pass the map

```go
update(m)
```

The map value is copied into the function parameter. That value refers to the same underlying map data.

### Step 2: Modify an existing entry

```go
m["A"] = 100
```

This updates the shared underlying map. The caller observes the change.

### Step 3: Assign a new map to the local parameter

```go
m = make(map[string]int)
```

This makes the local parameter refer to a different, newly initialized map.

### Step 4: Insert `"B"`

```go
m["B"] = 200
```

This inserts into the new map, not the caller's original map.

### Mental model

```text
Caller map value ───┐
                    ├──→ Original underlying map
Function parameter ┘
                         A → 100

After local reassignment:

Caller map value ───────→ Original underlying map
                            A → 100

Function parameter ─────→ New underlying map
                            B → 200
```

**Golden rule:** Copying a map value does not clone its contents. The copied value still refers to the same underlying map until that local variable is reassigned.

---

## 15. Maps and Reference-Like Behavior

In Go, a map value is commonly described as having reference-like behavior because copies can access the same underlying map.

Example:

```go
m1 := map[string]int{
    "A": 10,
}

m2 := m1

m2["A"] = 100

fmt.Println(m1["A"])
```

Output:

```text
100
```

Both map variables can access the same underlying map.

However:

```go
m2 = make(map[string]int)
m2["B"] = 200
```

reassigns `m2` to a different map. It doesn't reassign `m1`.

**Be precise:** A map is a Go map value that refers to underlying map storage. It isn't correct to treat assignment as an automatic deep copy of the contents.

---

## 16. Maps Are Not Safe for Concurrent Writes

Go maps are not safe for unsynchronized concurrent writes.

Example:

```go
m := make(map[string]int)

go func() {
    m["A"] = 10
}()

go func() {
    m["B"] = 20
}()
```

Concurrent access like this is unsafe. It may cause a data race, and concurrent map read/write activity can trigger a runtime failure.

### Safer option: `sync.Mutex`

```go
var (
    m  = make(map[string]int)
    mu sync.Mutex
)

func set(key string, value int) {
    mu.Lock()
    defer mu.Unlock()

    m[key] = value
}

func get(key string) int {
    mu.Lock()
    defer mu.Unlock()

    return m[key]
}
```

Use the same synchronization strategy for every access that needs protection. Locking writes alone is insufficient if other goroutines read or iterate over the map concurrently without synchronization.

For read-heavy workloads, `sync.RWMutex` may be worth considering.

Go also provides `sync.Map`, a concurrent map intended for particular use cases. A regular map with a mutex is often simpler when the application already needs coordinated access to related state.

---

## 17. Maps Inside Structs

Maps are often used as struct fields.

```go
type User struct {
    ID       int
    Name     string
    Metadata map[string]string
}
```

Usage:

```go
u := User{
    ID:   1,
    Name: "Amjad",
    Metadata: map[string]string{
        "role": "developer",
        "team": "backend",
    },
}

fmt.Println(u.Metadata["role"])
```

Output:

```text
developer
```

Remember that copying the `User` struct copies its map value, not all of the map's entries. The copied structs can still refer to the same underlying map.

This matters when deciding whether mutations should be shared or isolated.

---

## 18. A Common Map Pitfall: Struct Values

Consider:

```go
type User struct {
    Age int
}

users := map[string]User{
    "Amjad": {Age: 24},
}
```

You cannot directly do this:

```go
// users["Amjad"].Age = 25 // Compile error
```

The map lookup returns a value; you cannot modify a field of that returned struct value directly.

One solution is to retrieve, modify, and assign it back:

```go
user := users["Amjad"]
user.Age = 25
users["Amjad"] = user
```

Another option is to store pointers:

```go
users := map[string]*User{
    "Amjad": {Age: 24},
}

users["Amjad"].Age = 25
```

Use pointers only when their shared-mutation and lifetime behavior are appropriate for your design.

---

## 19. Practical Backend Use Cases

Maps are helpful in real-world Go services.

### Lookup table

```go
statusCodes := map[string]int{
    "success": 200,
    "created": 201,
    "missing": 404,
}
```

### Frequency counting

```go
counts := make(map[string]int)

items := []string{"A", "B", "A", "C", "A"}

for _, item := range items {
    counts[item]++
}

fmt.Println(counts["A"])
```

Output:

```text
3
```

Because the zero value of `int` is `0`, `counts[item]++` works even if the key doesn't exist yet.

### Grouping values

```go
groups := make(map[string][]string)

groups["backend"] = append(groups["backend"], "Go")
groups["backend"] = append(groups["backend"], "SQL")
```

### Metadata

```go
metadata := map[string]string{
    "request_id": "abc123",
    "source":     "api",
}
```

Make sure to validate untrusted input and protect shared maps when concurrent access is involved.

---

## 20. Common Mistakes

### Mistake 1: Writing to a nil map

```go
var m map[string]int
m["A"] = 10 // Panic
```

Fix:

```go
m := make(map[string]int)
m["A"] = 10
```

### Mistake 2: Treating a missing key as definitely absent

A missing key and an existing key storing the zero value look identical with a single-value lookup.

Use comma-ok if the distinction matters.

### Mistake 3: Assuming map iteration order is stable

Never depend on the order returned by `range`.

### Mistake 4: Comparing two maps with `==`

Maps can only be compared directly with `nil`.

### Mistake 5: Assuming a function gets a deep copy

A copied map value can modify the original underlying map.

### Mistake 6: Sharing a map across goroutines without synchronization

Protect concurrent accesses with a mutex or use an appropriate concurrent data structure.

### Mistake 7: Assuming assignment makes an independent map

```go
m2 := m1
```

does not create an independent copy of all the entries.

---

## 21. Day 5 Interview Questions

**1. What is a map in Go?**

A built-in data structure that associates unique keys with values.

**2. How do you create a map?**

Using `make` or a map literal.

```go
m := make(map[string]int)
```

**3. What happens when you read a missing key?**

The lookup returns the value type's zero value.

**4. What is the comma-ok idiom?**

A two-value lookup used to determine whether a key exists.

```go
value, ok := m[key]
```

**5. What happens when you assign to a nil map?**

The program panics.

**6. Is map iteration order guaranteed?**

No.

**7. Can slices be map keys?**

No, because slices are not comparable types.

**8. Can two maps be compared using `==`?**

No. A map can only be compared directly with `nil`.

**9. What happens when a map is passed to a function?**

The map value is copied, but both copies can refer to the same underlying map.

**10. Why doesn't reassignment of a local map parameter replace the caller's map?**

The reassignment changes only the local map variable.

**11. Are Go maps safe for concurrent writes?**

No. Use synchronization, such as a mutex, or an appropriate concurrent data structure.

**12. What is the difference between a nil map and an empty map?**

Both have length zero, but a nil map cannot be written to, whereas an initialized empty map can.

---

## 22. Day 5 Golden Mental Model

```text
Map
 ├── Key → unique, comparable type
 ├── Value → any type
 ├── Missing lookup → zero value
 ├── Comma-ok → checks key existence
 ├── Nil map → readable but not writable
 ├── range → unordered iteration
 └── Concurrent writes → require synchronization
```

For map copying:

```text
m2 := m1
    ↓
Map value copied
    ↓
Underlying map storage shared

m2 = make(map[K]V)
    ↓
m2 now refers to a new map
    ↓
m1 remains unchanged
```

### Final takeaways

- Use `make` or a map literal to initialize a map.
- Use comma-ok when the existence of a key matters.
- Never assume map iteration order.
- Maps require comparable key types.
- Maps cannot be directly compared with other maps.
- Copying a map value does not clone the entries.
- Reassigning a local map variable doesn't replace the caller's map.
- Synchronize concurrent access to shared maps.
