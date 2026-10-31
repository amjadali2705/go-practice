# Day 14 – Database Access in Go

## 1. Introduction to Database Access

Go provides the `database/sql` package for interacting with relational databases.

It supports common database operations such as:
- Connecting to a database.
- Inserting records.
- Fetching records.
- Updating records.
- Deleting records.
- Managing transactions.
- Handling database errors.

The `database/sql` package provides a common interface. A database driver handles communication with the actual database.

Examples:
- SQLite: `modernc.org/sqlite`
- MySQL/MariaDB: `github.com/go-sql-driver/mysql`
- PostgreSQL: `github.com/jackc/pgx`

## 2. Database Connection in Go

Import the required packages:

```go
import (
    "database/sql"
    _ "modernc.org/sqlite"
)
```

The underscore (`_`) is a blank import. It loads the driver so that it can register itself with `database/sql`, without directly using its package name.

### Opening a database

```go
db, err := sql.Open("sqlite", "go_practice.db")
if err != nil {
    return err
}
defer db.Close()
```

Important:

`sql.Open()` does not necessarily establish a connection immediately. It creates a database handle that manages a connection pool.

### Verify database connectivity

```go
if err := db.PingContext(ctx); err != nil {
    return err
}
```

`PingContext()` verifies that the database is reachable.

### Important distinction

| Function | Purpose |
|---|---|
| `sql.Open()` | Creates a database handle |
| `db.PingContext()` | Checks database connectivity |
| `db.Close()` | Closes the database handle and its pool |

## 3. What is *sql.DB?

`*sql.DB` represents a database connection pool, not one individual database connection.

It manages connections and makes them available to application requests.

Example:

```go
func GetUser(db *sql.DB, id int64) {
    // Execute database operations
}
```

You should create the database handle once when your application starts and reuse it across handlers and repositories.

Do not create a new database handle for every HTTP request.

### Connection pool configuration

```go
db.SetMaxOpenConns(10)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(30 * time.Minute)
db.SetConnMaxIdleTime(5 * time.Minute)
```

| Function | Meaning |
|---|---|
| `SetMaxOpenConns()` | Maximum number of open connections |
| `SetMaxIdleConns()` | Maximum number of idle connections |
| `SetConnMaxLifetime()` | Maximum lifetime of a connection |
| `SetConnMaxIdleTime()` | Maximum idle time of a connection |

Too few connections can cause requests to wait. Too many connections can overload the database.

For our SQLite exercise, we used one open connection to keep the setup simple. Database pool settings should be chosen according to the database engine and workload.

## 4. Executing SQL Queries

Go provides different functions for different SQL operations.

### A. ExecContext()

Use `ExecContext()` for operations that do not return rows.

Examples:
- INSERT
- UPDATE
- DELETE
- CREATE TABLE

```go
query := `INSERT INTO users (name, email)
          VALUES (?, ?)`

_, err := db.ExecContext(ctx, query, "Amjad", "amjad@example.com")
if err != nil {
    return err
}
```

### B. QueryContext()

Use `QueryContext()` when a query can return multiple rows.

```go
rows, err := db.QueryContext(ctx,
    "SELECT id, name FROM users")
if err != nil {
    return err
}
defer rows.Close()

for rows.Next() {
    var id int64
    var name string

    if err := rows.Scan(&id, &name); err != nil {
        return err
    }

    fmt.Println(id, name)
}

if err := rows.Err(); err != nil {
    return err
}
```

Important:
- `rows.Next()` moves to the next row.
- `rows.Scan()` copies column values into Go variables.
- `rows.Close()` releases row-related resources.
- `rows.Err()` checks for errors encountered during iteration.

### C. QueryRowContext()

Use `QueryRowContext()` when you expect one row.

```go
var name string

err := db.QueryRowContext(ctx,
    "SELECT name FROM users WHERE id = ?",
    101,
).Scan(&name)

if errors.Is(err, sql.ErrNoRows) {
    return errors.New("user not found")
}

if err != nil {
    return err
}
```

`QueryRowContext()` returns a `*sql.Row`. The query error is generally reported by calling `Scan()`.

### Quick comparison

| Function | Use case |
|---|---|
| `ExecContext()` | INSERT, UPDATE, DELETE |
| `QueryContext()` | Multiple rows |
| `QueryRowContext()` | Single row |

## 5. Parameterized Queries and SQL Injection

Never construct SQL queries by directly concatenating user input.

Unsafe example:

```go
query := "SELECT * FROM users WHERE name = '" + name + "'"
```

