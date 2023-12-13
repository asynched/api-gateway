package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app *fiber.App
}

func NewServer() *Server {
	return &Server{
		app: fiber.New(fiber.Config{
			DisableStartupMessage: true,
		}),
	}
}

func (server *Server) Run(host string, port int) error {
	return server.app.Listen(fmt.Sprintf("%s:%d", host, port))
}

type Controller interface {
	Setup(router fiber.Router)
}

func (setup *Server) Setup(group string, controller Controller) {
	controller.Setup(setup.app.Group(group))
}
