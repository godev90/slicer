package slicer_test

import (
	"testing"

	"github.com/godev90/slicer"
)

const (
	slicePageError = "SlicePage failed: %v"
)

// Test error handling and edge cases
func TestSearchAndErrorHandling(t *testing.T) {
	tests := []struct {
		name      string
		searchAnd *slicer.SearchQueryAnd
		expectErr bool
	}{
		{
			name:      "Nil SearchQueryAnd",
			searchAnd: nil,
			expectErr: false, // Should handle gracefully
		},
		{
			name: "SearchQueryAnd with nil Fields slice",
			searchAnd: &slicer.SearchQueryAnd{
				Fields: nil,
			},
			expectErr: false,
		},
		{
			name: "SearchQueryAnd with empty SearchField values",
			searchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "", Keyword: ""},
					{Field: "status", Keyword: "active"},
				},
			},
			expectErr: false, // Should filter out empty fields gracefully
		},
		{
			name: "Mixed empty and valid fields",
			searchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "", Keyword: "ignored"},
					{Field: "status", Keyword: "active"},
					{Field: "department", Keyword: ""},
					{Field: "role", Keyword: "senior"},
				},
			},
			expectErr: false, // Should use only valid fields
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &slicer.QueryOptions{
				Page:      1,
				Limit:     10,
				SearchAnd: tt.searchAnd,
			}

			// Test proto conversion doesn't panic
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("ToProto() panicked: %v", r)
					}
				}()
				proto := opts.ToProto()
				if proto != nil {
					slicer.QueryFromProto(proto)
				}
			}()

			// Test slice pagination doesn't panic
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("SlicePage() panicked: %v", r)
					}
				}()
				testData := []TestUser{
					{1, "John", "john@test.com", "active", "eng", "senior", "NYC", "USA"},
				}
				allowedFields := map[string]string{"status": "status"}
				paginator := slicer.NewSlicePaginator(testData, allowedFields)
				_, _ = slicer.SlicePage(paginator, *opts)
			}()
		})
	}
}

// Test memory and resource handling
func TestSearchAndResourceHandling(t *testing.T) {
	tests := []struct {
		name        string
		fieldsCount int
	}{
		{"Small number of fields", 5},
		{"Medium number of fields", 50},
		{"Large number of fields", 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create SearchQueryAnd with many fields
			fields := make([]*slicer.SearchField, tt.fieldsCount)
			for i := 0; i < tt.fieldsCount; i++ {
				fields[i] = &slicer.SearchField{
					Field:   "field" + string(rune(i)),
					Keyword: "keyword" + string(rune(i)),
				}
			}

			searchAnd := &slicer.SearchQueryAnd{Fields: fields}
			opts := &slicer.QueryOptions{
				Page:      1,
				Limit:     10,
				SearchAnd: searchAnd,
			}

			// Test proto conversion with many fields
			proto := opts.ToProto()
			if proto == nil {
				t.Error("ToProto() returned nil")
			}

			converted := slicer.QueryFromProto(proto)
			if converted.SearchAnd == nil {
				t.Error("QueryFromProto() returned nil SearchAnd")
			}

			if len(converted.SearchAnd.Fields) != tt.fieldsCount {
				t.Errorf("Expected %d fields after conversion, got %d", tt.fieldsCount, len(converted.SearchAnd.Fields))
			}
		})
	}
}

// Test concurrent access safety
func TestSearchAndConcurrency(t *testing.T) {
	opts := &slicer.QueryOptions{
		Page:  1,
		Limit: 10,
		SearchAnd: &slicer.SearchQueryAnd{
			Fields: []*slicer.SearchField{
				{Field: "status", Keyword: "active"},
				{Field: "department", Keyword: "engineering"},
			},
		},
	}

	// Test concurrent proto conversions
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Concurrent access panicked: %v", r)
				}
				done <- true
			}()

			for j := 0; j < 100; j++ {
				proto := opts.ToProto()
				slicer.QueryFromProto(proto)
			}
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Test data integrity and immutability
func TestSearchAndDataIntegrity(t *testing.T) {
	original := &slicer.QueryOptions{
		Page:  1,
		Limit: 10,
		SearchAnd: &slicer.SearchQueryAnd{
			Fields: []*slicer.SearchField{
				{Field: "status", Keyword: "active"},
			},
		},
	}

	// Store original values
	originalField := original.SearchAnd.Fields[0].Field
	originalKeyword := original.SearchAnd.Fields[0].Keyword

	// Convert to proto and back
	proto := original.ToProto()
	converted := slicer.QueryFromProto(proto)

	// Modify the converted version
	if converted.SearchAnd != nil && len(converted.SearchAnd.Fields) > 0 {
		converted.SearchAnd.Fields[0].Field = "modified"
		converted.SearchAnd.Fields[0].Keyword = "modified"
	}

	// Verify original is unchanged
	if original.SearchAnd.Fields[0].Field != originalField {
		t.Error("Original data was modified during conversion")
	}
	if original.SearchAnd.Fields[0].Keyword != originalKeyword {
		t.Error("Original data was modified during conversion")
	}
}

// Test edge cases with malformed data
func TestSearchAndMalformedData(t *testing.T) {
	tests := []struct {
		name      string
		searchAnd *slicer.SearchQueryAnd
	}{
		{
			name: "Field with only spaces",
			searchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "   ", Keyword: "active"},
				},
			},
		},
		{
			name: "Keyword with only spaces",
			searchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "status", Keyword: "   "},
				},
			},
		},
		{
			name: "Very long field name",
			searchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: string(make([]byte, 1000)), Keyword: "value"},
				},
			},
		},
		{
			name: "Very long keyword",
			searchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "field", Keyword: string(make([]byte, 1000))},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &slicer.QueryOptions{
				Page:      1,
				Limit:     10,
				SearchAnd: tt.searchAnd,
			}

			// Should not panic with malformed data
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Processing malformed data panicked: %v", r)
					}
				}()

				proto := opts.ToProto()
				if proto != nil {
					slicer.QueryFromProto(proto)
				}
			}()
		})
	}
}
