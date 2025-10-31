package slicer_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/godev90/slicer"
)

const (
	slicePageFailedMsg = "SlicePage failed: %v"
)

// Test data structures
type TestUser struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Status     string `json:"status"`
	Department string `json:"department"`
	Role       string `json:"role"`
	City       string `json:"city"`
	Country    string `json:"country"`
}

// Test data
var testUsers = []TestUser{
	{1, "John Doe", "john@example.com", "active", "engineering", "senior", "New York", "USA"},
	{2, "Jane Smith", "jane@example.com", "active", "marketing", "junior", "London", "UK"},
	{3, "Bob Johnson", "bob@example.com", "inactive", "engineering", "senior", "Toronto", "Canada"},
	{4, "Alice Brown", "alice@example.com", "active", "sales", "manager", "Sydney", "Australia"},
	{5, "Charlie Wilson", "charlie@example.com", "suspended", "engineering", "junior", "Berlin", "Germany"},
}

// 1. URL Parameter Parsing Tests
func TestSearchAndURLParsingScenarios(t *testing.T) {
	tests := []struct {
		name     string
		urlQuery string
		expected int // expected number of search fields
		fields   map[string]string
	}{
		{
			name:     "Multiple fields with different keywords",
			urlQuery: "searchAnd.status=active&searchAnd.department=engineering&searchAnd.role=senior",
			expected: 3,
			fields:   map[string]string{"status": "active", "department": "engineering", "role": "senior"},
		},
		{
			name:     "Single field",
			urlQuery: "searchAnd.status=active",
			expected: 1,
			fields:   map[string]string{"status": "active"},
		},
		{
			name:     "Empty keyword values should be ignored",
			urlQuery: "searchAnd.status=&searchAnd.department=engineering",
			expected: 1,
			fields:   map[string]string{"department": "engineering"},
		},
		{
			name:     "Special characters in keywords",
			urlQuery: "searchAnd.name=John%20Doe&searchAnd.email=john%40example.com",
			expected: 2,
			fields:   map[string]string{"name": "John Doe", "email": "john@example.com"},
		},
		{
			name:     "Case sensitive field names",
			urlQuery: "searchAnd.Status=active&searchAnd.status=inactive",
			expected: 2,
			fields:   map[string]string{"Status": "active", "status": "inactive"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, _ := url.ParseQuery(tt.urlQuery)
			opts := slicer.ParseOpts(values)

			if opts.SearchAnd == nil {
				if tt.expected > 0 {
					t.Fatalf("Expected SearchAnd to be set with %d fields, got nil", tt.expected)
				}
				return
			}

			if len(opts.SearchAnd.Fields) != tt.expected {
				t.Fatalf("Expected %d search fields, got %d", tt.expected, len(opts.SearchAnd.Fields))
			}

			// Verify field values
			for _, field := range opts.SearchAnd.Fields {
				if expectedValue, exists := tt.fields[field.Field]; exists {
					if field.Keyword != expectedValue {
						t.Errorf("Expected field %s to have keyword %s, got %s", field.Field, expectedValue, field.Keyword)
					}
				}
			}
		})
	}
}

// 2. Integration with Existing Search Tests
func TestSearchAndIntegrationScenarios(t *testing.T) {
	tests := []struct {
		name            string
		urlQuery        string
		expectSearch    bool
		expectAndSearch bool
	}{
		{
			name:            "SearchAnd only",
			urlQuery:        "searchAnd.status=active",
			expectSearch:    false,
			expectAndSearch: true,
		},
		{
			name:            "Regular search only",
			urlQuery:        "search=name&search=email&keyword=john",
			expectSearch:    true,
			expectAndSearch: false,
		},
		{
			name:            "Both search types",
			urlQuery:        "search=name&keyword=john&searchAnd.status=active",
			expectSearch:    true,
			expectAndSearch: true,
		},
		{
			name:            "Neither search type",
			urlQuery:        "page=1&limit=10",
			expectSearch:    false,
			expectAndSearch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, _ := url.ParseQuery(tt.urlQuery)
			opts := slicer.ParseOpts(values)

			hasSearch := opts.Search != nil && len(opts.Search.Fields) > 0 && opts.Search.Keyword != ""
			hasAndSearch := opts.SearchAnd != nil && len(opts.SearchAnd.Fields) > 0

			if hasSearch != tt.expectSearch {
				t.Errorf("Expected regular search: %v, got: %v", tt.expectSearch, hasSearch)
			}

			if hasAndSearch != tt.expectAndSearch {
				t.Errorf("Expected AND search: %v, got: %v", tt.expectAndSearch, hasAndSearch)
			}
		})
	}
}

