package jwt

import (
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Maker interface {
	GenerateAccessToken(claims jwtgo.MapClaims) (string, error)
	GenerateRefreshToken(claims jwtgo.MapClaims) (string, error)
}

type maker struct {
	config MakerConfig
}

func NewJwtMaker(config MakerConfig) Maker {
	return &maker{
		config: config,
	}
}

func (m *maker) GenerateAccessToken(claims jwtgo.MapClaims) (string, error) {
	return m.generateJwt(claims, m.config.AccessTokenDuration, m.config.AccessTokenSecret)
}

func (m *maker) GenerateRefreshToken(claims jwtgo.MapClaims) (string, error) {
	return m.generateJwt(claims, m.config.RefreshTokenDuration, m.config.RefreshTokenSecret)
}

func (m *maker) generateJwt(claims jwtgo.MapClaims, exp time.Duration, secret string) (string, error) {
	id, _ := gonanoid.New()
	now := time.Now()

	claims["id"] = id
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(exp).Unix()

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
