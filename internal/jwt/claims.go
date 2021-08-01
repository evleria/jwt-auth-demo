package jwt

import (
	jwtgo "github.com/dgrijalva/jwt-go"
	"time"
)

type Claims interface {
	jwtgo.Claims
	SetId(id string)
	SetIssuedAt(issuedAt time.Time)
	SetExpiresAt(expiresAt time.Time)
}

type AccessTokenClaims struct {
	UserId int    `json:"userId"`
	Email  string `json:"email"`
	jwtgo.StandardClaims
}

func (t *AccessTokenClaims) SetId(id string) {
	t.Id = id
}

func (t *AccessTokenClaims) SetIssuedAt(issuedAt time.Time) {
	t.IssuedAt = issuedAt.Unix()
}

func (t *AccessTokenClaims) SetExpiresAt(expiresAt time.Time) {
	t.ExpiresAt = expiresAt.Unix()
}

type RefreshTokenClaims struct {
	UserId int `json:"userId"`
	jwtgo.StandardClaims
}

func (t *RefreshTokenClaims) SetId(id string) {
	t.Id = id
}

func (t *RefreshTokenClaims) SetIssuedAt(issuedAt time.Time) {
	t.IssuedAt = issuedAt.Unix()
}

func (t *RefreshTokenClaims) SetExpiresAt(expiresAt time.Time) {
	t.ExpiresAt = expiresAt.Unix()
}
