package main

import (
	"testing"
)

// func TestAdd(t *testing.T) {
// 	got := Add(2, 3)

// 	if got != 5 {
// 		t.Errorf("got %d, want %d", got, 5)
// 	}
// }

// func TestExample(t *testing.T) {
// 	t.Error("something went wrong") // This will mark the test as failed but continue execution

// 	fmt.Println("hello") // This will be printed even though the test failed
// }

// func TestExample1(t *testing.T) {
// 	t.Fatal("something went wrong") // This will mark the test as failed and stop execution

// 	fmt.Println("hello") // This will not be printed because the test failed and execution stopped
// }

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(10, 20)
	}
}