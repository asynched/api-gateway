package repositories

import (
	"database/sql"

	"github.com/asynched/api-gateway/internal/domain/entities"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type ServiceRepository interface {
	FindAll() []entities.Service
	FindById(id string) (entities.Service, error)
	FindByHost(host string) []entities.Service
	Delete(id string) error
	Register(service entities.Service) (entities.Service, error)
	UpdateStatus(serviceId string, status entities.ServiceStatus) error
}

type sqlServiceRepositoryQueries struct {
	findAll      *sql.Stmt
	findById     *sql.Stmt
	findByHost   *sql.Stmt
	register     *sql.Stmt
	updateStatus *sql.Stmt
	delete       *sql.Stmt
}

type SQLServiceRepository struct {
	queries *sqlServiceRepositoryQueries
}

func initDatabase(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS "services" (
			"id" VARCHAR(32) NOT NULL PRIMARY KEY,
			"tag" VARCHAR(255) NOT NULL,
			"host" VARCHAR(255) NOT NULL,
			"address" VARCHAR(255) NOT NULL,
			"name" VARCHAR(255) NOT NULL,
			-- healthy, unhealthy, unknown
			"status" VARCHAR(12) NOT NULL,
			"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			"updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS "services_tag" ON "services" ("tag");
		
		CREATE INDEX IF NOT EXISTS "services_host" ON "services" ("host");
		
		CREATE INDEX IF NOT EXISTS "services_address" ON "services" ("address");
	`)

	return err
}

func NewSQLServiceRepository(db *sql.DB) (*SQLServiceRepository, error) {
	if err := initDatabase(db); err != nil {
		return nil, err
	}

	queries := &sqlServiceRepositoryQueries{}

	var err error

	queries.findAll, err = db.Prepare("SELECT * FROM services")
	if err != nil {
		return nil, err
	}

	queries.findByHost, err = db.Prepare("SELECT * FROM services WHERE host = $1")
	if err != nil {
		return nil, err
	}

	queries.register, err = db.Prepare("INSERT INTO services (id, tag, host, address, name, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING *")
	if err != nil {
		return nil, err
	}

	queries.updateStatus, err = db.Prepare("UPDATE services SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2")
	if err != nil {
		return nil, err
	}

	queries.findById, err = db.Prepare("SELECT * FROM services WHERE id = $1")
	if err != nil {
		return nil, err
	}

	queries.delete, err = db.Prepare("DELETE FROM services WHERE id = $1")
	if err != nil {
		return nil, err
	}

	return &SQLServiceRepository{queries}, nil
}

func (r *SQLServiceRepository) FindAll() []entities.Service {
	rows, err := r.queries.findAll.Query()
	if err != nil {
		return nil
	}
	defer rows.Close()

	services := make([]entities.Service, 0)

	for rows.Next() {
		service := entities.Service{}
		err := rows.Scan(&service.Id, &service.Tag, &service.Host, &service.Address, &service.Name, &service.Status, &service.CreatedAt, &service.UpdatedAt)
		if err != nil {
			return nil
		}
		services = append(services, service)
	}

	return services
}

func (r *SQLServiceRepository) FindById(id string) (entities.Service, error) {
	service := entities.Service{}
	rows, err := r.queries.findAll.Query()

	if err != nil {
		return service, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&service.Id, &service.Tag, &service.Host, &service.Address, &service.Name, &service.Status, &service.CreatedAt, &service.UpdatedAt)
		if err != nil {
			return service, err
		}
	}

	return service, nil
}

func (r *SQLServiceRepository) FindByHost(host string) []entities.Service {
	rows, err := r.queries.findByHost.Query(host)

	if err != nil {
		return nil
	}

	defer rows.Close()

	services := make([]entities.Service, 0)

	for rows.Next() {
		service := entities.Service{}
		err := rows.Scan(&service.Id, &service.Tag, &service.Host, &service.Address, &service.Name, &service.Status, &service.CreatedAt, &service.UpdatedAt)
		if err != nil {
			return nil
		}
		services = append(services, service)
	}

	return services
}

func (r *SQLServiceRepository) Register(service entities.Service) (entities.Service, error) {
	service.Id = uuid.NewString()

	_, err := r.queries.register.Exec(service.Id, service.Tag, service.Host, service.Address, service.Name, service.Status)

	if err != nil {
		return service, err
	}

	return r.FindById(service.Id)
}

func (r *SQLServiceRepository) UpdateStatus(serviceId string, status entities.ServiceStatus) error {
	_, err := r.queries.updateStatus.Exec(status, serviceId)

	if err != nil {
		return err
	}

	return nil
}

func (r *SQLServiceRepository) Delete(id string) error {
	_, err := r.queries.delete.Exec(id)

	if err != nil {
		return err
	}

	return nil
}
