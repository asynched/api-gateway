package entities

type Server struct {
	// Server id
	Id string `json:"id"`
	// Server tag, e.g: "users"
	Tag string `json:"tag"`
	// Server host to check on headers e.g: "localhost:3000"
	Host string `json:"host"`
	// Server address, e.g: "http://localhost:3000"
	Address string `json:"address"`
	// Server name, e.g: "Users API"
	Name string `json:"name"`
	// Server status, e.g: "healthy"
	Status ServerStatus `json:"status"`
	// Server created at, e.g: "2020-01-01T00:00:00Z"
	CreatedAt string `json:"createdAt"`
	// Server updated at, e.g: "2020-01-01T00:00:00Z"
	UpdatedAt string `json:"updatedAt"`
}

type ServerStatus string

const (
	ServerStatusHealthy   ServerStatus = "healthy"
	ServerStatusUnhealthy ServerStatus = "unhealthy"
	ServerStatusUnknown   ServerStatus = "unknown"
)
