package dummydb 

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

// Global DB connection — accessible from any package
var DB *sql.DB

// Init connects to the database
func Init() {

	wd, _ := os.Getwd()
log.Println("WORKING DIR:", wd)


	var err error
	DB, err = sql.Open("sqlite3", "./dummydb/socialnetwork.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Database connected")
}