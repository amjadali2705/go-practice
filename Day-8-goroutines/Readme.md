# Go Day 8 — Goroutines, Channels & Advanced Concurrency

## 1. What Is Concurrency?

**Concurrency** means managing multiple tasks whose execution overlaps. **Parallelism** means executing multiple tasks at the same time.

Example:

```text
Concurrency
Task A ────────┐
      Task B ───┴────────

Parallelism
CPU 1: Task A ──────────
CPU 2: Task B ──────────
```

A program can be concurrent without running all tasks simultaneously.

Concurrency is especially important in Go backend systems handling HTTP requests, database operations, network calls, and background jobs.

---

## 2. Goroutines

A goroutine is a lightweight execution unit managed by the Go runtime.

Example:

```go
package main

import (
    "fmt"
    "time"
)

func sayHello() {
    fmt.Println("Hello")
}

func main() {
    go sayHello()

    time.Sleep(time.Second)
    fmt.Println("Main")
}
```

Possible output:

```text
Hello
Main
```

The `go` keyword starts a goroutine.

```go
go sayHello()
```

However, the `time.Sleep` is only a simple demonstration, not a reliable synchronization technique. The scheduler does not guarantee that a goroutine will finish before `main` exits.

### Important rules

- Goroutines are managed by the Go runtime.
- Creating a goroutine does not automatically wait for its completion.
- The `main` function does not wait for other goroutines when it returns.
- Goroutines can share memory, which makes synchronization important.

---

## 3. Waiting for Goroutines with `sync.WaitGroup`

Use a `WaitGroup` when one goroutine needs to wait for a group of goroutines to finish.

```go
package main

import (
    "fmt"
    "sync"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Println("Worker", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }

    wg.Wait()
    fmt.Println("All workers finished")
}
```

The worker messages can appear in any order, but `All workers finished` prints only after all three workers call `Done()`.

### Important methods

```go
wg.Add(1) // Increment the counter
wg.Done() // Decrement the counter
wg.Wait() // Wait until the counter reaches zero
```

Prefer adding to the counter before starting the goroutine.

**Remember:** `WaitGroup` waits for completion. It does not protect shared variables against data races.

---

## 4. Race Conditions

A data race occurs when multiple goroutines access the same memory concurrently, at least one access is a write, and the accesses are not properly synchronized.

Example:

```go
var counter int

func increment() {
    counter++
}
```

If many goroutines call `increment()` at the same time, they can race.

Even though `counter++` looks like one operation, it involves reading, incrementing, and writing a value. Those steps are not automatically atomic.

### Fix with a mutex

```go
var (
    counter int
    mu      sync.Mutex
)

func increment() {
    mu.Lock()
    counter++
    mu.Unlock()
}
```

You can also use `defer` for unlocking:

```go
func increment() {
    mu.Lock()
    defer mu.Unlock()

    counter++
}
```

A mutex ensures that only one goroutine at a time executes the protected section.

---

## 5. `sync.Mutex` and `sync.RWMutex`

### `sync.Mutex`

A mutex provides exclusive access to shared data.

```go
mu.Lock()
counter++
mu.Unlock()
```

Use it when multiple goroutines need to read or modify shared state safely.

### `sync.RWMutex`

An `RWMutex` supports multiple concurrent readers or one exclusive writer.

```go
var mu sync.RWMutex

func readData() {
    mu.RLock()
    defer mu.RUnlock()

    // Read shared data
}

func updateData() {
    mu.Lock()
    defer mu.Unlock()

    // Modify shared data
}
```

Use `RLock()`/`RUnlock()` for read access and `Lock()`/`Unlock()` for write access.

An `RWMutex` is not automatically faster than a `Mutex`; it depends on the workload.

---

## 6. Channels

Channels allow goroutines to communicate by sending and receiving values.

Create a channel:

```go
ch := make(chan int)
```

Send:

```go
ch <- 10
```

Receive:

```go
value := <-ch
```

For an unbuffered channel, a send and receive synchronize with each other. The sender waits until a receiver is ready, and the receiver waits until a sender is ready.

