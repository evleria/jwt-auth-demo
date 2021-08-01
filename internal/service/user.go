package service

import (
	"context"
	"github.com/evleria/jwt-auth-demo/internal/repository"
	"github.com/evleria/jwt-auth-demo/internal/repository/entities"
)

type User interface {
	Me(ctx context.Context, userId int) (*entities.User, error)
}

type user struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) User {
	return &user{
		userRepository: userRepository,
	}
}

func (u *user) Me(ctx context.Context, userId int) (*entities.User, error) {
	return u.userRepository.GetUserById(ctx, userId)
}
