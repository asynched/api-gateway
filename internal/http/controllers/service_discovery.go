package controllers

import (
	"github.com/asynched/api-gateway/internal/domain/repositories"
	"github.com/asynched/api-gateway/internal/dto"
	"github.com/gofiber/fiber/v2"
)

type ServiceDiscoveryController struct {
	serverRepository repositories.ServerRepository
}

func NewServiceDiscoveryController(serverRepository repositories.ServerRepository) *ServiceDiscoveryController {
	return &ServiceDiscoveryController{
		serverRepository: serverRepository,
	}
}

func (controller *ServiceDiscoveryController) HandleFindAll(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"data": controller.serverRepository.FindAll(),
	})
}

func (controller *ServiceDiscoveryController) HandleRegister(c *fiber.Ctx) error {
	request := dto.CreateServerDto{}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"cause": err.Error(),
		})
	}

	service, err := controller.serverRepository.Register(request.ToServer())

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to register service",
			"cause": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": service,
	})
}

func (controller *ServiceDiscoveryController) HandleDelete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := controller.serverRepository.Delete(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to delete service",
			"cause": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Service deleted successfully",
	})
}

func (controller *ServiceDiscoveryController) HandleFindById(c *fiber.Ctx) error {
	id := c.Params("id")

	service, err := controller.serverRepository.FindById(id)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Failed to find service",
			"cause": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": service,
	})
}

func (controller *ServiceDiscoveryController) Setup(router fiber.Router) {
	router.Post("/", controller.HandleRegister)
	router.Get("/", controller.HandleFindAll)
	router.Get("/:id", controller.HandleFindById)
	router.Delete("/:id", controller.HandleDelete)
}
