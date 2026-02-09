package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seedub/example-api/models"
)

func TestHealthCheckHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	HealthCheckHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}
}

func TestCreateItemHandler(t *testing.T) {
	// Reset store
	store = NewStore()

	item := models.Item{
		Name:        "Test Item",
		Description: "Test Description",
	}

	body, _ := json.Marshal(item)
	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	CreateItemHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var createdItem models.Item
	if err := json.NewDecoder(w.Body).Decode(&createdItem); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if createdItem.Name != item.Name {
		t.Errorf("Expected name '%s', got '%s'", item.Name, createdItem.Name)
	}

	if createdItem.ID == "" {
		t.Error("Expected ID to be generated")
	}
}

func TestCreateItemHandlerMissingName(t *testing.T) {
	// Reset store
	store = NewStore()

	item := models.Item{
		Description: "Test Description",
	}

	body, _ := json.Marshal(item)
	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	CreateItemHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetItemsHandler(t *testing.T) {
	// Reset store and add test data
	store = NewStore()
	testItem := &models.Item{
		ID:          "test-id",
		Name:        "Test Item",
		Description: "Test Description",
	}
	store.items["test-id"] = testItem

	req := httptest.NewRequest(http.MethodGet, "/api/items", nil)
	w := httptest.NewRecorder()

	GetItemsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var items []*models.Item
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items))
	}
}

func TestGetItemHandler(t *testing.T) {
	// Reset store and add test data
	store = NewStore()
	testItem := &models.Item{
		ID:          "test-id",
		Name:        "Test Item",
		Description: "Test Description",
	}
	store.items["test-id"] = testItem

	req := httptest.NewRequest(http.MethodGet, "/api/items/test-id", nil)
	w := httptest.NewRecorder()

	GetItemHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var item models.Item
	if err := json.NewDecoder(w.Body).Decode(&item); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if item.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", item.ID)
	}
}

func TestGetItemHandlerNotFound(t *testing.T) {
	// Reset store
	store = NewStore()

	req := httptest.NewRequest(http.MethodGet, "/api/items/nonexistent", nil)
	w := httptest.NewRecorder()

	GetItemHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUpdateItemHandler(t *testing.T) {
	// Reset store and add test data
	store = NewStore()
	testItem := &models.Item{
		ID:          "test-id",
		Name:        "Test Item",
		Description: "Test Description",
	}
	store.items["test-id"] = testItem

	updates := models.Item{
		Name:        "Updated Item",
		Description: "Updated Description",
	}

	body, _ := json.Marshal(updates)
	req := httptest.NewRequest(http.MethodPut, "/api/items/test-id", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	UpdateItemHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var updatedItem models.Item
	if err := json.NewDecoder(w.Body).Decode(&updatedItem); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if updatedItem.Name != updates.Name {
		t.Errorf("Expected name '%s', got '%s'", updates.Name, updatedItem.Name)
	}
}

func TestUpdateItemHandlerNotFound(t *testing.T) {
	// Reset store
	store = NewStore()

	updates := models.Item{
		Name: "Updated Item",
	}

	body, _ := json.Marshal(updates)
	req := httptest.NewRequest(http.MethodPut, "/api/items/nonexistent", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	UpdateItemHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteItemHandler(t *testing.T) {
	// Reset store and add test data
	store = NewStore()
	testItem := &models.Item{
		ID:          "test-id",
		Name:        "Test Item",
		Description: "Test Description",
	}
	store.items["test-id"] = testItem

	req := httptest.NewRequest(http.MethodDelete, "/api/items/test-id", nil)
	w := httptest.NewRecorder()

	DeleteItemHandler(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Verify item was deleted
	if _, exists := store.items["test-id"]; exists {
		t.Error("Expected item to be deleted")
	}
}

func TestDeleteItemHandlerNotFound(t *testing.T) {
	// Reset store
	store = NewStore()

	req := httptest.NewRequest(http.MethodDelete, "/api/items/nonexistent", nil)
	w := httptest.NewRecorder()

	DeleteItemHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
