package migrations

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Migration struct {
	Name string
	SQL  string
}

func Run(db *pgxpool.Pool, migrationsDir string) error {
	if err := ensureMigrationsTable(db); err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	applied, err := getApplied(db)
	if err != nil {
		return fmt.Errorf("get applied migrations: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("glob migrations: %w", err)
	}
	slices.Sort(files)

	for _, f := range files {
		name := filepath.Base(f)
		if slices.Contains(applied, name) {
			continue
		}

		log.Info().Str("name", name).Msg("applying migration")

		sql, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}

		if _, err := db.Exec(context.Background(), string(sql)); err != nil {
			return fmt.Errorf("execute %s: %w", name, err)
		}

		if _, err := db.Exec(context.Background(),
			`INSERT INTO schema_migrations (name, applied_at) VALUES ($1, $2)`, name, time.Now()); err != nil {
			return fmt.Errorf("record %s: %w", name, err)
		}

		log.Info().Str("name", name).Msg("migration applied")
	}

	return nil
}

func ensureMigrationsTable(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			name VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func getApplied(db *pgxpool.Pool) ([]string, error) {
	rows, err := db.Query(context.Background(),
		`SELECT name FROM schema_migrations ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

func SetupLogger() {
	level := zerolog.InfoLevel
	if os.Getenv("DEV_MODE") == "true" {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)
}
