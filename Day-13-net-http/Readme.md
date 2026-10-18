# Go Day 13 — HTTP Servers, Clients & Context

HTTP handling is a fundamental skill for Go backend developers. Go's standard library provides the `net/http` package to build HTTP servers, create APIs, call external services, handle requests, and write responses.

These concepts form the foundation of frameworks such as Gin.

## 1. What is `net/http`?

The `net/http` package provides HTTP server and client functionality.

```go
import "net/http"
```

A typical backend request flow looks like this:

```text
HTTP Client
     ↓
HTTP Request
     ↓
Router
     ↓
Handler
     ↓
Service
     ↓
Repository / DAO
     ↓
Database
     ↓
HTTP Response
```

Go can handle HTTP requests directly without using Gin or another framework.

## 2. Your First HTTP Server

Example:

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello from Go!")
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/hello", helloHandler)

    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    log.Fatal(server.ListenAndServe())
}
```

Run:

```powershell
go run .
```

Open:

```text
http://localhost:8080/hello
```

Response:

```text
Hello from Go!
```

Important components:

- `http.NewServeMux()` creates a router.
- `HandleFunc()` registers a handler function.
- `http.ResponseWriter` is used to construct the response.
- `*http.Request` contains information about the request.
- `ListenAndServe()` starts the HTTP server.

For production servers, configure appropriate timeouts rather than relying solely on the default settings.

## 3. `http.Handler` and `http.HandlerFunc`

The core HTTP handler interface is:

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

A type implementing `ServeHTTP` satisfies `http.Handler`.

Example:

```go
type HelloHandler struct{}

func (h HelloHandler) ServeHTTP(
    w http.ResponseWriter,
    r *http.Request,
) {
    fmt.Fprintln(w, "Hello")
}
```

You can pass it directly to the server:

```go
http.ListenAndServe(":8080", HelloHandler{})
```

You can also register ordinary functions:

```go
func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello")
}
```

```go
mux.HandleFunc("/hello", helloHandler)
```

`http.HandlerFunc` is an adapter that allows a compatible function to implement `http.Handler`.

Remember:

```text
http.Handler
    → interface with ServeHTTP()

http.HandlerFunc
    → adapts a function into an HTTP handler
```

## 4. Understanding `http.ResponseWriter`

`http.ResponseWriter` is the interface used to write an HTTP response.

Example:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusCreated)

    fmt.Fprintln(w, "Resource created")
}
```

The usual order is:

```text
1. Set response headers
2. Write status code
3. Write response body
```

### Important rule

The first operation that commits the response determines the status code. If you write the body first, Go implicitly commits `200 OK`.

Incorrect:

```go
fmt.Fprintln(w, "Created")
w.WriteHeader(http.StatusCreated)
```

Correct:

```go
w.WriteHeader(http.StatusCreated)
fmt.Fprintln(w, "Created")
```

Once the response has been committed, you cannot change its status code by calling `WriteHeader()` again.

## 5. Routing with `http.ServeMux`

Create a router:

```go
mux := http.NewServeMux()

mux.HandleFunc("/health", healthHandler)
mux.HandleFunc("/hello", helloHandler)
```

Go 1.22 and later support method-aware patterns and path wildcards.

Example:

```go
mux.HandleFunc("GET /users", getUsers)
mux.HandleFunc("POST /users", createUser)
mux.HandleFunc("GET /users/{id}", getUser)
```

For a request to:

```text
GET /users/123
```

you can retrieve the ID using:

```go
func getUser(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    fmt.Fprintln(w, "User ID:", id)
}
```

Output:

```text
User ID: 123
```

The value is a string. Validate and convert it if your application requires an integer or another type.

## 6. HTTP Methods

Common HTTP methods and their typical purposes:

| Method | Purpose |
|---|---|
| `GET` | Retrieve a resource |
| `POST` | Create a resource or submit data |
| `PUT` | Replace a resource |
| `PATCH` | Partially update a resource |
| `DELETE` | Delete a resource |

