package models

import (
	"testing"
	"time"
)

func TestItemCreation(t *testing.T) {
	item := Item{
		ID:          "test-id",
		Name:        "Test Item",
		Description: "Test Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if item.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", item.ID)
	}

	if item.Name != "Test Item" {
		t.Errorf("Expected Name 'Test Item', got '%s'", item.Name)
	}

	if item.Description != "Test Description" {
		t.Errorf("Expected Description 'Test Description', got '%s'", item.Description)
	}
}
