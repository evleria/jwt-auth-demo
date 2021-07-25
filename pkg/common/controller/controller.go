package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strings"
)

type Base struct {
	validate *validator.Validate
}

func NewBase() *Base {
	return &Base{
		validate: validator.New(),
	}
}

func (b *Base) Validate(input interface{}) error {
	err := b.validate.Struct(input)
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
