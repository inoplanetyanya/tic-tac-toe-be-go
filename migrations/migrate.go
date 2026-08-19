package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %s", err.Error())
	}

	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dbsslmode := os.Getenv("DB_SSLMODE")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbuser, dbpassword, dbhost, dbport, dbname, dbsslmode)

	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)

	if err != nil {
		log.Fatalf("Error creating migration instance: %v", err)
	}

	defer func() {
		if _, err := m.Close(); err != nil {
			log.Printf("Error closing migration: %v", err)
		}
	}()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Error applying migrations: %v", err)
	}

	log.Println("Migrations successfully applied!")
}
