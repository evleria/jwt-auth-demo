// Package service encapsulates usecases
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/evleria/jwt-auth-demo/internal/config"
	"github.com/evleria/jwt-auth-demo/internal/jwt"
	"github.com/evleria/jwt-auth-demo/internal/repository"
)

// Auth contains usecase logic for authentication
type Auth interface {
	Register(ctx context.Context, firstName, lastName, email, password string) error
	Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, err error)
	Refresh(ctx context.Context, refreshToken string) (string, error)
	Logout(ctx context.Context, refreshToken string) error
	ValidateAccessToken(ctx context.Context, accessToken string) (*jwt.AccessTokenClaims, error)
}

type auth struct {
	userRepository  repository.User
	tokenRepository repository.Token
	jwtMaker        jwt.Maker
}

// NewAuthService creates auth service
func NewAuthService(userRepository repository.User, tokenRepository repository.Token, jwtMaker jwt.Maker) Auth {
	return &auth{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		jwtMaker:        jwtMaker,
	}
}

func (s *auth) Register(ctx context.Context, firstName, lastName, email, password string) error {
	var hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("cannot register: %v", err)
	}
	err = s.userRepository.CreateNewUser(ctx, firstName, lastName, email, string(hash))
	if err != nil {
		return errors.New("cannot register a new user")
	}
	return nil
}

func (s *auth) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("cannot find user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid password provided")
	}

	accessToken, err = s.jwtMaker.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", "", errors.New("cannot generate access token")
	}

	refreshToken, err = s.jwtMaker.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", errors.New("cannot generate refresh token")
	}

	return accessToken, refreshToken, nil
}

func (s *auth) Refresh(ctx context.Context, refreshToken string) (string, error) {
	userID, err := s.checkRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", err
	}

	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return "", errors.New("cannot find user")
	}

	accessToken, err := s.jwtMaker.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", errors.New("cannot generate access token")
	}
	return accessToken, err
}

func (s *auth) Logout(ctx context.Context, refreshToken string) error {
	userID, err := s.checkRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	return s.tokenRepository.Blacklist(ctx, userID, time.Now(), config.GetDuration("REFRESH_TOKEN_DURATION", time.Hour*24*7))
}

func (s *auth) ValidateAccessToken(ctx context.Context, accessToken string) (*jwt.AccessTokenClaims, error) {
	claims, err := s.jwtMaker.VerifyAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	blacklisted, err := s.isBlacklisted(ctx, claims.UserID, claims.IssuedAt)
	if err != nil {
		return nil, err
	}
	if blacklisted {
		return nil, errors.New("token is blacklisted")
	}

	return &claims, nil
}

func (s *auth) checkRefreshToken(ctx context.Context, refreshToken string) (int, error) {
	claims, err := s.jwtMaker.VerifyRefreshToken(refreshToken)
	if err != nil {
		return 0, err
	}

	blacklisted, err := s.isBlacklisted(ctx, claims.UserID, claims.IssuedAt)
	if err != nil {
		return 0, err
	}
	if blacklisted {
		return 0, errors.New("token is blacklisted")
	}

	return claims.UserID, nil
}

func (s *auth) isBlacklisted(ctx context.Context, userID int, issuedAt int64) (bool, error) {
	t, inBlacklist, err := s.tokenRepository.IsBlacklisted(ctx, userID)
	if err != nil {
		return false, err
	}

	if inBlacklist {
		iat := time.Unix(issuedAt, 0)
		if t.After(iat) {
			return true, nil
		}
	}
	return false, nil
}
