package main

import "fmt"

func main() {
	m := make(map[string]int)
	m["A"] = 10
	fmt.Println(m["A"]) //output - 10

	var m1 map[string]int
	fmt.Println(m1 == nil) //output - true

	var m2 map[string]int
	fmt.Println(m2["A"]) //output - 0 (zero value)

	//var m3 map[string]int
	//m3["A"] = 10 // This will cause a runtime panic: assignment to entry in nil map

	m4 := map[string]int{
		"A": 10,
	}
	value, ok := m4["B"]
	fmt.Println(value) //output - 0 (zero value)
	fmt.Println(ok)    //output - false

	m5 := map[string]int{
		"A": 10,
	}
	delete(m5, "A")
	fmt.Println(len(m5)) //output - 0

	m6 := map[string]int{}
	update(m6)
	fmt.Println(m6) //output - map[Go:100]

	m7 := map[string]int{
		"A": 10,
	}
	update1(m7)
	fmt.Println(m7) //output - map[A:10]

	m8 := map[string]int{
		"A": 10,
		"B": 20,
	}
	update2(m8)
	fmt.Println(m8) //output - map[B:20]

	m9 := map[string]int{}
	m9["A"] = 10
	m9["A"] = 20
	fmt.Println(m9) //output - map[A:20]

	m10 := map[string]int{
		"A": 10,
	}
	update3(m10)
	fmt.Println(m10) //output - map[A:100]
}

func update(m map[string]int) {
	m["Go"] = 100
}

func update1(m map[string]int) {
	m = make(map[string]int)
	m["A"] = 100
}

func update2(m map[string]int) {
	delete(m, "A")
}

func update3(m map[string]int) {
	m["A"] = 100
	m = make(map[string]int)
	m["B"] = 200
}

// Maps store key-value pairs.
// make() initializes a map.
// A nil map can be read but cannot be written to.
// Use value, ok := m[key] to distinguish missing keys from zero values.
// delete() is safe even if the key doesn't exist.
// Maps are backed by a hash table.
// A map variable is a map header.
// Go is always pass-by-value; the map header is copied, but the underlying hash table is shared.
// Maps cannot be compared (except with nil).
// Map iteration order is random.
// Maps are not safe for concurrent writes.
