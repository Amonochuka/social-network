package sqlite

import (
	"database/sql"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func NewDB() *sql.DB {
	// Always use the database inside the backend folder
	dbPath, err := filepath.Abs("app.db")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		log.Fatal(err)
	}

	runMigrations(db)

	return db
}

func runMigrations(db *sql.DB) {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal("migrate driver: ", err)
	}

	migrationsPath, err := filepath.Abs("pkg/db/migrations/sqlite")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatal("migrate init: ", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("migrate up: ", err)
	}
}