Example:

```go
mux.HandleFunc("GET /users", getUsers)
mux.HandleFunc("POST /users", createUser)
mux.HandleFunc("DELETE /users/{id}", deleteUser)
```

Keep each handler focused on its method and responsibility.

Test unsupported methods as part of your API contract. An unsupported method normally results in `405 Method Not Allowed`.

## 7. Path Parameters vs Query Parameters

### Path parameters

Request:

```text
GET /users/123
```

Retrieve the parameter:

```go
id := r.PathValue("id")
```

Path parameters commonly identify the resource being accessed.

### Query parameters

Request:

```text
GET /transactions?page=2&limit=50
```

Retrieve them:

```go
page := r.URL.Query().Get("page")
limit := r.URL.Query().Get("limit")
```

Both values are strings.

If you need integers:

```go
pageNum, err := strconv.Atoi(page)
if err != nil {
    http.Error(w, "Invalid page", http.StatusBadRequest)
    return
}
```

Query parameters are often used for pagination, sorting, filtering, and optional behavior.

## 8. Returning JSON

Go provides the `encoding/json` package.

For API response data, structs are often clearer than generic maps.

```go
type User struct {
    ID     int    `json:"id"`
    Name   string `json:"name"`
    Active bool   `json:"active"`
}
```

Use a JSON encoder:

```go
func getUser(w http.ResponseWriter, r *http.Request) {
    user := User{
        ID:     123,
        Name:   "Amjad",
        Active: true,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    if err := json.NewEncoder(w).Encode(user); err != nil {
        log.Println("Failed to encode response:", err)
    }
}
```

Response:

```json
{
  "id": 123,
  "name": "Amjad",
  "active": true
}
```

The JSON tags determine the field names used in the JSON representation.

### Important

Setting:

```go
w.Header().Set("Content-Type", "application/json")
```

tells the client that the response body is JSON. It does not convert arbitrary text into valid JSON; you must encode the response correctly.

## 9. Decoding JSON Requests

Suppose a client sends:

```json
{
  "name": "Amjad",
  "email": "amjad@example.com"
}
```

Define the request structure:

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

Decode the body:

```go
func createUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(
            w,
            "Invalid JSON request",
            http.StatusBadRequest,
        )
        return
    }

    fmt.Fprintln(w, "Creating user:", req.Name)
}
```

Valid JSON doesn't necessarily mean valid business input. After decoding, validate fields such as email, name, amount, or user ID according to your rules.

### Limit request body size

For untrusted input, apply an appropriate size limit:

```go
r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MiB
```

You can also use `json.Decoder.DisallowUnknownFields()` if your API needs to reject unknown JSON fields.

## 10. HTTP Status Codes

| Status | Meaning |
|---|---|
| `200 OK` | Request succeeded |
| `201 Created` | Resource created |
| `204 No Content` | Success without a response body |
| `400 Bad Request` | Invalid request |
| `401 Unauthorized` | Authentication missing or invalid |
| `403 Forbidden` | Access isn't permitted |
| `404 Not Found` | Resource doesn't exist |
| `405 Method Not Allowed` | Method isn't supported |
| `409 Conflict` | Conflict with existing state |
| `500 Internal Server Error` | Unexpected server error |
| `503 Service Unavailable` | Service temporarily unavailable |

Use meaningful status codes instead of returning `200 OK` for every response.

## 11. JSON Error Responses

A consistent error response makes an API easier for clients to consume.

```go
type ErrorResponse struct {
    Error string `json:"error"`
}
```

Create a helper:

```go
func writeJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    if err := json.NewEncoder(w).Encode(data); err != nil {
        log.Println("JSON encoding error:", err)
    }
}
```

Use it:

```go
writeJSON(w, http.StatusNotFound, ErrorResponse{
    Error: "user not found",
})
```

Response:

```json
{
  "error": "user not found"
}
```