Example:

```go
package main

import "fmt"

func main() {
    ch := make(chan int)

    go func() {
        ch <- 100
    }()

    value := <-ch
    fmt.Println(value)
}
```

Output:

```text
100
```

### Channel direction

Use directional channel types to clarify ownership and responsibilities:

```go
func producer(ch chan<- int) {
    ch <- 10
}

func consumer(ch <-chan int) {
    value := <-ch
    fmt.Println(value)
}
```

- `chan<- int` means send-only.
- `<-chan int` means receive-only.

These restrictions are checked at compile time.

---

## 7. Buffered vs Unbuffered Channels

### Unbuffered channel

```go
ch := make(chan int)
```

A send needs a corresponding receiver to proceed.

```go
ch <- 10
```

If no receiver becomes ready, the send blocks.

### Buffered channel

```go
ch := make(chan int, 2)
```

The channel can hold two values without a receiver immediately receiving them.

```go
ch <- 10
ch <- 20
```

Both sends can complete because the buffer has capacity two.

A third send blocks until space becomes available.

### Mental model

```text
Unbuffered channel
→ sender and receiver synchronize

Buffered channel
→ values can wait in the buffer
```

Buffering does not remove the need for synchronization in the rest of your program.

---

## 8. Closing Channels

Use `close(ch)` to signal that no more values will be sent on a channel.

```go
ch := make(chan int)

go func() {
    defer close(ch)

    for i := 1; i <= 3; i++ {
        ch <- i
    }
}()

for value := range ch {
    fmt.Println(value)
}
```

Output:

```text
1
2
3
```

When the channel is closed and all buffered values have been received, the `range` loop ends.

### Receive from a closed channel

```go
value, ok := <-ch
```

When a closed, drained channel is received from:

```text
value → zero value of the element type
ok    → false
```

For `int`, the zero value is `0`.

### Important mistakes

```go
ch <- 10
close(ch)
```

Sending after closing a channel panics.

Closing a channel that is already closed also panics. Closing a nil channel panics.

**Golden rule:** Usually, the component that knows no more values will be sent should close the channel. A receiver should not close a channel that another component is still sending on.

---

## 9. Channel `range`

Example:

```go
for value := range ch {
    fmt.Println(value)
}
```

This receives values until the channel is closed and drained.

A common bug is forgetting to close the channel.

```go
ch := make(chan int)

go func() {
    for i := 0; i < 5; i++ {
        ch <- i
    }
}()

for value := range ch {
    fmt.Println(value)
}
```

The loop prints `0` through `4`, then waits forever for another value or closure. If all remaining goroutines are blocked, the program may terminate with a deadlock error.

Correct:

```go
go func() {
    defer close(ch)

    for i := 0; i < 5; i++ {
        ch <- i
    }
}()
```

---

## 10. The `select` Statement

`select` lets a goroutine wait on multiple channel operations.

```go
select {
case value := <-ch:
    fmt.Println("Received:", value)

case ch <- 10:
    fmt.Println("Sent value")
}
```

Go executes one case whose communication can proceed.

If multiple cases are ready, Go selects one pseudo-randomly. There is no guarantee that the first listed case wins.

### Timeout pattern

```go
select {
case result := <-resultCh:
    fmt.Println("Result:", result)

case <-time.After(2 * time.Second):
    fmt.Println("Timed out")
}
```

For repeated operations in production code, `context.WithTimeout` is often useful because it can propagate cancellation throughout a request.

---

## 11. Deadlocks

A deadlock occurs when goroutines cannot make progress because they're waiting on events that cannot happen.

Example:

```go
func main() {
    ch := make(chan int)

    ch <- 10

    fmt.Println(<-ch)
}
```

Why does it deadlock?

1. The channel is unbuffered.
2. The send waits for a receiver.
3. The receiver is never reached because the send blocks.
4. The program cannot progress.

Fix by receiving concurrently or using an appropriate buffered channel.

A buffer isn't always the best fix; choose a design that matches the intended communication and lifecycle.

