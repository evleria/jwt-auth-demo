package auth

import (
	ctrl "github.com/evleria/jwt-auth-demo/pkg/common/controller"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Controller interface {
	Register(context echo.Context) error
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
// @Summary Register a new user
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
