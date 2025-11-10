package slicer_test

import (
	"testing"

	"github.com/godev90/slicer"
)

// TestInt64CaseSpecifically - test specifically designed to hit the int64 case in toInt64
func TestInt64CaseSpecifically(t *testing.T) {
	// Create a struct that explicitly uses int64 types
	type Int64User struct {
		ID       int64  `json:"id"`        // Explicitly int64
		Score    int64  `json:"score"`     // Explicitly int64
		Age      int64  `json:"age"`       // Explicitly int64
		Name     string `json:"name"`      // For control
		BadStr   string `json:"bad_str"`   // For testing string parse error
	}

	users := []Int64User{
		{ID: int64(1), Score: int64(100), Age: int64(25), Name: "Alice", BadStr: "not-a-number"},
		{ID: int64(2), Score: int64(200), Age: int64(30), Name: "Bob", BadStr: "invalid"},
		{ID: int64(3), Score: int64(150), Age: int64(35), Name: "Charlie", BadStr: "abc123"},
	}

	allowedFields := map[string]string{
		"id":     "id",
		"score":  "score",
		"age":    "age",
		"name":   "name",
		"badStr": "bad_str",
	}

	paginator := slicer.NewSlicePaginator(users, allowedFields)

	t.Run("int64 ID comparison", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "id", Op: "eq", Value: "2"}, // Should hit int64 case
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
		items := paginator.Items()
		if len(items) != 1 {
			t.Errorf("Expected 1 item for int64 ID comparison, got %d", len(items))
		}
	})

	t.Run("int64 score comparison", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "score", Op: "gt", Value: "125"}, // Should hit int64 case
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
		items := paginator.Items()
		if len(items) != 2 {
			t.Errorf("Expected 2 items for int64 score comparison, got %d", len(items))
		}
	})

	t.Run("int64 age sorting", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Sort: []slicer.SortField{
				{Field: "age", Desc: false}, // Should hit int64 case in sorting - CRITICAL for toInt64
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
		items := paginator.Items()
		if len(items) != 3 {
			t.Errorf("Expected 3 items for int64 age sorting, got %d", len(items))
		}
		// Should be sorted by age: Alice (25), Bob (30), Charlie (35)
		if len(items) >= 3 {
			if items[0].Age != 25 || items[1].Age != 30 || items[2].Age != 35 {
				t.Errorf("Items not sorted correctly by int64 age")
			}
		}
	})

	t.Run("int64 ID sorting specifically", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Sort: []slicer.SortField{
				{Field: "id", Desc: true}, // Should hit int64 case - this is the KEY test
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
		items := paginator.Items()
		if len(items) != 3 {
			t.Errorf("Expected 3 items for int64 ID sorting, got %d", len(items))
		}
		// Should be sorted by ID desc: Charlie (3), Bob (2), Alice (1)
		if len(items) >= 3 {
			if items[0].ID != 3 || items[1].ID != 2 || items[2].ID != 1 {
				t.Errorf("Items not sorted correctly by int64 ID desc")
			}
		}
	})

	t.Run("int64 score sorting specifically", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Sort: []slicer.SortField{
				{Field: "score", Desc: false}, // Should hit int64 case - another KEY test
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
		items := paginator.Items()
		if len(items) != 3 {
			t.Errorf("Expected 3 items for int64 score sorting, got %d", len(items))
		}
		// Should be sorted by score asc: Alice (100), Charlie (150), Bob (200)
		if len(items) >= 3 {
			if items[0].Score != 100 || items[1].Score != 150 || items[2].Score != 200 {
				t.Errorf("Items not sorted correctly by int64 score asc")
			}
		}
	})

	t.Run("int64 multiple field test", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "id", Op: "gt", Value: "1"}, // int64 comparison
				{Field: "score", Op: "lt", Value: "250"}, // int64 comparison
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
		items := paginator.Items()
		if len(items) != 2 {
			t.Errorf("Expected 2 items for multiple int64 comparisons, got %d", len(items))
		}
	})

	t.Run("string parsing in toInt64 with invalid strings", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Sort: []slicer.SortField{
				{Field: "badStr", Desc: false}, // This should hit string case in toInt64 with parsing errors
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
		items := paginator.Items()
		if len(items) != 3 {
			t.Errorf("Expected 3 items for string parsing test, got %d", len(items))
		}
	})
}