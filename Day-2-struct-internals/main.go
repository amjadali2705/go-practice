package main

import (
	"fmt"
	"unsafe"
)

// 1. Primitive Type Sizes (Must Memorize)

// On a 64-bit machine:

// Type	Size
// bool	1 byte
// int8	1 byte
// int16	2 bytes
// int32	4 bytes
// int64	8 bytes
// byte	1 byte
// rune	4 bytes
// pointer	8 bytes
// Important aliases:
// byte = uint8
// rune = int32
// 2. Alignment Rule

// A type likes to start at addresses divisible by its alignment.

// Examples:

// Type	Can Start At
// bool	0,1,2,3…
// int16	0,2,4,6…
// int32	0,4,8,12…
// int64	0,8,16…
// 3. Padding

// Padding = empty bytes inserted by compiler so next field starts at correct alignment.

// Example:

// type T struct {
//     A bool
//     B int32
// }

// Memory:

// [A][P][P][P][B][B][B][B]
// A = 1
// Padding = 3
// B = 4

// Total = 8 bytes

// 4. Field Ordering Rule
// ❌ Bad
// type S1 struct {
//     A bool
//     B int64
//     C bool
// }
// ✅ Better
// type S2 struct {
//     B int64
//     A bool
//     C bool
// }
// Rule:

// Put larger fields first, smaller fields later.

// 5. Empty Struct
// type Signal struct{}

// Size:

// 0 bytes

// Useful for:

// sets
// signaling channels
// 6. Check size in Go

// Use:

// unsafe.Sizeof(variable)

// Example:

// fmt.Println(unsafe.Sizeof(S1{}))
// 7. Interview One-Liners
// What is padding?

// Extra bytes inserted by compiler for proper alignment.

// Why alignment?
// To allow CPU-efficient memory access.

// How to optimize struct memory?
// Place large aligned fields first.

// Does field order matter?
// Yes, it affects padding and total struct size.

func main() {
	type A struct {
		X bool
		Y int16
	}

	fmt.Println(unsafe.Sizeof(A{})) // Output: 4 bytes (1 byte for X + 1 byte padding + 2 bytes for Y)

	type B struct {
		X bool
		Y int32
	}
	fmt.Println(unsafe.Sizeof(B{})) // Output: 8 bytes (1 byte for X + 3 bytes padding + 4 bytes for Y)

	type C struct {
		X bool
		Y int16
		Z int32
	}
	fmt.Println(unsafe.Sizeof(C{})) // Output: 8 bytes (1 byte for X + 1 byte padding + 2 bytes for Y + 4 bytes for Z)

	type D1 struct {
		A bool
		B int64
		C bool
	}
	fmt.Println(unsafe.Sizeof(D1{})) // Output: 24 bytes (1 byte for A + 7 bytes padding + 8 bytes for B + 1 byte for C + 7 bytes padding)

	type D2 struct {
		B int64
		A bool
		C bool
	}
	fmt.Println(unsafe.Sizeof(D2{})) // Output: 16 bytes (8 bytes for B + 1 byte for A + 1 byte for C + 6 bytes padding)

	type S1 struct {
		A bool
		B int64
		C bool
	}

	type S2 struct {
		B int64
		A bool
		C bool
	}
	fmt.Println("S1:", unsafe.Sizeof(S1{}))
	fmt.Println("S2:", unsafe.Sizeof(S2{}))

}
