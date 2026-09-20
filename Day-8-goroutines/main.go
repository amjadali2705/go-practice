package main

import (
	"context"
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Println(id, job)
	}
}

func worker1(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopped")
			return

		default:
			fmt.Println("Working")
		}
	}
}

func main() {
	// go fmt.Println("Hello") // This line runs in a separate goroutine, which may not complete before the main function exits
	// fmt.Println("World")    // This line may execute before the goroutine above, leading to "World" being printed before "Hello"

	// wg.Add(2) // This line increments the WaitGroup counter by 2, indicating that two goroutines will be launched
	// go func() {
	// 	defer wg.Done()  // This line ensures that the WaitGroup counter is decremented when the goroutine completes
	// 	fmt.Println("A") // This line runs in a separate goroutine and prints "A"
	// }()
	// go func() {
	// 	defer wg.Done()  // This line ensures that the WaitGroup counter is decremented when the goroutine completes
	// 	fmt.Println("B") // This line runs in a separate goroutine and prints "B"
	// }()
	// wg.Wait()           // This line blocks the main function until the WaitGroup counter reaches zero, ensuring that both goroutines complete before proceeding
	// fmt.Println("Done") // This line prints "Done" after both goroutines have completed

	// counter := 0
	// for i := 0; i < 100; i++ {
	// 	go func() {
	// 		counter++ // This line increments the counter variable, but it is not safe for concurrent access and may lead to a race condition
	// 	}()
	// }

	// ch := make(chan int) // This line creates a new channel of type int, which will be used for communication between goroutines
	// ch <- 10             // This line sends the value 10 into the channel, but it will block because there is no goroutine receiving from the channel yet
	// fmt.Println("Done")  // This line will not be reached because the previous line blocks indefinitely, leading to a deadlock situation

	// ch1 := make(chan int, 2) // This line creates a buffered channel of type int with a capacity of 2, allowing it to hold up to 2 values without blocking
	// ch1 <- 10                // This line sends the value 10 into the buffered channel, which will not block because the channel has enough capacity to hold the value
	// ch1 <- 20                // This line sends the value 20 into the buffered channel, which will not block because the channel has enough capacity to hold the value
	// fmt.Println("Done")      // This line prints "Done" after both values have been sent into the buffered channel without blocking

	// ch2 := make(chan int) // This line creates an unbuffered channel of type int, which will be used for communication between goroutines
	// close(ch2)            // This line closes the channel, indicating that no more values will be sent into it. Any subsequent sends to the channel will cause a panic, and any receives from the channel will return the zero value of the channel's type along with a false boolean value.
	// value, ok := <-ch2    // This line attempts to receive a value from the closed channel. Since the channel is closed, it will return the zero value of the channel's type (0 for int) and a boolean value of false, indicating that the channel is closed and no more values can be received.
	// fmt.Println(value)    // This line prints the value received from the closed channel, which will be 0 since the channel is closed and no values were sent into it.
	// fmt.Println(ok)       // This line prints the boolean value received from the closed channel, which will be false since the channel is closed and no more values can be received.

	// ch3 := make(chan int) // This line creates an unbuffered channel of type int, which will be used for communication between goroutines
	// close(ch3)            // This line closes the channel, indicating that no more values will be sent into it. Any subsequent sends to the channel will cause a panic, and any receives from the channel will return the zero value of the channel's type along with a false boolean value.
	// ch3 <- 10             // This line attempts to send the value 10 into the closed channel. Since the channel is closed, this will cause a panic, as sending to a closed channel is not allowed in Go.

	// ch4 := make(chan string)
	// ch5 := make(chan string)
	// go func() {
	// 	ch4 <- "A" // This line sends the string "A" into the channel ch4, which will block until another goroutine receives from the channel.
	// }()
	// go func() {
	// 	ch5 <- "B" // This line sends the string "B" into the channel ch5, which will block until another goroutine receives from the channel.
	// }()
	// select {
	// case v := <-ch4: // This line attempts to receive a value from the channel ch4. If a value is available, it will be assigned to the variable v, and the corresponding case block will execute.
	// 	fmt.Println(v) // This line prints the value received from the channel ch4, which will be "A" if the case block executes.

	// case v := <-ch5: // This line attempts to receive a value from the channel ch5. If a value is available, it will be assigned to the variable v, and the corresponding case block will execute.
	// 	fmt.Println(v) // This line prints the value received from the channel ch5, which will be "B" if the case block executes.
	// }

	// counter1 := 0
	// for j := 0; j < 1000; j++ {
	// 	go func() {
	// 		counter1++
	// 	}()
	// }
	// fmt.Println(counter1) // This line prints the value of counter1, which may not be 1000 due to the race condition caused by concurrent access to the counter1 variable. The final value of counter1 is unpredictable and may vary between runs, as multiple goroutines are incrementing the same variable without synchronization mechanisms in place.

	// jobs := make(chan int)
	// var wg sync.WaitGroup
	// for i := 1; i <= 2; i++ {
	// 	wg.Add(1)
	// 	go worker(i, jobs, &wg) // This line launches a worker goroutine that will process jobs from the jobs channel. The worker function takes an ID, the jobs channel, and a pointer to the WaitGroup as arguments. The WaitGroup is used to wait for all worker goroutines to finish processing before the main function exits.
	// }
	// for i := 1; i <= 4; i++ {
	// 	jobs <- i // This line sends the job (an integer) into the jobs channel. The worker goroutines will receive these jobs and process them concurrently. The jobs channel is used to distribute the jobs among the worker goroutines.
	// }
	// close(jobs)
	// wg.Wait()
	// fmt.Println("Done")

	// ctx, cancel := context.WithCancel(context.Background()) // This line creates a new context with a cancel function. The context is used to manage the lifecycle of the worker goroutine, allowing it to be stopped when the cancel function is called.
	// go worker1(ctx)
	// time.Sleep(time.Second)
	// cancel() // This line calls the cancel function, which signals the worker goroutine to stop. The worker goroutine will receive the cancellation signal through the context and exit gracefully.
	// time.Sleep(time.Second)

	// var counter atomic.Int64
	// for i := 0; i < 1000; i++ {
	// 	go func() {
	// 		counter.Add(1)
	// 	}()
	// }
	// fmt.Println(counter.Load())

	// ch6 := make(chan int)
	// go func() {
	// 	ch6 <- 10
	// 	fmt.Println("sent")
	// }()
	// fmt.Println(<-ch6)

	ch7 := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			ch7 <- i
		}
	}()
	for v := range ch7 {
		fmt.Println(v)
	}

}

