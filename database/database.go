package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// DB variable holds database property
var DB *sql.DB

// DBerr is the error of the database
var DBerr error

// InitDB initialized the database and creates a WEB_URL table
func InitDB() (*sql.DB, error) {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	var err error

	databaseURI := os.Getenv("DATABASE_URL")

	DB, DBerr = sql.Open("postgres", databaseURI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	stmt, err := DB.Prepare(`CREATE TABLE IF NOT EXISTS web_url(
		id SERIAL,
		original_url varchar not null,
		short_url varchar PRIMARY KEY not null,
		created_at time not null,
		expires_at time,
		has_expired bool default false);`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	_, err = stmt.Exec()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	DB.Exec("create unique index shortUrl_original_index on web_url(original_url);")
	DB.Exec("create unique index shortUrl_short_index on web_url(short_url);")

	return DB, nil

}

// TestPing test function for the package database
func TestPing() {

	// Pings the global database
	pingError := DB.Ping()

	if pingError != nil {
		// An error was returned while pinging the database
		log.Fatal(pingError)
	} else {
		// Database Ping successful
		log.Println("Database: Ping successful.")
	}

}
