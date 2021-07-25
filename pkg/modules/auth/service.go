package auth

import (
	"errors"
	"fmt"
	"github.com/evleria/jwt-auth-demo/pkg/common/bcrypt"
)

type Service interface {
	Register(firstName, lastName, email, password string) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Register(firstName, lastName, email, password string) error {
	hash, err := bcrypt.Hash(password)
	if err != nil {
		return fmt.Errorf("cannot register: %v", err)
	}

	err = s.repository.CreateNewUser(firstName, lastName, email, hash)
	if err != nil {
		return errors.New("cannot register a new user")
	}
	return nil
}
