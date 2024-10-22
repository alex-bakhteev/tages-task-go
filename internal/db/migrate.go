package db

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
)

func RunMigrations(databaseURL string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	sourceURL := "file://" + wd + "/internal/db/migrations"
	m, err := migrate.New(
		sourceURL,
		databaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database migrations successfully applied")
}
