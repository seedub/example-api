package api

import (
	"net/http"
	"strings"
)

// NewRouter creates and configures the HTTP router
func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", HealthCheckHandler)

	// API endpoints
	mux.HandleFunc("/api/items", handleItems)
	mux.HandleFunc("/api/items/", handleItemByID)

	// Apply middleware
	return CORSMiddleware(LoggingMiddleware(mux))
}

// handleItems routes to the appropriate handler based on HTTP method
func handleItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetItemsHandler(w, r)
	case http.MethodPost:
		CreateItemHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleItemByID routes to the appropriate handler based on HTTP method
func handleItemByID(w http.ResponseWriter, r *http.Request) {
	// Check if there's an ID in the path
	id := strings.TrimPrefix(r.URL.Path, "/api/items/")
	if id == "" || id == "/" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		GetItemHandler(w, r)
	case http.MethodPut:
		UpdateItemHandler(w, r)
	case http.MethodDelete:
		DeleteItemHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