Do not expose raw database errors, connection details, credentials, or tokens to clients. Log safe technical details internally.

Also remember that after the response status has been committed, an encoding error can't change the status already sent.

## 12. Middleware

Middleware wraps another handler to add behavior before or after it runs.

Common uses include:

- Logging
- Authentication
- Authorization
- Request IDs
- Metrics
- CORS
- Recovering from panics at a suitable request boundary

Example:

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(
        w http.ResponseWriter,
        r *http.Request,
    ) {
        start := time.Now()

        log.Println("Request:", r.Method, r.URL.Path)

        next.ServeHTTP(w, r)

        log.Println("Duration:", time.Since(start))
    })
}
```

Register it:

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /health", healthHandler)

server := &http.Server{
    Addr:    ":8080",
    Handler: loggingMiddleware(mux),
}
```

### What does `next.ServeHTTP(w, r)` do?

It forwards the current request and response writer to the next handler.

Flow:

```text
Request
   ↓
Logging Middleware
   ↓
Next Handler
   ↓
Response
   ↓
Middleware continues after next returns
```

Middleware can execute code both before and after `next.ServeHTTP()`.

The example above doesn't capture the final status code. To log status codes, a middleware can wrap `ResponseWriter` and record the status when it is committed.

## 13. Making HTTP Requests

Go can also act as an HTTP client.

```go
resp, err := http.Get("https://example.com")
if err != nil {
    log.Println("Request failed:", err)
    return
}
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
if err != nil {
    log.Println("Failed to read response:", err)
    return
}

fmt.Println(string(body))
```

Remember to close the response body after a successful call.

### Important distinction

`http.Get()` returning `err == nil` does not guarantee `200 OK`.

A server may return `404`, `401`, or `500` while the HTTP request itself completes successfully.

Check the status:

```go
if resp.StatusCode != http.StatusOK {
    // Handle the unexpected response status.
}
```

## 14. HTTP Client Timeouts

Use a configured HTTP client for outbound requests.

```go
client := &http.Client{
    Timeout: 5 * time.Second,
}

resp, err := client.Get("https://example.com")
if err != nil {
    log.Println("Request failed:", err)
    return
}
defer resp.Body.Close()
```

The client timeout limits the total duration of the request, including reading the response body.

Choose timeouts according to the external service and your application's latency requirements.

## 15. What Is `context.Context`?

A context allows code to communicate cancellation, deadlines, and request-scoped values.

```go
import "context"
```

Contexts are important when a request triggers work that may outlive the caller's interest in the result, such as a database query or another HTTP request.

### `context.WithCancel`

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

Calling:

```go
cancel()
```

signals cancellation.

It does **not** forcibly terminate a goroutine. The goroutine or operation must cooperate.

### `context.WithTimeout`

```go
ctx, cancel := context.WithTimeout(
    context.Background(),
    2*time.Second,
)
defer cancel()
```

This sets a deadline approximately two seconds from when the context is created.

### `context.WithDeadline`

```go
deadline := time.Now().Add(2 * time.Second)

ctx, cancel := context.WithDeadline(
    context.Background(),
    deadline,
)
defer cancel()
```

This sets a specific deadline time.

Remember:

```text
WithCancel
→ manually signal cancellation

WithTimeout
→ cancel when a duration expires

WithDeadline
→ cancel when a specified time arrives
```

A child context is also canceled when its parent is canceled. A child cannot extend an earlier parent deadline.

## 16. Important Context Methods

### `ctx.Done()`

Returns a channel that closes when the context is canceled or its deadline expires.

```go
select {
case <-ctx.Done():
    fmt.Println("Stopped:", ctx.Err())
}
```

### `ctx.Err()`

Returns the reason the context ended:

```text
context.Canceled
context.DeadlineExceeded
```

Before cancellation, it normally returns `nil`.

### `ctx.Deadline()`

Returns the deadline and a boolean indicating whether one exists.

