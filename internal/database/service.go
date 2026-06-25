// Package database provides functions to initialize and manage the database connection for the AIUB Companion application.
package database

import (
	"context"
	"database/sql"

	"github.com/wailsapp/wails/v3/pkg/application"
	_ "modernc.org/sqlite"
)

// Service holds the database connection and queries
type Service struct {
	db *sql.DB
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	instance, err := open()
	if err != nil {
		return err
	}
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
