package services

import (
	"errors"
	"fmt"
	"github.com/evleria/jwt-auth-demo/internal/common/jwt"
	"github.com/evleria/jwt-auth-demo/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService interface {
	Register(firstName, lastName, email, password string) error
	Login(email, password string) (string, string, error)
	Refresh(refreshToken string) (string, error)
}

type authService struct {
	userRepository  repositories.UserRepository
	tokenRepository repositories.TokenRepository
	jwtMaker        jwt.Maker
}

func NewAuthService(userRepository repositories.UserRepository, tokenRepository repositories.TokenRepository, jwtMaker jwt.Maker) *authService {
	return &authService{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		jwtMaker:        jwtMaker,
	}
}

func (s *authService) Register(firstName, lastName, email, password string) error {
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

func (s *authService) Login(email, password string) (string, string, error) {
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

func (s *authService) Refresh(refreshToken string) (string, error) {
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
