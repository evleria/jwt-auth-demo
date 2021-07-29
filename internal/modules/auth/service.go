package auth

import (
	"errors"
	"fmt"
	"github.com/evleria/jwt-auth-demo/internal/common/bcrypt"
	"github.com/evleria/jwt-auth-demo/internal/common/jwt"
	"time"
)

type Service interface {
	Register(firstName, lastName, email, password string) error
	Login(email, password string) (string, string, error)
	Refresh(refreshToken string) (string, error)
}

type service struct {
	userRepository  UserRepository
	tokenRepository TokenRepository
	jwtMaker        jwt.Maker
}

func NewService(userRepository UserRepository, tokenRepository TokenRepository, jwtMaker jwt.Maker) *service {
	return &service{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		jwtMaker:        jwtMaker,
	}
}

func (s *service) Register(firstName, lastName, email, password string) error {
	hash, err := bcrypt.Hash(password)
	if err != nil {
		return fmt.Errorf("cannot register: %v", err)
	}

	err = s.userRepository.CreateNewUser(firstName, lastName, email, hash)
	if err != nil {
		return errors.New("cannot register a new user")
	}
	return nil
}

func (s *service) Login(email, password string) (string, string, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return "", "", errors.New("cannot find user")
	}

	if !bcrypt.Compare(user.PassHash, password) {
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

func (s *service) Refresh(refreshToken string) (string, error) {
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
