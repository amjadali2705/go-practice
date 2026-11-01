# Day 15 – Building REST APIs with Database in Go

## 1. Introduction

On Day 15, we combined concepts from:

- Day 13: HTTP server and REST APIs.
- Day 14: Database access using `database/sql`.

We learned how to build a REST API that fetches user information from an SQLite database.

We followed this architecture:

**Handler → Service → Repository → Database**

The main objective was to separate HTTP handling, business logic, and database operations.

---

## 2. Handler–Service–Repository Architecture

A typical Go backend application separates responsibilities into different layers.

### A. Handler Layer

The handler communicates with the HTTP client.

Responsibilities:
- Receive HTTP requests.
- Extract path parameters and query parameters.
- Perform basic request validation.
- Call the service layer.
- Return HTTP responses in JSON format.
- Convert application errors into HTTP status codes.

Example:

```go
func (h *UserHandler) GetAll(
    w http.ResponseWriter,
    r *http.Request,
) {
    users, err := h.service.GetAll(r.Context())

    if err != nil {
        writeJSON(w, http.StatusInternalServerError,
            map[string]string{
                "error": "internal server error",
            })
        return
    }

    writeJSON(w, http.StatusOK, users)
}
```

### B. Service Layer

The service contains the application's business logic.

Responsibilities:
- Implement business rules.
- Validate business inputs.
- Coordinate repository operations.
- Propagate meaningful errors.

Example:

```go
func (s *UserService) GetUserByID(
    ctx context.Context,
    id int64,
) (User, error) {

    if id <= 0 {
        return User{}, ErrInvalidUserID
    }

    return s.repo.GetByID(ctx, id)
}
```

### C. Repository Layer

The repository interacts with the database.

Responsibilities:
- Execute SQL queries.
- Insert, update, delete and retrieve records.
- Scan database results into Go structures.
- Handle database-specific errors.
- Translate database errors into application errors when appropriate.

Example:

```go
func (r *SQLiteUserRepository) GetByID(
    ctx context.Context,
    id int64,
) (User, error) {

    var user User

    query := `
        SELECT id, name, email
        FROM users
        WHERE id = ?
    `

    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
    )

    if errors.Is(err, sql.ErrNoRows) {
        return User{}, ErrUserNotFound
    }

    if err != nil {
        return User{}, fmt.Errorf("get user: %w", err)
    }

    return user, nil
}
```

### Responsibility comparison

| Layer | Responsibility |
|---|---|
| Handler | HTTP requests and responses |
| Service | Business logic |
| Repository | Database operations |
| Database | Data storage |

**Important:** Avoid writing SQL queries directly inside HTTP handlers.

---

## 3. Model Layer

A model represents the structure of application data.

Example:

```go
type User struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

JSON tags control the names of fields when converting Go structures into JSON.

Example response:

```json
{
    "id": 101,
    "name": "Amjad",
    "email": "amjad@example.com"
}
```

---

## 4. Dependency Injection in Go

Dependency Injection means providing an object with its required dependencies from outside rather than making it create those dependencies itself.

Example:

```go
userRepo := NewSQLiteUserRepository(db)

userService := NewUserService(userRepo)

userHandler := NewUserHandler(userService)
```

Here:

- The repository receives the database handle.
- The service receives the repository.
- The handler receives the service.

### Why use Dependency Injection?

1. Reduces tight coupling.
2. Makes unit testing easier.
3. Allows replacing implementations.
4. Improves maintainability.
5. Makes dependencies explicit.

### Repository interface

```go
type UserRepository interface {
    GetByID(ctx context.Context, id int64) (User, error)
    GetAll(ctx context.Context) ([]User, error)
    GetByUserName(ctx context.Context, name string) ([]User, error)
}
```

The interface defines the behavior that the service expects from a repository.

SQLite implements these methods.

A different repository implementation could use MariaDB or PostgreSQL while satisfying the same interface.

---

## 5. Creating REST API Routes

In Go 1.22+, `http.ServeMux` supports method-aware route patterns and path parameters.

Example:

```go
mux := http.NewServeMux()

mux.HandleFunc("GET /users", userHandler.GetAll)
mux.HandleFunc("GET /users/{id}", userHandler.GetUserByID)
```

The route registration specifies the HTTP method and path.

### Path parameters

Example request:

```http
GET /users/101
```

Extract the user ID:

```go
id := r.PathValue("id")
```

### Query parameters

Example request:

```http
GET /users?name=Amjad
```

Extract the name:

```go
name := r.URL.Query().Get("name")
```

Important:

- Path parameters identify a resource or part of a route.
- Query parameters provide optional filtering, searching, sorting or pagination information.

Do not register query parameters as separate routes.

Correct:

```go
mux.HandleFunc("GET /users", userHandler.GetAll)
```

The same route supports:

```http
GET /users
GET /users?name=Amjad
GET /users?name=Rahul
```

### Conditional handler logic

```go
name := r.URL.Query().Get("name")

