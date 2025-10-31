package slicer_test

import (
	"net/url"
	"testing"

	"github.com/godev90/slicer"
)

func TestSearchAndParsing(t *testing.T) {
	// Test URL values with searchAnd parameters
	values := url.Values{}
	values.Add("searchAnd.status", "active")
	values.Add("searchAnd.category", "backend")
	values.Add("page", "1")
	values.Add("limit", "10")

	opts := slicer.ParseOpts(values)

	// Check that SearchAnd is populated correctly
	if opts.SearchAnd == nil {
		t.Fatal("SearchAnd should not be nil")
	}

	if len(opts.SearchAnd.Fields) != 2 {
		t.Fatalf("Expected 2 search fields, got %d", len(opts.SearchAnd.Fields))
	}

	// Convert to map for easier testing since order might vary
	fieldMap := make(map[string]string)
	for _, field := range opts.SearchAnd.Fields {
		fieldMap[field.Field] = field.Keyword
	}

	// Check fields
	if fieldMap["status"] != "active" {
		t.Errorf("Expected status=active, got status=%s", fieldMap["status"])
	}
	if fieldMap["category"] != "backend" {
		t.Errorf("Expected category=backend, got category=%s", fieldMap["category"])
	}
}

func TestSearchAndWithRegularSearch(t *testing.T) {
	// Test combining regular search with searchAnd
	values := url.Values{}
	values.Add("search", "name,description")
	values.Add("keyword", "golang")
	values.Add("searchAnd.status", "active")
	values.Add("searchAnd.category", "backend")

	opts := slicer.ParseOpts(values)

	// Check regular search is still working
	if opts.Search == nil {
		t.Fatal("Search should not be nil")
	}
	if len(opts.Search.Fields) != 2 {
		t.Fatalf("Expected 2 search fields, got %d", len(opts.Search.Fields))
	}
	if opts.Search.Keyword != "golang" {
		t.Errorf("Expected keyword 'golang', got '%s'", opts.Search.Keyword)
	}

	// Check SearchAnd is also working
	if opts.SearchAnd == nil {
		t.Fatal("SearchAnd should not be nil")
	}
	if len(opts.SearchAnd.Fields) != 2 {
		t.Fatalf("Expected 2 searchAnd fields, got %d", len(opts.SearchAnd.Fields))
	}
}

func TestProtoConversion(t *testing.T) {
	// Test proto conversion with SearchAnd
	opts := slicer.QueryOptions{
		Page:  1,
		Limit: 10,
		SearchAnd: &slicer.SearchQueryAnd{
			Fields: []*slicer.SearchField{
				{Field: "status", Keyword: "active"},
				{Field: "category", Keyword: "backend"},
			},
		},
	}

	// Convert to proto
	proto := opts.ToProto()
	if proto.SearchAnd == nil {
		t.Fatal("Proto SearchAnd should not be nil")
	}
	if len(proto.SearchAnd.Fields) != 2 {
		t.Fatalf("Expected 2 proto search fields, got %d", len(proto.SearchAnd.Fields))
	}

	// Convert back from proto
	converted := slicer.QueryFromProto(proto)
	if converted.SearchAnd == nil {
		t.Fatal("Converted SearchAnd should not be nil")
	}
	if len(converted.SearchAnd.Fields) != 2 {
		t.Fatalf("Expected 2 converted search fields, got %d", len(converted.SearchAnd.Fields))
	}

	// Check specific field
	fieldMap := make(map[string]string)
	for _, field := range converted.SearchAnd.Fields {
		fieldMap[field.Field] = field.Keyword
	}
	if fieldMap["status"] != "active" {
		t.Errorf("Proto conversion failed for status field")
	}
	if fieldMap["category"] != "backend" {
		t.Errorf("Proto conversion failed for category field")
	}
}

func TestSlicePaginatorSearchAnd(t *testing.T) {
	// Test data
	type TestUser struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Status   string `json:"status"`
		Category string `json:"category"`
	}

	users := []TestUser{
		{ID: 1, Name: "John", Status: "active", Category: "admin"},
		{ID: 2, Name: "Jane", Status: "active", Category: "user"},
		{ID: 3, Name: "Bob", Status: "inactive", Category: "user"},
		{ID: 4, Name: "Alice", Status: "active", Category: "admin"},
	}

	allowedFields := map[string]string{
		"id":       "id",
		"name":     "name",
		"status":   "status",
		"category": "category",
	}

	paginator := slicer.NewSlicePaginator(users, allowedFields)

	opts := slicer.QueryOptions{
		Page:  1,
		Limit: 10,
		SearchAnd: &slicer.SearchQueryAnd{
			Fields: []*slicer.SearchField{
				{Field: "status", Keyword: "active"},
				{Field: "category", Keyword: "admin"},
			},
		},
	}

	result, err := slicer.SlicePage(paginator, opts)
	if err != nil {
		t.Fatalf("SlicePage returned error: %v", err)
	}

	resultUsers, ok := result.Items.([]TestUser)
	if !ok {
		t.Fatalf("Expected []TestUser, got %T", result.Items)
	}

	// Should match John and Alice (both have status=active AND category=admin)
	expectedCount := 2
	if len(resultUsers) != expectedCount {
		t.Errorf("Expected %d results, got %d", expectedCount, len(resultUsers))
	}

	// Verify total count
	if result.Total != int64(expectedCount) {
		t.Errorf("Expected total %d, got %d", expectedCount, result.Total)
	}
}
