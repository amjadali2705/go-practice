# Go Day 12 — Packages, Modules & Project Structure

## 1. What is a Package?

A **package** is a collection of Go source files in the same directory that are compiled together as a unit.

Example:

```text
calculator/
├── add.go
├── subtract.go
└── multiply.go
```

`add.go`:

```go
package calculator

func Add(a, b int) int {
    return a + b
}
```

`subtract.go`:

```go
package calculator

func Subtract(a, b int) int {
    return a - b
}
```

Both files belong to the `calculator` package.

### Important rules

- Ordinary Go files in the same directory generally declare the same package name.
- Different packages should normally be placed in different directories.
- Test files can use the package's name or an external test package name such as `calculator_test`.

## 2. What is a Module?

A **module** is a collection of Go packages that are versioned together.

Example:

```text
myapp/
├── go.mod
├── main.go
├── calculator/
│   └── calculator.go
└── user/
    └── user.go
```

Conceptually:

```text
Module: myapp
├── Package: main
├── Package: calculator
└── Package: user
```

Remember:

- **Package** → organizes Go source code.
- **Module** → groups packages and manages their dependencies and versions.

## 3. Creating a Go Module

Run:

```powershell
mkdir day12
cd day12
go mod init example.com/day12
```

This creates a `go.mod` file.

Example:

```go
module example.com/day12

go 1.23
```

The module path is `example.com/day12`. The Go version shown is illustrative; use the version appropriate for your environment.

Useful commands:

```powershell
go run .
go build .
go test ./...
```

- `go run .` runs the application.
- `go build .` builds the current package.
- `go test ./...` tests packages under the current module.

## 4. Understanding Import Paths

Suppose your project looks like this:

```text
day12/
├── go.mod
├── main.go
└── calculator/
    └── calculator.go
```

Your `go.mod` contains:

```go
module example.com/day12
```

The package file `calculator/calculator.go` contains:

```go
package calculator

func Add(a, b int) int {
    return a + b
}
```

You can import it from `main.go`:

```go
package main

import (
    "fmt"

    "example.com/day12/calculator"
)

func main() {
    res := calculator.Add(10, 20)
    fmt.Println("Result:", res)
}
```

Output:

```text
Result: 30
```

### How is the import path formed?

```text
Module path:
example.com/day12

        +

Package directory:
calculator

        ↓

Import path:
example.com/day12/calculator
```

The import path is based on the module path in `go.mod` plus the package's directory path.

The package name declared in the source file determines the default local name used to access it.

## 5. Exported vs Unexported Identifiers

Go uses capitalization to control whether identifiers are exported.

### Exported identifiers

Names beginning with an uppercase letter can be accessed from other packages.

```go
package calculator

func Add(a, b int) int {
    return a + b
}
```

Other packages can call:

```go
calculator.Add(10, 20)
```

### Unexported identifiers

Names beginning with a lowercase letter are accessible only within the same package.

```go
func add(a, b int) int {
    return a + b
}
```

Another package cannot call:

```go
calculator.add(10, 20) // Compile error
```

The rule also applies to types, variables, constants, methods, and struct fields.

Example:

```go
type User struct {
    ID    int
    name  string
    Email string
}
```

Another package can access `ID` and `Email`, but not `name`.

**Remember:**

```text
Uppercase → exported
lowercase → unexported
```

Exported identifiers form part of a package's public API.

## 6. Package Name vs Module Name

The package name does not have to match the final part of the module path.

For example:

```go
module example.com/myapp
```

Suppose the directory is:

```text
myapp/
└── calculator/
    └── calculator.go
```

The source file could declare:

```go
package calc
```

You still import it using the directory-based path:

```go
import "example.com/myapp/calculator"
```

The default local name is `calc`, so you would use:

```go
calc.Add(1, 2)
```

You can also specify an alias explicitly:

```go
import calc "example.com/myapp/calculator"
```

Avoid naming your own package or module after a standard-library package when it could cause confusion or import-resolution conflicts. Use distinct, meaningful names.

## 7. Understanding `go.mod` and `go.sum`

### `go.mod`

The `go.mod` file records information such as:

- Module path
- Go version
- Dependency requirements

Example:

```go
module example.com/myapp

go 1.23
```

### `go.sum`

The `go.sum` file contains cryptographic checksums used to verify module content and module metadata.

| File | Purpose |
|---|---|
| `go.mod` | Module identity, Go version, dependency requirements |
| `go.sum` | Dependency checksums |

Both files should normally be committed to Git so developers and CI environments can resolve and verify dependencies consistently.

## 8. Important Go Module Commands

Initialize a module:

```powershell
go mod init example.com/myapp
```

Clean up dependency requirements and checksums:

```powershell
go mod tidy
```

Download required modules into the module cache:

```powershell
go mod download
```

Display the dependency graph:

```powershell
go mod graph
```

List module dependencies:

```powershell
go list -m all
```

Run all package tests:

```powershell
go test ./...
```

Build all packages:

```powershell
go build ./...
```

### What does `go mod tidy` do?

