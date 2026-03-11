package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// DB wrapper struct if needed, or just return *sql.DB
type DB struct {
	*sql.DB
}

func Connect(dbPath string) (*DB, error) {
	// Ensure parent directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// modernc.org/sqlite has PRAGMAs passed in via DSN sometimes, but we'll execute them
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// PRAGMAs for performance and foreign keys
	pragmas := []string{
		"PRAGMA foreign_keys = ON;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA temp_store = MEMORY;",
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return nil, fmt.Errorf("failed to execute pragma %q: %w", p, err)
		}
	}

	return &DB{db}, nil
}

func (db *DB) Migrate(ctx context.Context, logger *slog.Logger) error {
	// Create migrations table if not exists
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}

	sort.Strings(files)

	for _, file := range files {
		var exists string
		err := db.QueryRowContext(ctx, "SELECT version FROM schema_migrations WHERE version = ?", file).Scan(&exists)
		if err == sql.ErrNoRows {
			// Needs migration
			logger.Info("applying migration", "file", file)

			content, err := migrationsFS.ReadFile(filepath.Join("migrations", file))
			if err != nil {
				return fmt.Errorf("failed to read migration %s: %w", file, err)
			}

			tx, err := db.BeginTx(ctx, nil)
			if err != nil {
				return fmt.Errorf("failed to begin tx for migration %s: %w", file, err)
			}

			if _, err := tx.ExecContext(ctx, string(content)); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to execute migration %s: %w", file, err)
			}

			if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES (?)", file); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to record migration %s: %w", file, err)
			}

			if err := tx.Commit(); err != nil {
				return fmt.Errorf("failed to commit migration %s: %w", file, err)
			}
			logger.Info("applied migration", "file", file)
		} else if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", file, err)
		}
	}

	return nil
}
