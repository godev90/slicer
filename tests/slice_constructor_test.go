package slicer_test

import (
	"testing"

	"github.com/godev90/slicer"
)

// Test slice paginator constructor and basic functions
func TestSlicePaginatorConstructor(t *testing.T) {
	t.Run("NewSlicePaginator", func(t *testing.T) {
		type User struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		users := []User{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
		}

		allowedFields := map[string]string{
			"id":   "id",
			"name": "name",
		}

		paginator := slicer.NewSlicePaginator(users, allowedFields)

		if paginator == nil {
			t.Fatal("NewSlicePaginator should not return nil")
		}

		// Test initial state
		items := paginator.Items()
		if items == nil {
			t.Errorf("Initial items should be empty slice, got nil")
		}
		if len(items) != 0 {
			t.Errorf("Initial items should be empty, got %d items", len(items))
		}
	})

	t.Run("Items and SetItems", func(t *testing.T) {
		type User struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		users := []User{
			{ID: 1, Name: "Alice"},
		}

		allowedFields := map[string]string{
			"id":   "id",
			"name": "name",
		}

		paginator := slicer.NewSlicePaginator(users, allowedFields)

		// Test SetItems
		testItems := []User{
			{ID: 100, Name: "Test"},
		}

		paginator.SetItems(testItems)

		// Test Items
		retrievedItems := paginator.Items()
		if len(retrievedItems) != 1 {
			t.Errorf("Expected 1 item, got %d", len(retrievedItems))
		}

		if retrievedItems[0].ID != 100 {
			t.Errorf("Expected ID 100, got %d", retrievedItems[0].ID)
		}

		if retrievedItems[0].Name != "Test" {
			t.Errorf("Expected name 'Test', got %s", retrievedItems[0].Name)
		}
	})
}
