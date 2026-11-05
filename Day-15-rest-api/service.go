package main

import (
	"context"
	"errors"
	"strings"
)

var ErrInvalidUserID = errors.New("invalid user ID")
var (
	ErrInvalidUserName    = errors.New("invalid username")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(
	ctx context.Context,
	id int64,
) (User, error) {

	if id <= 0 {
		return User{}, ErrInvalidUserID
	}

	return s.repo.GetByID(ctx, id)
}

func (s *UserService) GetAll(ctx context.Context) ([]User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) GetUserByUserName(
	ctx context.Context,
	name string,
) (User, error) {
	if name == "" {
		return User{}, ErrInvalidUserID
	}
	return s.repo.GetByUserName(ctx, name)
}

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

	_, err := s.repo.GetByEmail(ctx, user.Email)

	if err == nil {
		return User{}, ErrEmailAlreadyExists
	}

	if !errors.Is(err, ErrUserNotFound) {
		return User{}, err
	}

	return s.repo.Create(ctx, user)
}

func (s *UserService) UpdateUser(
	ctx context.Context,
	id int64,
	user User,
) error {

	if id <= 0 {
		return ErrInvalidUserID
	}

	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	if user.Name == "" {
		return ErrInvalidUserName
	}

	if user.Email == "" {
		return ErrInvalidEmail
	}

	existing, err := s.repo.GetByEmail(ctx, user.Email)

	if err == nil && existing.ID != id {
		return ErrEmailAlreadyExists
	}

	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return err
	}

	return s.repo.Update(ctx, id, user)
}

func (s *UserService) DeleteUser(
	ctx context.Context,
	id int64,
) error {

	if id <= 0 {
		return ErrInvalidUserID
	}

	return s.repo.Delete(ctx, id)
}
