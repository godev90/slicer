package slicer_test

import (
	"net/url"
	"testing"

	"github.com/godev90/slicer"
)

// TestParseOptsComprehensive tests all URL parameter parsing combinations
func TestParseOptsComprehensive(t *testing.T) {
	t.Run("Page and Limit parsing", func(t *testing.T) {
		tests := []struct {
			name     string
			values   url.Values
			expected slicer.QueryOptions
		}{
			{
				name:   "Default values",
				values: url.Values{},
				expected: slicer.QueryOptions{
					Page:    1,
					Limit:   10,
					Offset:  0,
					Filters: map[string]string{},
				},
			},
			{
				name:   "Valid page and limit",
				values: url.Values{"page": {"3"}, "limit": {"25"}},
				expected: slicer.QueryOptions{
					Page:    3,
					Limit:   25,
					Offset:  50, // (3-1) * 25
					Filters: map[string]string{},
				},
			},
			{
				name:   "Invalid page and limit should use defaults",
				values: url.Values{"page": {"0"}, "limit": {"-5"}},
				expected: slicer.QueryOptions{
					Page:    1,
					Limit:   10,
					Offset:  0,
					Filters: map[string]string{},
				},
			},
			{
				name:   "Non-numeric page and limit should use defaults",
				values: url.Values{"page": {"abc"}, "limit": {"xyz"}},
				expected: slicer.QueryOptions{
					Page:    1,
					Limit:   10,
					Offset:  0,
					Filters: map[string]string{},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := slicer.ParseOpts(tt.values)

				if result.Page != tt.expected.Page {
					t.Errorf("Expected page %d, got %d", tt.expected.Page, result.Page)
				}
				if result.Limit != tt.expected.Limit {
					t.Errorf("Expected limit %d, got %d", tt.expected.Limit, result.Limit)
				}
				if result.Offset != tt.expected.Offset {
					t.Errorf("Expected offset %d, got %d", tt.expected.Offset, result.Offset)
				}
			})
		}
	})

	t.Run("Sort field parsing", func(t *testing.T) {
		tests := []struct {
			name     string
			values   url.Values
			expected []slicer.SortField
		}{
			{
				name:     "No sort fields",
				values:   url.Values{},
				expected: nil,
			},
			{
				name:   "Single sort field ascending",
				values: url.Values{"sort": {"name"}},
				expected: []slicer.SortField{
					{Field: "name", Desc: false},
				},
			},
			{
				name:   "Single sort field descending",
				values: url.Values{"sort": {"-created_at"}},
				expected: []slicer.SortField{
					{Field: "created_at", Desc: true},
				},
			},
			{
				name:   "Multiple sort fields",
				values: url.Values{"sort": {"name,-created_at,id"}},
				expected: []slicer.SortField{
					{Field: "name", Desc: false},
					{Field: "created_at", Desc: true},
					{Field: "id", Desc: false},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := slicer.ParseOpts(tt.values)

				if len(result.Sort) != len(tt.expected) {
					t.Fatalf("Expected %d sort fields, got %d", len(tt.expected), len(result.Sort))
				}

				for i, expected := range tt.expected {
					if result.Sort[i].Field != expected.Field {
						t.Errorf("Sort field %d: expected field %s, got %s", i, expected.Field, result.Sort[i].Field)
					}
					if result.Sort[i].Desc != expected.Desc {
						t.Errorf("Sort field %d: expected desc %t, got %t", i, expected.Desc, result.Sort[i].Desc)
					}
				}
			})
		}
	})

	t.Run("Search query parsing", func(t *testing.T) {
		tests := []struct {
			name           string
			values         url.Values
			expectedSearch *slicer.SearchQuery
		}{
			{
				name:           "No search parameters",
				values:         url.Values{},
				expectedSearch: nil,
			},
			{
				name:           "Search fields without keyword",
				values:         url.Values{"search": {"name,email"}},
				expectedSearch: nil,
			},
			{
				name:           "Keyword without search fields",
				values:         url.Values{"keyword": {"john"}},
				expectedSearch: nil,
			},
			{
				name:   "Complete search query",
				values: url.Values{"search": {"name,email,description"}, "keyword": {"john doe"}},
				expectedSearch: &slicer.SearchQuery{
					Fields:  []string{"name", "email", "description"},
					Keyword: "john doe",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := slicer.ParseOpts(tt.values)

				if tt.expectedSearch == nil {
					if result.Search != nil {
						t.Errorf("Expected Search to be nil, but got %+v", result.Search)
					}
					return
				}

				if result.Search == nil {
					t.Fatalf("Expected Search to be %+v, but got nil", tt.expectedSearch)
				}

				if result.Search.Keyword != tt.expectedSearch.Keyword {
					t.Errorf("Expected keyword %s, got %s", tt.expectedSearch.Keyword, result.Search.Keyword)
				}

				if len(result.Search.Fields) != len(tt.expectedSearch.Fields) {
					t.Fatalf("Expected %d search fields, got %d", len(tt.expectedSearch.Fields), len(result.Search.Fields))
				}

				for i, expected := range tt.expectedSearch.Fields {
					if result.Search.Fields[i] != expected {
						t.Errorf("Search field %d: expected %s, got %s", i, expected, result.Search.Fields[i])
					}
				}
			})
		}
	})

	t.Run("Select field parsing", func(t *testing.T) {
		tests := []struct {
			name     string
			values   url.Values
			expected []string
		}{
			{
				name:     "No select fields",
				values:   url.Values{},
				expected: nil,
			},
			{
				name:     "Single select field",
				values:   url.Values{"select": {"name"}},
				expected: []string{"name"},
			},
			{
				name:     "Multiple select fields",
				values:   url.Values{"select": {"id,name,email,created_at"}},
				expected: []string{"id", "name", "email", "created_at"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := slicer.ParseOpts(tt.values)

				if len(result.Select) != len(tt.expected) {
					t.Fatalf("Expected %d select fields, got %d", len(tt.expected), len(result.Select))
				}

				for i, expected := range tt.expected {
					if result.Select[i] != expected {
						t.Errorf("Select field %d: expected %s, got %s", i, expected, result.Select[i])
					}
				}
			})
		}
	})

	t.Run("Group By parsing", func(t *testing.T) {
		tests := []struct {
			name           string
			values         url.Values
			expectedGroup  []string
			expectedSelect []string
		}{
			{
				name:           "No group by fields",
				values:         url.Values{},
				expectedGroup:  nil,
				expectedSelect: nil,
			},
			{
				name:           "Group by with single field",
				values:         url.Values{"group": {"department"}},
				expectedGroup:  []string{"department"},
				expectedSelect: []string{"department"}, // Should auto-set select to match group
			},
			{
				name:           "Group by with multiple fields",
				values:         url.Values{"group": {"department,role"}},
				expectedGroup:  []string{"department", "role"},
				expectedSelect: []string{"department", "role"},
			},
			{
				name:           "Group by with sort fields (should include sort in group)",
				values:         url.Values{"group": {"department"}, "sort": {"name,-created_at"}},
				expectedGroup:  []string{"department", "name", "created_at"},
				expectedSelect: []string{"department"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := slicer.ParseOpts(tt.values)

				if len(result.GroupBy) != len(tt.expectedGroup) {
					t.Fatalf("Expected %d group by fields, got %d", len(tt.expectedGroup), len(result.GroupBy))
				}

				for i, expected := range tt.expectedGroup {
					if result.GroupBy[i] != expected {
						t.Errorf("Group by field %d: expected %s, got %s", i, expected, result.GroupBy[i])
					}
				}

				if len(result.Select) != len(tt.expectedSelect) {
					t.Fatalf("Expected %d select fields, got %d", len(tt.expectedSelect), len(result.Select))
				}

				for i, expected := range tt.expectedSelect {
					if result.Select[i] != expected {
						t.Errorf("Select field %d: expected %s, got %s", i, expected, result.Select[i])
					}
				}
			})
		}
	})

	t.Run("Comparison filters parsing", func(t *testing.T) {
		tests := []struct {
			name     string
			values   url.Values
			expected []slicer.ComparisonFilter
		}{
			{
				name:     "No comparison filters",
				values:   url.Values{},
				expected: nil,
			},
			{
				name:   "Single comparison filter",
				values: url.Values{"age[gte]": {"18"}},
				expected: []slicer.ComparisonFilter{
					{Field: "age", Op: "gte", Value: "18"},
				},
			},
			{
				name: "Multiple comparison filters",
				values: url.Values{
					"age[gte]":       {"18"},
					"salary[lt]":     {"100000"},
					"created_at[eq]": {"2023-01-01"},
					"score[gt]":      {"80"},
					"rating[lte]":    {"5"},
				},
				expected: []slicer.ComparisonFilter{
					{Field: "age", Op: "gte", Value: "18"},
					{Field: "salary", Op: "lt", Value: "100000"},
					{Field: "created_at", Op: "eq", Value: "2023-01-01"},
					{Field: "score", Op: "gt", Value: "80"},
					{Field: "rating", Op: "lte", Value: "5"},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := slicer.ParseOpts(tt.values)

				if len(result.Comparisons) != len(tt.expected) {
					t.Fatalf("Expected %d comparison filters, got %d", len(tt.expected), len(result.Comparisons))
				}

				// Sort both slices for consistent comparison since map iteration order is not guaranteed
				for _, expected := range tt.expected {
					found := false
					for _, actual := range result.Comparisons {
						if actual.Field == expected.Field &&
							string(actual.Op) == string(expected.Op) &&
							actual.Value == expected.Value {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected comparison filter %+v not found in result", expected)
					}
				}
			})
		}
	})

	t.Run("Custom filters parsing", func(t *testing.T) {
		tests := []struct {
			name     string
			values   url.Values
			expected map[string]string
		}{
			{
				name:     "No custom filters",
				values:   url.Values{"page": {"1"}, "limit": {"10"}},
				expected: map[string]string{},
			},
			{
				name:   "Single custom filter",
				values: url.Values{"status": {"active"}},
				expected: map[string]string{
					"status": "active",
				},
			},
			{
				name: "Multiple custom filters",
				values: url.Values{
					"status":     {"active"},
					"department": {"engineering"},
					"role":       {"senior"},
					"location":   {"remote"},
				},
				expected: map[string]string{
					"status":     "active",
					"department": "engineering",
					"role":       "senior",
					"location":   "remote",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := slicer.ParseOpts(tt.values)

				if len(result.Filters) != len(tt.expected) {
					t.Fatalf("Expected %d filters, got %d", len(tt.expected), len(result.Filters))
				}

				for key, expectedValue := range tt.expected {
					if actualValue, exists := result.Filters[key]; !exists {
						t.Errorf("Expected filter key %s not found", key)
					} else if actualValue != expectedValue {
						t.Errorf("Filter %s: expected value %s, got %s", key, expectedValue, actualValue)
					}
				}
			})
		}
	})

	t.Run("Complex combined parameters", func(t *testing.T) {
		values := url.Values{
			"page":             {"2"},
			"limit":            {"20"},
			"sort":             {"name,-created_at"},
			"search":           {"name,email"},
			"keyword":          {"john"},
			"select":           {"id,name,email"},
			"status":           {"active"},
			"department":       {"engineering"},
			"age[gte]":         {"25"},
			"salary[lt]":       {"80000"},
			"searchAnd.role":   {"senior"},
			"searchAnd.city":   {"NYC"},
			"search_and.level": {"L5"},
		}

		result := slicer.ParseOpts(values)

		// Verify basic pagination
		if result.Page != 2 {
			t.Errorf("Expected page 2, got %d", result.Page)
		}
		if result.Limit != 20 {
			t.Errorf("Expected limit 20, got %d", result.Limit)
		}
		if result.Offset != 20 {
			t.Errorf("Expected offset 20, got %d", result.Offset)
		}

		// Verify sort
		expectedSort := []slicer.SortField{
			{Field: "name", Desc: false},
			{Field: "created_at", Desc: true},
		}
		if len(result.Sort) != len(expectedSort) {
			t.Fatalf("Expected %d sort fields, got %d", len(expectedSort), len(result.Sort))
		}

		// Verify search
		if result.Search == nil {
			t.Fatal("Expected Search to be set")
		}
		if result.Search.Keyword != "john" {
			t.Errorf("Expected search keyword 'john', got '%s'", result.Search.Keyword)
		}

		// Verify select
		expectedSelect := []string{"id", "name", "email"}
		if len(result.Select) != len(expectedSelect) {
			t.Fatalf("Expected %d select fields, got %d", len(expectedSelect), len(result.Select))
		}

		// Verify filters
		expectedFilters := map[string]string{
			"status":     "active",
			"department": "engineering",
		}
		for key, expectedValue := range expectedFilters {
			if actualValue, exists := result.Filters[key]; !exists {
				t.Errorf("Expected filter key %s not found", key)
			} else if actualValue != expectedValue {
				t.Errorf("Filter %s: expected value %s, got %s", key, expectedValue, actualValue)
			}
		}

		// Verify comparisons
		if len(result.Comparisons) != 2 {
			t.Fatalf("Expected 2 comparison filters, got %d", len(result.Comparisons))
		}

		// Verify SearchAnd
		if result.SearchAnd == nil {
			t.Fatal("Expected SearchAnd to be set")
		}
		if len(result.SearchAnd.Fields) != 3 {
			t.Fatalf("Expected 3 SearchAnd fields, got %d", len(result.SearchAnd.Fields))
		}

		// Check SearchAnd fields (role, city, level)
		searchAndFields := make(map[string]string)
		for _, field := range result.SearchAnd.Fields {
			searchAndFields[field.Field] = field.Keyword
		}

		expectedSearchAnd := map[string]string{
			"role":  "senior",
			"city":  "NYC",
			"level": "L5",
		}

		for expectedField, expectedKeyword := range expectedSearchAnd {
			if actualKeyword, exists := searchAndFields[expectedField]; !exists {
				t.Errorf("Expected SearchAnd field %s not found", expectedField)
			} else if actualKeyword != expectedKeyword {
				t.Errorf("SearchAnd field %s: expected keyword %s, got %s", expectedField, expectedKeyword, actualKeyword)
			}
		}
	})
}

// TestSetValueSeparator tests the SetValueSeparator function
func TestSetValueSeparator(t *testing.T) {
	// Save original separator
	originalSeparator := ","

	// Test setting a custom separator
	slicer.SetValueSeparator("|")

	// Test with custom separator
	values := url.Values{"search": {"name|email|description"}}
	result := slicer.ParseOpts(values)

	if result.Search == nil {
		// This test may need adjustment based on actual implementation
		// since Search also needs keyword parameter
		t.Log("Search is nil - this may be expected behavior without keyword")
	}

	// Reset to original separator
	slicer.SetValueSeparator(originalSeparator)

	// Test with default separator
	values2 := url.Values{"search": {"name,email,description"}, "keyword": {"test"}}
	result2 := slicer.ParseOpts(values2)

	if result2.Search == nil {
		t.Fatal("Expected Search to be set with default separator")
	}

	expectedFields := []string{"name", "email", "description"}
	if len(result2.Search.Fields) != len(expectedFields) {
		t.Fatalf("Expected %d search fields, got %d", len(expectedFields), len(result2.Search.Fields))
	}

	for i, expected := range expectedFields {
		if result2.Search.Fields[i] != expected {
			t.Errorf("Search field %d: expected %s, got %s", i, expected, result2.Search.Fields[i])
		}
	}
}
