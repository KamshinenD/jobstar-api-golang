package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	 _ "github.com/lib/pq"
	 "github.com/joho/godotenv"
)

var DB *sql.DB


func InitDB() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		panic("Could not connect to database")
	}

	DB.SetMaxOpenConns(10) //number of connections that can be opened simultaneausly
	DB.SetMaxIdleConns(5)  //no. of connections to be allowed when idle

	createTables()
}

func createTables() {
	// Enable the pgcrypto extension
	DB.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`)

	createUsersTable:=`
	CREATE TABLE IF NOT EXISTS users(
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		firstName TEXT NOT NULL,
		lastName TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		location TEXT NOT NULL
	)
	`
	_, err:= DB.Exec(createUsersTable)
	if err != nil {
		// panic("Could not create users table")
		log.Fatal("Could not create users table:", err)
	}
}