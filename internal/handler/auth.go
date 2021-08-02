// Package handler encapsulates work with HTTP
package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"

	"github.com/evleria/jwt-auth-demo/internal/service"
)

// Auth contains http handlers for each endpoint in auth group
type Auth interface {
	Register(context echo.Context) error
	Login(context echo.Context) error
	Refresh(context echo.Context) error
	Logout(context echo.Context) error
}

type auth struct {
	validate *validator.Validate
	service  service.Auth
}

// NewAuthHandler creates auth handler
func NewAuthHandler(svc service.Auth) Auth {
	return &auth{
		validate: validator.New(),
		service:  svc,
	}
}

// Register godoc
// @Tags Auth
// @Summary Registers a new user
// @Param registerData body RegisterRequest true "Registration information"
// @Success 201 "Created"
// @Failure 400 {object} DefaultHTTPError
// @Failure 500 {object} DefaultHTTPError
// @Router /auth/register [post]
func (c *auth) Register(ctx echo.Context) error {
	request := new(RegisterRequest)
	err := ctx.Bind(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err = c.Validate(request)
	if err != nil {
		return err
	}

	err = c.service.Register(ctx.Request().Context(), request.FirstName, request.LastName, request.Email, request.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusCreated)
}

// Login godoc
// @Tags Auth
// @Summary Logins a user
// @Param loginData body LoginRequest true "Login information"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} DefaultHTTPError
// @Failure 500 {object} DefaultHTTPError
// @Router /auth/login [post]
func (c *auth) Login(ctx echo.Context) error {
	request := new(LoginRequest)
	err := ctx.Bind(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err = c.Validate(request)
	if err != nil {
		return err
	}
	accessToken, refreshToken, err := c.service.Login(ctx.Request().Context(), request.Email, request.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	response := &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return ctx.JSON(http.StatusOK, response)
}

// Refresh godoc
// @Tags Auth
// @Summary Refreshes access token
// @Param refreshData body RefreshRequest true "Refresh token"
// @Success 200 {object} RefreshResponse
// @Failure 400 {object} DefaultHTTPError
// @Failure 500 {object} DefaultHTTPError
// @Router /auth/refresh [post]
func (c *auth) Refresh(ctx echo.Context) error {
	request := new(RefreshRequest)
	err := ctx.Bind(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	accessToken, err := c.service.Refresh(ctx.Request().Context(), request.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := RefreshResponse{
		AccessToken: accessToken,
	}
	return ctx.JSON(http.StatusOK, response)
}

// Logout godoc
// @Tags Auth
// @Summary Logouts a user
// @Param logoutData body LogoutRequest true "Refresh token"
// @Success 200 "Logged out"
// @Failure 400 {object} DefaultHTTPError
// @Failure 500 {object} DefaultHTTPError
// @Router /auth/logout [post]
func (c *auth) Logout(ctx echo.Context) error {
	request := new(LogoutRequest)
	err := ctx.Bind(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err = c.service.Logout(ctx.Request().Context(), request.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}

func (c *auth) Validate(input interface{}) error {
	err := c.validate.Struct(input)
	if err != nil {
		validationErrs := err.(validator.ValidationErrors)
		fields := make([]string, 0, len(validationErrs))
		for _, e := range validationErrs {
			fields = append(fields, e.Field())
		}
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid fields: %s", strings.Join(fields, ", ")))
	}
	return nil
}

// DefaultHTTPError represents default http error
type DefaultHTTPError struct {
	Message string `json:"message"`
}

// RegisterRequest represents request of register endpoint
type RegisterRequest struct {
	FirstName string `json:"firstName" validate:"required,min=2,max=20"`
	LastName  string `json:"lastName" validate:"required,min=2,max=20"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=30"`
}

// LoginRequest represents request of login endpoint
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}

// LoginResponse represents response of login endpoint
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// RefreshRequest represents request of refresh endpoint
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// RefreshResponse represents response of refresh endpoint
type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}

// LogoutRequest represents request of logout endpoint
type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}
