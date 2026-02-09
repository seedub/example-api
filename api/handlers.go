package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/seedub/example-api/models"
)

// Store holds the in-memory data
type Store struct {
	mu    sync.RWMutex
	items map[string]*models.Item
}

// NewStore creates a new Store
func NewStore() *Store {
	return &Store{
		items: make(map[string]*models.Item),
	}
}

var store = NewStore()

// HealthCheckHandler returns a simple health check
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// GetItemsHandler returns all items
func GetItemsHandler(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	items := make([]*models.Item, 0, len(store.items))
	for _, item := range store.items {
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// CreateItemHandler creates a new item
func CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if item.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	// Generate a simple ID
	item.ID = generateID()
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	store.items[item.ID] = &item

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// GetItemHandler returns a single item by ID
func GetItemHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/items/"):]
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	store.mu.RLock()
	defer store.mu.RUnlock()

	item, exists := store.items[id]
	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// UpdateItemHandler updates an existing item
func UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/items/"):]
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var updates models.Item
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	item, exists := store.items[id]
	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	if updates.Name != "" {
		item.Name = updates.Name
	}
	if updates.Description != "" {
		item.Description = updates.Description
	}
	item.UpdatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// DeleteItemHandler deletes an item
func DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/items/"):]
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.items[id]; !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	delete(store.items, id)

	w.WriteHeader(http.StatusNoContent)
}

// generateID generates a simple ID based on current time
func generateID() string {
	return time.Now().Format("20060102150405.000000")
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed in %v", time.Since(start))
	})
}
