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




# Day 16 – Complete CRUD APIs in Go

## 1. Introduction to CRUD

CRUD stands for:

- **C – Create:** Add new records.
- **R – Read:** Retrieve existing records.
- **U – Update:** Modify existing records.
- **D – Delete:** Remove records.

In Day 16, we extended our Day 15 REST API project to implement all four CRUD operations using Go, SQLite and `database/sql`.

### CRUD and HTTP methods

| Operation | HTTP method | Purpose |
|---|---|---|
| Create | POST | Create a new user |
| Read | GET | Retrieve users |
| Update | PUT | Update an existing user |
| Delete | DELETE | Delete an existing user |

Our API endpoints:

```text
POST   /users
GET    /users
GET    /users/{id}
PUT    /users/{id}
DELETE /users/{id}
```

## 2. Application Architecture

We continue to follow the Handler → Service → Repository architecture.

```text
HTTP Client
     |
     v
Handler
     |
     v
Service
     |
     v
Repository
     |
     v
SQLite Database
```

### Handler Layer

Responsibilities:
- Receive HTTP requests.
- Decode JSON request bodies.
- Extract path and query parameters.
- Perform request validation.
- Call service methods.
- Return HTTP responses and status codes.

### Service Layer

Responsibilities:
- Implement business logic.
- Validate business rules.
- Check duplicate email addresses.
- Coordinate repository operations.
- Return meaningful application errors.

### Repository Layer

Responsibilities:
- Execute SQL queries.
- Insert, retrieve, update and delete records.
- Scan database results.
- Check affected rows.
- Handle database-specific errors.

**Important:** The handler should not execute SQL queries directly.

## 3. Model and Request Structures

A model represents the application's data.

```go
type User struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

For API requests, we can create separate structures.

### CreateUserRequest

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

### UpdateUserRequest

```go
type UpdateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

### Why use separate request structures?

1. The client should not control the database-generated ID.
2. We control which fields the client can provide.
3. Request validation becomes easier.
4. API input structures remain separate from database models.

## 4. POST – Create a User

POST is used to create a new record.

Example request:

```http
POST /users
Content-Type: application/json
```

Request body:

```json
{
    "name": "Amjad",
    "email": "amjad@example.com"
}
```

Expected response:

```json
{
    "id": 103,
    "name": "Amjad",
    "email": "amjad@example.com"
}
```

HTTP status: `201 Created`.

### A. Decode JSON request body

Go provides `json.NewDecoder()` to decode JSON into a Go structure.

```go
var req CreateUserRequest

err := json.NewDecoder(r.Body).Decode(&req)
if err != nil {
    writeJSON(w, http.StatusBadRequest,
        map[string]string{
            "error": "invalid request body",
        })
    return
}
```

The decoded request can be converted into a User model:

```go
user := User{
    Name:  req.Name,
    Email: req.Email,
}
```

### B. Service validation

```go
func (s *UserService) CreateUser(
    ctx context.Context,
    user User,
) (User, error) {

    user.Name = strings.TrimSpace(user.Name)
    user.Email = strings.TrimSpace(user.Email)

    if user.Name == "" {
        return User{}, ErrInvalidUserName
    }

    if user.Email == "" {
        return User{}, ErrInvalidEmail
    }

    return s.repo.Create(ctx, user)
}
```

The service validates business inputs before asking the repository to create the record.

Email syntax should also be validated. Duplicate checking can be performed using the repository.

### C. Repository INSERT query

```go
func (r *SQLiteUserRepository) Create(
    ctx context.Context,
    user User,
) (User, error) {

    query := `
        INSERT INTO users (name, email)
        VALUES (?, ?)
    `

    result, err := r.db.ExecContext(
        ctx, query, user.Name, user.Email,
    )
    if err != nil {
        return User{}, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return User{}, err
    }

    user.ID = id

    return user, nil
}
```

### LastInsertId()

`LastInsertId()` retrieves the ID generated by the most recent successful INSERT operation on that database connection.

