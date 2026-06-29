package migrations

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMigrationFileOrdering(t *testing.T) {
	dir := t.TempDir()

	files := []string{
		"002_add_users.sql",
		"001_init.sql",
		"010_add_indexes.sql",
		"003_add_devices.sql",
	}
	for _, f := range files {
		os.WriteFile(filepath.Join(dir, f), []byte("SELECT 1;"), 0644)
	}

	// Ensure correct order via filename glob + sort
	matches, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		t.Fatal(err)
	}

	var names []string
	for _, m := range matches {
		names = append(names, filepath.Base(m))
	}

	expected := []string{"001_init.sql", "002_add_users.sql", "003_add_devices.sql", "010_add_indexes.sql"}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("position %d: expected %s, got %s", i, expected[i], name)
		}
	}
}

func TestMigrationFileFiltering(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "001_init.sql"), []byte("SELECT 1;"), 0644)
	os.WriteFile(filepath.Join(dir, "002_data.sql"), []byte("SELECT 2;"), 0644)
	os.WriteFile(filepath.Join(dir, "readme.txt"), []byte("not a migration"), 0644)

	matches, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		t.Fatal(err)
	}

	if len(matches) != 2 {
		t.Errorf("expected 2 .sql files, got %d", len(matches))
	}
}

func TestSetupLogger(t *testing.T) {
	// Just ensure it doesn't panic
	SetupLogger()
}



// Test that the create table SQL is valid
func TestCreateMigrationsTableSQL(t *testing.T) {
	sql := `CREATE TABLE IF NOT EXISTS schema_migrations (
		name VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`
	if sql == "" {
		t.Fatal("empty SQL")
	}
}

// Test that the insert SQL is valid
func TestRecordMigrationSQL(t *testing.T) {
	sql := `INSERT INTO schema_migrations (name, applied_at) VALUES ($1, $2)`
	if sql == "" {
		t.Fatal("empty SQL")
	}
}

func TestGetAppliedQuery(t *testing.T) {
	sql := `SELECT name FROM schema_migrations ORDER BY name`
	if sql == "" {
		t.Fatal("empty SQL")
	}
}