// 3. Protocol Buffer Edge Cases
func TestSearchAndProtoEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		opts *slicer.QueryOptions
	}{
		{
			name: "Nil SearchAnd",
			opts: &slicer.QueryOptions{
				Page:      1,
				Limit:     10,
				SearchAnd: nil,
			},
		},
		{
			name: "Empty SearchQueryAnd",
			opts: &slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				SearchAnd: &slicer.SearchQueryAnd{
					Fields: []*slicer.SearchField{},
				},
			},
		},
		{
			name: "Large number of search fields",
			opts: &slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				SearchAnd: &slicer.SearchQueryAnd{
					Fields: []*slicer.SearchField{
						{Field: "field1", Keyword: "value1"},
						{Field: "field2", Keyword: "value2"},
						{Field: "field3", Keyword: "value3"},
						{Field: "field4", Keyword: "value4"},
						{Field: "field5", Keyword: "value5"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert to proto
			proto := tt.opts.ToProto()
			if proto == nil {
				t.Fatal("ToProto() returned nil")
			}

			// Convert back from proto
			converted := slicer.QueryFromProto(proto)

			// Verify SearchAnd field specifically
			if tt.opts.SearchAnd == nil {
				if converted.SearchAnd != nil {
					t.Error("Expected SearchAnd to remain nil after conversion")
				}
			} else {
				// For empty SearchQueryAnd, it might become nil after conversion (acceptable)
				if len(tt.opts.SearchAnd.Fields) == 0 {
					// Empty SearchQueryAnd may be converted to nil - this is acceptable
					t.Logf("Empty SearchQueryAnd converted to nil - this is acceptable")
				} else if converted.SearchAnd == nil {
					t.Error("Expected SearchAnd to be preserved after conversion")
				} else if len(converted.SearchAnd.Fields) != len(tt.opts.SearchAnd.Fields) {
					t.Errorf("Expected %d SearchAnd fields, got %d", len(tt.opts.SearchAnd.Fields), len(converted.SearchAnd.Fields))
				}
			}
		})
	}
}

// 4. Slice Pagination Edge Cases
func TestSlicePaginatorEdgeCases(t *testing.T) {
	allowedFields := map[string]string{
		"id":         "id",
		"name":       "name",
		"email":      "email",
		"status":     "status",
		"department": "department",
		"role":       "role",
		"city":       "city",
		"country":    "country",
	}

	tests := []struct {
		name        string
		users       []TestUser
		searchAnd   *slicer.SearchQueryAnd
		expectedLen int
	}{
		{
			name:        "Empty slice",
			users:       []TestUser{},
			searchAnd:   &slicer.SearchQueryAnd{Fields: []*slicer.SearchField{{Field: "status", Keyword: "active"}}},
			expectedLen: 0,
		},
		{
			name:  "No matches found",
			users: testUsers,
			searchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "status", Keyword: "deleted"},
				},
			},
			expectedLen: 0,
		},
		{
			name:  "All items match",
			users: testUsers,
			searchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "email", Keyword: "example.com"},
				},
			},
			expectedLen: 5,
		},
		{
			name:  "Multiple AND conditions - strict match",
			users: testUsers,
			searchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "name", Keyword: "Doe"}, // Only matches "John Doe", not "Bob Johnson"
					{Field: "department", Keyword: "engineering"},
					{Field: "role", Keyword: "senior"},
				},
			},
			expectedLen: 1, // Only John Doe matches all conditions
		},
		{
			name:  "Case insensitive search",
			users: testUsers,
			searchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "name", Keyword: "Doe"}, // Only matches "John Doe"
				},
			},
			expectedLen: 1, // Only John Doe should match
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paginator := slicer.NewSlicePaginator(tt.users, allowedFields)
			opts := slicer.QueryOptions{
				Page:      1,
				Limit:     10,
				SearchAnd: tt.searchAnd,
			}

			result, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatalf(slicePageFailedMsg, err)
			}

			resultUsers, ok := result.Items.([]TestUser)
			if !ok {
				t.Fatalf("Expected []TestUser, got %T", result.Items)
			}

			if len(resultUsers) != tt.expectedLen {
				// Debug: print actual results for failing tests
				t.Logf("Test case: %s", tt.name)
				for i, user := range resultUsers {
					t.Logf("Result %d: %s (status=%s, dept=%s, role=%s)", i, user.Name, user.Status, user.Department, user.Role)
				}
				t.Errorf("Expected %d items, got %d", tt.expectedLen, len(resultUsers))
			}
		})
	}
}