if name != "" {
    h.GetUserByUserName(w, r, name)
    return
}

// Otherwise fetch all users.
```

The handler decides which service operation should be called.

Note: In this implementation, an empty `name` value is treated as no filter. An application can instead explicitly reject `?name=` as invalid input.

---

## 6. QueryRowContext() vs QueryContext()

This was one of the important concepts covered during the practical exercise.

### A. QueryRowContext()

Use `QueryRowContext()` when you expect a single row.

Example:

```go
err := db.QueryRowContext(
    ctx,
    "SELECT id, name FROM users WHERE id = ?",
    101,
).Scan(&id, &name)
```

It returns a `*sql.Row`.

If no matching row exists, `Scan()` returns `sql.ErrNoRows`.

### B. QueryContext()

Use `QueryContext()` when you want multiple rows.

Example:

```go
rows, err := db.QueryContext(
    ctx,
    "SELECT id, name FROM users",
)
if err != nil {
    return nil, err
}
defer rows.Close()

var users []User

for rows.Next() {
    var user User

    if err := rows.Scan(&user.ID, &user.Name); err != nil {
        return nil, err
    }

    users = append(users, user)
}

if err := rows.Err(); err != nil {
    return nil, err
}
```

### Comparison

| Function | Use case | Result |
|---|---|---|
| QueryRowContext() | Single row | `*sql.Row` |
| QueryContext() | Multiple rows | `*sql.Rows` |

---

## 7. Searching Users by Name

We implemented an API to search users by name.

Example:

```http
GET /users?name=Amjad
```

### Repository query

```sql
SELECT id, name, email
FROM users
WHERE name = ?
```

Always use parameterized queries instead of concatenating user input into SQL.

### Important: Can multiple users have the same name?

Suppose the database contains:

| ID | Name | Email |
|---:|---|---|
| 101 | Amjad | amjad1@example.com |
| 102 | Amjad | amjad2@example.com |

If we use `QueryRowContext()`, only one matching row is returned.

If we want all matching users, we must use `QueryContext()`.

The return type must also change.

Before:

```go
(User, error)
```

After:

```go
([]User, error)
```

Update the repository interface, repository implementation and service method consistently.

### Example repository implementation

```go
func (r *SQLiteUserRepository) GetByUserName(
    ctx context.Context,
    name string,
) ([]User, error) {

    query := `
        SELECT id, name, email
        FROM users
        WHERE name = ?
    `

    rows, err := r.db.QueryContext(ctx, query, name)
    if err != nil {
        return nil, fmt.Errorf("get users by name: %w", err)
    }
    defer rows.Close()

    users := make([]User, 0)

    for rows.Next() {
        var user User

        if err := rows.Scan(
            &user.ID,
            &user.Name,
            &user.Email,
        ); err != nil {
            return nil, err
        }

        users = append(users, user)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return users, nil
}
```

### Why initialize the slice?

```go
users := make([]User, 0)
```

An empty initialized slice becomes `[]` in JSON.

A nil slice becomes `null`.

For a collection endpoint, returning an empty array when no users match is generally preferable.

---

## 8. Context Propagation

Context carries cancellation signals, deadlines and request-scoped values.

In an HTTP backend, use the incoming request context.

Example:

```go
func (h *UserHandler) GetAll(
    w http.ResponseWriter,
    r *http.Request,
) {
    users, err := h.service.GetAll(r.Context())
}
```

Pass it through the service:

```go
func (s *UserService) GetAll(
    ctx context.Context,
) ([]User, error) {
    return s.repo.GetAll(ctx)
}
```

Then use it in the repository:

```go
rows, err := r.db.QueryContext(ctx, query)
```

### Why is this important?

Suppose a client disconnects while a database query is running.

The HTTP request context is cancelled.

That cancellation signal can propagate through the service to the repository and database operation.

The driver and database operation must support cancellation for it to take effect.

Recommended flow:

**Handler → Service → Repository → Database**

Use `r.Context()` rather than creating an unrelated `context.Background()` inside every layer of a request.

---

## 9. Error Handling and HTTP Status Codes

The repository should translate database errors into meaningful application errors.

Example:

```go
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidUserID = errors.New("invalid user ID")
var ErrInvalidUserName = errors.New("invalid username")
```

Use `errors.Is()` in the handler to identify errors.

Example:

```go
if errors.Is(err, ErrUserNotFound) {
    writeJSON(w, http.StatusNotFound,
        map[string]string{
            "error": err.Error(),
        })
    return
}
```

### Common HTTP status codes

| Status | Meaning |
|---|---|
| 200 | OK |
| 201 | Created |
| 204 | No Content |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 409 | Conflict |
| 500 | Internal Server Error |

For a collection search such as `GET /users?name=Unknown`, returning `200 OK` with `[]` is generally appropriate because the collection request succeeded but found no matching records.

Never expose raw database errors or sensitive internal details to API clients.

---

## 10. Testing the REST API

We used PowerShell and `curl.exe` to test the endpoints.

### Get all users

```powershell
curl.exe -i "http://localhost:8080/users"
```

### Get one user by ID

```powershell
curl.exe -i "http://localhost:8080/users/101"
```

### Search users by name

```powershell
curl.exe -i "http://localhost:8080/users?name=Amjad"
```

### Search for a nonexistent name

```powershell
curl.exe -i "http://localhost:8080/users?name=Unknown"
```

### Search with spaces

```powershell
curl.exe -i "http://localhost:8080/users?name=Md%20Amjad%20Ali"
```

Expected response for multiple matching users:

```json
[
    {
        "id": 101,
        "name": "Amjad",
        "email": "amjad1@example.com"
    },
    {
        "id": 102,
        "name": "Amjad",
        "email": "amjad2@example.com"
    }
]
```

---

## 11. Important Interview Questions

**Q1. What is the responsibility of the handler layer?**

It receives HTTP requests, validates request input, calls the service and returns HTTP responses.

**Q2. What is the responsibility of the service layer?**

It handles business logic and business validations.

**Q3. What is the responsibility of the repository layer?**

It executes database operations and retrieves or modifies data.

**Q4. What is Dependency Injection?**

Providing an object with its required dependencies from outside instead of making it create those dependencies itself.

**Q5. Why use interfaces in Go?**

Interfaces define behavior and allow different implementations. They also make testing easier by allowing fake implementations.

**Q6. What is the difference between QueryContext() and QueryRowContext()?**

QueryContext() retrieves multiple rows, while QueryRowContext() is used when expecting a single row.

**Q7. How do you retrieve a query parameter?**

Use:

```go
name := r.URL.Query().Get("name")
```

**Q8. How do you retrieve a path parameter?**

Use:

```go
id := r.PathValue("id")
```

**Q9. Why use r.Context()?**

To propagate request cancellation, deadlines and request-scoped values through application layers.

**Q10. Why shouldn't a handler directly execute SQL queries?**

It separates HTTP responsibilities from database operations and makes the application easier to test and maintain.

**Q11. Why use QueryContext() when multiple users can have the same name?**

Because it returns multiple rows, allowing us to collect every matching user into a slice.

**Q12. What is the difference between a nil slice and an empty slice in JSON?**

A nil slice serializes as `null`, while an initialized empty slice serializes as `[]`.

---

## 12. Day 15 Practical Assignment

We implemented the following functionality:

- Database connection using SQLite.
- User model.
- Repository interface.
- Repository implementation.
- Service layer.
- HTTP handler.
- GET /users endpoint.
- GET /users/{id} endpoint.
- GET /users?name=Amjad endpoint.
- Error handling.
- Context propagation.
- Testing with curl.

### Additional practice

Build the following endpoints using the same architecture:

1. `POST /users` – Create a user.
2. `PUT /users/{id}` – Update user information.
3. `DELETE /users/{id}` – Delete a user.
4. `GET /users?name=Amjad` – Search users by name.

Also write unit tests for your service and handler layers.

---

## 13. Day 15 – Key Takeaways

1. A clean Go backend separates Handler, Service and Repository responsibilities.
2. The handler handles HTTP requests and responses.
3. The service contains business logic.
4. The repository executes database queries.
5. Dependency Injection reduces tight coupling.
6. Interfaces make code easier to test and extend.
7. Use `r.URL.Query().Get()` for query parameters.
8. Use `r.PathValue()` for path parameters.
9. Use `QueryRowContext()` when expecting one row.
10. Use `QueryContext()` when retrieving multiple rows.
11. Pass request context through every layer.
12. Use `errors.Is()` to identify application errors.
13. Convert internal errors into appropriate HTTP status codes.
14. Return an empty JSON array `[]` when a collection search has no matches.

**Day 15 Completed: Building REST APIs with Database in Go.**