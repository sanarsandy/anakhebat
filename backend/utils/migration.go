package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

func ApplyMigration(db *sqlx.DB, migrationFile string) error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}

	// Construct path to migration file
	migrationPath := filepath.Join(cwd, "migrations", migrationFile)

	// Read migration file
	content, err := ioutil.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file at %s: %v", migrationPath, err)
	}

	// Execute migration
	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %v", err)
	}

	log.Printf("Migration %s applied successfully!", migrationFile)
	return nil
}