It is useful when the database generates IDs automatically.

### Duplicate email handling

We should enforce email uniqueness in the database:

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE
);
```

The service can check for an existing email, but the database UNIQUE constraint is essential because two concurrent requests might pass the service-level check at the same time.

Translate the database's unique-constraint violation into an application error such as:

```go
ErrEmailAlreadyExists = errors.New("email already exists")
```

Return HTTP `409 Conflict` for this condition.

## 5. PUT – Update an Existing User

PUT is generally used to replace the complete resource representation.

Example:

```http
PUT /users/101
Content-Type: application/json
```

Request body:

```json
{
    "name": "Amjad Ali",
    "email": "amjad.ali@example.com"
}
```

Expected response:

```json
{
    "message": "user updated successfully"
}
```

HTTP status: `200 OK`.

### A. Extract the user ID

```go
id, err := strconv.ParseInt(
    r.PathValue("id"),
    10,
    64,
)

if err != nil || id <= 0 {
    writeJSON(w, http.StatusBadRequest,
        map[string]string{
            "error": "invalid user ID",
        })
    return
}
```

`r.PathValue()` extracts the path parameter, and `strconv.ParseInt()` converts it into an int64.

### B. Repository UPDATE query

```go
query := `
    UPDATE users
    SET name = ?, email = ?
    WHERE id = ?
`
```

Execute it:

```go
result, err := r.db.ExecContext(
    ctx,
    query,
    user.Name,
    user.Email,
    id,
)
if err != nil {
    return err
}
```

### C. Check RowsAffected()

```go
affected, err := result.RowsAffected()
if err != nil {
    return err
}

if affected == 0 {
    return ErrUserNotFound
}
```

Why?

Suppose we execute:

```http
PUT /users/999
```

If user 999 does not exist, the UPDATE can execute without returning an SQL error, but zero rows may be affected.

`RowsAffected()` helps us detect this situation.

Note: A unique email constraint can also cause an UPDATE to fail when another user already owns that email.

### D. Update service validation

Before updating, the service should check:

1. User ID is valid.
2. Name is not empty.
3. Email is not empty.
4. The email is not already owned by another user.

The same user should be allowed to retain their existing email.

## 6. DELETE – Delete an Existing User

DELETE removes a record from the database.

Example:

```http
DELETE /users/101
```

Successful response:

```text
HTTP/1.1 204 No Content
```

### Repository DELETE query

```go
query := `DELETE FROM users WHERE id = ?`

result, err := r.db.ExecContext(ctx, query, id)
if err != nil {
    return err
}
```

Check affected rows:

```go
affected, err := result.RowsAffected()
if err != nil {
    return err
}

if affected == 0 {
    return ErrUserNotFound
}

return nil
```

### Handler response

```go
w.WriteHeader(http.StatusNoContent)
```

HTTP 204 means the operation succeeded and there is no response body.

## 7. PUT vs PATCH

These two HTTP methods are commonly used for updating resources.

| PUT | PATCH |
|---|---|
| Generally replaces the complete resource | Updates selected fields |
| Usually expects all required fields | Can accept only changed fields |
| Example: Update name and email | Example: Update only name |

Example PUT:

```json
{
    "name": "Amjad Ali",
    "email": "amjad@example.com"
}
```

Example PATCH:

```json
{
    "name": "Amjad"
}
```

For our Day 16 exercise, we implemented PUT.

## 8. HTTP Status Codes

| Status | Meaning | Example |
|---|---|---|
| 200 | OK | User retrieved or updated |
| 201 | Created | User created |
| 204 | No Content | User deleted |
| 400 | Bad Request | Invalid JSON or user ID |
| 404 | Not Found | User does not exist |
| 409 | Conflict | Duplicate email |
| 500 | Internal Server Error | Unexpected database failure |

### Important error handling

We use custom errors:

```go
var (
    ErrUserNotFound      = errors.New("user not found")
    ErrInvalidUserID     = errors.New("invalid user ID")
    ErrInvalidUserName   = errors.New("invalid username")
    ErrInvalidEmail      = errors.New("invalid email")
    ErrEmailAlreadyExists = errors.New("email already exists")
)
```

In the handler, use:

```go
if errors.Is(err, ErrUserNotFound) {
    writeJSON(w, http.StatusNotFound,
        map[string]string{
            "error": err.Error(),
        })
    return
}
```

`errors.Is()` checks whether an error matches a particular target, including through wrapped errors.

Never expose raw database errors to clients.

## 9. JSON Request Validation

Go provides `json.NewDecoder()` for decoding JSON request bodies.

```go
var req CreateUserRequest

