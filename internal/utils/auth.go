package utils

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"shambashare/internal/database"
	"shambashare/internal/templates"
)

// SignupHandler handles user signup
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templates.RenderTemplate(w, "signup.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("first-name")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		password := r.FormValue("password")

		// Validate input
		if name == "" || email == "" || phone == "" || password == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// Log the attempt
		fmt.Printf("Signup attempt for email: %s, name: %s, phone: %s\n", email, name, phone)
		log.Printf("Signup attempt for email: %s, name: %s, phone: %s", email, name, phone)

		// Check if user already exists
		exists, err := database.CheckUserExists(email)
		if err != nil {
			http.Error(w, "Database check error: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Database check error: %v", err)
			return
		}
		if exists {
			http.Error(w, "User with this email already exists", http.StatusConflict)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			log.Printf("Password hash error: %v", err)
			return
		}

		// Insert user into the database
		err = database.CreateUser(name, email, phone, hashedPassword)
		if err != nil {
			http.Error(w, "Error saving user to the database: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Signup database insert error: %v", err)
			return
		}

		// Successful signup
		log.Printf("User %s signed up successfully", email)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templates.RenderTemplate(w, "login.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validate input
		if email == "" || password == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// Check if user exists and get hashed password
		hashedPassword, err := database.GetUserPasswordByEmail(email)
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			log.Println("Login query error:", err)
			return
		}

		// Compare passwords
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Successful login
		// TODO: Implement proper session management

		name, err := database.GetUserNamByEmail(email)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			log.Println("Username query error:", err)
			return
		}
		CurrentUser.FullName = name
		CurrentUser.Logo = string(name[0])

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}