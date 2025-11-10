package slicer_test

import (
	"testing"

	"github.com/godev90/slicer"
)

// TestDirectTypeConversion tests conversion functions by ensuring field types hit all conversion branches
func TestDirectTypeConversion(t *testing.T) {
	// Data with mixed types that should exercise all branches
	type MixedTypeData struct {
		IntField     int     `json:"int_field"`     // int case for toInt64
		Int64Field   int64   `json:"int64_field"`   // int64 case for toInt64  
		StringField  string  `json:"string_field"`  // string case for toInt64/toFloat64
		Float32Field float32 `json:"float32_field"` // float32 case for toFloat64
		Float64Field float64 `json:"float64_field"` // float64 case for toFloat64
		BoolField    bool    `json:"bool_field"`    // default case for both
	}

	data := []MixedTypeData{
		{
			IntField:     100,
			Int64Field:   200,
			StringField:  "300",
			Float32Field: 1.5,
			Float64Field: 2.5,
			BoolField:    true,
		},
		{
			IntField:     400,
			Int64Field:   500,
			StringField:  "600",
			Float32Field: 3.5,
			Float64Field: 4.5,
			BoolField:    false,
		},
	}

	allowedFields := map[string]string{
		"intField":     "int_field",
		"int64Field":   "int64_field",
		"stringField":  "string_field",
		"float32Field": "float32_field",
		"float64Field": "float64_field",
		"boolField":    "bool_field",
	}

	paginator := slicer.NewSlicePaginator(data, allowedFields)

	// Test int → int64 conversion path
	t.Run("int to int64 conversion", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "intField", Op: "gt", Value: "50"},
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal(err)
		}
	})

	// Test int64 → int64 (direct) path
	t.Run("int64 direct path", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "int64Field", Op: "lt", Value: "300"},
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal(err)
		}
	})

	// Test string → int64 conversion path
	t.Run("string to int64 conversion", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "stringField", Op: "eq", Value: "300"},
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal(err)
		}
	})

	// Test float32 → float64 conversion path
	t.Run("float32 to float64 conversion", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "float32Field", Op: "gt", Value: "1.0"},
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal(err)
		}
	})

	// Test float64 → float64 (direct) path
	t.Run("float64 direct path", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "float64Field", Op: "lt", Value: "5.0"},
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal(err)
		}
	})

	// Test string → float64 conversion path
	t.Run("string to float64 conversion", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "stringField", Op: "gt", Value: "200.5"},
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal(err)
		}
	})

	// Test default case for toInt64 (bool type)
	t.Run("bool triggers toInt64 default", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "boolField", Op: "gt", Value: "0"}, // bool gets converted via toInt64 default (returns 0)
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal(err)
		}
	})

	// Test default case for toFloat64 (bool type)
	t.Run("bool triggers toFloat64 default", func(t *testing.T) {
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "boolField", Op: "gt", Value: "0.5"}, // bool gets converted via toFloat64 default (returns 0.0)
			},
		}
		_, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatal(err)
		}
	})
}