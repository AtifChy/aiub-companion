// Package database provides functions to initialize and manage the database connection for the AIUB Companion application.
package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/wailsapp/wails/v3/pkg/application"
	_ "modernc.org/sqlite"
)

// Service holds the database connection and queries
type Service struct {
	db   *sql.DB
	path string
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	path, err := dbPath()
	if err != nil {
		return fmt.Errorf("failed to get database path: %w", err)
	}

	instance, err := open(path)
	if err != nil {
		return err
	}

	s.path = path
	s.db = instance.db
	return nil
}

func (s *Service) ServiceShutdown() error {
	return s.Close()
}

func (s *Service) DB() *sql.DB {
	return s.db
}

// Close closes the database connection
func (s *Service) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
