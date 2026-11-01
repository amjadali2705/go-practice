package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (User, error)
	GetAll(ctx context.Context) ([]User, error)
	GetByUserName(ctx context.Context, name string) (User, error)
}

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

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

func (r *SQLiteUserRepository) GetAll(
	ctx context.Context,
) ([]User, error) {

	query := `
        SELECT id, name, email
        FROM users
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}
	defer rows.Close()

	// Initialize an empty slice.
	users := make([]User, 0)

	for rows.Next() {
		var user User

		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
		); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}

		users = append(users, user)
	}

	// Check errors encountered during iteration.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

func (r *SQLiteUserRepository) GetByUserName(
	ctx context.Context,
	name string,
) (User, error) {
	var user User

	query := `
        SELECT id, name, email
        FROM users
        WHERE name = ?
    `

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUserNotFound
	}

	if err != nil {
		return User{}, fmt.Errorf("get user by name: %w", err)
	}

	return user, nil
}