err := json.NewDecoder(r.Body).Decode(&req)
```

Invalid JSON should produce HTTP 400.

For stronger request handling, consider:

```go
decoder := json.NewDecoder(r.Body)
decoder.DisallowUnknownFields()
```

This rejects unknown JSON fields.

Also consider limiting the request body size and checking that the request contains only one JSON value.

### Common validations

| Field | Validation |
|---|---|
| ID | Positive integer |
| Name | Required, not whitespace-only |
| Email | Required and valid format |
| Email uniqueness | Must not belong to another user |
| Request body | Valid JSON |

## 10. Register CRUD Routes

In `main.go`:

```go
mux.HandleFunc("POST /users", userHandler.Create)

mux.HandleFunc("GET /users", userHandler.GetAll)

mux.HandleFunc("GET /users/{id}", userHandler.GetUserByID)

mux.HandleFunc("PUT /users/{id}", userHandler.Update)

mux.HandleFunc("DELETE /users/{id}", userHandler.Delete)
```

These are method-aware routes supported by Go 1.22 and later.

Query parameters are not registered as separate routes.

For example:

```http
GET /users
GET /users?name=Amjad
```

Both use the same registered route.

## 11. Testing APIs Using PowerShell

We used `Invoke-RestMethod` to avoid PowerShell argument-escaping problems when sending JSON.

### A. Create user

```powershell
$body = @{
    name  = "Amjad"
    email = "amjad@example.com"
} | ConvertTo-Json

Invoke-RestMethod `
    -Uri "http://localhost:8080/users" `
    -Method POST `
    -ContentType "application/json" `
    -Body $body
```

### B. Get all users

```powershell
Invoke-RestMethod `
    -Uri "http://localhost:8080/users" `
    -Method GET
```

### C. Get user by ID

```powershell
Invoke-RestMethod `
    -Uri "http://localhost:8080/users/101" `
    -Method GET
```

### D. Update user

```powershell
$body = @{
    name  = "Amjad Ali"
    email = "amjad.ali@example.com"
} | ConvertTo-Json -Compress

Invoke-RestMethod `
    -Uri "http://localhost:8080/users/101" `
    -Method PUT `
    -ContentType "application/json" `
    -Body $body
```

### E. Delete user

```powershell
Invoke-WebRequest `
    -Uri "http://localhost:8080/users/101" `
    -Method DELETE
