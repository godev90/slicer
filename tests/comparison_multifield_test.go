package slicer_test

import (
	"testing"

	"github.com/godev90/slicer"
)

// Test to isolate the multi-field struct issue
func TestMultiFieldComparison(t *testing.T) {
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

	// Create paginator
	paginator := slicer.NewSlicePaginator(users, allowedFields)

	// Test simple equality comparison
	result, err := slicer.SlicePage(paginator, slicer.QueryOptions{
		Page:  1,
		Limit: 10,
		Comparisons: []slicer.ComparisonFilter{
			{Field: "name", Op: slicer.EQ, Value: "Alice"},
		},
	})

	if err != nil {
		t.Fatalf("SlicePage returned error: %v", err)
	}

	resultUsers, ok := result.Items.([]User)
	if !ok {
		t.Fatalf("Expected []User, got %T", result.Items)
	}

	t.Logf("Found %d users", len(resultUsers))
	for i, user := range resultUsers {
		t.Logf("User %d: %+v", i, user)
	}

	if len(resultUsers) != 1 {
		t.Errorf("Expected 1 user, got %d", len(resultUsers))
	}

	if len(resultUsers) > 0 && resultUsers[0].Name != "Alice" {
		t.Errorf("Expected Alice, got %s", resultUsers[0].Name)
	}
}