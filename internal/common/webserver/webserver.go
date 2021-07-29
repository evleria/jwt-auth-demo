package webserver

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

type WebServer interface {
	Listen(port int) error
	Engine() *echo.Echo
}

type webserver struct {
	e *echo.Echo
}

func New() WebServer {
	return &webserver{
		e: echo.New(),
	}
}

func (ws *webserver) Listen(port int) error {
	return ws.e.Start(fmt.Sprintf(":%d", port))
}

func (ws *webserver) Engine() *echo.Echo {
	return ws.e
}
