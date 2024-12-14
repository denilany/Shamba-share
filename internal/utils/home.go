package utils

import (
	"net/http"

	"shambashare/internal/templates"
)

type UserName struct {
	FullName string
	Logo string
}

var CurrentUser UserName

// HomeHandler handles the home page route
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	templates.RenderTemplate(w, "index.html", nil)
}

// FindLandHandler handles the findland page route
func FindLandHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/findland" {
		http.NotFound(w, r)
		return
	}

	templates.RenderTemplate(w, "findland.html", nil)
}

// DashboardHandler handles the dashboard page route
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/dashboard" {
		http.NotFound(w, r)
		return
	}

	templates.RenderTemplate(w, "dashboard.html", CurrentUser)
}
