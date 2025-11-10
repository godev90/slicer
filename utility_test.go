package slicer

import (
	"errors"
	"testing"
)

// Test utility functions to improve coverage
func TestUtilityFunctions(t *testing.T) {
	t.Run("SetValueSeparator", func(t *testing.T) {
		// Test that SetValueSeparator works by using it in parsing
		originalOpts := ParseOpts(map[string][]string{
			"sort": {"name,age"},
		})

		SetValueSeparator("|")

		// Test parsing with new separator
		newOpts := ParseOpts(map[string][]string{
			"sort": {"name|age"},
		})

		// Reset to default
		SetValueSeparator(",")

		// Verify the change worked
		if len(originalOpts.Sort) != 2 {
			t.Errorf("Expected 2 sort fields with comma separator, got %d", len(originalOpts.Sort))
		}
		if len(newOpts.Sort) != 2 {
			t.Errorf("Expected 2 sort fields with pipe separator, got %d", len(newOpts.Sort))
		}
	})

	t.Run("ErrorPage", func(t *testing.T) {
		testErr := errors.New("test error message")
		opts := QueryOptions{Page: 2, Limit: 5}

		err := ErrorPage(testErr, opts)

		if err.Items == nil {
			t.Error("ErrorPage should initialize Items")
		}

		if len(err.Items.([]string)) != 0 {
			t.Errorf("Expected empty items slice, got %v", err.Items)
		}

		if err.Total != 0 {
			t.Errorf("Expected total 0, got %d", err.Total)
		}

		if err.Page != 1 {
			t.Errorf("Expected page 1 (ErrorPage always returns page 1), got %d", err.Page)
		}

		if err.Limit != opts.Limit {
			t.Errorf("Expected limit %d, got %d", opts.Limit, err.Limit)
		}

		if err.LastError == nil {
			t.Error("Expected LastError to be set")
		}
	})

	t.Run("DefaultFilterByJson", func(t *testing.T) {
		// Test with a simple struct type
		type TestStruct struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
			City string `json:"city"`
		}

		result := DefaultFilterByJson[TestStruct]()

		expected := map[string]string{
			"name": "name",
			"age":  "age",
			"city": "city",
		}

		if len(result) != len(expected) {
			t.Errorf("Expected %d fields, got %d", len(expected), len(result))
		}

		for field, tag := range expected {
			if result[field] != tag {
				t.Errorf("Expected field %s to map to %s, got %s", field, tag, result[field])
			}
		}
	})
}
