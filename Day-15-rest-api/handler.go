package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *UserHandler) GetUserByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": "invalid user ID"})
		return
	}

	user, err := h.service.GetUserByID(r.Context(), id)

	if errors.Is(err, ErrInvalidUserID) {
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": err.Error()})
		return
	}

	if errors.Is(err, ErrUserNotFound) {
		writeJSON(w, http.StatusNotFound,
			map[string]string{"error": err.Error()})
		return
	}

	if err != nil {
		log.Printf("Get user failed: %v", err)

		writeJSON(w, http.StatusInternalServerError,
			map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetUserByUserName(w http.ResponseWriter, r *http.Request, name string) {

	user, err := h.service.GetUserByUserName(r.Context(), name)
	if errors.Is(err, ErrUserNotFound) {
		writeJSON(w, http.StatusNotFound,
			map[string]string{"error": err.Error()})
		return
	}

	if err != nil {
		log.Printf("Get user failed: %v", err)

		writeJSON(w, http.StatusInternalServerError,
			map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	if name != "" {
		h.GetUserByUserName(w, r, name)
		return
	}

	users, err := h.service.GetAll(r.Context())
	if err != nil {
		log.Printf("Get all users failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, users)
}

func (h *UserHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": "invalid request body"})
		return
	}

	user := User{
		Name:  req.Name,
		Email: req.Email,
	}

	createdUser, err := h.service.CreateUser(r.Context(), user)

	if errors.Is(err, ErrInvalidUserName) ||
		errors.Is(err, ErrInvalidEmail) {
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": err.Error()})
		return
	}

	if errors.Is(err, ErrEmailAlreadyExists) {
		writeJSON(w, http.StatusConflict,
			map[string]string{"error": err.Error()})
		return
	}

	if err != nil {
		log.Printf("Create user failed: %v", err)
		writeJSON(w, http.StatusInternalServerError,
			map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, createdUser)
}

func (h *UserHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": "invalid user ID"})
		return
	}

	var req UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": "invalid request body"})
		return
	}

	user := User{
		Name:  req.Name,
		Email: req.Email,
	}

	err = h.service.UpdateUser(r.Context(), id, user)

	if errors.Is(err, ErrUserNotFound) {
		writeJSON(w, http.StatusNotFound,
			map[string]string{"error": err.Error()})
		return
	}

	if errors.Is(err, ErrInvalidUserID) ||
		errors.Is(err, ErrInvalidUserName) ||
		errors.Is(err, ErrInvalidEmail) {
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": err.Error()})
		return
	}

	if errors.Is(err, ErrEmailAlreadyExists) {
		writeJSON(w, http.StatusConflict,
			map[string]string{"error": err.Error()})
		return
	}

	if err != nil {
		log.Printf("Update user failed: %v", err)
		writeJSON(w, http.StatusInternalServerError,
			map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK,
		map[string]string{"message": "user updated successfully"})
}

func (h *UserHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": "invalid user ID"})
		return
	}

	err = h.service.DeleteUser(r.Context(), id)

	if errors.Is(err, ErrUserNotFound) {
		writeJSON(w, http.StatusNotFound,
			map[string]string{"error": err.Error()})
		return
	}

	if err != nil {
		log.Printf("Delete user failed: %v", err)
		writeJSON(w, http.StatusInternalServerError,
			map[string]string{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
