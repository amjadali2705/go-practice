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

	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, user User) (User, error)
	Update(ctx context.Context, id int64, user User) error
	Delete(ctx context.Context, id int64) error
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
		return User{}, fmt.Errorf("create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return User{}, err
	}

	user.ID = id

	return user, nil
}

func (r *SQLiteUserRepository) Update(
	ctx context.Context,
	id int64,
	user User,
) error {

	query := `
        UPDATE users
        SET name = ?, email = ?
        WHERE id = ?
    `

	result, err := r.db.ExecContext(
		ctx, query, user.Name, user.Email, id,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *SQLiteUserRepository) Delete(
	ctx context.Context,
	id int64,
) error {

	query := `DELETE FROM users WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *SQLiteUserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (User, error) {

	var user User

	query := `
        SELECT id, name, email
        FROM users
        WHERE email = ?
    `

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUserNotFound
	}

	if err != nil {
		return User{}, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}
