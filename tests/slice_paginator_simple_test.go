package slicer_test

import (
	"testing"

	"github.com/godev90/slicer"
)

// Simple test user structure for slice paginator tests
type SlicePaginatorUser struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Age    int    `json:"age"`
	Active bool   `json:"active"`
}

func TestSlicePaginatorSimple(t *testing.T) {
	// Test data
	users := []SlicePaginatorUser{
		{ID: 1, Name: "Alice", Email: "alice@test.com", Age: 25, Active: true},
		{ID: 2, Name: "Bob", Email: "bob@test.com", Age: 30, Active: true},
		{ID: 3, Name: "Charlie", Email: "charlie@test.com", Age: 35, Active: false},
		{ID: 4, Name: "Diana", Email: "diana@test.com", Age: 28, Active: true},
		{ID: 5, Name: "Eve", Email: "eve@test.com", Age: 32, Active: false},
	}

	allowedFields := map[string]string{
		"id":     "id",
		"name":   "name",
		"email":  "email",
		"age":    "age",
		"active": "active",
	}

	t.Run("Constructor", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)
		if paginator == nil {
			t.Fatal("NewSlicePaginator returned nil")
		}

		items := paginator.Items()
		if len(items) != 0 {
			t.Errorf("Expected empty items initially, got %d items", len(items))
		}
	})

	t.Run("SetItems and Items", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)
		testItems := users[:2]
		paginator.SetItems(testItems)

		items := paginator.Items()
		if len(items) != 2 {
			t.Errorf("Expected 2 items, got %d", len(items))
		}
		if items[0].Name != "Alice" || items[1].Name != "Bob" {
			t.Error("Items not set correctly")
		}
	})

	t.Run("Basic pagination", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 2,
		}

		result, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage returned error:", err)
		}

		if result.Total != 5 {
			t.Errorf("Expected total 5, got %d", result.Total)
		}
		if result.Page != 1 {
			t.Errorf("Expected page 1, got %d", result.Page)
		}
		if result.Limit != 2 {
			t.Errorf("Expected limit 2, got %d", result.Limit)
		}

		items := paginator.Items()
		if len(items) != 2 {
			t.Errorf("Expected 2 items on first page, got %d", len(items))
		}
	})

	t.Run("Comparison filters", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "age", Op: "gt", Value: "30"},
			},
		}

		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage returned error:", err)
		}

		items := paginator.Items()
		if len(items) != 2 {
			t.Errorf("Expected 2 users with age > 30, got %d", len(items))
		}
		
		// Verify the correct users are returned
		for _, item := range items {
			if item.Age <= 30 {
				t.Errorf("Found user with age %d, expected > 30", item.Age)
			}
		}
	})

	t.Run("Custom filters", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Filters: map[string]string{
				"active": "true",
			},
		}

		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage returned error:", err)
		}

		items := paginator.Items()
		if len(items) != 3 {
			t.Errorf("Expected 3 active users, got %d", len(items))
		}
		
		// Verify all returned users are active
		for _, item := range items {
			if !item.Active {
				t.Errorf("Found inactive user %s", item.Name)
			}
		}
	})

	t.Run("Search functionality", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Search: &slicer.SearchQuery{
				Fields:  []string{"name"},
				Keyword: "alice",
			},
		}

		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage returned error:", err)
		}

		items := paginator.Items()
		if len(items) != 1 {
			t.Errorf("Expected 1 user matching 'alice', got %d", len(items))
		}
		if len(items) > 0 && items[0].Name != "Alice" {
			t.Errorf("Expected Alice, got %s", items[0].Name)
		}
	})

	t.Run("SearchAnd functionality", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			SearchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "name", Keyword: "alice"},
				},
			},
		}

		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage returned error:", err)
		}

		items := paginator.Items()
		if len(items) != 1 {
			t.Errorf("Expected 1 user matching SearchAnd, got %d", len(items))
		}
		if len(items) > 0 && items[0].Name != "Alice" {
			t.Errorf("Expected Alice, got %s", items[0].Name)
		}
	})

	t.Run("Sorting functionality", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Sort: []slicer.SortField{
				{Field: "age", Desc: true},
			},
		}

		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage returned error:", err)
		}

		items := paginator.Items()
		if len(items) != 5 {
			t.Errorf("Expected 5 users, got %d", len(items))
		}
		
		// Verify sorted by age descending
		if len(items) > 0 && items[0].Name != "Charlie" {
			t.Errorf("Expected first user to be Charlie (oldest), got %s", items[0].Name)
		}
		
		// Check order is correct
		for i := 1; i < len(items); i++ {
			if items[i-1].Age < items[i].Age {
				t.Error("Items not sorted by age descending")
			}
		}
	})

	t.Run("Combined operations", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Filters: map[string]string{
				"active": "true",
			},
			Sort: []slicer.SortField{
				{Field: "age", Desc: false},
			},
		}

		result, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage returned error:", err)
		}

		items := paginator.Items()
		if len(items) != 3 {
			t.Errorf("Expected 3 active users, got %d", len(items))
		}
		if result.Total != 3 {
			t.Errorf("Expected total 3, got %d", result.Total)
		}
		
		// Check they're sorted by age ascending and all active
		for i := 0; i < len(items); i++ {
			if !items[i].Active {
				t.Errorf("Found inactive user %s", items[i].Name)
			}
			if i > 0 && items[i-1].Age > items[i].Age {
				t.Error("Items not sorted by age ascending")
			}
		}
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Run("Empty slice", func(t *testing.T) {
			emptyUsers := []SlicePaginatorUser{}
			paginator := slicer.NewSlicePaginator(emptyUsers, allowedFields)
			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
			}

			result, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			items := paginator.Items()
			if len(items) != 0 {
				t.Errorf("Expected 0 items for empty slice, got %d", len(items))
			}
			if result.Total != 0 {
				t.Errorf("Expected total 0 for empty slice, got %d", result.Total)
			}
		})

		t.Run("Page beyond available data", func(t *testing.T) {
			paginator := slicer.NewSlicePaginator(users, allowedFields)
			opts := slicer.QueryOptions{
				Page:  10,
				Limit: 2,
			}

			result, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			items := paginator.Items()
			if len(items) != 0 {
				t.Errorf("Expected 0 items for page beyond data, got %d", len(items))
			}
			if result.Total != 5 {
				t.Errorf("Total should still be 5, got %d", result.Total)
			}
		})

		t.Run("Invalid comparison values", func(t *testing.T) {
			paginator := slicer.NewSlicePaginator(users, allowedFields)
			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Comparisons: []slicer.ComparisonFilter{
					{Field: "age", Op: "gt", Value: "invalid_number"},
				},
			}

			_, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			items := paginator.Items()
			if len(items) != 0 {
				t.Errorf("Expected 0 items for invalid comparison value, got %d", len(items))
			}
		})
	})
}