This can expose your application to SQL injection.

Safe example:

```go
query := "SELECT * FROM users WHERE name = ?"

row := db.QueryRowContext(ctx, query, name)
```

Parameterized queries separate SQL instructions from user-provided values.

Note: `?` is used by SQLite and MySQL drivers. PostgreSQL commonly uses placeholders such as `$1`, `$2`.

## 6. Handling Database Results

### RowsAffected()

`RowsAffected()` tells us how many rows were affected by an INSERT, UPDATE, or DELETE operation.

```go
result, err := db.ExecContext(ctx,
    "UPDATE accounts SET balance_paise = ? WHERE id = ?",
    50000, 101,
)
if err != nil {
    return err
}

affected, err := result.RowsAffected()
if err != nil {
    return err
}

fmt.Println("Rows affected:", affected)
```

Why is this important?

A query can execute without an SQL error and still affect zero rows.

For example, updating an account that does not exist might affect zero rows.

### sql.ErrNoRows

When `QueryRowContext()` does not find a matching record, `Scan()` returns `sql.ErrNoRows`.

```go
if errors.Is(err, sql.ErrNoRows) {
    return errors.New("record not found")
}
```

Do not treat every database error as a missing record. Handle unexpected errors separately.

## 7. Context in Database Operations

Context provides cancellation and deadline propagation.

It is useful when database operations take too long or an HTTP client disconnects.

Example:

```go
ctx, cancel := context.WithTimeout(
    context.Background(),
    5*time.Second,
)
defer cancel()

err := db.PingContext(ctx)
if err != nil {
    return err
}
```

Important database functions with context support:

- `db.PingContext()`
- `db.ExecContext()`
- `db.QueryContext()`
- `db.QueryRowContext()`
- `db.BeginTx()`
- `tx.ExecContext()`

### Why pass context to database operations?

Suppose a client sends an API request and disconnects before the database query finishes.

Passing the request context allows the cancellation signal to reach the database operation.

The driver and database operation must support cancellation for it to take effect.

For a backend application, the recommended flow is:

Handler → Service → Repository → Database

Pass the same request context through these layers rather than creating unrelated background contexts inside them.

## 8. Database Transactions

A transaction groups multiple database operations into one logical unit.

A transaction helps ensure that related changes succeed or fail together.

Example: Money transfer.

Suppose:

| Account | Balance |
|---|---:|
| 101 | ₹1,000 |
| 102 | ₹500 |

We want to transfer ₹100 from account 101 to account 102.

Expected result:

| Account | New Balance |
|---|---:|
| 101 | ₹900 |
| 102 | ₹600 |

Both operations must succeed.

If the debit succeeds but the credit fails, the debit must be rolled back.

### Transaction functions

| Function | Purpose |
|---|---|
| `db.BeginTx()` | Starts a transaction |
| `tx.ExecContext()` | Executes a query inside the transaction |
| `tx.Commit()` | Commits the transaction |
| `tx.Rollback()` | Rolls back uncommitted changes |

### Example

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil {
    return err
}

defer tx.Rollback()

_, err = tx.ExecContext(ctx, debitQuery, amount, fromID)
if err != nil {
    return err
}

_, err = tx.ExecContext(ctx, creditQuery, amount, toID)
if err != nil {
    return err
}

