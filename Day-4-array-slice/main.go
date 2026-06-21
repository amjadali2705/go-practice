package main

import "fmt"

func main() {
	arr := [3]int{1, 2, 3}
	fmt.Println(len(arr)) //output - 3

	s := []int{1, 2, 3, 4}
	fmt.Println(len(s)) //output - 4

	arr1 := [5]int{1, 2, 3, 4, 5}
	s1 := arr1[1:3]
	fmt.Println(len(s1)) // 2
	fmt.Println(cap(s1)) // 4

	arr2 := [3]int{1, 2, 3}
	s2 := arr2[:]
	s2[1] = 100
	fmt.Println(arr2) // [1, 100, 3]

	var s3 []int
	fmt.Println(s3 == nil) // true
	fmt.Println(len(s3))   // 0
	fmt.Println(cap(s3))   // 0

	s4 := make([]int, 2, 5)
	fmt.Println(len(s4)) // 2
	fmt.Println(cap(s4)) // 5

	arr3 := [5]int{1, 2, 3, 4, 5}
	s5 := arr3[1:4]
	s6 := arr3[2:5]
	s5[1] = 100
	fmt.Println(s5)   // [2, 100, 4]
	fmt.Println(s6)   // [100, 4, 5]
	fmt.Println(arr3) // [1, 2, 100, 4, 5]

	arr4 := [5]int{1, 2, 3, 4, 5}
	s7 := arr4[1:3]
	s8 := append(s7, 100)
	fmt.Println(arr4) // [1, 2, 3, 100, 5]
	fmt.Println(s7)   // [2, 3]
	fmt.Println(s8)   // [2, 3, 100]

	arr5 := [5]int{1, 2, 3, 4, 5}
	s9 := arr5[1:3]
	s10 := append(s9, 100)
	s11 := append(s9, 200)
	fmt.Println(arr5) // [1, 2, 3, 200, 5]
	fmt.Println(s10)  // [2, 3, 200]
	fmt.Println(s11)  // [2, 3, 200]

	arr6 := [5]int{1, 2, 3, 4, 5}
	s12 := arr6[1:3]
	s13 := append(s12, 100)
	arr6[3] = 999
	fmt.Println(s13) // [2, 3, 999]

	arr7 := [3]int{1, 2, 3}
	s14 := arr7[:]
	s15 := append(s14, 4)
	s15[0] = 100
	fmt.Println(arr7) // [1, 2, 3]
	fmt.Println(s14)  // [1, 2, 3]
	fmt.Println(s15)  // [100, 2, 3, 4]

	a := []int{1, 2, 3}
	modify(a)
	fmt.Println(a) // [100, 2, 3]

	a1 := []int{1, 2, 3}
	modify1(a1)
	fmt.Println(a1) // [1, 2, 3]

	b := make([]int, 3, 5)
	b[0] = 1
	b[1] = 2
	b[2] = 3
	b = modify3(b)
	fmt.Println(b) // [100, 2, 3, 4]
}

func modify(s []int) {
	s[0] = 100
}

func modify1(s []int) {
	s = append(s, 4)
}

func modify3(s []int) []int {
	s[0] = 100
	s = append(s, 4)
	return s
}

// Arrays
// [3]int
// Fixed size
// Passed by value

// Slices
// []int
// Internally:
// Pointer

// Length
// Capacity
// len()
// len(s)
// Number of visible elements.

// cap()
// cap(s)
// Available storage from slice start to end of backing array.

// Slices Share Memory
// s := arr[:]
// Changing s may change arr.

// append()
// s = append(s, x)
// Reuses backing array if capacity exists.
// Allocates a new backing array if capacity is full.

// Always Capture append()
// s = append(s, x)

// Never:
// append(s, x)
// without using the returned slice.

// Pass-by-Value + Slices
// Go is still pass-by-value.
// The slice header is copied.
// The underlying array is shared.

// If you want to modify the slice itself, return it from the function.
// If you want to modify the underlying array, just modify the slice.
