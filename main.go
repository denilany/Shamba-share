package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"shambashare/internal/database"
	"shambashare/internal/utils"
)

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Serve static files
	fs := http.FileServer(http.Dir("frontend/templates"))
	http.Handle("/frontend/templates/", http.StripPrefix("/frontend/templates/", fs))

	// Set up routes
	http.HandleFunc("/", utils.HomeHandler)
	http.HandleFunc("/signup", utils.SignupHandler)
	http.HandleFunc("/login", utils.LoginHandler)
	http.HandleFunc("/findland", utils.FindLandHandler)
	http.HandleFunc("/dashboard", utils.DashboardHandler)
	http.HandleFunc("/landlease", utils.LandLeaseHandler)
	http.HandleFunc("/about", utils.AboutHandler)
	http.HandleFunc("/contact", utils.ContactHandler)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	log.Printf("Server running on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil))
}
