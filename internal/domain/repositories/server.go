package repositories

import (
	"database/sql"

	"github.com/asynched/api-gateway/internal/domain/entities"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type ServerRepository interface {
	FindAll() []entities.Server
	FindById(id string) (entities.Server, error)
	FindByHost(host string) []entities.Server
	Delete(id string) error
	Register(server entities.Server) (entities.Server, error)
	UpdateStatus(serverId string, status entities.ServerStatus) error
}

type sqlServerRepositoryQueries struct {
	findAll      *sql.Stmt
	findById     *sql.Stmt
	findByHost   *sql.Stmt
	register     *sql.Stmt
	updateStatus *sql.Stmt
	delete       *sql.Stmt
}

type SQLServerRepository struct {
	queries *sqlServerRepositoryQueries
}

func initDatabase(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS "servers" (
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
		
		CREATE INDEX IF NOT EXISTS "servers_tag" ON "servers" ("tag");
		
		CREATE INDEX IF NOT EXISTS "servers_host" ON "servers" ("host");
		
		CREATE INDEX IF NOT EXISTS "servers_address" ON "servers" ("address");
	`)

	return err
}

func NewSQLServerRepository(db *sql.DB) (*SQLServerRepository, error) {
	if err := initDatabase(db); err != nil {
		return nil, err
	}

	queries := &sqlServerRepositoryQueries{}

	var err error

	queries.findAll, err = db.Prepare("SELECT * FROM servers")
	if err != nil {
		return nil, err
	}

	queries.findByHost, err = db.Prepare("SELECT * FROM servers WHERE host = $1")
	if err != nil {
		return nil, err
	}

	queries.register, err = db.Prepare("INSERT INTO servers (id, tag, host, address, name, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING *")
	if err != nil {
		return nil, err
	}

	queries.updateStatus, err = db.Prepare("UPDATE servers SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2")
	if err != nil {
		return nil, err
	}

	queries.findById, err = db.Prepare("SELECT * FROM servers WHERE id = $1")
	if err != nil {
		return nil, err
	}

	queries.delete, err = db.Prepare("DELETE FROM servers WHERE id = $1")
	if err != nil {
		return nil, err
	}

	return &SQLServerRepository{queries}, nil
}

func (r *SQLServerRepository) FindAll() []entities.Server {
	rows, err := r.queries.findAll.Query()
	if err != nil {
		return nil
	}
	defer rows.Close()

	servers := make([]entities.Server, 0)

	for rows.Next() {
		server := entities.Server{}
		err := rows.Scan(&server.Id, &server.Tag, &server.Host, &server.Address, &server.Name, &server.Status, &server.CreatedAt, &server.UpdatedAt)

		if err != nil {
			return nil
		}

		servers = append(servers, server)
	}

	return servers
}

func (r *SQLServerRepository) FindById(id string) (entities.Server, error) {
	server := entities.Server{}
	rows, err := r.queries.findAll.Query()

	if err != nil {
		return server, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&server.Id, &server.Tag, &server.Host, &server.Address, &server.Name, &server.Status, &server.CreatedAt, &server.UpdatedAt)
		if err != nil {
			return server, err
		}
	}

	return server, nil
}

func (r *SQLServerRepository) FindByHost(host string) []entities.Server {
	rows, err := r.queries.findByHost.Query(host)

	if err != nil {
		return nil
	}

	defer rows.Close()

	servers := make([]entities.Server, 0)

	for rows.Next() {
		server := entities.Server{}
		err := rows.Scan(&server.Id, &server.Tag, &server.Host, &server.Address, &server.Name, &server.Status, &server.CreatedAt, &server.UpdatedAt)

		if err != nil {
			return nil
		}

		servers = append(servers, server)
	}

	return servers
}

func (r *SQLServerRepository) Register(server entities.Server) (entities.Server, error) {
	server.Id = uuid.NewString()

	_, err := r.queries.register.Exec(server.Id, server.Tag, server.Host, server.Address, server.Name, server.Status)

	if err != nil {
		return server, err
	}

	return r.FindById(server.Id)
}

func (r *SQLServerRepository) UpdateStatus(serverId string, status entities.ServerStatus) error {
	_, err := r.queries.updateStatus.Exec(status, serverId)

	if err != nil {
		return err
	}

	return nil
}

func (r *SQLServerRepository) Delete(id string) error {
	_, err := r.queries.delete.Exec(id)

	if err != nil {
		return err
	}

	return nil
}
