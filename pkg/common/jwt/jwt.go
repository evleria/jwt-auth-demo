package jwt

import (
	"errors"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Maker interface {
	GenerateAccessToken(userId int, email string) (string, error)
	GenerateRefreshToken(userId int) (string, error)
	VerifyRefreshToken(refreshToken string) (jwtgo.MapClaims, error)
}

type maker struct {
	config MakerConfig
}

func NewJwtMaker(config MakerConfig) Maker {
	return &maker{
		config: config,
	}
}

func (m *maker) GenerateAccessToken(userId int, email string) (string, error) {
	return m.generateJwt(jwtgo.MapClaims{
		"sub":   userId,
		"email": email,
	}, m.config.AccessTokenDuration, m.config.AccessTokenSecret)
}

func (m *maker) GenerateRefreshToken(userId int) (string, error) {
	return m.generateJwt(jwtgo.MapClaims{
		"sub": userId,
	}, m.config.RefreshTokenDuration, m.config.RefreshTokenSecret)
}

func (m *maker) VerifyRefreshToken(refreshToken string) (jwtgo.MapClaims, error) {
	return m.verifyJwt(refreshToken, m.config.RefreshTokenSecret)
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

func (m *maker) verifyJwt(token string, secret string) (jwtgo.MapClaims, error) {
	t, err := jwtgo.Parse(token, func(t *jwtgo.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := t.Claims.(jwtgo.Claims); !ok && !t.Valid {
		return nil, errors.New("token is invalid")
	}
	return t.Claims.(jwtgo.MapClaims), err
}
