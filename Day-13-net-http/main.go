package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var users = map[string]User{
	"1": {
		ID:     "1",
		Name:   "Amjad",
		Active: true,
	},
	"2": {
		ID:     "2",
		Name:   "Rahul",
		Active: true,
	},
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /users/{id}", getUser)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("Server starting on :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server error:", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user, ok := users[id]
	if !ok {
		writeJSON(w, http.StatusNotFound, ErrorResponse{
			Error: "user not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Failed to encode response:", err)
	}
}
