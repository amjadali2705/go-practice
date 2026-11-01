package main

import (
	"errors"
	"log"
	"net/http"
	"time"
)

func main() {
	// 1. Database connection
	db, err := ConnectDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	// 2. Create users table
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY,
            name TEXT NOT NULL,
            email TEXT NOT NULL UNIQUE
        )
    `

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	// Insert sample user
	_, err = db.Exec(`
        INSERT OR IGNORE INTO users (id, name, email)
        VALUES (101, 'Amjad', 'amjad@example.com')
    `)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Initialize repository
	userRepo := NewSQLiteUserRepository(db)

	// 4. Initialize service
	userService := NewUserService(userRepo)

	// 5. Initialize handler
	userHandler := NewUserHandler(userService)

	// 6. Register route
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", userHandler.GetUserByID)
	mux.HandleFunc("GET /users", userHandler.GetAll)
	// 7. Start HTTP server
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Server running on port 8080")

	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
