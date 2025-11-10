package slicer_test

import (
	"testing"

	"github.com/godev90/slicer"
)

// Test comparison functions systematically to understand failure
func TestComparisonFunctionDebug(t *testing.T) {
	type SimpleUser struct {
		Name string `json:"name"`
	}

	users := []SimpleUser{
		{Name: "Alice"},
		{Name: "Bob"},
	}

	allowedFields := map[string]string{
		"name": "name",
	}

	paginator := slicer.NewSlicePaginator(users, allowedFields)

	// First verify basic functionality works
	t.Run("No filters - should return all", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
		}
		result, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatalf("SlicePage error: %v", err)
		}

		items, ok := result.Items.([]SimpleUser)
		if !ok {
			t.Fatalf("Expected []SimpleUser, got %T", result.Items)
		}

		if len(items) != 2 {
			t.Errorf("Expected 2 users, got %d", len(items))
		}
		t.Logf("Basic test passed: %+v", items)
	})

	// Test with comparison filter
	t.Run("EQ comparison", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "name", Op: slicer.EQ, Value: "Alice"},
			},
		}

		result, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatalf("SlicePage error: %v", err)
		}

		items, ok := result.Items.([]SimpleUser)
		if !ok {
			t.Fatalf("Expected []SimpleUser, got %T", result.Items)
		}

		t.Logf("Comparison test result: %+v", items)
		t.Logf("Expected 1 Alice, got %d items", len(items))

		if len(items) == 1 && items[0].Name == "Alice" {
			t.Log("SUCCESS: Comparison working!")
		} else {
			t.Errorf("FAILED: Expected 1 Alice, got %+v", items)
		}
	})

	// Test comparison with different operators
	t.Run("Different operators", func(t *testing.T) {
		operators := []struct {
			name     string
			op       slicer.ComparisonOp
			val      string
			expected int
		}{
			{"EQ Alice", slicer.EQ, "Alice", 1},
			{"EQ Bob", slicer.EQ, "Bob", 1},
			{"GT Alice", slicer.GT, "Alice", 1}, // Bob > Alice
			{"LT Bob", slicer.LT, "Bob", 1},     // Alice < Bob
		}

		for _, tt := range operators {
			t.Run(tt.name, func(t *testing.T) {
				opts := slicer.QueryOptions{
					Page:  1,
					Limit: 10,
					Comparisons: []slicer.ComparisonFilter{
						{Field: "name", Op: tt.op, Value: tt.val},
					},
				}

				result, err := slicer.SlicePage(paginator, opts)
				if err != nil {
					t.Fatalf("SlicePage error: %v", err)
				}

				items, ok := result.Items.([]SimpleUser)
				if !ok {
					t.Fatalf("Expected []SimpleUser, got %T", result.Items)
				}

				t.Logf("%s: got %d items, expected %d", tt.name, len(items), tt.expected)
				if len(items) != tt.expected {
					t.Errorf("%s failed: expected %d, got %d items: %+v", tt.name, tt.expected, len(items), items)
				}
			})
		}
	})
}