---

## 12. Goroutine Leaks

A goroutine leak occurs when a goroutine remains alive but can no longer make useful progress, often because it is blocked indefinitely.

Example:

```go
func process() <-chan int {
    ch := make(chan int)

    go func() {
        result := 100
        ch <- result
    }()

    return ch
}
```

If the caller never receives from `ch`, the goroutine can stay blocked forever at the send.

One possible improvement for a single result is a buffer of size one:

```go
func process() <-chan int {
    ch := make(chan int, 1)

    go func() {
        ch <- 100
    }()

    return ch
}
```

This permits the goroutine to send its single result without waiting for an immediate receiver. It does not solve every kind of leak, especially when the work itself can block.

For cancellable or long-running work, use a context and ensure goroutines observe it and return.

---

## 13. Context Cancellation

The `context` package communicates cancellation, deadlines, and request-scoped values.

Example:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

A worker can listen for cancellation:

```go
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("Worker stopped:", ctx.Err())
            return

        default:
            // Do a unit of work.
            return // Example only; real work would loop as needed.
        }
    }
}
```

A clearer long-running worker pattern is:

```go
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return

        case job, ok := <-jobs:
            if !ok {
                return
            }

            _ = job // Process the job.
        }
    }
}
```

Here, `jobs` represents a channel supplied to the worker.

### Important rules

- `cancel()` signals that work should stop.
- Cancellation is cooperative; Go does not forcibly kill the goroutine.
- A goroutine must observe `ctx.Done()` and return.
- Operations that accept a context should receive the relevant context whenever possible.

For a backend request, propagate context from the HTTP handler through the service to the repository and database.

For example, with GORM:

```go
db.WithContext(ctx).Find(&users)
```

This allows database work to respect request cancellation and deadlines when supported by the underlying driver and operation.

---

## 14. Context Timeouts

Use `context.WithTimeout` to set a deadline.

```go
ctx, cancel := context.WithTimeout(
    context.Background(),
    2*time.Second,
)
defer cancel()
```

Then:

```go
select {
case result := <-resultCh:
    fmt.Println(result)

case <-ctx.Done():
    fmt.Println("Operation stopped:", ctx.Err())
}
```

`ctx.Err()` commonly returns `context.Canceled` or `context.DeadlineExceeded`.

### Production lesson

If a caller times out, its work may still continue unless the worker or downstream operation observes cancellation.

A timeout on the caller does not magically interrupt arbitrary Go code.

---

## 15. Worker Pools

A worker pool uses a fixed number of workers to process jobs from a shared jobs channel.

This bounds concurrency rather than launching an unbounded number of goroutines.

Example:

```go
func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
    defer wg.Done()

    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job)
    }
}

func main() {
    jobs := make(chan int)

    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go worker(i, jobs, &wg)
    }

    for i := 1; i <= 10; i++ {
        jobs <- i
    }

    close(jobs)
    wg.Wait()
}
```

### What does this guarantee?

- All ten jobs are sent to the channel.
- Each job is received by one worker in this design.
- Workers process jobs concurrently.
- The job-processing order across workers is not guaranteed.
- `wg.Wait()` returns after the workers finish.

### Where are worker pools useful?

- Database operations with controlled concurrency
- Processing batches of transactions
- Calling external APIs
- File processing
- Background report generation
- Kafka message-processing workloads

Choose the worker count according to workload, resource limits, and downstream capacity.

---

## 16. Fan-Out and Fan-In

### Fan-out

One source distributes work to multiple workers.

```text
Producer
   ↓
 jobs channel
   ├── Worker 1
   ├── Worker 2
   └── Worker 3
```

### Fan-in

Multiple workers send results to one channel.

```text
Worker 1 ──┐
Worker 2 ──┼──→ results channel → Consumer
Worker 3 ──┘
```

Combined pipeline:

```text
Producer
   ↓
 jobs
   ↓
Workers
   ↓
 results
   ↓
Consumer
```

