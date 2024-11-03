package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	defaultVersion = -1
)

func main() {
	version, force := parseFlags()
	migrationsPath, databaseURL := parseEnv()

	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		log.Fatalf("Could not initialize migrations: %v", err)
	}
	defer m.Close()

	if err := applyMigrations(m, version, force); err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	log.Println("Migrations completed successfully!")
}

func parseEnv() (string, string) {
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		log.Fatalf("Missing MIGRATIONS_PATH environment variable")
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatalf("Missing DATABASE_URL environment variable")
	}
	return migrationsPath, databaseURL
}

func parseFlags() (int, bool) {
	version := flag.Int("version", defaultVersion, "Version migrations to apply")
	force := flag.Bool("force", false, "Force change migrations version")

	flag.Parse()

	if *version < -1 {
		log.Fatalf("invalid version: %d", *version)
	}

	return *version, *force
}

func applyMigrations(m *migrate.Migrate, version int, force bool) error {
	currentVersion, _, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("could not retrieve current version: %v", err)
	}

	if version == defaultVersion {
		return m.Up()
	}

	steps := version - int(currentVersion)
	if currentVersion == uint(version) {
		log.Println("All migrations have already been applied.")
		return nil
	}

	if force {
		if err := m.Force(version); err != nil {
			return fmt.Errorf("failed to force set migration version: %v", err)
		}
	} else {
		if err := m.Steps(steps); err != nil {
			return fmt.Errorf("could not apply migrations: %v", err)
		}
		log.Printf("Applied %d migrations successfully.", steps)
	}

	return nil
}
