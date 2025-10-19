package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func InitDB() {
	// Get environment variables or use defaults
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "password")
	dbname := getEnv("POSTGRES_DB", "events_db")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
	initSchema()
}

//go:embed sql/schema.sql
var schemaCreateSQL string

func initSchema() {
	_, err := DB.Exec(schemaCreateSQL)
	if err != nil {
		log.Fatal("Error creating events table:", err)
	}

	_, err = DB.Exec(schemaCreateSQL)
	if err != nil {
		log.Fatal("Error creating registrations table:", err)
	}

	fmt.Println("Tables created successfully!")
}