```go
deadline, ok := ctx.Deadline()
if ok {
    fmt.Println("Deadline:", deadline)
}
```

### Cancellation is cooperative

```text
cancel()
   ↓
ctx.Done() closes
   ↓
Goroutine observes cancellation
   ↓
Goroutine returns or stops its operation
```

If a goroutine never checks the context and runs arbitrary work indefinitely, cancellation alone doesn't stop it.

## 17. Context Propagation in Backend Code

For request-scoped work, propagate the incoming context through the layers.

```text
HTTP Handler
     ↓
Service
     ↓
Repository
     ↓
Database
```

### Handler

```go
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    user, err := h.service.GetUser(ctx, 123)
    if err != nil {
        http.Error(
            w,
            "Failed to get user",
            http.StatusInternalServerError,
        )
        return
    }

    _ = user
}
```

### Service

```go
func (s *UserService) GetUser(
    ctx context.Context,
    id int,
) (*User, error) {
    return s.repo.GetUser(ctx, id)
}
```

### Repository using GORM

```go
func (r *UserRepository) GetUser(
    ctx context.Context,
    id int,
) (*User, error) {
    var user User

    err := r.db.WithContext(ctx).First(&user, id).Error
    if err != nil {
        return nil, err
    }

    return &user, nil
}
```

The context must be propagated and supported by the downstream operation. Simply accepting a context parameter doesn't automatically make a function cancellable.

Avoid replacing the incoming request context with `context.Background()` during request processing. Doing so can break cancellation and deadline propagation.

## 18. Context Cancellation vs HTTP Server Timeouts

These concepts complement each other.

**HTTP server timeouts** configure limits on reading headers, reading request bodies, writing responses, and keeping idle connections open.

**Context deadlines** provide a cancellation signal to work such as database operations or outgoing HTTP requests.

Example server configuration:

```go
server := &http.Server{
    Addr:              ":8080",
    Handler:           mux,
    ReadHeaderTimeout: 5 * time.Second,
    ReadTimeout:       10 * time.Second,
    WriteTimeout:      15 * time.Second,
    IdleTimeout:       60 * time.Second,
}
```

Example downstream operation:

```go
ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
defer cancel()
```

These timeout values are illustrative, not universal. Streaming, uploads, service latency, and deployment requirements can require different settings.

## 19. Testing HTTP Handlers

The `net/http/httptest` package allows you to test handlers without starting a real network server.

Example:

```go
func helloHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "Hello")
}
```

Test:

```go
func TestHelloHandler(t *testing.T) {
    req := httptest.NewRequest(
        http.MethodGet,
        "/hello",
        nil,
    )

    rec := httptest.NewRecorder()

    helloHandler(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf(
            "expected status %d, got %d",
            http.StatusOK,
            rec.Code,
        )
    }

    if got := strings.TrimSpace(rec.Body.String()); got != "Hello" {
        t.Errorf("expected %q, got %q", "Hello", got)
    }
}
```

Required imports:

```go
import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)
```

Useful recorder fields:

```go
rec.Code
rec.Body.String()
rec.Header()
```

You can also test the router by calling:

```go
mux.ServeHTTP(rec, req)
```

This tests routing and handler behavior together.

## 20. Production Server Configuration

For a real application, use an explicit `http.Server`:

```go
server := &http.Server{
    Addr:              ":8080",
    Handler:           mux,
    ReadHeaderTimeout: 5 * time.Second,
    ReadTimeout:       10 * time.Second,
    WriteTimeout:      15 * time.Second,
    IdleTimeout:       60 * time.Second,
}

log.Fatal(server.ListenAndServe())
```

The timeouts help reduce resource exhaustion from stalled connections and slow clients.

Also consider graceful shutdown:

```text
Receive SIGTERM
      ↓
Stop accepting new work
      ↓
Cancel background work when appropriate
      ↓
Wait for active work
      ↓
Close resources
```

Use facilities such as `signal.NotifyContext` and `http.Server.Shutdown()` to implement controlled shutdown.

