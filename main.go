package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

// renderTemplate is a helper function to render HTML templates
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(filepath.Join("frontend/templates", tmpl))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// homeHandler handles the home page route
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	renderTemplate(w, "index.html", nil)
}

// signupHandler handles user signup
func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "signup.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		// Collect all form values
		name := r.FormValue("first-name")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		password := r.FormValue("password")

		fmt.Printf("Signup attempt for email: %s, name: %s, phone: %s", email, name, phone)

		// Validate input
		if name == "" || email == "" || phone == "" || password == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// Log the attempt
		log.Printf("Signup attempt for email: %s, name: %s, phone: %s", email, name, phone)

		// Check if user already exists
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
		if err != nil {
			http.Error(w, "Database check error: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Database check error: %v", err)
			return
		}
		if count > 0 {
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
		result, err := db.Exec(
			"INSERT INTO users (name, email, phone, password) VALUES (?, ?, ?, ?)",
			name, email, phone, hashedPassword,
		)
		if err != nil {
			http.Error(w, "Error saving user to the database: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Signup database insert error: %v", err)
			return
		}

		// Log the result of the insertion
		rowsAffected, _ := result.RowsAffected()
		log.Printf("User inserted successfully. Rows affected: %d", rowsAffected)

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// loginHandler handles user login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "login.html", nil)
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

		// Check if user exists
		var hashedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&hashedPassword)
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
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func findLandHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/findland" {
		http.NotFound(w, r)
		return
	}

	renderTemplate(w, "findland.html", nil)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/dashboard" {
		http.NotFound(w, r)
		return
	}

	renderTemplate(w, "dashboard.html", nil)
}

func leaselandHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/leaseland" {
		http.NotFound(w, r)
		return
	}

	renderTemplate(w, "leaseland.html", nil)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
		http.NotFound(w, r)
		return
	}

	renderTemplate(w, "about.html", nil)
}
func contactHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/contact" {
		http.NotFound(w, r)
		return
	}

	renderTemplate(w, "contact.html", nil)
}

func main() {
	// Set up logging to a file
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Determine the absolute path for the database
	dbPath, err := filepath.Abs("./users.db")
	if err != nil {
		log.Fatal("Failed to get absolute path:", err)
	}
	log.Printf("Database path: %s", dbPath)

	// Open the database with additional parameters
	db, err = sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Verify database connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
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
		log.Fatal("Failed to create table:", err)
	}

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Other handlers remain the same
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/findland", findLandHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/leaseland", leaselandHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/contact", contactHandler)

	// Start the server
	port := 8080
	fmt.Printf("Server starting on port %d...\n", port)
	log.Printf("Server starting on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil))
}
