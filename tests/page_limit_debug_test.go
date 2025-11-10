package slicer_test

import (
	"testing"

	"github.com/godev90/slicer"
)

// Test to isolate Page/Limit issue
func TestPageLimitIssue(t *testing.T) {
	type SimpleUser struct {
		Name string `json:"name"`
	}

	users := []SimpleUser{
		{Name: "Alice"},
		{Name: "Bob"},
	}

	fields := map[string]string{
		"name": "name",
	}

	paginator := slicer.NewSlicePaginator(users, fields)

	// Test 1: Empty QueryOptions (like failing tests)
	t.Run("Empty QueryOptions", func(t *testing.T) {
		result, err := slicer.SlicePage(paginator, slicer.QueryOptions{})
		if err != nil {
			t.Fatalf("SlicePage error: %v", err)
		}

		items, ok := result.Items.([]SimpleUser)
		if !ok {
			t.Fatalf("Expected []SimpleUser, got %T", result.Items)
		}

		t.Logf("Empty QueryOptions result: %d items", len(items))
		for i, user := range items {
			t.Logf("User %d: %+v", i, user)
		}
	})

	// Test 2: With Page/Limit specified (like working tests)
	t.Run("With Page and Limit", func(t *testing.T) {
		result, err := slicer.SlicePage(paginator, slicer.QueryOptions{
			Page:  1,
			Limit: 10,
		})
		if err != nil {
			t.Fatalf("SlicePage error: %v", err)
		}

		items, ok := result.Items.([]SimpleUser)
		if !ok {
			t.Fatalf("Expected []SimpleUser, got %T", result.Items)
		}

		t.Logf("With Page/Limit result: %d items", len(items))
		for i, user := range items {
			t.Logf("User %d: %+v", i, user)
		}
	})

	// Test 3: Empty QueryOptions with comparison
	t.Run("Empty QueryOptions with comparison", func(t *testing.T) {
		result, err := slicer.SlicePage(paginator, slicer.QueryOptions{
			Comparisons: []slicer.ComparisonFilter{
				{Field: "name", Op: slicer.EQ, Value: "Alice"},
			},
		})
		if err != nil {
			t.Fatalf("SlicePage error: %v", err)
		}

		items, ok := result.Items.([]SimpleUser)
		if !ok {
			t.Fatalf("Expected []SimpleUser, got %T", result.Items)
		}

		t.Logf("Empty QueryOptions + comparison result: %d items", len(items))
		for i, user := range items {
			t.Logf("User %d: %+v", i, user)
		}
	})

	// Test 4: With Page/Limit and comparison (like working tests)
	t.Run("With Page/Limit and comparison", func(t *testing.T) {
		result, err := slicer.SlicePage(paginator, slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "name", Op: slicer.EQ, Value: "Alice"},
			},
		})
		if err != nil {
			t.Fatalf("SlicePage error: %v", err)
		}

		items, ok := result.Items.([]SimpleUser)
		if !ok {
			t.Fatalf("Expected []SimpleUser, got %T", result.Items)
		}

		t.Logf("With Page/Limit + comparison result: %d items", len(items))
		for i, user := range items {
			t.Logf("User %d: %+v", i, user)
		}
	})
}
