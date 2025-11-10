package slicer_test

import (
	"testing"

	slicerpb "github.com/godev90/slicer/pb"
)

// Test protocol buffer functions to boost coverage
func TestProtocolBufferFunctions(t *testing.T) {
	t.Run("QueryOptions protobuf methods", func(t *testing.T) {
		// Create a QueryOptions protobuf message
		opts := &slicerpb.QueryOptions{
			Page:  1,
			Limit: 10,
		}

		// Test getter methods
		if opts.GetPage() != 1 {
			t.Errorf("Expected page 1, got %d", opts.GetPage())
		}

		if opts.GetLimit() != 10 {
			t.Errorf("Expected limit 10, got %d", opts.GetLimit())
		}

		// Test String() method
		str := opts.String()
		if str == "" {
			t.Error("String() should not return empty string")
		}

		// Test ProtoMessage() method
		opts.ProtoMessage()

		// Test ProtoReflect() method
		reflection := opts.ProtoReflect()
		if reflection == nil {
			t.Error("ProtoReflect() should not return nil")
		}

		// Test Descriptor() method
		desc := opts.ProtoReflect().Descriptor()
		if desc == nil {
			t.Error("Descriptor() should not return nil")
		}

		// Test Reset() method
		opts.Reset()
		if opts.GetPage() != 0 || opts.GetLimit() != 0 {
			t.Error("Reset() should clear all fields")
		}
	})

	t.Run("SortField protobuf methods", func(t *testing.T) {
		sort := &slicerpb.SortField{
			Field: "name",
			Desc:  true,
		}

		// Test getter methods
		if sort.GetField() != "name" {
			t.Errorf("Expected field 'name', got %s", sort.GetField())
		}

		if !sort.GetDesc() {
			t.Error("Expected desc to be true")
		}

		// Test String(), ProtoMessage(), ProtoReflect(), Reset()
		str := sort.String()
		if str == "" {
			t.Error("String() should not return empty string")
		}

		sort.ProtoMessage()

		reflection := sort.ProtoReflect()
		if reflection == nil {
			t.Error("ProtoReflect() should not return nil")
		}

		sort.Reset()
		if sort.GetField() != "" || sort.GetDesc() {
			t.Error("Reset() should clear all fields")
		}
	})

	t.Run("SearchQuery protobuf methods", func(t *testing.T) {
		search := &slicerpb.SearchQuery{
			Fields:  []string{"name", "email"},
			Keyword: "test",
		}

		// Test getter methods
		fields := search.GetFields()
		if len(fields) != 2 || fields[0] != "name" || fields[1] != "email" {
			t.Errorf("Expected fields ['name', 'email'], got %v", fields)
		}

		if search.GetKeyword() != "test" {
			t.Errorf("Expected keyword 'test', got %s", search.GetKeyword())
		}

		// Test other methods
		search.String()
		search.ProtoMessage()
		search.ProtoReflect()
		search.Reset()
	})

	t.Run("ComparisonFilter protobuf methods", func(t *testing.T) {
		comp := &slicerpb.ComparisonFilter{
			Field: "age",
			Op:    "gt",
			Value: "18",
		}

		// Test getter methods
		if comp.GetField() != "age" {
			t.Errorf("Expected field 'age', got %s", comp.GetField())
		}

		if comp.GetOp() != "gt" {
			t.Errorf("Expected op 'gt', got %s", comp.GetOp())
		}

		if comp.GetValue() != "18" {
			t.Errorf("Expected value '18', got %s", comp.GetValue())
		}

		// Test other methods
		comp.String()
		comp.ProtoMessage()
		comp.ProtoReflect()
		comp.Reset()
	})

	t.Run("PageData protobuf methods", func(t *testing.T) {
		page := &slicerpb.PageData{
			Total: 100,
			Page:  1,
			Limit: 10,
		}

		// Test getter methods
		if page.GetTotal() != 100 {
			t.Errorf("Expected total 100, got %d", page.GetTotal())
		}

		if page.GetPage() != 1 {
			t.Errorf("Expected page 1, got %d", page.GetPage())
		}

		if page.GetLimit() != 10 {
			t.Errorf("Expected limit 10, got %d", page.GetLimit())
		}

		// Test other methods
		page.String()
		page.ProtoMessage()
		page.ProtoReflect()
		page.Reset()
	})

	t.Run("SearchQueryAnd protobuf methods", func(t *testing.T) {
		searchAnd := &slicerpb.SearchQueryAnd{
			Fields: []*slicerpb.SearchField{
				{Field: "name", Keyword: "john"},
				{Field: "department", Keyword: "engineering"},
			},
		}

		// Test getter methods
		fields := searchAnd.GetFields()
		if len(fields) != 2 {
			t.Errorf("Expected 2 fields, got %d", len(fields))
		}

		if fields[0].GetField() != "name" || fields[0].GetKeyword() != "john" {
			t.Errorf("Expected first field to be name:john, got %s:%s", fields[0].GetField(), fields[0].GetKeyword())
		}

		// Test other methods
		searchAnd.String()
		searchAnd.ProtoMessage()
		searchAnd.ProtoReflect()
		searchAnd.Reset()
	})

	t.Run("SearchField protobuf methods", func(t *testing.T) {
		field := &slicerpb.SearchField{
			Field:   "email",
			Keyword: "gmail",
		}

		// Test getter methods
		if field.GetField() != "email" {
			t.Errorf("Expected field 'email', got %s", field.GetField())
		}

		if field.GetKeyword() != "gmail" {
			t.Errorf("Expected keyword 'gmail', got %s", field.GetKeyword())
		}

		// Test other methods
		field.String()
		field.ProtoMessage()
		field.ProtoReflect()
		field.Reset()
	})
}
