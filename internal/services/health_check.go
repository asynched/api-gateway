package services

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/asynched/api-gateway/internal/domain/entities"
	"github.com/asynched/api-gateway/internal/domain/repositories"
)

type HealthCheckService struct {
	serverRepository repositories.ServerRepository
	client           *http.Client
}

func NewHealthCheckService(serverRepository repositories.ServerRepository) *HealthCheckService {
	return &HealthCheckService{
		serverRepository: serverRepository,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *HealthCheckService) Check(server entities.Server) error {
	request, err := http.NewRequest("GET", fmt.Sprintf("http://%s/healthcheck", server.Address), nil)

	if err != nil {
		return err
	}

	response, err := s.client.Do(request)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("server %s is not healthy", server.Name)
	}

	return nil
}

func (healthCheck *HealthCheckService) Run() {
	ticker := time.NewTicker(5 * time.Second)

	<-ticker.C

	for {
		for _, server := range healthCheck.serverRepository.FindAll() {
			if err := healthCheck.Check(server); err != nil {
				log.Printf("server='%s' status='%s' error='%s'\n", server.Id, entities.ServerStatusUnhealthy, err)
				healthCheck.serverRepository.UpdateStatus(server.Id, entities.ServerStatusUnhealthy)
				continue
			}

			if server.Status != entities.ServerStatusHealthy {
				log.Printf("server='%s' status='%s'\n", server.Id, entities.ServerStatusHealthy)
				healthCheck.serverRepository.UpdateStatus(server.Id, entities.ServerStatusHealthy)
			}
		}

		<-ticker.C
	}
}
