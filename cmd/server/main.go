package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/asynched/api-gateway/internal/database"
	"github.com/asynched/api-gateway/internal/domain/repositories"
	"github.com/asynched/api-gateway/internal/http"
	"github.com/asynched/api-gateway/internal/http/controllers"
	"github.com/asynched/api-gateway/internal/services"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)
	log.SetPrefix(fmt.Sprintf("[%d] [api-gateway] ", os.Getpid()))
}

var (
	host = flag.String("host", "0.0.0.0", "Host to run the server")
	port = flag.Int("port", 9190, "Port to run the server")
	file = flag.String("file", "api-gateway.db", "Database file")
)

func main() {
	flag.Parse()

	server := http.NewServer()

	db, err := database.CreateClient(*file)

	if err != nil {
		log.Fatalf("Failed to create database client: %s", err.Error())
	}

	defer db.Close()

	serviceRepository, err := repositories.NewSQLServiceRepository(db)

	if err != nil {
		log.Fatalf("Failed to create service repository: %s", err.Error())
	}

	healthCheck := services.NewHealthCheckService(serviceRepository)
	go healthCheck.Run()

	health := controllers.NewHealthController()
	server.Setup("/health", health)

	serviceDiscovery := controllers.NewServiceDiscoveryController(serviceRepository)
	server.Setup("/services", serviceDiscovery)

	proxy := controllers.NewProxyController(serviceRepository)
	server.Setup("/", proxy)

	log.Printf("Server is running on: http://%s:%d\n", *host, *port)
	log.Printf("Check health status at: http://%s:%d/health/check\n", *host, *port)

	if err := server.Run(*host, *port); err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}
}
