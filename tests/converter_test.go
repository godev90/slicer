package slicer_test

import (
	"encoding/json"
	"testing"

	"github.com/godev90/slicer"
)

// TestConverterFunctions tests ToProto/FromProto conversions
func TestConverterFunctions(t *testing.T) {
	t.Run("QueryOptions ToProto conversion", func(t *testing.T) {
		// Test complete QueryOptions conversion
		opts := slicer.QueryOptions{
			Page:  2,
			Limit: 25,
			Sort: []slicer.SortField{
				{Field: "name", Desc: false},
				{Field: "created_at", Desc: true},
			},
			Search: &slicer.SearchQuery{
				Fields:  []string{"name", "email", "description"},
				Keyword: "test keyword",
			},
			SearchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "status", Keyword: "active"},
					{Field: "department", Keyword: "engineering"},
				},
			},
			Filters: map[string]string{
				"status":     "active",
				"department": "engineering",
				"role":       "senior",
			},
			Select:  []string{"id", "name", "email"},
			GroupBy: []string{"department", "role"},
			Comparisons: []slicer.ComparisonFilter{
				{Field: "age", Op: "gte", Value: "25"},
				{Field: "salary", Op: "lt", Value: "100000"},
			},
		}

		proto := opts.ToProto()

		// Verify basic fields
		if proto.Page != uint32(opts.Page) {
			t.Errorf("Expected page %d, got %d", opts.Page, proto.Page)
		}
		if proto.Limit != uint32(opts.Limit) {
			t.Errorf("Expected limit %d, got %d", opts.Limit, proto.Limit)
		}

		// Verify sort fields
		if len(proto.Sort) != len(opts.Sort) {
			t.Fatalf("Expected %d sort fields, got %d", len(opts.Sort), len(proto.Sort))
		}
		for i, sortField := range opts.Sort {
			if proto.Sort[i].Field != sortField.Field {
				t.Errorf("Sort field %d: expected field %s, got %s", i, sortField.Field, proto.Sort[i].Field)
			}
			if proto.Sort[i].Desc != sortField.Desc {
				t.Errorf("Sort field %d: expected desc %t, got %t", i, sortField.Desc, proto.Sort[i].Desc)
			}
		}

		// Verify search query
		if proto.Search == nil {
			t.Fatal("Expected search to be converted")
		}
		if proto.Search.Keyword != opts.Search.Keyword {
			t.Errorf("Expected search keyword %s, got %s", opts.Search.Keyword, proto.Search.Keyword)
		}
		if len(proto.Search.Fields) != len(opts.Search.Fields) {
			t.Fatalf("Expected %d search fields, got %d", len(opts.Search.Fields), len(proto.Search.Fields))
		}

		// Verify SearchAnd
		if proto.SearchAnd == nil {
			t.Fatal("Expected SearchAnd to be converted")
		}
		if len(proto.SearchAnd.Fields) != len(opts.SearchAnd.Fields) {
			t.Fatalf("Expected %d SearchAnd fields, got %d", len(opts.SearchAnd.Fields), len(proto.SearchAnd.Fields))
		}

		// Verify filters
		if len(proto.Filters) != len(opts.Filters) {
			t.Fatalf("Expected %d filters, got %d", len(opts.Filters), len(proto.Filters))
		}

		// Verify select fields
		if len(proto.Select) != len(opts.Select) {
			t.Fatalf("Expected %d select fields, got %d", len(opts.Select), len(proto.Select))
		}

		// Verify group by fields
		if len(proto.GroupBy) != len(opts.GroupBy) {
			t.Fatalf("Expected %d group by fields, got %d", len(opts.GroupBy), len(proto.GroupBy))
		}

		// Verify comparisons
		if len(proto.Comparisons) != len(opts.Comparisons) {
			t.Fatalf("Expected %d comparisons, got %d", len(opts.Comparisons), len(proto.Comparisons))
		}
	})

	t.Run("QueryFromProto conversion", func(t *testing.T) {
		// Test nil proto
		result := slicer.QueryFromProto(nil)
		if result.Page != 1 {
			t.Errorf("Expected default page 1, got %d", result.Page)
		}
		if result.Limit != 10 {
			t.Errorf("Expected default limit 10, got %d", result.Limit)
		}

		// Test proto with zero values that should be defaulted
		opts := slicer.QueryOptions{
			Page:  0, // Should default to 1
			Limit: 0, // Should default to 10
		}
		proto := opts.ToProto()
		result = slicer.QueryFromProto(proto)

		if result.Page != 1 {
			t.Errorf("Expected page 1 (defaulted), got %d", result.Page)
		}
		if result.Limit != 10 {
			t.Errorf("Expected limit 10 (defaulted), got %d", result.Limit)
		}
	})

	t.Run("Round-trip conversion", func(t *testing.T) {
		// Test that converting to proto and back preserves data
		original := slicer.QueryOptions{
			Page:  3,
			Limit: 50,
			Sort: []slicer.SortField{
				{Field: "name", Desc: true},
			},
			Search: &slicer.SearchQuery{
				Fields:  []string{"title", "content"},
				Keyword: "search term",
			},
			SearchAnd: &slicer.SearchQueryAnd{
				Fields: []*slicer.SearchField{
					{Field: "category", Keyword: "tech"},
				},
			},
			Filters: map[string]string{
				"published": "true",
			},
			Select: []string{"id", "title"},
			Comparisons: []slicer.ComparisonFilter{
				{Field: "views", Op: "gt", Value: "1000"},
			},
		}

		proto := original.ToProto()
		converted := slicer.QueryFromProto(proto)

		// Verify basic fields
		if converted.Page != original.Page {
			t.Errorf("Round-trip page: expected %d, got %d", original.Page, converted.Page)
		}
		if converted.Limit != original.Limit {
			t.Errorf("Round-trip limit: expected %d, got %d", original.Limit, converted.Limit)
		}

		// Verify sort fields
		if len(converted.Sort) != len(original.Sort) {
			t.Fatalf("Round-trip sort fields: expected %d, got %d", len(original.Sort), len(converted.Sort))
		}

		// Verify search query
		if converted.Search == nil && original.Search != nil {
			t.Fatal("Round-trip search: expected search to be preserved")
		}
		if converted.Search != nil && original.Search != nil {
			if converted.Search.Keyword != original.Search.Keyword {
				t.Errorf("Round-trip search keyword: expected %s, got %s", original.Search.Keyword, converted.Search.Keyword)
			}
		}

		// Verify SearchAnd
		if converted.SearchAnd == nil && original.SearchAnd != nil {
			t.Fatal("Round-trip SearchAnd: expected SearchAnd to be preserved")
		}

		// Verify filters
		if len(converted.Filters) != len(original.Filters) {
			t.Fatalf("Round-trip filters: expected %d, got %d", len(original.Filters), len(converted.Filters))
		}
	})

	t.Run("PageData ToProto conversion", func(t *testing.T) {
		// Test data for PageData conversion
		type TestItem struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		items := []TestItem{
			{ID: 1, Name: "Item 1"},
			{ID: 2, Name: "Item 2"},
		}

		pageData := slicer.PageData{
			Items: items,
			Total: 100,
			Page:  2,
			Limit: 20,
		}

		proto, err := pageData.ToProto()
		if err != nil {
			t.Fatalf("ToProto returned error: %v", err)
		}

		if proto.Total != pageData.Total {
			t.Errorf("Expected total %d, got %d", pageData.Total, proto.Total)
		}
		if proto.Page != int32(pageData.Page) {
			t.Errorf("Expected page %d, got %d", pageData.Page, proto.Page)
		}
		if proto.Limit != int32(pageData.Limit) {
			t.Errorf("Expected limit %d, got %d", pageData.Limit, proto.Limit)
		}

		// Verify items were marshaled to JSON
		if len(proto.Items) == 0 {
			t.Error("Expected items to be marshaled to JSON bytes")
		}

		// Verify we can unmarshal the items back
		var unmarshaledItems []TestItem
		err = json.Unmarshal(proto.Items, &unmarshaledItems)
		if err != nil {
			t.Fatalf("Failed to unmarshal items: %v", err)
		}

		if len(unmarshaledItems) != len(items) {
			t.Errorf("Expected %d items after unmarshal, got %d", len(items), len(unmarshaledItems))
		}
	})

	t.Run("PageFromProto conversion", func(t *testing.T) {
		// Test nil proto
		result, err := slicer.PageFromProto(nil, nil)
		if err != nil {
			t.Errorf("PageFromProto with nil should not error, got: %v", err)
		}
		if result != nil {
			t.Errorf("PageFromProto with nil should return nil")
		}

		// Test with actual data
		type TestItem struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		items := []TestItem{
			{ID: 1, Name: "Item 1"},
			{ID: 2, Name: "Item 2"},
		}

		// Create a PageData and convert to proto
		pageData := slicer.PageData{
			Items: items,
			Total: 50,
			Page:  3,
			Limit: 15,
		}

		proto, err := pageData.ToProto()
		if err != nil {
			t.Fatalf("ToProto failed: %v", err)
		}

		// Convert back from proto
		var destItems []TestItem
		result, err = slicer.PageFromProto(proto, &destItems)
		if err != nil {
			t.Fatalf("PageFromProto failed: %v", err)
		}

		if result == nil {
			t.Fatal("PageFromProto returned nil")
		}

		// Verify page metadata
		if result.Total != pageData.Total {
			t.Errorf("Expected total %d, got %d", pageData.Total, result.Total)
		}
		if result.Page != pageData.Page {
			t.Errorf("Expected page %d, got %d", pageData.Page, result.Page)
		}
		if result.Limit != pageData.Limit {
			t.Errorf("Expected limit %d, got %d", pageData.Limit, result.Limit)
		}

		// Verify items were correctly deserialized
		resultItems, ok := result.Items.(*[]TestItem)
		if !ok {
			t.Fatalf("Expected *[]TestItem, got %T", result.Items)
		}

		if len(*resultItems) != len(items) {
			t.Errorf("Expected %d items, got %d", len(items), len(*resultItems))
		}

		for i, item := range items {
			if (*resultItems)[i].ID != item.ID {
				t.Errorf("Item %d: expected ID %d, got %d", i, item.ID, (*resultItems)[i].ID)
			}
			if (*resultItems)[i].Name != item.Name {
				t.Errorf("Item %d: expected name %s, got %s", i, item.Name, (*resultItems)[i].Name)
			}
		}
	})

	t.Run("PageFromProto with default values", func(t *testing.T) {
		// Test proto with zero page/limit values that should be defaulted
		pageData := slicer.PageData{
			Items: []string{"item1", "item2"},
			Total: 25,
			Page:  0, // Should default to 1
			Limit: 0, // Should default to 10
		}

		proto, err := pageData.ToProto()
		if err != nil {
			t.Fatalf("ToProto failed: %v", err)
		}

		var destItems []string
		result, err := slicer.PageFromProto(proto, &destItems)
		if err != nil {
			t.Fatalf("PageFromProto failed: %v", err)
		}

		if result.Page != 1 {
			t.Errorf("Expected page 1 (defaulted), got %d", result.Page)
		}
		if result.Limit != 10 {
			t.Errorf("Expected limit 10 (defaulted), got %d", result.Limit)
		}
	})

	t.Run("Error handling in conversions", func(t *testing.T) {
		// Test PageData ToProto with nil items (should work)
		pageData := slicer.PageData{
			Items: nil,
			Total: 0,
			Page:  1,
			Limit: 10,
		}

		proto, err := pageData.ToProto()
		if err != nil {
			t.Errorf("ToProto with nil items should not error, got: %v", err)
		}
		if proto == nil {
			t.Error("ToProto should not return nil proto")
		}
	})
}
