package service

import (
	"context"

	"github.com/evleria/jwt-auth-demo/internal/repository"
	"github.com/evleria/jwt-auth-demo/internal/repository/entities"
)

// User contains usecase logic for authentication
type User interface {
	Me(ctx context.Context, userID int) (*entities.User, error)
}

type user struct {
	userRepository repository.User
}

// NewUserService creates user service
func NewUserService(userRepository repository.User) User {
	return &user{
		userRepository: userRepository,
	}
}

func (u *user) Me(ctx context.Context, userID int) (*entities.User, error) {
	return u.userRepository.GetUserByID(ctx, userID)
}
