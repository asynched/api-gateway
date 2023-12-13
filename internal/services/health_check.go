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
	serviceRepository repositories.ServiceRepository
	client            *http.Client
}

func NewHealthCheckService(serviceRepository repositories.ServiceRepository) *HealthCheckService {
	return &HealthCheckService{
		serviceRepository: serviceRepository,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *HealthCheckService) Check(service entities.Service) error {
	request, err := http.NewRequest("GET", fmt.Sprintf("http://%s/healthcheck", service.Address), nil)

	if err != nil {
		return err
	}

	response, err := s.client.Do(request)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("service %s is not healthy", service.Name)
	}

	return nil
}

func (healthCheck *HealthCheckService) Run() {
	ticker := time.NewTicker(5 * time.Second)

	<-ticker.C

	for {
		for _, service := range healthCheck.serviceRepository.FindAll() {
			if err := healthCheck.Check(service); err != nil {
				log.Printf("service='%s' status='%s' error='%s'\n", service.Id, entities.ServiceStatusUnhealthy, err)
				healthCheck.serviceRepository.UpdateStatus(service.Id, entities.ServiceStatusUnhealthy)
				continue
			}

			if service.Status != entities.ServiceStatusHealthy {
				log.Printf("service='%s' status='%s'\n", service.Id, entities.ServiceStatusHealthy)
				healthCheck.serviceRepository.UpdateStatus(service.Id, entities.ServiceStatusHealthy)
			}
		}

		<-ticker.C
	}
}
