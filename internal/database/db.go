package database

import (
	"database/sql"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB sets up and returns a database connection
func InitDB() (*sql.DB, error) {
	dbPath, err := filepath.Abs("./users.db")
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Verify database connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Create the users table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			phone TEXT NOT NULL,
			password TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}

// CheckUserExists checks if a user with the given email already exists
func CheckUserExists(email string) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	return count > 0, err
}

// CreateUser inserts a new user into the database
func CreateUser(name, email, phone string, hashedPassword []byte) error {
	_, err := DB.Exec(
		"INSERT INTO users (name, email, phone, password) VALUES (?, ?, ?, ?)",
		name, email, phone, hashedPassword,
	)
	return err
}

// GetUserPasswordByEmail retrieves the hashed password for a given email
func GetUserPasswordByEmail(email string) (string, error) {
	var hashedPassword string
	err := DB.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&hashedPassword)
	return hashedPassword, err
}

// GetUserPasswordByEmail retrieves the hashed password for a given email
func GetUserNamByEmail(email string) (string, error) {
	var userName string
	err := DB.QueryRow("SELECT name FROM users WHERE email = ?", email).Scan(&userName)
	return userName, err
}
