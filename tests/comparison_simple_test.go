package slicer_test

import (
	"testing"

	"github.com/godev90/slicer"
)

// Very simple test to isolate the comparison issue
func TestSimpleComparison(t *testing.T) {
	type User struct {
		Name string `json:"name"`
	}

	users := []User{
		{Name: "Alice"},
		{Name: "Bob"},
	}

	allowedFields := map[string]string{
		"name": "name",
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
