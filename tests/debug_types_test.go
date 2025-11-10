package slicer_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/godev90/slicer"
)

// TestDebugFieldTypes - Debug what types are actually extracted from struct fields
func TestDebugFieldTypes(t *testing.T) {
	type DebugStruct struct {
		IntField    int     `json:"int_field"`
		Int64Field  int64   `json:"int64_field"`
		StringField string  `json:"string_field"`
		Float32     float32 `json:"float32_field"`
		Float64     float64 `json:"float64_field"`
	}

	data := []DebugStruct{
		{IntField: 42, Int64Field: 100, StringField: "25", Float32: 1.5, Float64: 2.5},
		{IntField: 10, Int64Field: 200, StringField: "30", Float32: 3.5, Float64: 4.5},
	}

	// Debug: Check what types Go reflection gives us
	val := reflect.ValueOf(data[0])
	typ := reflect.TypeOf(data[0])
	
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldValue := field.Interface()
		
		t.Logf("Field %s: declared_type=%v, reflect_type=%v, value=%v, kind=%v", 
			fieldType.Name, fieldType.Type, reflect.TypeOf(fieldValue), fieldValue, field.Kind())
	}

	allowedFields := map[string]string{
		"intField":    "int_field",
		"int64Field":  "int64_field",
		"stringField": "string_field",
		"float32":     "float32_field",
		"float64":     "float64_field",
	}

	paginator := slicer.NewSlicePaginator(data, allowedFields)

	// Test sorting each field to see what happens
	testCases := []struct {
		name      string
		fieldName string
	}{
		{"int field", "intField"},
		{"int64 field", "int64Field"},
		{"string field", "stringField"},
		{"float32 field", "float32"},
		{"float64 field", "float64"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Sort by %s", tc.name), func(t *testing.T) {
			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Sort: []slicer.SortField{
					{Field: tc.fieldName, Desc: false},
				},
			}
			_, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Logf("Error sorting by %s: %v", tc.fieldName, err)
			} else {
				t.Logf("Successfully sorted by %s", tc.fieldName)
			}
		})
	}
}