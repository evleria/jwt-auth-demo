package handler

import (
	"github.com/evleria/jwt-auth-demo/internal/jwt"
	"github.com/evleria/jwt-auth-demo/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type User interface {
	Me(context echo.Context) error
}

type user struct {
	userService service.User
}

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
// @Failure 400 {object} DefaultHttpError
// @Failure 500 {object} DefaultHttpError
// @Router /api/user/me [get]
func (u *user) Me(ctx echo.Context) error {
	claims := ctx.Get("user").(*jwt.AccessTokenClaims)

	user, err := u.userService.Me(ctx.Request().Context(), claims.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := MeResponse{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	return ctx.JSON(http.StatusOK, response)
}

type MeResponse struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}