//                  CONCURRENCY
//                       │
//           ┌───────────┴───────────┐
//           ↓                       ↓
//       Goroutines              Communication
//           │                       │
//           ↓                       ↓
//      WaitGroup                 Channels
//           │                       │
//           │                 ┌─────┴─────┐
//           │                 ↓           ↓
//           │             Buffered    Unbuffered
//           │
//           ↓
//    Shared State
//           │
//       ┌───┴────┐
//       ↓        ↓
//     Mutex    Atomic
//       │
//       ↓
//  Synchronization
//       │
//       ↓
//   Race Detector

// A goroutine is a lightweight concurrent execution unit managed by the Go runtime.

// Goroutine execution order is not guaranteed.

// sync.WaitGroup waits for a collection of goroutines to finish; it does not guarantee their execution order.

// A data race occurs when multiple goroutines access the same memory concurrently, at least one access is a write, without proper synchronization.

// A mutex provides mutual exclusion for protecting shared state.

// Channels provide a mechanism for communication and synchronization between goroutines.

// An unbuffered channel synchronizes sender and receiver directly; a buffered channel can hold values up to its capacity before a sender blocks.

// Sending on a closed channel panics, while receiving from a closed channel is allowed.

// select waits on multiple channel operations and executes a ready case.

// The race detector can be used with go test -race or go run -race to detect data races during execution.

// Worker Pool
// A worker pool limits concurrency by using a fixed number of goroutines to process many jobs.

// Channel
// An unbuffered channel requires sender and receiver synchronization.

// Channel Ownership
// The component that knows no more values will be sent should generally close the channel.

// WaitGroup
// WaitGroup synchronizes goroutine completion; it does not protect shared data.

// Mutex
// A mutex protects a critical section so that concurrent goroutines don't unsafely access shared mutable state.

// Atomic
// Atomic operations safely perform individual operations on shared values, but they don't wait for goroutines to finish.

// Context
// Context cancellation is cooperative; it signals goroutines to stop rather than forcibly terminating them.

// Goroutine Leak
// A goroutine leak occurs when a goroutine remains blocked indefinitely and cannot terminate.

// Channel Range
// A range over a channel terminates only when the channel is closed.

// Results Channel
// When multiple workers produce results, the results channel should be closed only after all result-producing workers have finished.
