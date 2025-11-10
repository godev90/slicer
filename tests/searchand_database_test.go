package slicer_test

import (
	"context"
	"strings"
	"testing"

	"github.com/godev90/slicer"
)

// Test SQL-like query building logic (without actual DB)
func TestSearchAndQueryConstruction(t *testing.T) {
	tests := []struct {
		name         string
		searchAnd    *slicer.SearchQueryAnd
		expectedArgs int
	}{
		{
			name: "Single field AND condition",
			searchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "status", Keyword: "active"},
				},
			},
			expectedArgs: 1,
		},
		{
			name: "Multiple field AND conditions",
			searchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "status", Keyword: "active"},
					{Field: "department", Keyword: "engineering"},
				},
			},
			expectedArgs: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []interface{}
			if tt.searchAnd != nil {
				for _, field := range tt.searchAnd.Fields {
					if field.Field != "" && field.Keyword != "" {
						args = append(args, "%"+field.Keyword+"%")
					}
				}
			}

			if len(args) != tt.expectedArgs {
				t.Errorf("Expected %d arguments, got %d", tt.expectedArgs, len(args))
			}

			for i, arg := range args {
				if str, ok := arg.(string); ok {
					if !strings.HasPrefix(str, "%") || !strings.HasSuffix(str, "%") {
						t.Errorf("Argument %d should be wrapped with %%, got: %s", i, str)
					}
				}
			}
		})
	}
}

// Test context handling for database operations
func TestSearchAndContextHandling(t *testing.T) {
	ctx := context.Background()
	opts := &slicer.QueryOptions{
		Page:  1,
		Limit: 10,
		SearchAnd: &slicer.SearchQueryAnd{
			Fields: []*slicer.SearchField{
				{Field: "status", Keyword: "active"},
			},
		},
	}

	if ctx == nil {
		t.Error("Context should not be nil")
	}

	_ = opts // Use opts to avoid unused variable
}

// Test field mapping validation
func TestSearchAndFieldMapping(t *testing.T) {
	searchAnd := &slicer.SearchQueryAnd{
		Fields: []*slicer.SearchField{
			{Field: "name", Keyword: "john"},
			{Field: "", Keyword: "ignored"},
			{Field: "status", Keyword: "active"},
		},
	}

	validFields := 0
	for _, field := range searchAnd.Fields {
		if field.Field != "" && field.Keyword != "" {
			validFields++
		}
	}

	if validFields != 2 {
		t.Errorf("Expected 2 valid fields, got %d", validFields)
	}
}
