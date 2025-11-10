package slicer_test

import (
	"testing"
	"time"

	"github.com/godev90/slicer"
)

func TestSlicePaginatorUtilityFunctions(t *testing.T) {
	baseTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	users := []AdvancedUser{
		{ID: 1, Name: "Alice", Age: 25, Salary: 50000.50, CreatedAt: baseTime},
		{ID: 2, Name: "Bob", Age: 30, Salary: 60000.75, CreatedAt: baseTime.Add(24 * time.Hour)},
		{ID: 3, Name: "Charlie", Age: 35, Salary: 70000.25, CreatedAt: baseTime.Add(48 * time.Hour)},
	}

	allowedFields := map[string]string{
		"id":         "id",
		"name":       "name",
		"age":        "age",
		"salary":     "salary",
		"created_at": "created_at",
	}

	// Test comparison functions through ComparisonFilter
	t.Run("String comparisons", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		// Test EQ comparison (compareString function)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "name", Op: slicer.EQ, Value: "Alice"},
			},
		}

		result, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatalf("SlicePage returned error: %v", err)
		}

		resultUsers, ok := result.Items.([]AdvancedUser)
		if !ok {
			t.Fatalf("Expected []AdvancedUser, got %T", result.Items)
		}

		if len(resultUsers) != 1 || resultUsers[0].Name != "Alice" {
			t.Errorf("String EQ comparison failed: expected 1 Alice, got %d users", len(resultUsers))
		}
	})

	t.Run("Integer comparisons", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		// Test GT comparison (compareInt64 function)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "age", Op: slicer.GT, Value: "30"},
			},
		}

		result, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatalf("SlicePage returned error: %v", err)
		}

		resultUsers, ok := result.Items.([]AdvancedUser)
		if !ok {
			t.Fatalf("Expected []AdvancedUser, got %T", result.Items)
		}

		if len(resultUsers) != 1 || resultUsers[0].Age != 35 {
			t.Errorf("Integer GT comparison failed: expected 1 user with age 35, got %d users", len(resultUsers))
		}
	})

	t.Run("Float comparisons", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		// Test GTE comparison (compareFloat64 function)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "salary", Op: slicer.GTE, Value: "60000.75"},
			},
		}

		result, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatalf("SlicePage returned error: %v", err)
		}

		resultUsers, ok := result.Items.([]AdvancedUser)
		if !ok {
			t.Fatalf("Expected []AdvancedUser, got %T", result.Items)
		}

		if len(resultUsers) != 2 {
			t.Errorf("Float GTE comparison failed: expected 2 users, got %d", len(resultUsers))
		}
	})

	t.Run("Time comparisons", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		// Test LT comparison (compareTime function)
		opts := slicer.QueryOptions{
			Page:  1,
			Limit: 10,
			Comparisons: []slicer.ComparisonFilter{
				{Field: "created_at", Op: slicer.LT, Value: "2023-01-02"},
			},
		}

		result, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatalf("SlicePage returned error: %v", err)
		}

		resultUsers, ok := result.Items.([]AdvancedUser)
		if !ok {
			t.Fatalf("Expected []AdvancedUser, got %T", result.Items)
		}

		if len(resultUsers) != 1 || resultUsers[0].Name != "Alice" {
			t.Errorf("Time LT comparison failed: expected 1 Alice, got %d users", len(resultUsers))
		}
	})

	// Test sorting functions through Sort fields
	t.Run("String sorting", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		// Test sortString function
		opts := slicer.QueryOptions{
			Sort: []slicer.SortField{
				{Field: "name", Desc: false},
			},
		}

		result, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatalf("SlicePage returned error: %v", err)
		}

		resultUsers, ok := result.Items.([]AdvancedUser)
		if !ok {
			t.Fatalf("Expected []AdvancedUser, got %T", result.Items)
		}

		expectedOrder := []string{"Alice", "Bob", "Charlie"}
		for i, user := range resultUsers {
			if user.Name != expectedOrder[i] {
				t.Errorf("String sorting failed: expected %s at position %d, got %s", expectedOrder[i], i, user.Name)
			}
		}
	})

	t.Run("Integer sorting", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		// Test sortInt64 function
		opts := slicer.QueryOptions{
			Sort: []slicer.SortField{
				{Field: "age", Desc: true},
			},
		}

		result, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatalf("SlicePage returned error: %v", err)
		}

		resultUsers, ok := result.Items.([]AdvancedUser)
		if !ok {
			t.Fatalf("Expected []AdvancedUser, got %T", result.Items)
		}

		expectedOrder := []int{35, 30, 25}
		for i, user := range resultUsers {
			if user.Age != expectedOrder[i] {
				t.Errorf("Integer sorting failed: expected %d at position %d, got %d", expectedOrder[i], i, user.Age)
			}
		}
	})

	t.Run("Float sorting", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		// Test sortFloat64 function
		opts := slicer.QueryOptions{
			Sort: []slicer.SortField{
				{Field: "salary", Desc: false},
			},
		}

		result, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatalf("SlicePage returned error: %v", err)
		}

		resultUsers, ok := result.Items.([]AdvancedUser)
		if !ok {
			t.Fatalf("Expected []AdvancedUser, got %T", result.Items)
		}

		expectedOrder := []float64{50000.50, 60000.75, 70000.25}
		for i, user := range resultUsers {
			if user.Salary != expectedOrder[i] {
				t.Errorf("Float sorting failed: expected %f at position %d, got %f", expectedOrder[i], i, user.Salary)
			}
		}
	})

	t.Run("Time sorting", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		// Test sortTime function
		opts := slicer.QueryOptions{
			Sort: []slicer.SortField{
				{Field: "created_at", Desc: true},
			},
		}

		result, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatalf("SlicePage returned error: %v", err)
		}

		resultUsers, ok := result.Items.([]AdvancedUser)
		if !ok {
			t.Fatalf("Expected []AdvancedUser, got %T", result.Items)
		}

		// Check that times are in descending order
		for i := 0; i < len(resultUsers)-1; i++ {
			if resultUsers[i].CreatedAt.Before(resultUsers[i+1].CreatedAt) {
				t.Errorf("Time sorting failed: user %d created at %v should be after user %d created at %v",
					i, resultUsers[i].CreatedAt, i+1, resultUsers[i+1].CreatedAt)
			}
		}
	})

	// Test findFieldByColumn function through field access
	t.Run("Field lookup by JSON tag", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		// Test findFieldByColumn function by accessing fields via JSON tags
		opts := slicer.QueryOptions{
			Sort: []slicer.SortField{
				{Field: "id", Desc: false}, // Should work via JSON tag
			},
		}

		result, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatalf("SlicePage returned error: %v", err)
		}

		resultUsers, ok := result.Items.([]AdvancedUser)
		if !ok {
			t.Fatalf("Expected []AdvancedUser, got %T", result.Items)
		}

		// Should be sorted by ID: 1, 2, 3
		expectedIDs := []int{1, 2, 3}
		for i, user := range resultUsers {
			if user.ID != expectedIDs[i] {
				t.Errorf("Field lookup failed: expected ID %d at position %d, got %d", expectedIDs[i], i, user.ID)
			}
		}
	})

	// Test compareSort function through multiple comparisons
	t.Run("Multiple field sorting", func(t *testing.T) {
		// Add users with same age to test compareSort
		usersWithDuplicates := []AdvancedUser{
			{ID: 1, Name: "Alice", Age: 30, Salary: 50000.50},
			{ID: 2, Name: "Bob", Age: 30, Salary: 60000.75},
			{ID: 3, Name: "Charlie", Age: 25, Salary: 70000.25},
		}

		paginator := slicer.NewSlicePaginator(usersWithDuplicates, allowedFields)

		opts := slicer.QueryOptions{
			Sort: []slicer.SortField{
				{Field: "age", Desc: false},    // Primary sort
				{Field: "salary", Desc: false}, // Secondary sort
			},
		}

		result, err := slicer.SlicePage(paginator, opts)
		if err != nil {
			t.Fatalf("SlicePage returned error: %v", err)
		}

		resultUsers, ok := result.Items.([]AdvancedUser)
		if !ok {
			t.Fatalf("Expected []AdvancedUser, got %T", result.Items)
		}

		// Charlie (25) should be first, then Alice (30, 50000.50), then Bob (30, 60000.75)
		expectedOrder := []string{"Charlie", "Alice", "Bob"}
		for i, user := range resultUsers {
			if user.Name != expectedOrder[i] {
				t.Errorf("Multiple field sorting failed: expected %s at position %d, got %s", expectedOrder[i], i, user.Name)
			}
		}
	})
}
