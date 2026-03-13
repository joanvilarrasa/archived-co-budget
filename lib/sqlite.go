package lib

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"

	_ "modernc.org/sqlite"
)

func OpenSQLite(databasePath string, migrationsPath string) (*sql.DB, error) {
	db, dbErr := sql.Open("sqlite", databasePath)
	if dbErr != nil {
		return nil, dbErr
	}

	if pingErr := db.Ping(); pingErr != nil {
		db.Close()
		return nil, pingErr
	}

	if migrateErr := migrate(db, migrationsPath); migrateErr != nil {
		db.Close()
		return nil, migrateErr
	}

	return db, nil
}

func migrate(db *sql.DB, migrationsPath string) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename TEXT PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return err
	}

	appliedMigrations, err := loadAppliedMigrations(db)
	if err != nil {
		return err
	}

	entries, entriesErr := os.ReadDir(migrationsPath)
	if entriesErr != nil {
		return entriesErr
	}

	migrationFiles := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".sql" {
			continue
		}

		migrationFiles = append(migrationFiles, entry.Name())
	}

	sort.Strings(migrationFiles)

	for _, migrationFile := range migrationFiles {
		if appliedMigrations[migrationFile] {
			continue
		}

		fullPath := filepath.Join(migrationsPath, migrationFile)
		migrationSQL, readErr := os.ReadFile(fullPath)
		if readErr != nil {
			return readErr
		}

		tx, txErr := db.Begin()
		if txErr != nil {
			return txErr
		}

		if _, execErr := tx.Exec(string(migrationSQL)); execErr != nil {
			tx.Rollback()
			return execErr
		}

		if _, execErr := tx.Exec(
			`INSERT INTO schema_migrations(filename) VALUES(?)`,
			migrationFile,
		); execErr != nil {
			tx.Rollback()
			return execErr
		}

		if commitErr := tx.Commit(); commitErr != nil {
			return commitErr
		}
	}

	return nil
}

func loadAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query(`SELECT filename FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	appliedMigrations := make(map[string]bool)
	for rows.Next() {
		var filename string
		if scanErr := rows.Scan(&filename); scanErr != nil {
			return nil, scanErr
		}

		appliedMigrations[filename] = true
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}

	return appliedMigrations, nil
}
