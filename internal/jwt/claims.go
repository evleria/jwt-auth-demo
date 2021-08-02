// Package jwt encapsulates work with access and refresh JWT tokens
package jwt

import (
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
)

// Claims represents common jwt claims
type Claims interface {
	jwtgo.Claims
	SetID(id string)
	SetIssuedAt(issuedAt time.Time)
	SetExpiresAt(expiresAt time.Time)
}

// AccessTokenClaims contains claims of access token
type AccessTokenClaims struct {
	UserID int    `json:"userId"`
	Email  string `json:"email"`
	jwtgo.StandardClaims
}

// SetID sets id claim of access token
func (t *AccessTokenClaims) SetID(id string) {
	t.Id = id
}

// SetIssuedAt sets issuedAt claim of access token
func (t *AccessTokenClaims) SetIssuedAt(issuedAt time.Time) {
	t.IssuedAt = issuedAt.Unix()
}

// SetExpiresAt sets expiresAt claim of access token
func (t *AccessTokenClaims) SetExpiresAt(expiresAt time.Time) {
	t.ExpiresAt = expiresAt.Unix()
}

// RefreshTokenClaims contains claims of refresh token
type RefreshTokenClaims struct {
	UserID int `json:"userId"`
	jwtgo.StandardClaims
}

// SetID sets id claim of refresh token
func (t *RefreshTokenClaims) SetID(id string) {
	t.Id = id
}

// SetIssuedAt sets issuedAt claim of refresh token
func (t *RefreshTokenClaims) SetIssuedAt(issuedAt time.Time) {
	t.IssuedAt = issuedAt.Unix()
}

// SetExpiresAt set expiresAt claim of refresh token
func (t *RefreshTokenClaims) SetExpiresAt(expiresAt time.Time) {
	t.ExpiresAt = expiresAt.Unix()
}