It adds missing module requirements needed by the code and removes requirements that are no longer needed. It also updates `go.sum` as necessary.

## 9. The `internal` Directory

Go supports a special directory named `internal`.

Example:

```text
myapp/
├── go.mod
├── main.go
└── internal/
    └── auth/
        └── auth.go
```

The `internal/auth` package can be imported by code within its permitted parent directory tree, but unrelated modules cannot import it directly.

This is useful for application-specific implementation details, such as:

- Authentication logic
- Internal business rules
- Private service components
- Application-specific helpers

A possible backend layout:

```text
myapp/
├── go.mod
├── cmd/
│   └── api/
│       └── main.go
└── internal/
    ├── user/
    ├── transaction/
    └── config/
```

This is a convention, not a mandatory structure. Introduce packages when they provide clear boundaries.

## 10. Import Cycles

Go does not allow cyclic package imports.

Problematic dependency design:

```text
Handler → Service → Repository
    ↑                       |
    └───────────────────────┘
```

Here, the repository eventually imports the handler, creating a cycle.

Instead, keep dependencies flowing in one direction:

```text
Handler
   ↓
Service
   ↓
Repository
   ↓
Database
```

Responsibilities:

- **Handler** → receives requests and writes responses.
- **Service** → implements business logic.
- **Repository** → handles persistence.

A repository should return data or errors rather than importing a handler to send an HTTP response.

If packages need a shared type, consider placing that type in an appropriate common package or redesigning the dependency. Avoid creating a large shared package just to hide an import cycle.

## 11. The `init()` Function

Go has a special initialization function:

```go
package main

import "fmt"

func init() {
    fmt.Println("init")
}

func main() {
    fmt.Println("main")
}
```

Output:

```text
init
main
```

Important properties:

- `init()` has no parameters and no return values.
- Go runs it automatically during package initialization.
- A package can have multiple `init()` functions across its files.
- Imported packages are initialized before the packages that import them.
- `init()` is not a substitute for explicit dependency setup.

Avoid hiding complicated business logic inside `init()`. Explicit setup functions are generally easier to test.

## 12. Blank Imports

A blank import uses `_`:

```go
import _ "example.com/some/driver"
```

This imports the package for its initialization side effects without making its exported identifiers available under a local package name.

A common use is registering a database driver with another library.

Use blank imports only when those side effects are intentional.

## 13. A Practical Go Backend Structure

One possible layout for a small-to-medium API:

```text
payment-api/
├── go.mod
├── go.sum
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── handler/
│   │   └── user_handler.go
│   ├── service/
│   │   └── user_service.go
│   ├── repository/
│   │   └── user_repository.go
│   ├── model/
│   │   └── user.go
│   └── config/
│       └── config.go
└── migrations/
```

The typical request flow is:

```text
HTTP Request
     ↓
Handler
     ↓
Service
     ↓
Repository
     ↓
Database
```

Keep each package focused on a clear responsibility. For small applications, a simpler layout may be better than introducing many packages immediately.

## 14. Day 12 Coding Exercise — Completed

You created a module and a separate calculator package.

`calculator/calculator.go`:

```go
package calculator

func Add(a, b int) int {
    return a + b
}
```

`main.go`:

```go
package main

import (
    "fmt"

    "example.com/day12/calculator"
)

func main() {
    res := calculator.Add(10, 20)
    fmt.Println("Result:", res)
}
```

Output:

```text
Result: 30
```

You successfully practised:

- Creating a module
- Declaring a package
- Importing your own package
- Calling an exported function
- Running, testing, and building a Go project

## 15. Day 12 Interview Questions

**1. What is the difference between a package and a module?**

A package groups Go source files, while a module groups packages and manages their versions and dependencies.

**2. What makes an identifier exported?**

Its name begins with an uppercase letter.

**3. How is an import path formed?**

It is generally the module path plus the package's directory path.

**4. What is the purpose of `go.mod`?**

It records module information, the Go version, and dependency requirements.

**5. What is the purpose of `go.sum`?**

It stores checksums used to verify module content and metadata.

**6. What is an import cycle?**

A circular dependency between packages, which Go rejects.

**7. What is the `internal` directory used for?**

Restricting package imports to code within the permitted parent directory tree.

**8. When does `init()` run?**

During package initialization, before the importing program reaches `main()`.

**9. What is a blank import?**

An import using `_`, commonly used when a package must perform initialization side effects.

**10. Why separate handlers, services, and repositories?**

To keep HTTP handling, business logic, and database operations independent and easier to test and maintain.

## 16. Day 12 Golden Rules

```text
Package       → groups source files
Module        → groups and versions packages
go.mod        → module and dependency requirements
go.sum        → dependency checksums
Uppercase     → exported identifier
lowercase     → unexported identifier
internal/     → restricted import boundary
init()        → automatic package initialization
Import cycle  → not allowed
```

For backend architecture:

```text
Handler
   ↓
Service
   ↓
Repository
   ↓
Database
```

The most important takeaway from Day 12 is that a good Go project is not just about writing correct functions. It is also about **organizing packages, controlling dependencies, and maintaining clear boundaries**.
