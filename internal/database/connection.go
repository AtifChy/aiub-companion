// Package database provides functions to initialize and manage the database connection for the AIUB Companion application.
package database

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"aiub-companion/internal/app"

	_ "modernc.org/sqlite"
)

//go:embed sql/schemas/*.sql
var schemasFS embed.FS

// DBInstance holds the database connection and queries
type DBInstance struct {
	SQLDB *sql.DB
}

// InitDB initializes the database connection and creates the schema if it doesn't exist
func InitDB() (*DBInstance, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config dir: %w", err)
	}

	appDir := filepath.Join(configDir, app.Name)
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create app dir: %w", err)
	}

	dbPath := filepath.Join(appDir, app.Name+".db")

	// Open the database connection
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	conn.SetMaxOpenConns(1)

	// Enable foreign keys and WAL journal mode
	if _, err := conn.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}
	if _, err := conn.Exec("PRAGMA journal_mode = WAL;"); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to set journal mode: %w", err)
	}

	// Create the schema if it doesn't exist
	if err := runSchemas(conn); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to run schemas: %w", err)
	}

	return &DBInstance{SQLDB: conn}, nil
}

// runSchemas reads and executes all SQL schema files from the embedded filesystem
func runSchemas(conn *sql.DB) error {
	entries, err := schemasFS.ReadDir("sql/schemas")
	if err != nil {
		return fmt.Errorf("failed to read schema dir: %w", err)
	}

	for _, entry := range entries {
		sql, err := schemasFS.ReadFile("sql/schemas/" + entry.Name())
		if err != nil {
			return fmt.Errorf("failed to read schema file %s: %w", entry.Name(), err)
		}
		if _, err := conn.Exec(string(sql)); err != nil {
			return fmt.Errorf("failed to execute schema file %s: %w", entry.Name(), err)
		}
	}

	return nil
}

// Close closes the database connection
func (db *DBInstance) Close() error {
	if db.SQLDB != nil {
		return db.SQLDB.Close()
	}
	return nil
}
