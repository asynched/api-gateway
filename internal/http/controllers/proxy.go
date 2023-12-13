package controllers

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/asynched/api-gateway/internal/domain/repositories"
	"github.com/asynched/api-gateway/internal/services"
	"github.com/gofiber/fiber/v2"
)

type ProxyController struct {
	serverRepository repositories.ServerRepository
	cacheService     *services.CacheService
}

func NewProxyController(serverRepository repositories.ServerRepository) *ProxyController {
	return &ProxyController{
		serverRepository: serverRepository,
		cacheService:     services.NewCacheService(),
	}
}

func (controller *ProxyController) HandleRequest(c *fiber.Ctx) error {
	host := c.Get("host")

	if host == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	cacheKey := fmt.Sprintf("http://%s%s", host, c.OriginalURL())
	if response, ok := controller.cacheService.Get(cacheKey); ok {
		for key := range response.Headers {
			c.Set(key, response.Headers[key])
		}

		return c.Status(response.StatusCode).Send(response.Body)
	}

	servers := controller.serverRepository.FindByHost(host)

	if len(servers) == 0 {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "No services available for this host",
		})
	}

	server := servers[rand.Int()%len(servers)]
	url := fmt.Sprintf("http://%s%s", server.Address, c.OriginalURL())

	reader := bytes.NewReader(c.Body())
	request, err := http.NewRequest(c.Method(), url, reader)

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

	defer response.Body.Close()

	for key := range response.Header {
		c.Set(key, response.Header.Get(key))
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Couldn't read response body",
			"cause": err.Error(),
		})
	}

	if c.Method() == "GET" && c.Get("Cache-Control") != "no-cache" {
		cachedResponse := services.CacheResponse{
			StatusCode: response.StatusCode,
			Body:       body,
			Headers:    make(map[string]string),
			Ttl:        time.Second * 10,
		}

		for key := range response.Header {
			cachedResponse.Headers[key] = response.Header.Get(key)
		}

		controller.cacheService.Set(cacheKey, cachedResponse)
	}

	return c.Status(response.StatusCode).Send(body)
}

func (controller *ProxyController) Setup(router fiber.Router) {
	router.All("*", controller.HandleRequest)
}
