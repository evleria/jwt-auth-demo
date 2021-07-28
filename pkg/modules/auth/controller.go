package auth

import (
	ctrl "github.com/evleria/jwt-auth-demo/pkg/common/controller"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Controller interface {
	Register(context echo.Context) error
	Login(context echo.Context) error
	Refresh(context echo.Context) error
}

type controller struct {
	*ctrl.Base
	service Service
}

func NewController(service Service) Controller {
	return &controller{
		Base:    ctrl.NewBase(),
		service: service,
	}
}

// Register godoc
// @Tags Auth
// @Summary Registers a new user
// @Param registerData body RegisterRequest true "Registration information"
// @Success 201 "Created"
// @Failure 400 {object} ctrl.DefaultHttpError
// @Failure 500 {object} ctrl.DefaultHttpError
// @Router /auth/register [post]
func (c *controller) Register(ctx echo.Context) error {
	request := new(RegisterRequest)
	err := ctx.Bind(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err = c.Validate(request)
	if err != nil {
		return err
	}

	err = c.service.Register(request.FirstName, request.LastName, request.Email, request.Password)
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
// @Failure 400 {object} ctrl.DefaultHttpError
// @Failure 500 {object} ctrl.DefaultHttpError
// @Router /auth/login [post]
func (c *controller) Login(ctx echo.Context) error {
	request := new(LoginRequest)
	err := ctx.Bind(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err = c.Validate(request)
	if err != nil {
		return err
	}
	accessToken, refreshToken, err := c.service.Login(request.Email, request.Password)
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
// @Summary Refresh a user
// @Param refreshData body RefreshRequest true "Refresh information"
// @Success 200 {object} RefreshResponse
// @Failure 400 {object} ctrl.DefaultHttpError
// @Failure 500 {object} ctrl.DefaultHttpError
// @Router /auth/refresh [post]
func (c *controller) Refresh(ctx echo.Context) error {
	request := new(RefreshRequest)
	err := ctx.Bind(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	accessToken, err := c.service.Refresh(request.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	response := RefreshResponse{
		AccessToken: accessToken,
	}
	return ctx.JSON(http.StatusOK, response)
}
