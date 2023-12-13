package entities

type Service struct {
	// Service id
	Id string `json:"id"`
	// Service tag, e.g: "users"
	Tag string `json:"tag"`
	// Service host to check on headers e.g: "localhost:3000"
	Host string `json:"host"`
	// Service address, e.g: "http://localhost:3000"
	Address string `json:"address"`
	// Service name, e.g: "Users API"
	Name string `json:"name"`
	// Service status, e.g: "healthy"
	Status ServiceStatus `json:"status"`
	// Service created at, e.g: "2020-01-01T00:00:00Z"
	CreatedAt string `json:"createdAt"`
	// Service updated at, e.g: "2020-01-01T00:00:00Z"
	UpdatedAt string `json:"updatedAt"`
}

type ServiceStatus string

const (
	ServiceStatusHealthy   ServiceStatus = "healthy"
	ServiceStatusUnhealthy ServiceStatus = "unhealthy"
	ServiceStatusUnknown   ServiceStatus = "unknown"
)
