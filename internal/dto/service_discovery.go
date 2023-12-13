package dto

import "github.com/asynched/api-gateway/internal/domain/entities"

type ServiceDiscoveryRequest struct {
	Tag     string `json:"tag"`
	Host    string `json:"host"`
	Address string `json:"address"`
	Name    string `json:"name"`
}

func (r *ServiceDiscoveryRequest) ToService() entities.Service {
	return entities.Service{
		Tag:     r.Tag,
		Host:    r.Host,
		Address: r.Address,
		Name:    r.Name,
		Status:  entities.ServiceStatusUnknown,
	}
}