Fan-out improves parallel processing where tasks can be performed independently. Fan-in gives the application a single stream of results to consume.

### Closing the results channel safely

If multiple workers send to `results`, don't let one worker close it while others may still send.

Instead:

```go
go func() {
    wg.Wait()
    close(results)
}()
```

This assumes the `WaitGroup` tracks all result-producing workers and no other sender will write to `results` after they finish.

The consumer can then use:

```go
for result := range results {
    process(result)
}
```

---

## 17. Channel Ownership

Channel ownership is a design rule for making concurrency safer.

Typical worker-pool ownership:

```text
Producer:
- sends jobs
- closes jobs when no more jobs will be sent

Workers:
- receive jobs
- produce results

Results closer:
- waits for all workers to finish
- closes results

Consumer:
- receives results
```

This division prevents common errors such as:

```text
send on closed channel
double close
range waiting forever
```

Make sure all potential senders have stopped before closing a channel they send on.

---

## 18. `sync.Once`

`sync.Once` makes sure a function is executed at most once for a particular `Once` instance.

Example:

```go
var once sync.Once

func initialize() {
    once.Do(func() {
        fmt.Println("Initialize once")
    })
}
```

Even if many goroutines call `initialize()`, the function passed to `Do` executes only once for that `sync.Once`.

Typical use cases include one-time initialization of shared resources.

It does not mean "run once per request" or "once globally for every possible `Once` value." The guarantee belongs to the specific instance.

---

## 19. Atomic Operations

The `sync/atomic` package provides atomic operations for certain values.

Example with a modern typed atomic:

```go
var counter atomic.Int64

counter.Add(1)

value := counter.Load()
```

Atomic operations make individual operations safe for concurrent access.

However, atomic operations do not automatically wait for goroutines to finish.

If 1,000 goroutines increment an atomic counter, the counter operations can be race-free, yet `main` can still return before all goroutines finish if it doesn't synchronize with them.

Use both when necessary:

```text
Atomic → safe counter updates
WaitGroup → wait for completion
```

Use a mutex when a larger critical section or multiple related values need to be protected together.

---

## 20. Mutex vs Atomic vs Channel

Use the tool that best fits the problem.

### Mutex

Best when you need to protect shared mutable state or several related operations.

```go
mu.Lock()
balance -= amount
mu.Unlock()
```

### Atomic

Best for simple shared values where the required operation is supported atomically.

```go
counter.Add(1)
```

### Channel

Best when goroutines need to communicate, pass work, or transfer ownership of values.

```go
jobs <- job
```

Mental model:

```text
Mutex  → protect shared state
Atomic → safe individual atomic operations
Channel → communicate or transfer work/data
```

Channels can also be used for synchronization, but they should not be added merely because the program uses goroutines.

---

## 21. Race Detector

Go includes a race detector.

Run tests with:

```powershell
go test -race ./...
```

Run an application with:

```powershell
go run -race .
```

The race detector identifies many data races that occur during the execution being tested.

It does not guarantee the absence of:

- Deadlocks
- Goroutine leaks
- Logical concurrency bugs
- Starvation
- Races on code paths that were never executed

Run race detection during testing, especially for code that accesses shared state.

---

## 22. Graceful Shutdown

Backend applications should shut down safely when they receive a termination signal.

A typical sequence:

```text
Receive SIGTERM
      ↓
Stop accepting new work
      ↓
Cancel background processing as appropriate
      ↓
Wait for active work
      ↓
Shut down resources
```

Useful tools include:

- `signal.NotifyContext`
- `context.Context`
- `sync.WaitGroup`
- `http.Server.Shutdown`

Example of an HTTP server shutdown:

```go
ctx, cancel := context.WithTimeout(
    context.Background(),
    10*time.Second,
)
defer cancel()

if err := server.Shutdown(ctx); err != nil {
    log.Printf("HTTP shutdown failed: %v", err)
}
```

The timeout bounds how long the shutdown process waits for HTTP server work to finish.

A complete shutdown also needs deliberate handling of background workers, producers, and database resources.

---

