package slicer_test

import (
	"reflect"
	"testing"

	"github.com/godev90/slicer"
)

// TestFunctionCoverageComplete ensures 100% coverage of toInt64 and toFloat64 functions
func TestFunctionCoverageComplete(t *testing.T) {
	// Test data structures to exercise all type conversion branches
	type TestData struct {
		ID       int     `json:"id"`
		IntValue int     `json:"int_value"`
		Int64Val int64   `json:"int64_val"`
		StrInt   string  `json:"str_int"`
		Float32  float32 `json:"float32"`
		Float64  float64 `json:"float64"`
		StrFloat string  `json:"str_float"`
		Other    bool    `json:"other"` // For default case
	}

	testData := []TestData{
		{
			ID:       1,
			IntValue: 42,
			Int64Val: 12345,
			StrInt:   "678",
			Float32:  3.14,
			Float64:  2.718,
			StrFloat: "1.23",
			Other:    true,
		},
		{
			ID:       2,
			IntValue: 100,
			Int64Val: 54321,
			StrInt:   "999",
			Float32:  1.41,
			Float64:  9.876,
			StrFloat: "4.56",
			Other:    false,
		},
	}

	allowedFields := map[string]string{
		"id":        "id",
		"intValue":  "int_value",
		"int64Val":  "int64_val",
		"strInt":    "str_int",
		"float32":   "float32",
		"float64":   "float64",
		"strFloat":  "str_float",
		"other":     "other",
	}

	t.Run("toInt64 coverage - int type", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(testData, allowedFields)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "intValue", Op: "gt", Value: "30"}, // Will convert int to int64
			},
		}

		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
	})

	t.Run("toInt64 coverage - int64 type", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(testData, allowedFields)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "int64Val", Op: "lt", Value: "20000"}, // Will use int64 directly
			},
		}

		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
	})

	t.Run("toInt64 coverage - string type", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(testData, allowedFields)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "strInt", Op: "eq", Value: "678"}, // Will parse string to int64
			},
		}

		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
	})

	t.Run("toInt64 coverage - default case (bool)", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(testData, allowedFields)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "other", Op: "gt", Value: "0"}, // Bool will trigger default case
			},
		}

		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
	})

	t.Run("toFloat64 coverage - float32 type", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(testData, allowedFields)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "float32", Op: "gt", Value: "2.0"}, // Will convert float32 to float64
			},
		}

		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
	})

	t.Run("toFloat64 coverage - float64 type", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(testData, allowedFields)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "float64", Op: "lt", Value: "5.0"}, // Will use float64 directly
			},
		}

		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
	})

	t.Run("toFloat64 coverage - string type", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(testData, allowedFields)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "strFloat", Op: "eq", Value: "1.23"}, // Will parse string to float64
			},
		}

		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
	})

	t.Run("toFloat64 coverage - default case (bool)", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(testData, allowedFields)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "other", Op: "gt", Value: "0.5"}, // Bool will trigger default case in toFloat64
			},
		}

		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal("SlicePage failed:", err)
		}
	})

	t.Run("Debug function calls - check what types we're actually hitting", func(t *testing.T) {
		// Let's examine what types our struct fields actually are
		val := reflect.ValueOf(testData[0])
		typ := reflect.TypeOf(testData[0])
		
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := typ.Field(i)
			t.Logf("Field %s: type=%v, value=%v", fieldType.Name, field.Type(), field.Interface())
		}
	})
}