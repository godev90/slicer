package slicer_test

import (
	"testing"

	"github.com/godev90/slicer"
)

func TestSlicePaginatorBasicFunctions(t *testing.T) {
	// Test data
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	users := []User{
		{ID: 1, Name: "Alice", Age: 25},
		{ID: 2, Name: "Bob", Age: 30},
		{ID: 3, Name: "Charlie", Age: 35},
	}

	allowedFields := map[string]string{
		"id":   "id",
		"name": "name",
		"age":  "age",
	}

	t.Run("NewSlicePaginator", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		if paginator == nil {
			t.Error("NewSlicePaginator should not return nil")
		}

		// Check initial state
		items := paginator.Items()
		if items == nil {
			t.Errorf("Initial items should be empty slice, got nil")
		}
		if len(items) != 0 {
			t.Errorf("Initial items should be empty, got %d items", len(items))
		}
	})

	t.Run("SetItems and Items", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		// Test setting items
		testItems := []User{
			{ID: 100, Name: "Test", Age: 40},
		}

		paginator.SetItems(testItems)

		// Test getting items
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

	t.Run("SlicePage with empty options", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
		}

		pageData, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Errorf("SlicePage failed: %v", err)
		}

		if pageData.Total != 3 {
			t.Errorf("Expected total 3, got %d", pageData.Total)
		}

		if pageData.Page != 1 {
			t.Errorf("Expected page 1, got %d", pageData.Page)
		}

		if pageData.Limit != 10 {
			t.Errorf("Expected limit 10, got %d", pageData.Limit)
		}

		// Check if items were set
		items := paginator.Items()
		if len(items) != 3 {
			t.Errorf("Expected 3 items in paginator, got %d", len(items))
		}
	})

	t.Run("SlicePage with pagination", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		opts := slicer.QueryOptions{
			Page:  2,
			Limit: 2,
		}

		pageData, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Errorf("SlicePage failed: %v", err)
		}

		if pageData.Total != 3 {
			t.Errorf("Expected total 3, got %d", pageData.Total)
		}

		// Page 2 with limit 2 should give 1 item (Charlie)
		items := paginator.Items()
		if len(items) != 1 {
			t.Errorf("Expected 1 item on page 2, got %d", len(items))
		}

		if len(items) > 0 && items[0].Name != "Charlie" {
			t.Errorf("Expected Charlie on page 2, got %s", items[0].Name)
		}
	})

	t.Run("SlicePage with simple filter", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Filters: map[string]string{
				"name": "Alice",
			},
		}

		pageData, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Errorf("SlicePage failed: %v", err)
		}

		if pageData.Total != 1 {
			t.Errorf("Expected total 1, got %d", pageData.Total)
		}

		items := paginator.Items()
		if len(items) != 1 {
			t.Errorf("Expected 1 Alice, got %d items", len(items))
		}

		if len(items) > 0 && items[0].Name != "Alice" {
			t.Errorf("Expected Alice, got %s", items[0].Name)
		}
	})

	t.Run("SlicePage with empty slice", func(t *testing.T) {
		emptyUsers := []User{}
		paginator := slicer.NewSlicePaginator(emptyUsers, allowedFields)

		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
		}

		pageData, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Errorf("SlicePage failed: %v", err)
		}

		if pageData.Total != 0 {
			t.Errorf("Expected total 0, got %d", pageData.Total)
		}

		items := paginator.Items()
		if len(items) != 0 {
			t.Errorf("Expected 0 items, got %d", len(items))
		}
	})

	t.Run("SlicePage with invalid field filter", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Filters: map[string]string{
				"invalid_field": "value",
			},
		}

		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Errorf("SlicePage failed: %v", err)
		}

		// Should return all items since invalid field is ignored
		// Check via the paginator items since invalid fields are ignored
		items := paginator.Items()
		if len(items) != 3 {
			t.Errorf("Expected 3 items, got %d", len(items))
		}
	})
}
