package models

import (
	"errors"
	"fmt"
	"log"

	"jobstar.com/api/db"
	"jobstar.com/api/utils"
)

type User struct {
	ID                string `json:"id"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	Location          string `json:"location"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	IsAdmin           bool   `json:"isAdmin"`
	VerificationToken string `json:"verificationToken"`
	IsVerified        bool   `json:"isVerified"`
}

type UserLogin struct {
	ID       string `json:"id"` //id is skipped so that swagger does notpick it
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdate struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Location  string `json:"location"`
}

func (u *User) Save() error {
	// Check if DB is initialized
	if db.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	query := `INSERT INTO users(firstName, lastName, email, password, location, verificationToken ) 
	VALUES($1, $2, $3, $4, $5, $6) RETURNING id`

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
	err = stmt.QueryRow(u.FirstName, u.LastName, u.Email, hashedPassword, u.Location, u.VerificationToken).Scan(&u.ID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}

	return nil
}

func (u *UserLogin) ValidateCredentials() error {
	query := "SELECT id, password, isVerified FROM users WHERE email=$1"

	var retrievedPassword string
	var isVerified bool
	row := db.DB.QueryRow(query, u.Email)
	// err := db.DB.QueryRow(query, u.Email).Scan(&u.ID, &retrievedPassword) //binding the password. Weare also binding the UserID so that we can acccess it to generate jwt token during login
	err := row.Scan(&u.ID, &retrievedPassword, &isVerified) //binding the password. We are also binding the UserID so that we can acccess it to generate jwt token during login

	if err != nil {
		fmt.Println(err)
		return err
	}

	if !isVerified {
		return errors.New("account not verified")
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)

	if !passwordIsValid {
		return errors.New("credentials Invalid")
	}

	return nil
}

func (u *UserUpdate) Update() error {
	// Check if DB is initialized
	if db.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	query := `UPDATE users SET firstName=$1, lastName=$2, location=$3 WHERE id=$4`

	// Prepare the query
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	// Execute the query
	_, err = stmt.Exec(u.FirstName, u.LastName, u.Location, u.ID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}

	return nil
}

func Verify(email, token string) error {
	// Check if DB is initialized
	if db.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	// First, retrieve the stored verification token from the database for the given email
	var storedToken string
	query := `SELECT verificationToken FROM users WHERE email=$1`

	err := db.DB.QueryRow(query, email).Scan(&storedToken)

	if err != nil {
		return err
	}

	// Check if the provided token matches the stored token
	if token != storedToken {
		return fmt.Errorf("invalid verification token")
	}

	// If the token is valid, update the user's verification status
	updateQuery := `UPDATE users SET verificationToken=$1, isVerified=$2 WHERE email=$3`

	// Prepare the query
	stmt, err := db.DB.Prepare(updateQuery)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	// Execute the query
	_, err = stmt.Exec("", true, email)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}

	return nil
}