## 23. Common Concurrency Mistakes

### Mistake 1: Assuming `main` waits for goroutines

It doesn't. Use proper synchronization when completion matters.

### Mistake 2: Assuming `counter++` is atomic

It isn't. Protect shared updates with a mutex or use an appropriate atomic operation.

### Mistake 3: Closing a channel from the wrong component

If a producer might still send, closing the channel causes a panic when a later send occurs.

### Mistake 4: Forgetting to close a producer-owned channel

Consumers ranging over the channel may wait indefinitely.

### Mistake 5: Assuming `WaitGroup` protects shared variables

It only coordinates completion; it does not prevent data races by itself.

### Mistake 6: Thinking `cancel()` kills a goroutine

Cancellation is a signal. The goroutine must respond to it.

### Mistake 7: Using `time.Sleep()` to synchronize goroutines

Sleeping doesn't guarantee another goroutine finished. Use a `WaitGroup`, channel, mutex, or another appropriate synchronization mechanism.

### Mistake 8: Forgetting about slow or blocked downstream operations

Workers can become stuck on database calls, external APIs, or channel operations. Use cancellation, deadlines, and sound lifecycle management.

---

## 24. Interview Questions

**1. What is the difference between concurrency and parallelism?**

Concurrency deals with coordinating multiple overlapping tasks. Parallelism means executing tasks at the same time.

**2. What is a goroutine?**

A lightweight execution unit managed by the Go runtime.

**3. Why do we use a `WaitGroup`?**

To wait until a group of goroutines completes.

**4. Does a `WaitGroup` prevent data races?**

No. It coordinates completion but doesn't protect concurrent access to shared data.

**5. What is a data race?**

Concurrent access to the same memory location where at least one access is a write and the accesses aren't properly synchronized.

**6. Difference between a buffered and unbuffered channel?**

An unbuffered channel synchronizes sends and receives. A buffered channel can hold a limited number of values without an immediate receiver.

**7. What happens when you send on a closed channel?**

The program panics.

**8. What happens when you receive from a closed, drained channel?**

The receive returns the element type's zero value; the two-value form also returns `false`.

**9. What is a goroutine leak?**

A goroutine that remains alive but cannot make useful progress, often because it is blocked indefinitely.

**10. What does `context.CancelFunc` do?**

It signals cancellation to operations using that context. It doesn't forcefully terminate goroutines.

**11. What is a worker pool?**

A design that uses a bounded number of workers to process jobs from a shared queue or channel.

**12. What is fan-out and fan-in?**

Fan-out distributes work to multiple workers. Fan-in combines results from multiple workers into one stream.

**13. Mutex vs atomic?**

Mutexes protect critical sections. Atomics support particular race-free operations on individual values.

**14. When should a channel be closed?**

Normally, when the component responsible for sending knows no more values will be sent.

**15. What does `go test -race ./...` do?**

Runs package tests with the race detector enabled to identify many data races.

---

## 25. Day 8 Golden Mental Model

```text
Goroutine
→ concurrent execution

WaitGroup
→ wait for completion

Mutex
→ protect shared mutable state

Atomic
→ race-free individual atomic operations

Channel
→ communication and data transfer

Worker pool
→ bounded concurrency

Fan-out
→ distribute work to multiple workers

Fan-in
→ combine results from multiple workers

Context
→ propagate cancellation and deadlines

Channel close
→ signal that no more values will be sent

sync.Once
→ run initialization at most once per instance

Race detector
→ detect many data races
```

### Backend concurrency lifecycle

```text
HTTP Request / Job Producer
          ↓
       Jobs Channel
          ↓
       Worker Pool
          ↓
     Results Channel
          ↓
        Consumer
```

Throughout that lifecycle, use synchronization to protect shared state, close channels only when it is safe to do so, and propagate cancellation so work can stop when it is no longer needed.

**Final takeaway:** Go concurrency isn't just about starting goroutines. It's about coordinating their work, sharing data safely, controlling resource usage, preventing leaks, and ensuring the application can shut down cleanly.