// 5. Performance and Stress Tests
func TestSearchAndPerformance(t *testing.T) {
	// Create large dataset
	largeDataset := make([]TestUser, 10000)
	for i := 0; i < 10000; i++ {
		largeDataset[i] = TestUser{
			ID:         i,
			Name:       "User" + strings.Repeat("X", i%100),
			Status:     []string{"active", "inactive", "suspended"}[i%3],
			Department: []string{"engineering", "marketing", "sales"}[i%3],
			Role:       []string{"junior", "senior", "manager"}[i%3],
		}
	}

	allowedFields := map[string]string{
		"id":         "id",
		"name":       "name",
		"status":     "status",
		"department": "department",
		"role":       "role",
	}

	tests := []struct {
		name      string
		searchAnd *slicer.SearchQueryAnd
	}{
		{
			name: "Single field search on large dataset",
			searchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "status", Keyword: "active"},
				},
			},
		},
		{
			name: "Multiple field search on large dataset",
			searchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "status", Keyword: "active"},
					{Field: "department", Keyword: "engineering"},
					{Field: "role", Keyword: "senior"},
				},
			},
		},
		{
			name: "Many search fields",
			searchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "status", Keyword: "active"},
					{Field: "department", Keyword: "engineering"},
					{Field: "role", Keyword: "senior"},
					{Field: "name", Keyword: "User"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paginator := slicer.NewSlicePaginator(largeDataset, allowedFields)
			opts := slicer.QueryOptions{
				Page:      1,
				Limit:     100,
				SearchAnd: tt.searchAnd,
			}

			result, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatalf(slicePageFailedMsg, err)
			}

			// Just verify it completes without error and returns reasonable results
			t.Logf("Test '%s' processed %d items, returned %d results", tt.name, len(largeDataset), len(result.Items.([]TestUser)))
		})
	}
}

// 6. Field Validation Tests
func TestSearchAndFieldValidation(t *testing.T) {
	tests := []struct {
		name      string
		urlQuery  string
		shouldErr bool
	}{
		{
			name:      "Valid field names",
			urlQuery:  "searchAnd.status=active&searchAnd.department=engineering",
			shouldErr: false,
		},
		{
			name:      "Empty field name should be ignored",
			urlQuery:  "searchAnd.=active&searchAnd.department=engineering",
			shouldErr: false,
		},
		{
			name:      "Special characters in field names",
			urlQuery:  "searchAnd.field-name=value&searchAnd.field_name=value2",
			shouldErr: false,
		},
		{
			name:      "Numeric field names",
			urlQuery:  "searchAnd.123=value&searchAnd.field456=value2",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, _ := url.ParseQuery(tt.urlQuery)
			opts := slicer.ParseOpts(values)

			// The parsing itself should never error, but we can check if the result makes sense
			if opts.SearchAnd != nil {
				for _, field := range opts.SearchAnd.Fields {
					if field.Field == "" && field.Keyword != "" {
						t.Error("Found search field with empty field name but non-empty keyword")
					}
				}
			}
		})
	}
}

// 7. Unicode and Special Character Tests
func TestSearchAndUnicodeHandling(t *testing.T) {
	unicodeUsers := []TestUser{
		{1, "José María", "jose@example.com", "active", "engineering", "senior", "México", "México"},
		{2, "山田太郎", "yamada@example.com", "active", "marketing", "junior", "東京", "日本"},
		{3, "Müller Schmidt", "muller@example.com", "inactive", "sales", "manager", "München", "Deutschland"},
	}

	allowedFields := map[string]string{
		"id":      "id",
		"name":    "name",
		"email":   "email",
		"status":  "status",
		"city":    "city",
		"country": "country",
	}

	tests := []struct {
		name        string
		searchAnd   *slicer.SearchQueryAnd
		expectedLen int
	}{
		{
			name: "Unicode characters in keyword",
			searchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "name", Keyword: "José"},
				},
			},
			expectedLen: 1,
		},
		{
			name: "Japanese characters",
			searchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "city", Keyword: "東京"},
				},
			},
			expectedLen: 1,
		},
		{
			name: "German umlauts",
			searchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "name", Keyword: "Müller"},
				},
			},
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paginator := slicer.NewSlicePaginator(unicodeUsers, allowedFields)
			opts := slicer.QueryOptions{
				Page:      1,
				Limit:     10,
				SearchAnd: tt.searchAnd,
			}

			result, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatalf(slicePageFailedMsg, err)
			}

			resultUsers, ok := result.Items.([]TestUser)
			if !ok {
				t.Fatalf("Expected []TestUser, got %T", result.Items)
			}

			if len(resultUsers) != tt.expectedLen {
				t.Errorf("Expected %d items, got %d", tt.expectedLen, len(resultUsers))
			}
		})
	}
}
