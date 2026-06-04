package dummydb 

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

// Global DB connection — accessible from any package
var DB *sql.DB

// Init connects to the database
func Init() {
	var err error
	DB, err = sql.Open("sqlite3", "./socialnetwork.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Database connected")
}