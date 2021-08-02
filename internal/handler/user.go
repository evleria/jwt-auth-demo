package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/evleria/jwt-auth-demo/internal/jwt"
	"github.com/evleria/jwt-auth-demo/internal/service"
)

// User contains method of current user
type User interface {
	Me(context echo.Context) error
}

type user struct {
	userService service.User
}

// NewUserHandler creates user handler
func NewUserHandler(userService service.User) User {
	return &user{
		userService: userService,
	}
}

// Me godoc
// @Tags User
// @Summary Shows information about current user
// @Param Authorization header string true "Authorization header"
// @Success 200 {object} MeResponse
// @Failure 400 {object} DefaultHTTPError
// @Failure 500 {object} DefaultHTTPError
// @Router /api/user/me [get]
func (u *user) Me(ctx echo.Context) error {
	claims := ctx.Get("user").(*jwt.AccessTokenClaims)

	user, err := u.userService.Me(ctx.Request().Context(), claims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := MeResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	return ctx.JSON(http.StatusOK, response)
}

// MeResponse contains information of response
type MeResponse struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}
