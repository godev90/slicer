package slicer

import (
	"net/url"
	"testing"
)

func TestSearchAndParsing(t *testing.T) {
	// Test URL values with search_and parameters
	values := url.Values{}
	values.Add("search_and.status", "active")
	values.Add("search_and.category", "backend")
	values.Add("page", "1")
	values.Add("limit", "10")

	opts := ParseOpts(values)

	// Check that SearchAnd is populated correctly
	if opts.SearchAnd == nil {
		t.Fatal("SearchAnd should not be nil")
	}

	if len(opts.SearchAnd.Fields) != 2 {
		t.Fatalf("Expected 2 search fields, got %d", len(opts.SearchAnd.Fields))
	}

	// Check first field
	if opts.SearchAnd.Fields[0].Field != "status" || opts.SearchAnd.Fields[0].Keyword != "active" {
		t.Errorf("Expected status=active, got %s=%s", opts.SearchAnd.Fields[0].Field, opts.SearchAnd.Fields[0].Keyword)
	}

	// Check second field
	if opts.SearchAnd.Fields[1].Field != "category" || opts.SearchAnd.Fields[1].Keyword != "backend" {
		t.Errorf("Expected category=backend, got %s=%s", opts.SearchAnd.Fields[1].Field, opts.SearchAnd.Fields[1].Keyword)
	}
}

func TestSearchAndWithRegularSearch(t *testing.T) {
	// Test combining regular search with search_and
	values := url.Values{}
	values.Add("search", "name,description")
	values.Add("keyword", "golang")
	values.Add("search_and.status", "active")
	values.Add("search_and.category", "backend")

	opts := ParseOpts(values)

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
		t.Fatalf("Expected 2 search_and fields, got %d", len(opts.SearchAnd.Fields))
	}
}

func TestProtoConversion(t *testing.T) {
	// Test proto conversion with SearchAnd
	opts := QueryOptions{
		Page:  1,
		Limit: 10,
		SearchAnd: &SearchQueryAnd{
			Fields: []*SearchField{
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
	converted := QueryFromProto(proto)
	if converted.SearchAnd == nil {
		t.Fatal("Converted SearchAnd should not be nil")
	}
	if len(converted.SearchAnd.Fields) != 2 {
		t.Fatalf("Expected 2 converted search fields, got %d", len(converted.SearchAnd.Fields))
	}
	if converted.SearchAnd.Fields[0].Field != "status" || converted.SearchAnd.Fields[0].Keyword != "active" {
		t.Errorf("Proto conversion failed for first field")
	}
}
