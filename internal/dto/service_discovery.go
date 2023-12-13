package dto

import "github.com/asynched/api-gateway/internal/domain/entities"

type CreateServerDto struct {
	Tag     string `json:"tag"`
	Host    string `json:"host"`
	Address string `json:"address"`
	Name    string `json:"name"`
}

func (dto *CreateServerDto) ToServer() entities.Server {
	return entities.Server{
		Tag:     dto.Tag,
		Host:    dto.Host,
		Address: dto.Address,
		Name:    dto.Name,
		Status:  entities.ServerStatusUnknown,
	}
}