## 21. `net/http` vs Gin

| Standard library | Gin |
|---|---|
| `http.Handler` | Gin handlers are adapted to the HTTP server |
| `http.ResponseWriter` | `c.Writer` |
| `*http.Request` | `c.Request` |
| `http.NewServeMux()` | `gin.New()` or `gin.Default()` |
| Custom middleware | `router.Use(...)` |
| `json.NewEncoder(w).Encode(data)` | `c.JSON(status, data)` |
| `httptest` | Can also test Gin handlers and routes |

Gin simplifies routing and common web tasks, but the underlying HTTP concepts remain relevant.

## 22. Common Mistakes

1. **Writing the body before the status code:** The response status may already have become `200 OK`.
2. **Ignoring the response status:** `err == nil` doesn't guarantee `200 OK`.
3. **Not closing response bodies:** This can interfere with resource release and connection reuse.
4. **Using no suitable timeouts:** Requests can consume resources for too long.
5. **Not validating decoded JSON:** Valid JSON can still contain invalid business input.
6. **Exposing internal errors:** Database and infrastructure details shouldn't automatically be returned to clients.
7. **Ignoring context cancellation:** A downstream operation may continue after its result is no longer needed.
8. **Using `context.Background()` unnecessarily:** This can break request-scoped cancellation propagation.
9. **Forgetting to test status codes:** Test both successful and error responses.

## 23. Day 13 Interview Questions

**1. What is `net/http`?**

Go's standard library package for HTTP servers, clients, handlers, and related functionality.

**2. What is `http.Handler`?**

An interface that requires a `ServeHTTP(ResponseWriter, *Request)` method.

**3. What is the difference between path and query parameters?**

Path parameters identify resources, while query parameters commonly specify filtering, pagination, or optional behavior.

**4. Why set `Content-Type`?**

It identifies the response body's media type, such as JSON.

**5. What does `next.ServeHTTP(w, r)` do?**

It passes the request to the next handler in the chain.

**6. Does `http.Get()` with no error guarantee a successful status?**

No. The server may return any HTTP status, including `404` or `500`.

**7. Why configure HTTP client timeouts?**

To bound how long a request may take and avoid waiting indefinitely.

**8. What does `context.WithCancel()` do?**

Creates a derived context that can be canceled using the returned cancel function.

**9. What is the difference between `WithTimeout` and `WithDeadline`?**

`WithTimeout` uses a relative duration; `WithDeadline` specifies an absolute deadline time.

**10. Does cancellation forcibly kill a goroutine?**

No. The goroutine or operation must cooperate.

**11. What does `ctx.Done()` return?**

A channel that closes when the context is canceled or its deadline expires.

**12. What is the purpose of `httptest`?**

To test HTTP handlers and routes without launching a real network server.

**13. Why should context flow from handler to repository?**

So cancellation and deadlines can propagate to downstream operations.

**14. What's the difference between an HTTP server timeout and a context deadline?**

Server timeouts configure HTTP connection and read/write behavior; context deadlines communicate cancellation to associated work.

## 24. Day 13 Golden Mental Model

```text
net/http
    ↓
Server + Router + Handler
    ↓
Parse request
    ↓
Validate input
    ↓
Call Service with context
    ↓
Repository / External API
    ↓
Handle errors
    ↓
Write status + JSON response
```

Keep these rules in mind:

```text
ResponseWriter
→ writes the response

Request
→ contains request information and context

Handler
→ handles an HTTP request

Middleware
→ wraps handlers

Client timeout
→ limits an outgoing request

Context
→ communicates cancellation and deadlines

httptest
→ tests HTTP behavior without a real network server
```

### Final takeaway

The key to Day 13 is understanding the complete HTTP lifecycle, not just registering routes. A robust Go API parses and validates inputs, calls appropriate service logic, propagates contexts, handles errors, uses correct status codes, and applies suitable timeout and shutdown behavior.
