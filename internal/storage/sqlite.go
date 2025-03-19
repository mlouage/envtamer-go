package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

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

// Init initializes the database tables
func (s *Storage) Init() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS directories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			path TEXT UNIQUE
		);
		CREATE TABLE IF NOT EXISTS env_vars (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			directory_id INTEGER,
			key TEXT,
			value TEXT,
			FOREIGN KEY (directory_id) REFERENCES directories(id),
			UNIQUE(directory_id, key)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to initialize database tables: %w", err)
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

	// Get or create directory
	var dirID int
	err = tx.QueryRow("SELECT id FROM directories WHERE path = ?", directory).Scan(&dirID)
	if err != nil {
		if err == sql.ErrNoRows {
			result, err := tx.Exec("INSERT INTO directories (path) VALUES (?)", directory)
			if err != nil {
				return fmt.Errorf("failed to insert directory: %w", err)
			}
			id, err := result.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get last insert ID: %w", err)
			}
			dirID = int(id)
		} else {
			return fmt.Errorf("failed to query directory: %w", err)
		}
	}

	// Delete existing env vars for this directory
	_, err = tx.Exec("DELETE FROM env_vars WHERE directory_id = ?", dirID)
	if err != nil {
		return fmt.Errorf("failed to delete existing env vars: %w", err)
	}

	// Insert new env vars
	stmt, err := tx.Prepare("INSERT INTO env_vars (directory_id, key, value) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for key, value := range envVars {
		_, err = stmt.Exec(dirID, key, value)
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
	var dirID int
	err := s.db.QueryRow("SELECT id FROM directories WHERE path = ?", directory).Scan(&dirID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("directory not found: %s", directory)
		}
		return nil, fmt.Errorf("failed to query directory: %w", err)
	}

	rows, err := s.db.Query("SELECT key, value FROM env_vars WHERE directory_id = ?", dirID)
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

	return envVars, nil
}

// ListDirectories lists all directories
func (s *Storage) ListDirectories() ([]string, error) {
	rows, err := s.db.Query("SELECT path FROM directories")
	if err != nil {
		return nil, fmt.Errorf("failed to query directories: %w", err)
	}
	defer rows.Close()

	var directories []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		directories = append(directories, path)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return directories, nil
}
