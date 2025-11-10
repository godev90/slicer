package slicer_test

import (
	"testing"

	"github.com/godev90/slicer"
)

// TestSlicePaginatorBasic tests basic slice paginator functionality
func TestSlicePaginatorBasic(t *testing.T) {
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

	// Test NewSlicePaginator
	paginator := slicer.NewSlicePaginator(users, allowedFields)
	if paginator == nil {
		t.Fatal("NewSlicePaginator returned nil")
	}

	// Test Items() method (initially empty)
	items := paginator.Items()
	if items != nil && len(items) != 0 {
		t.Error("Items() should be initially empty")
	}

	// Test SetItems() method
	newUsers := []User{{ID: 99, Name: "Test"}}
	paginator.SetItems(newUsers)

	retrievedItems := paginator.Items()
	if len(retrievedItems) != 1 || retrievedItems[0].ID != 99 {
		t.Error("SetItems() did not work correctly")
	}

	// Reset paginator to original data for further testing
	paginator = slicer.NewSlicePaginator(users, allowedFields)

	// Test SlicePage with basic pagination
	opts := slicer.QueryOptions{
		Page:  1,
		Limit: 2,
	}

	result, err := slicer.SlicePage(paginator, opts)
	if err != nil {
		t.Fatalf("SlicePage returned error: %v", err)
	}

	if result.Page != 1 {
		t.Errorf("Expected page 1, got %d", result.Page)
	}

	if result.Limit != 2 {
		t.Errorf("Expected limit 2, got %d", result.Limit)
	}

	resultUsers, ok := result.Items.([]User)
	if !ok {
		t.Fatalf("Expected []User, got %T", result.Items)
	}

	if len(resultUsers) != 2 {
		t.Errorf("Expected 2 users (first page), got %d", len(resultUsers))
	}

	// Test with filters
	paginator = slicer.NewSlicePaginator(users, allowedFields)
	opts = slicer.QueryOptions{
		Page:  1,
		Limit: 10,
		Filters: map[string]string{
			"name": "Alice",
		},
	}

	result, err = slicer.SlicePage(paginator, opts)
	if err != nil {
		t.Fatalf("SlicePage with filters returned error: %v", err)
	}

	resultUsers, ok = result.Items.([]User)
	if !ok {
		t.Fatalf("Expected []User, got %T", result.Items)
	}

	// Should find Alice
	if len(resultUsers) != 1 {
		t.Errorf("Expected 1 user matching filter, got %d", len(resultUsers))
	}

	// Test with comparisons
	opts = slicer.QueryOptions{
		Page:  1,
		Limit: 10,
		Comparisons: []slicer.ComparisonFilter{
			{Field: "age", Op: "gte", Value: "30"},
		},
	}

	result, err = slicer.SlicePage(paginator, opts)
	if err != nil {
		t.Fatalf("SlicePage with comparisons returned error: %v", err)
	}

	resultUsers, ok = result.Items.([]User)
	if !ok {
		t.Fatalf("Expected []User, got %T", result.Items)
	}

	// Should find Bob and Charlie (age >= 30)
	if len(resultUsers) != 2 {
		t.Errorf("Expected 2 users matching age >= 30, got %d", len(resultUsers))
	}

	// Test with search
	opts = slicer.QueryOptions{
		Page:  1,
		Limit: 10,
		Search: &slicer.SearchQuery{
			Fields:  []string{"name"},
			Keyword: "Bob",
		},
	}

	result, err = slicer.SlicePage(paginator, opts)
	if err != nil {
		t.Fatalf("SlicePage with search returned error: %v", err)
	}

	resultUsers, ok = result.Items.([]User)
	if !ok {
		t.Fatalf("Expected []User, got %T", result.Items)
	}

	// Should find Bob
	if len(resultUsers) != 1 {
		t.Errorf("Expected 1 user matching search for Bob, got %d", len(resultUsers))
	}

	// Test with sorting
	opts = slicer.QueryOptions{
		Page:  1,
		Limit: 10,
		Sort: []slicer.SortField{
			{Field: "age", Desc: true}, // Sort by age descending
		},
	}

	result, err = slicer.SlicePage(paginator, opts)
	if err != nil {
		t.Fatalf("SlicePage with sort returned error: %v", err)
	}

	resultUsers, ok = result.Items.([]User)
	if !ok {
		t.Fatalf("Expected []User, got %T", result.Items)
	}

	// Should be sorted by age descending, so Charlie (35) should be first
	if len(resultUsers) > 0 && resultUsers[0].Name != "Charlie" {
		t.Errorf("Expected Charlie to be first (highest age), got %s", resultUsers[0].Name)
	}
}
