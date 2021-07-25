package auth

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strings"
)

type Controller interface {
	Register(context echo.Context) error
}

type controller struct {
	validate *validator.Validate
	service  Service
}

func NewController(service Service) Controller {
	return &controller{
		validate: validator.New(),
		service:  service,
	}
}

// Register godoc
// @Summary Register a new user
// @Param registerData body RegisterRequest true "Registration information"
// @Success 201 "Created"
// @Failure 400 {object} DefaultHttpError
// @Failure 500 {object} DefaultHttpError
// @Router /auth/register [post]
func (c *controller) Register(ctx echo.Context) error {
	request := new(RegisterRequest)
	err := ctx.Bind(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err = c.validateStruct(request)
	if err != nil {
		return validationError(err)
	}

	err = c.service.Register(request.FirstName, request.LastName, request.Email, request.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusCreated)
}

func (c *controller) validateStruct(input interface{}) error {
	return c.validate.Struct(input)
}

func validationError(err error) *echo.HTTPError {
	errors := err.(validator.ValidationErrors)
	invalidFields := make([]string, 0, len(errors))
	for _, e := range errors {
		invalidFields = append(invalidFields, e.Field())
	}
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid fields: %s", strings.Join(invalidFields, ", ")))
}

type DefaultHttpError struct {
	Message string `json:"message"`
}
