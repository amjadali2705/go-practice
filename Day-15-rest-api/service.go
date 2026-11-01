package main

import (
	"context"
	"errors"
)

var ErrInvalidUserID = errors.New("invalid user ID")

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
