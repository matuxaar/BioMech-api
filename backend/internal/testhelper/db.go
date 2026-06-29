package testhelper

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TestDB struct {
	Pool   *pgxpool.Pool
	Cancel func()
}

func StartDB(t *testing.T, migrationsDir string) *TestDB {
	db, err := startContainer(migrationsDir)
	if err != nil {
		t.Fatalf("failed to start test db: %v", err)
	}
	t.Cleanup(func() {
		db.Pool.Close()
		if db.Cancel != nil {
			db.Cancel()
		}
	})
	return db
}

func StartDBForMain(migrationsDir string) (*TestDB, error) {
	return startContainer(migrationsDir)
}

func startContainer(migrationsDir string) (*TestDB, error) {
	containerName := fmt.Sprintf("biomech-test-pg-%d", time.Now().UnixNano())

	exec.Command("docker", "rm", "-f", containerName).Run()

	cmd := exec.Command("docker", "run", "-d",
		"--name", containerName,
		"-e", "POSTGRES_USER=test",
		"-e", "POSTGRES_PASSWORD=test",
		"-e", "POSTGRES_DB=desertacia",
		"-p", "5433:5432",
		"postgres:17-alpine",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("start container: %w\n%s", err, out)
	}
	containerID := strings.TrimSpace(string(out))

	cancel := func() {
		exec.Command("docker", "rm", "-f", containerID).Run()
	}

	dsn := "postgres://test:test@localhost:5433/desertacia?sslmode=disable"

	ctx, cancelWait := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelWait()

	var pool *pgxpool.Pool
	for {
		select {
		case <-ctx.Done():
			cancel()
			return nil, fmt.Errorf("timeout waiting for postgres")
		default:
			var err error
			pool, err = pgxpool.New(context.Background(), dsn)
			if err == nil {
				if pingErr := pool.Ping(context.Background()); pingErr == nil {
					goto connected
				}
				pool.Close()
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
connected:

	if err := runMigrations(pool, migrationsDir); err != nil {
		pool.Close()
		cancel()
		return nil, fmt.Errorf("migrations: %w", err)
	}

	return &TestDB{Pool: pool, Cancel: cancel}, nil
}

func runMigrations(pool *pgxpool.Pool, dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return err
	}

	for _, f := range files {
		sql, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("read %s: %w", f, err)
		}
		if _, err := pool.Exec(context.Background(), string(sql)); err != nil {
			return fmt.Errorf("exec %s: %w", filepath.Base(f), err)
		}
	}

	return nil
}