return tx.Commit()
```

Important:

- Use `tx.ExecContext()` for all operations that belong to the transaction.
- Do not accidentally use `db.ExecContext()` for an operation that must be part of that transaction.
- Use `defer tx.Rollback()` as a safety net when the function returns early.
- After a successful commit, the deferred rollback has no effect on the committed transaction.

## 9. ACID Properties

Transactions are commonly described using ACID.

| Property | Meaning |
|---|---|
| Atomicity | All operations succeed together or are rolled back |
| Consistency | Data remains valid according to its constraints and rules |
| Isolation | Concurrent transactions are isolated according to the database's isolation level |
| Durability | Committed changes are preserved according to the database's durability guarantees |

These properties are important for financial applications, including payment and banking systems.

## 10. Practical Exercise: Money Transfer

We created the following table:

```sql
CREATE TABLE accounts (
    id INTEGER PRIMARY KEY,
    balance_paise INTEGER NOT NULL
);
```

Initial records:

```sql
INSERT INTO accounts VALUES
(101, 100000),
(102, 50000);
```

We implemented:

```go
func TransferMoney(
    ctx context.Context,
    db *sql.DB,
    fromID int64,
    toID int64,
    amount int64,
) error
```

### Business rules

1. Reject zero or negative transfer amounts.
2. Start a database transaction.
3. Debit the sender only when sufficient balance exists.
4. Check the number of affected rows.
5. Credit the receiver.
6. Verify that the receiver account exists.
7. Commit only after both operations succeed.
8. Roll back on errors.

### Important SQL queries

Debit:

```sql
UPDATE accounts
SET balance_paise = balance_paise - ?
WHERE id = ? AND balance_paise >= ?
```

Credit:

```sql
UPDATE accounts
SET balance_paise = balance_paise + ?
WHERE id = ?
```

The balance condition prevents the sender from going into a negative balance through this operation.

### Important error cases

```go
var (
    ErrInvalidAmount     = errors.New("invalid transfer amount")
    ErrInsufficientFunds = errors.New("insufficient funds")
    ErrAccountNotFound   = errors.New("account not found")
)
```

We use these errors to represent meaningful business failures rather than exposing raw database errors to API clients.

For production financial systems, additional protection is necessary, such as concurrency-safe balance updates, appropriate database constraints, overflow checks, transfer validation, idempotency, and suitable transaction isolation.

## 11. SQLite vs MariaDB

For our practice, we used SQLite because it does not require a separate database server.

| Feature | SQLite | MariaDB |
|---|---|---|
| Installation | Go driver is sufficient | Requires database server |
| Database storage | Local file | Database server storage |
| Go package | `modernc.org/sqlite` | `go-sql-driver/mysql` |
| `database/sql` | Supported | Supported |
| Transactions | Supported | Supported |
| Connection pooling | Supported by Go API | Supported by Go API |

Most `database/sql` concepts transfer between these databases, but SQL syntax, locking, concurrency, and driver behavior can differ.

## 12. Handler → Service → Repository Architecture

A typical Go backend separates responsibilities.

**Handler**
- Receives HTTP requests.
- Validates basic request structure.
- Calls the service.
- Converts errors into HTTP responses.

**Service**
- Implements business logic.
- Validates business rules.
- Coordinates transactions when appropriate.

**Repository**
- Executes database queries.
- Maps database results to Go models.
- Handles database-specific errors.

**Database**
- Stores application data.

Example flow:

```text
HTTP Request
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
Database
```

For a missing user, a repository can translate `sql.ErrNoRows` into a domain error such as `ErrUserNotFound`. The handler can then use `errors.Is()` to return HTTP 404.

## 13. Important Interview Questions

**Q1. What is `database/sql`?**

It is Go's standard library package for accessing relational databases through database drivers.

**Q2. What is `*sql.DB`?**

It represents a database handle that manages a connection pool.

**Q3. Does `sql.Open()` immediately connect to the database?**

Not necessarily. Use `PingContext()` to verify connectivity.

**Q4. Difference between `ExecContext()` and `QueryContext()`?**

`ExecContext()` executes commands that do not return rows. `QueryContext()` returns rows.

**Q5. What is `sql.ErrNoRows`?**

It indicates that a `QueryRowContext()` operation did not find a matching row, reported when `Scan()` is called.

**Q6. Why use parameterized queries?**

To separate SQL instructions from supplied values and prevent SQL injection.

**Q7. Why use transactions?**

To ensure related operations are committed together or rolled back together.

**Q8. Difference between Commit and Rollback?**

Commit finalizes changes; Rollback cancels uncommitted changes.

**Q9. What does RowsAffected() do?**

It reports how many rows were affected by an SQL operation.

**Q10. Why use context in database operations?**

To propagate cancellation and deadlines to database operations.

**Q11. Why use defer rows.Close()?**

To release resources associated with the returned rows when the function exits.

**Q12. Why check rows.Err()?**

To detect errors that occurred during row iteration.

---

## Day 14 – Key Takeaways

1. `database/sql` is Go's standard database interface.
2. `*sql.DB` manages a connection pool.
3. `sql.Open()` creates the database handle; `PingContext()` verifies connectivity.
4. Use `ExecContext()` for commands that do not return rows.
5. Use `QueryContext()` for multiple rows.
6. Use `QueryRowContext()` for one row.
7. Always handle database errors.
8. Use parameterized queries to prevent SQL injection.
9. Use `RowsAffected()` when business logic depends on affected records.
10. Pass context through Handler → Service → Repository.
11. Use transactions for operations that must succeed or fail together.
12. Use `defer rows.Close()` and check `rows.Err()` when iterating through rows.
13. Keep database connections shared and properly configured.