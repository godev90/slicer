package slicer_test

import (
	"testing"

	"github.com/godev90/slicer"
)

// Minimal test to isolate comparison bug
func TestMinimalComparison(t *testing.T) {
	type SimpleUser struct {
		Name string `json:"name"`
	}

	users := []SimpleUser{
		{Name: "Alice"},
		{Name: "Bob"},
	}

	// Create paginator with field mapping
	fields := map[string]string{
		"name": "name",
	}

	paginator := slicer.NewSlicePaginator(users, fields)

	// Test most basic case - no filters first
	t.Run("No filters - should return all", func(t *testing.T) {
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

		t.Logf("No filters result: %d items", len(items))
		if len(items) != 2 {
			t.Errorf("Expected 2 users, got %d", len(items))
		}
	})

	// Test comparison
	t.Run("Single EQ comparison", func(t *testing.T) {
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

		t.Logf("EQ comparison result: %d items", len(items))
		for i, user := range items {
			t.Logf("User %d: %+v", i, user)
		}

		if len(items) != 1 {
			t.Errorf("Expected 1 user, got %d", len(items))
		}
		if len(items) > 0 && items[0].Name != "Alice" {
			t.Errorf("Expected Alice, got %s", items[0].Name)
		}
	})
}

// TestToInt64Simple - direct test to ensure we hit all branches
func TestToInt64Simple(t *testing.T) {
	// Create test data with all required types
	type AllTypesStruct struct {
		IntVal    int    `json:"int_val"`
		Int64Val  int64  `json:"int64_val"`
		StringVal string `json:"string_val"`
		BoolVal   bool   `json:"bool_val"`
	}

	data := []AllTypesStruct{
		{IntVal: 42, Int64Val: 1000, StringVal: "123", BoolVal: true},
		{IntVal: 50, Int64Val: 2000, StringVal: "456", BoolVal: false},
	}

	allowedFields := map[string]string{
		"intVal":    "int_val",
		"int64Val":  "int64_val", 
		"stringVal": "string_val",
		"boolVal":   "bool_val",
	}

	paginator := slicer.NewSlicePaginator(data, allowedFields)

	// Test 1: int case  
	opts1 := slicer.QueryOptions{
		Page:  1,
		Limit: 10,
		Comparisons: []slicer.ComparisonFilter{
			{Field: "intVal", Op: "gt", Value: "40"},
		},
	}
	_, _ = slicer.SlicePage(paginator, opts1)

	// Test 2: int64 case
	opts2 := slicer.QueryOptions{
		Page:  1,
		Limit: 10,
		Comparisons: []slicer.ComparisonFilter{
			{Field: "int64Val", Op: "lt", Value: "1500"},
		},
	}
	_, _ = slicer.SlicePage(paginator, opts2)

	// Test 3: string case
	opts3 := slicer.QueryOptions{
		Page:  1,
		Limit: 10,
		Comparisons: []slicer.ComparisonFilter{
			{Field: "stringVal", Op: "eq", Value: "123"},
		},
	}
	_, _ = slicer.SlicePage(paginator, opts3)

	// Test 4: default case (bool)
	opts4 := slicer.QueryOptions{
		Page:  1,
		Limit: 10,
		Comparisons: []slicer.ComparisonFilter{
			{Field: "boolVal", Op: "gt", Value: "0"},
		},
	}
	_, _ = slicer.SlicePage(paginator, opts4)
}
