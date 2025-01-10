package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"shambashare/internal/database"
	"shambashare/internal/utils"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Create a new router
	r := mux.NewRouter()

	// Serve static files
	fs := http.FileServer(http.Dir("frontend/templates"))
	r.PathPrefix("/frontend/templates/").Handler(http.StripPrefix("/frontend/templates/", fs))

	// Set up routes
	r.HandleFunc("/", utils.HomeHandler).Methods("GET")
	r.HandleFunc("/signup", utils.SignupHandler).Methods("GET", "POST")
	r.HandleFunc("/login", utils.LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/findland", utils.FindLandHandler).Methods("GET")
	r.HandleFunc("/dashboard", utils.DashboardHandler).Methods("GET")
	r.HandleFunc("/leaseland", utils.LandLeaseHandler).Methods("GET", "POST")
	r.HandleFunc("/viewlandforlease", utils.ViewLandLeaseHandler).Methods("GET")
	r.HandleFunc("/startleasing", utils.StartLeaseHandler).Methods("POST")
	r.HandleFunc("/about", utils.AboutHandler).Methods("GET")
	r.HandleFunc("/contact", utils.ContactHandler).Methods("GET", "POST")

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
