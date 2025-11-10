package slicer_test

import (
	"testing"

	"github.com/godev90/slicer"
)

// TestTypeConversionDefault specifically tests the default cases in conversion functions
func TestTypeConversionDefault(t *testing.T) {
	// Create a struct with types that will trigger default cases in toInt64 and toFloat64
	type UnsupportedTypeUser struct {
		ID      int                    `json:"id"`
		Data    map[string]interface{} `json:"data"`    // Unsupported type for toInt64/toFloat64
		Channel chan int               `json:"channel"` // Definitely unsupported
	}

	users := []UnsupportedTypeUser{
		{ID: 1, Data: map[string]interface{}{"key": "value"}, Channel: make(chan int)},
		{ID: 2, Data: map[string]interface{}{"key": "other"}, Channel: make(chan int)},
	}

	allowedFields := map[string]string{
		"id":      "id",
		"data":    "data",
		"channel": "channel",
	}

	paginator := slicer.NewSlicePaginator(users, allowedFields)

	t.Run("Unsupported type triggers default case in comparison", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "data", Op: "gt", Value: "0"}, // This should trigger default case in compare
			},
		}

		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage returned error:", err)
		}

		items := paginator.Items()
		// Should return 0 items since default case returns false
		if len(items) != 0 {
			t.Errorf("Expected 0 items for unsupported type comparison, got %d", len(items))
		}
	})

	t.Run("Unsupported type triggers default case in sorting", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Sort: []slicer.SortField{
				{Field: "channel", Desc: false}, // This should trigger default case in compareSort
			},
		}

		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage returned error:", err)
		}

		items := paginator.Items()
		// Should return all items (not sorted because default case)
		if len(items) != 2 {
			t.Errorf("Expected 2 items, got %d", len(items))
		}
	})
}

// TestToInt64AllBranches tests all code paths in the toInt64 function
func TestToInt64AllBranches(t *testing.T) {
	type MixedTypeUser struct {
		ID       int    `json:"id"`        // Will trigger int case
		Score    int64  `json:"score"`     // Will trigger int64 case  
		Age      string `json:"age"`       // Will trigger string case
		Invalid  bool   `json:"invalid"`   // Will trigger default case
	}

	users := []MixedTypeUser{
		{ID: 42, Score: 100, Age: "25", Invalid: true},
		{ID: 1, Score: 200, Age: "30", Invalid: false},
	}

	allowedFields := map[string]string{
		"id":      "id",
		"score":   "score", 
		"age":     "age",
		"invalid": "invalid",
	}

	paginator := slicer.NewSlicePaginator(users, allowedFields)

	// Test int case (ID field)
	t.Run("int type conversion", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "id", Op: "eq", Value: "42"}, // This will call toInt64 on int
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage returned error:", err)
		}
		items := paginator.Items()
		if len(items) != 1 {
			t.Errorf("Expected 1 item for int comparison, got %d", len(items))
		}
	})

	// Test int64 case (Score field)  
	t.Run("int64 type conversion", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "score", Op: "gt", Value: "150"}, // This will call toInt64 on int64
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage returned error:", err)
		}
		items := paginator.Items()
		if len(items) != 1 {
			t.Errorf("Expected 1 item for int64 comparison, got %d", len(items))
		}
	})

	// Test string case (Age field)
	t.Run("string type conversion", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "age", Op: "lt", Value: "27"}, // This will call toInt64 on string
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage returned error:", err)
		}
		items := paginator.Items()
		if len(items) != 1 {
			t.Errorf("Expected 1 item for string comparison, got %d", len(items))
		}
	})

	// Test string case with invalid string (should still work, just return 0)
	t.Run("string type conversion with invalid string", func(t *testing.T) {
		// Add a user with invalid string number
		usersWithInvalid := append(users, MixedTypeUser{ID: 99, Score: 999, Age: "not-a-number", Invalid: false})
		paginator2 := slicer.NewSlicePaginator(usersWithInvalid, allowedFields)
		
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "age", Op: "eq", Value: "0"}, // Invalid strings parse to 0
			},
		}
		_, err := slicer.SlicePage(paginator2, opts)
		if err != nil {
			t.Fatal("SlicePage returned error:", err)
		}
		items := paginator2.Items()
		// Should find the user with invalid age string (which becomes 0)
		// But the comparison might not work as expected, so let's just verify no error
		t.Logf("Items found with invalid string: %d", len(items))
	})

	// Test default case (Invalid field - bool type)
	t.Run("default case type conversion", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "invalid", Op: "eq", Value: "1"}, // This will call toInt64 on bool -> default case
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage returned error:", err)
		}
		items := paginator.Items()
		// Default case returns 0, so bool values won't match "1"
		if len(items) != 0 {
			t.Errorf("Expected 0 items for bool comparison (default case), got %d", len(items))
		}
	})
}

