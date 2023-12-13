package database

import (
	"database/sql"
	"fmt"
)

func CreateClient(filename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_jounal_mode=WAL&mode=rwc&cache=shared", filename))

	if err != nil {
		return nil, err
	}

	return db, nil
}
