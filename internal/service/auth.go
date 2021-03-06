package service

import (
	"errors"
	"fmt"
	"github.com/evleria/jwt-auth-demo/internal/jwt"
	"github.com/evleria/jwt-auth-demo/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Auth interface {
	Register(firstName, lastName, email, password string) error
	Login(email, password string) (string, string, error)
	Refresh(refreshToken string) (string, error)
}

type auth struct {
	userRepository  repository.UserRepository
	tokenRepository repository.Token
	jwtMaker        jwt.Maker
}

func NewAuthService(userRepository repository.UserRepository, tokenRepository repository.Token, jwtMaker jwt.Maker) *auth {
	return &auth{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		jwtMaker:        jwtMaker,
	}
}

func (s *auth) Register(firstName, lastName, email, password string) error {
	var hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("cannot register: %v", err)
	}
	err = s.userRepository.CreateNewUser(firstName, lastName, email, string(hash))
	if err != nil {
		return errors.New("cannot register a new user")
	}
	return nil
}

func (s *auth) Login(email, password string) (string, string, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return "", "", errors.New("cannot find user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid password provided")
	}

	accessToken, err := s.jwtMaker.GenerateAccessToken(user.Id, user.Email)
	if err != nil {
		return "", "", errors.New("cannot generate access token")
	}

	refreshToken, err := s.jwtMaker.GenerateRefreshToken(user.Id)
	if err != nil {
		return "", "", errors.New("cannot generate refresh token")
	}

	return accessToken, refreshToken, nil
}

func (s *auth) Refresh(refreshToken string) (string, error) {
	claims, err := s.jwtMaker.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	userId := int(claims["sub"].(float64))
	t, inBlacklist, err := s.tokenRepository.IsBlacklisted(userId)
	if err != nil {
		return "", err
	}
	if inBlacklist {
		iat := time.Unix(int64(claims["iat"].(int)), 0)
		if t.After(iat) {
			return "", errors.New("token is blacklisted")
		}
	}

	user, err := s.userRepository.GetUserById(userId)
	if err != nil {
		return "", errors.New("cannot find user")
	}

	accessToken, err := s.jwtMaker.GenerateAccessToken(user.Id, user.Email)
	if err != nil {
		return "", errors.New("cannot generate access token")
	}
	return accessToken, err
}
