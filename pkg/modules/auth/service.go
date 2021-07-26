package auth

import (
	"errors"
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/evleria/jwt-auth-demo/pkg/common/bcrypt"
	"github.com/evleria/jwt-auth-demo/pkg/common/jwt"
)

type Service interface {
	Register(firstName, lastName, email, password string) error
	Login(email, password string) (string, string, error)
}

type service struct {
	repository Repository
	jwtMaker   jwt.Maker
}

func NewService(repository Repository, jwtMaker jwt.Maker) Service {
	return &service{
		repository: repository,
		jwtMaker:   jwtMaker,
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

func (s *service) Login(email, password string) (string, string, error) {
	user, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return "", "", errors.New("cannot find user")
	}

	if !bcrypt.Compare(user.PassHash, password) {
		return "", "", errors.New("invalid password provided")
	}

	accessToken, err := s.jwtMaker.GenerateAccessToken(jwtgo.MapClaims{
		"sub":   user.Id,
		"email": user.Email,
	})
	if err != nil {
		return "", "", errors.New("cannot generate access token")
	}

	refreshToken, err := s.jwtMaker.GenerateRefreshToken(jwtgo.MapClaims{
		"sub": user.Id,
	})
	if err != nil {
		return "", "", errors.New("cannot generate refresh token")
	}

	return accessToken, refreshToken, nil
}
