package jwt

import (
	"errors"
	"github.com/evleria/jwt-auth-demo/internal/config"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Maker interface {
	GenerateAccessToken(userId int, email string) (string, error)
	GenerateRefreshToken(userId int) (string, error)
	VerifyAccessToken(accessToken string) (AccessTokenClaims, error)
	VerifyRefreshToken(refreshToken string) (RefreshTokenClaims, error)
}

type maker struct {
	accessTokenSecret    string
	accessTokenDuration  time.Duration
	refreshTokenSecret   string
	refreshTokenDuration time.Duration
}

func NewJwtMakerFromConfig() Maker {
	return &maker{
		accessTokenSecret:    config.GetString("ACCESS_TOKEN_SECRET", "access_secret"),
		accessTokenDuration:  config.GetDuration("ACCESS_TOKEN_DURATION", time.Minute*5),
		refreshTokenSecret:   config.GetString("REFRESH_TOKEN_SECRET", "refresh_secret"),
		refreshTokenDuration: config.GetDuration("REFRESH_TOKEN_DURATION", time.Hour*24*7),
	}
}

func (m *maker) GenerateAccessToken(userId int, email string) (string, error) {
	return m.generateJwt(&AccessTokenClaims{
		UserId: userId,
		Email:  email,
	}, m.accessTokenDuration, m.accessTokenSecret)
}

func (m *maker) GenerateRefreshToken(userId int) (string, error) {
	return m.generateJwt(&RefreshTokenClaims{
		UserId: userId,
	}, m.refreshTokenDuration, m.refreshTokenSecret)
}

func (m *maker) VerifyAccessToken(accessToken string) (AccessTokenClaims, error) {
	claims := AccessTokenClaims{}
	err := m.verifyJwt(accessToken, m.accessTokenSecret, &claims)
	return claims, err
}

func (m *maker) VerifyRefreshToken(refreshToken string) (RefreshTokenClaims, error) {
	claims := RefreshTokenClaims{}
	err := m.verifyJwt(refreshToken, m.refreshTokenSecret, &claims)
	return claims, err
}

func (m *maker) generateJwt(claims Claims, exp time.Duration, secret string) (string, error) {
	id, err := gonanoid.New()
	if err != nil {
		return "", err
	}

	now := time.Now()

	claims.SetId(id)
	claims.SetIssuedAt(now)
	claims.SetExpiresAt(now.Add(exp))

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (m *maker) verifyJwt(token string, secret string, claims Claims) error {
	t, err := jwtgo.ParseWithClaims(token, claims, func(t *jwtgo.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return err
	}

	if _, ok := t.Claims.(Claims); !ok && !t.Valid {
		return errors.New("token is invalid")
	}
	return nil
}
