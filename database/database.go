package database

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func ConnectDB() (*sql.DB, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	// Get database credentials from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Construct the database connection string
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Connect to the database
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	// Ping the database to check if the connection is successful
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	
	fmt.Println("Db connected")
	
	return db, nil
}
