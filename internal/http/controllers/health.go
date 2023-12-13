package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type HealthController struct {
	startup time.Time
}

func NewHealthController() *HealthController {
	return &HealthController{
		startup: time.Now(),
	}
}

func (controller *HealthController) HandleCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "up",
		"startup": controller.startup.Format(time.RFC3339),
		"uptime":  time.Since(controller.startup).String(),
	})
}

func (controller *HealthController) Setup(router fiber.Router) {
	router.Get("/check", controller.HandleCheck)
}
