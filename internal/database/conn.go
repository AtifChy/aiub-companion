package database

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"aiub-companion/internal/config"
)

//go:embed sql/schemas/*.sql
var schemasFS embed.FS

// open initializes the database connection and creates the schema if it doesn't exist
func open(path string) (*Service, error) {
	// Open the database connection
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Set connection pool settings
	conn.SetMaxOpenConns(1)

	// Enable foreign keys and WAL journal mode
	if _, err := conn.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}
	if _, err := conn.Exec("PRAGMA journal_mode = WAL;"); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("set journal mode: %w", err)
	}

	// Create the schema if it doesn't exist
	if err := runSchemas(conn); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("run schemas: %w", err)
	}

	return &Service{db: conn}, nil
}

// dbPath returns the path to the database file, creating the necessary directories if they don't exist
func dbPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user config dir: %w", err)
	}

	appDir := filepath.Join(configDir, config.AppName)
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return "", fmt.Errorf("create app dir: %w", err)
	}
	return filepath.Join(appDir, config.AppName+".db"), nil
}

// runSchemas reads and executes all SQL schema files from the embedded filesystem
func runSchemas(conn *sql.DB) error {
	entries, err := schemasFS.ReadDir("sql/schemas")
	if err != nil {
		return fmt.Errorf("read schema dir: %w", err)
	}

	for _, entry := range entries {
		sql, err := schemasFS.ReadFile("sql/schemas/" + entry.Name())
		if err != nil {
			return fmt.Errorf("read schema file %s: %w", entry.Name(), err)
		}
		if _, err := conn.Exec(string(sql)); err != nil {
			return fmt.Errorf("execute schema file %s: %w", entry.Name(), err)
		}
	}

	return nil
}