```

### curl.exe JSON formatting issue

In PowerShell, passing JSON through `curl.exe` can cause quoting and escaping problems.

For example, we encountered:

```text
invalid character 'n' looking for beginning of object key string
```

This means the JSON received by Go was malformed.

Using `ConvertTo-Json` with `Invoke-RestMethod` is a convenient alternative for PowerShell.

## 12. HTTP Testing with httptest

Go provides the `net/http/httptest` package.

It allows us to test HTTP handlers without starting a real HTTP server.

Important functions:

| Function | Purpose |
|---|---|
| `httptest.NewRequest()` | Creates an HTTP test request |
| `httptest.NewRecorder()` | Records the HTTP response |
| `recorder.Code` | Provides the HTTP status code |
| `recorder.Body` | Provides the response body |
| `mux.ServeHTTP()` | Runs the handler using a test request |

### Example

```go
func TestCreateUser(t *testing.T) {
    req := httptest.NewRequest(
        http.MethodPost,
        "/users",
        strings.NewReader(`{
            "name": "Amjad",
            "email": "amjad@example.com"
        }`),
    )

    req.Header.Set("Content-Type", "application/json")

    recorder := httptest.NewRecorder()

    mux.ServeHTTP(recorder, req)

    if recorder.Code != http.StatusCreated {
        t.Fatalf(
            "expected status %d, got %d",
            http.StatusCreated,
            recorder.Code,
        )
    }
}
```

Run tests:

```powershell
go test -v
```

An HTTP test that uses an in-memory SQLite database and exercises the handler, service and repository together is an integration-style test.

It is different from an isolated handler unit test that uses a fake service.

## 13. Important Database Concepts from Day 16

### LastInsertId()

Retrieves the ID generated by the last successful insert on the connection.

### RowsAffected()

Returns the number of rows affected by an SQL operation.

Useful for checking whether an UPDATE or DELETE actually matched a record.

### UNIQUE Constraint

Prevents duplicate values in a database column.

Example:

```sql
email TEXT NOT NULL UNIQUE
```

The database constraint protects against duplicate emails even when concurrent requests execute.

### Parameterized Queries

Always pass user-supplied values separately.

Good:

```go
db.ExecContext(ctx,
    "DELETE FROM users WHERE id = ?",
    id,
)
```

Avoid constructing SQL through string concatenation.

## 14. Day 16 Interview Questions

**Q1. What is CRUD?**

CRUD stands for Create, Read, Update and Delete.

**Q2. Difference between POST and PUT?**

POST commonly creates a resource. PUT generally replaces an existing resource representation.

**Q3. Why do we use LastInsertId()?**

To retrieve the automatically generated ID of an inserted record.

**Q4. Why do we check RowsAffected()?**

To determine how many records were modified and detect operations that matched no rows.

**Q5. What is the difference between HTTP 400 and HTTP 409?**

400 indicates a bad request or invalid input. 409 indicates a conflict with the current resource state, such as a duplicate email.

**Q6. Why add a UNIQUE constraint on email?**

To enforce uniqueness at the database level, including during concurrent requests.

**Q7. What is the purpose of httptest.NewRecorder()?**

It captures the HTTP response produced by a handler so we can inspect its status, headers and body.

**Q8. What happens if we delete a user that does not exist?**

The DELETE query can execute without an SQL error but affect zero rows. We should return a user-not-found error and HTTP 404.

**Q9. Difference between PUT and PATCH?**

PUT generally replaces the complete resource representation, while PATCH modifies selected fields.

**Q10. Why should business validation be in the service layer?**

It keeps business rules separate from HTTP-specific handling and database-specific operations.

## 15. Day 16 Practical Assignments

We practised these operations:

1. POST /users – Create a user.
2. GET /users – Retrieve all users.
3. GET /users/{id} – Retrieve a user by ID.
4. PUT /users/{id} – Update a user.
5. DELETE /users/{id} – Delete a user.

Additional practice:

- Handle duplicate email addresses.
- Handle invalid JSON.
- Handle nonexistent users.
- Validate names and email addresses.
- Write HTTP tests using `httptest`.
- Verify database records after every CRUD operation.

## 16. Day 16 Key Takeaways

1. CRUD means Create, Read, Update and Delete.
2. POST is commonly used to create resources.
3. GET retrieves resources.
4. PUT generally replaces a resource representation.
5. PATCH updates selected fields.
6. DELETE removes resources.
7. `json.NewDecoder()` parses JSON request bodies.
8. `LastInsertId()` retrieves a generated database ID.
9. `RowsAffected()` helps detect updates or deletes that matched no records.
10. Database UNIQUE constraints prevent duplicate values.
11. Use `errors.Is()` to identify meaningful application errors.
12. Use HTTP status codes appropriately.
13. Use `httptest` to test HTTP handlers.
14. Maintain Handler → Service → Repository separation.

**Day 16 Completed: Complete CRUD APIs in Go.**