// Package database provides functions to initialize and manage the database connection for the AIUB Companion application.
package database

import (
	"context"
	"database/sql"
	"embed"

	"github.com/wailsapp/wails/v3/pkg/application"
	_ "modernc.org/sqlite"
)

//go:embed sql/schemas/*.sql
var schemasFS embed.FS

// Service holds the database connection and queries
type Service struct {
	db *sql.DB
}

func NewService() *Service {
	return &Service{}
}

func (db *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	instance, err := open()
	if err != nil {
		return err
	}
	db.db = instance.db
	return nil
}

func (db *Service) ServiceShutdown() error {
	return db.Close()
}

func (db *Service) DB() *sql.DB {
	return db.db
}

// Close closes the database connection
func (db *Service) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}