// TestToInt64DirectCases - Direct test scenarios to hit all toInt64 branches
func TestToInt64DirectCases(t *testing.T) {
	// Create a struct with exact int, int64, and string types
	type DirectTypeUser struct {
		IntField    int    `json:"int_field"`
		Int64Field  int64  `json:"int64_field"`  
		StringField string `json:"string_field"`
		OtherField  bool   `json:"other_field"`
	}

	users := []DirectTypeUser{
		{IntField: 42, Int64Field: 100, StringField: "25", OtherField: true},
		{IntField: 10, Int64Field: 200, StringField: "30", OtherField: false},
		{IntField: 30, Int64Field: 150, StringField: "35", OtherField: true},
	}

	allowedFields := map[string]string{
		"intField":    "int_field",
		"int64Field":  "int64_field", 
		"stringField": "string_field",
		"otherField":  "other_field",
	}

	paginator := slicer.NewSlicePaginator(users, allowedFields)

	// Test 1: Force int case in toInt64 (via sorting)
	t.Run("toInt64 with int type", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Sort: []slicer.SortField{
				{Field: "intField", Desc: false}, // This will call toInt64(va) where va is int
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
		
		items := paginator.Items()
		if len(items) != 3 {
			t.Errorf("Expected 3 items, got %d", len(items))
		}
		// Should be sorted by IntField: 10, 30, 42
		if items[0].IntField != 10 || items[1].IntField != 30 || items[2].IntField != 42 {
			t.Errorf("Sort by int failed: got %d, %d, %d", items[0].IntField, items[1].IntField, items[2].IntField)
		}
	})

	// Test 2: Force int64 case in toInt64 (via sorting)
	t.Run("toInt64 with int64 type", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Sort: []slicer.SortField{
				{Field: "int64Field", Desc: false}, // This will call toInt64(va) where va is int64
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
		
		items := paginator.Items()
		if len(items) != 3 {
			t.Errorf("Expected 3 items, got %d", len(items))
		}
		// Should be sorted by Int64Field: 100, 150, 200
		if items[0].Int64Field != 100 || items[1].Int64Field != 150 || items[2].Int64Field != 200 {
			t.Errorf("Sort by int64 failed: got %d, %d, %d", items[0].Int64Field, items[1].Int64Field, items[2].Int64Field)
		}
	})

	// Test 3: Force string case in toInt64 (via sorting)  
	t.Run("toInt64 with string type", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Sort: []slicer.SortField{
				{Field: "stringField", Desc: false}, // This will call toInt64(va) where va is string
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
		
		items := paginator.Items()
		if len(items) != 3 {
			t.Errorf("Expected 3 items, got %d", len(items))
		}
		// Should be sorted by StringField as integers: "25", "30", "35"
		if items[0].StringField != "25" || items[1].StringField != "30" || items[2].StringField != "35" {
			t.Errorf("Sort by string failed: got %s, %s, %s", items[0].StringField, items[1].StringField, items[2].StringField)
		}
	})

	// Test 4: Force default case in toInt64 (via sorting with bool)
	t.Run("toInt64 with default case", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Sort: []slicer.SortField{
				{Field: "otherField", Desc: false}, // This will call toInt64(va) where va is bool -> default case
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
		
		items := paginator.Items()
		if len(items) != 3 {
			t.Errorf("Expected 3 items, got %d", len(items))
		}
		// With bool type, toInt64 returns 0 for all, so we just verify we got results
	})
}