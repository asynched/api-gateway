package controllers

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/asynched/api-gateway/internal/domain/repositories"
	"github.com/gofiber/fiber/v2"
)

type ProxyController struct {
	serviceRepository repositories.ServiceRepository
}

func NewProxyController(serviceRepository repositories.ServiceRepository) *ProxyController {
	return &ProxyController{
		serviceRepository: serviceRepository,
	}
}

func (controller *ProxyController) HandleRequest(c *fiber.Ctx) error {
	host := c.Get("host")

	if host == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	services := controller.serviceRepository.FindByHost(host)

	if len(services) == 0 {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "No services available for this host",
		})
	}

	service := services[rand.Int()%len(services)]
	requestUrl := fmt.Sprintf("http://%s%s", service.Address, c.OriginalURL())

	reader := bytes.NewReader(c.Body())
	request, err := http.NewRequest(c.Method(), requestUrl, reader)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
			"cause": err.Error(),
		})
	}

	for key := range c.GetReqHeaders() {
		request.Header.Set(key, c.Get(key))
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := client.Do(request)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
			"cause": err.Error(),
		})
	}

	for key := range response.Header {
		c.Set(key, response.Header.Get(key))
	}

	return c.Status(response.StatusCode).SendStream(response.Body)
}

func (controller *ProxyController) Setup(router fiber.Router) {
	router.All("*", controller.HandleRequest)
}
