package models

import (
	"fmt"
	"log"

	"jobstar.com/api/db"
	"jobstar.com/api/utils"
)

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Location  string `json:"location"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (u *User) Save() error {
	// Check if DB is initialized
	if db.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	query := `INSERT INTO users(firstName, lastName, email, password, location) 
	VALUES($1, $2, $3, $4, $5) RETURNING id`

	// Prepare and execute the query
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	// Use QueryRow to execute the query and retrieve the generated ID
	err = stmt.QueryRow(u.FirstName, u.LastName, u.Email, hashedPassword, u.Location).Scan(&u.ID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}

	return nil
}
