package slicer_test

import (
	"testing"

	"github.com/godev90/slicer"
)

// Simple test to debug comparison functionality
func TestComparisonDebug(t *testing.T) {
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	users := []User{
		{ID: 1, Name: "Alice", Age: 25},
		{ID: 2, Name: "Bob", Age: 30},
	}

	allowedFields := map[string]string{
		"id":   "id",
		"name": "name",
		"age":  "age",
	}

	paginator := slicer.NewSlicePaginator(users, allowedFields)

	// Test basic functionality first
	t.Run("Basic pagination", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
		}

		result, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatalf("SlicePage returned error: %v", err)
		}

		resultUsers, ok := result.Items.([]User)
		if !ok {
			t.Fatalf("Expected []User, got %T", result.Items)
		}

		if len(resultUsers) != 2 {
			t.Errorf("Expected 2 users, got %d", len(resultUsers))
		}
		t.Logf("Basic pagination works: %+v", resultUsers)
	})

	// Test comparison functionality
	t.Run("String comparison EQ", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "name", Op: slicer.EQ, Value: "Alice"},
			},
		}

		// Log what we're testing
		t.Logf("Testing comparison: field=%s, op=%v, value=%s", "name", slicer.EQ, "Alice")
		t.Logf("Allowed fields: %+v", allowedFields)

		for _, cmp := range opts.Comparisons {
			t.Logf("Processing comparison: field=%s, op=%v, value=%s", cmp.Field, cmp.Op, cmp.Value)
		}

		result, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatalf("SlicePage returned error: %v", err)
		}

		resultUsers, ok := result.Items.([]User)
		if !ok {
			t.Fatalf("Expected []User, got %T", result.Items)
		}

		t.Logf("Comparison result: %+v", resultUsers)
		t.Logf("Total found: %d", len(resultUsers))

		// Let's also check field mapping manually
		t.Logf("Manual field check:")
		for i, user := range users {
			t.Logf("User %d: %+v", i, user)
			if user.Name == "Alice" {
				t.Logf("Found Alice manually!")
			}
		}

		if len(resultUsers) != 1 {
			t.Errorf("Expected 1 user, got %d", len(resultUsers))
		} else if resultUsers[0].Name != "Alice" {
			t.Errorf("Expected Alice, got %s", resultUsers[0].Name)
		}
	})

	// Test different comparison
	t.Run("String comparison GT", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "name", Op: slicer.GT, Value: "Alice"},
			},
		}

		result, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatalf("SlicePage returned error: %v", err)
		}

		resultUsers, ok := result.Items.([]User)
		if !ok {
			t.Fatalf("Expected []User, got %T", result.Items)
		}

		t.Logf("GT comparison result: %+v", resultUsers)
		if len(resultUsers) != 1 {
			t.Errorf("Expected 1 user (Bob > Alice), got %d", len(resultUsers))
		} else if resultUsers[0].Name != "Bob" {
			t.Errorf("Expected Bob, got %s", resultUsers[0].Name)
		}
	})
}
