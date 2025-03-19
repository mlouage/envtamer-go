package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Storage represents a SQLite database connection
type Storage struct {
	db *sql.DB
}

// New creates a new Storage instance
func New() (*Storage, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	dbDir := filepath.Join(home, ".envtamer-go")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	dbPath := filepath.Join(dbDir, "envtamer-go.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Init() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS "EnvVariables" (
            "Directory" TEXT NOT NULL,
            "Key" TEXT NOT NULL,
            "Value" TEXT NOT NULL,
            CONSTRAINT "PK_EnvVariables" PRIMARY KEY ("Directory", "Key")
        );
	`)
	if err != nil {
		return fmt.Errorf("failed to initialize database table: %w", err)
	}
	return nil
}

// Close closes the database connection
func (s *Storage) Close() error {
	return s.db.Close()
}

// SaveEnvVars saves environment variables for a directory
func (s *Storage) SaveEnvVars(directory string, envVars map[string]string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete existing env vars for this directory
	_, err = tx.Exec("DELETE FROM EnvVariables WHERE Directory = ?", directory)
	if err != nil {
		return fmt.Errorf("failed to delete existing env vars: %w", err)
	}

	// Insert new env vars
	stmt, err := tx.Prepare("INSERT INTO EnvVariables (Directory, Key, Value) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for key, value := range envVars {
		_, err = stmt.Exec(directory, key, value)
		if err != nil {
			return fmt.Errorf("failed to insert env var: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetEnvVars retrieves environment variables for a directory
func (s *Storage) GetEnvVars(directory string) (map[string]string, error) {
	rows, err := s.db.Query("SELECT Key, Value FROM EnvVariables WHERE Directory = ?", directory)
	if err != nil {
		return nil, fmt.Errorf("failed to query env vars: %w", err)
	}
	defer rows.Close()

	envVars := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		envVars[key] = value
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	if len(envVars) == 0 {
		// Check if directory exists at all
		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM EnvVariables WHERE Directory = ?", directory).Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("failed to check directory existence: %w", err)
		}
		if count == 0 {
			return nil, fmt.Errorf("directory not found: %s", directory)
		}
	}

	return envVars, nil
}

// ListDirectories lists all directories
func (s *Storage) ListDirectories() ([]string, error) {
	rows, err := s.db.Query("SELECT DISTINCT Directory FROM EnvVariables")
	if err != nil {
		return nil, fmt.Errorf("failed to query directories: %w", err)
	}
	defer rows.Close()

	var directories []string
	for rows.Next() {
		var dir string
		if err := rows.Scan(&dir); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		directories = append(directories, dir)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return directories, nil
}
