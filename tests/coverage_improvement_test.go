package slicer_test

import (
	"testing"
	"time"

	"github.com/godev90/slicer"
)

// Test data structure for coverage improvement tests
type CoverageTestUser struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`        // int (supported)
	Score     float32   `json:"score"`      // float32 to test different float types
	Balance   float64   `json:"balance"`    // float64
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

// TestCoverageImprovements focuses on testing uncovered code paths
func TestCoverageImprovements(t *testing.T) {
	users := []CoverageTestUser{
		{ID: 1, Name: "Alice", Age: 25, Score: 85.5, Balance: 1000.0, Status: "active", CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), IsActive: true},
		{ID: 2, Name: "Bob", Age: 30, Score: 92.3, Balance: 1500.0, Status: "inactive", CreatedAt: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC), IsActive: false},
		{ID: 3, Name: "Charlie", Age: 35, Score: 78.9, Balance: 2000.0, Status: "pending", CreatedAt: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC), IsActive: true},
	}

	allowedFields := map[string]string{
		"id":         "id",
		"name":       "name",
		"age":        "age",
		"score":      "score",
		"balance":    "balance",
		"status":     "status",
		"created_at": "created_at",
		"is_active":  "is_active",
	}

	t.Run("Type conversion coverage tests", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		t.Run("int comparisons", func(t *testing.T) {
			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Comparisons: []slicer.ComparisonFilter{
					{Field: "age", Op: "eq", Value: "30"},
				},
			}

			_, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			items := paginator.Items()
			if len(items) != 1 {
				t.Errorf("Expected 1 user with age 30, got %d", len(items))
			}
		})

		t.Run("float32 comparisons", func(t *testing.T) {
			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Comparisons: []slicer.ComparisonFilter{
					{Field: "score", Op: "gt", Value: "80.0"},
				},
			}

			_, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			items := paginator.Items()
			if len(items) != 2 {
				t.Errorf("Expected 2 users with score > 80, got %d", len(items))
			}
		})

		t.Run("float64 comparisons", func(t *testing.T) {
			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Comparisons: []slicer.ComparisonFilter{
					{Field: "balance", Op: "gte", Value: "1500.0"},
				},
			}

			_, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			items := paginator.Items()
			if len(items) != 2 {
				t.Errorf("Expected 2 users with balance >= 1500, got %d", len(items))
			}
		})

		t.Run("time.Time comparisons", func(t *testing.T) {
			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Comparisons: []slicer.ComparisonFilter{
					{Field: "created_at", Op: "gt", Value: "2023-01-15T00:00:00Z"},
				},
			}

			_, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			items := paginator.Items()
			if len(items) != 2 {
				t.Errorf("Expected 2 users created after Jan 15, got %d", len(items))
			}
		})
	})

	t.Run("Sorting with different types", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		t.Run("Sort by int field", func(t *testing.T) {
			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Sort: []slicer.SortField{
					{Field: "age", Desc: false},
				},
			}

			_, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			items := paginator.Items()
			if len(items) != 3 {
				t.Errorf("Expected 3 users, got %d", len(items))
			}
			
			// Verify sort order
			for i := 1; i < len(items); i++ {
				if items[i-1].Age > items[i].Age {
					t.Error("Items not sorted by age ascending")
				}
			}
		})

		t.Run("Sort by float32 field", func(t *testing.T) {
			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Sort: []slicer.SortField{
					{Field: "score", Desc: true},
				},
			}

			_, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			items := paginator.Items()
			if len(items) != 3 {
				t.Errorf("Expected 3 users, got %d", len(items))
			}
			
			// Verify sort order (descending)
			for i := 1; i < len(items); i++ {
				if items[i-1].Score < items[i].Score {
					t.Error("Items not sorted by score descending")
				}
			}
		})

		t.Run("Sort by float64 field", func(t *testing.T) {
			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Sort: []slicer.SortField{
					{Field: "balance", Desc: false},
				},
			}

			_, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			items := paginator.Items()
			if len(items) != 3 {
				t.Errorf("Expected 3 users, got %d", len(items))
			}
			
			// Verify sort order
			for i := 1; i < len(items); i++ {
				if items[i-1].Balance > items[i].Balance {
					t.Error("Items not sorted by balance ascending")
				}
			}
		})

		t.Run("Sort by time.Time field", func(t *testing.T) {
			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Sort: []slicer.SortField{
					{Field: "created_at", Desc: true},
				},
			}

			_, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			items := paginator.Items()
			if len(items) != 3 {
				t.Errorf("Expected 3 users, got %d", len(items))
			}
			
			// Verify sort order (descending by date)
			for i := 1; i < len(items); i++ {
				if items[i-1].CreatedAt.Before(items[i].CreatedAt) {
					t.Error("Items not sorted by created_at descending")
				}
			}
		})
	})

	t.Run("Default case coverage tests", func(t *testing.T) {
		// Test struct with unsupported field types to trigger default cases
		type UnsupportedUser struct {
			ID      int                    `json:"id"`
			Age32   int32                  `json:"age32"`   // int32 - unsupported type
			Data    map[string]interface{} `json:"data"`    // Unsupported type
			Enabled chan bool              `json:"enabled"` // Definitely unsupported type
		}

		unsupportedUsers := []UnsupportedUser{
			{ID: 1, Age32: 25, Data: map[string]interface{}{"key": "value"}, Enabled: make(chan bool)},
			{ID: 2, Age32: 30, Data: map[string]interface{}{"key": "other"}, Enabled: make(chan bool)},
		}

		unsupportedFields := map[string]string{
			"id":      "id",
			"age32":   "age32",
			"data":    "data",
			"enabled": "enabled",
		}

		paginator := slicer.NewSlicePaginator(unsupportedUsers, unsupportedFields)

		t.Run("Unsupported int32 type comparison", func(t *testing.T) {
			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Comparisons: []slicer.ComparisonFilter{
					{Field: "age32", Op: "eq", Value: "25"},
				},
			}

			_, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			items := paginator.Items()
			// Should return empty since int32 type comparison hits default case and returns false
			if len(items) != 0 {
				t.Errorf("Expected 0 items for unsupported int32 type comparison, got %d", len(items))
			}
		})

		t.Run("Unsupported type comparison", func(t *testing.T) {
			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Comparisons: []slicer.ComparisonFilter{
					{Field: "data", Op: "eq", Value: "anything"},
				},
			}

			_, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			items := paginator.Items()
			// Should return empty since unsupported type comparison always returns false
			if len(items) != 0 {
				t.Errorf("Expected 0 items for unsupported type comparison, got %d", len(items))
			}
		})

		t.Run("Unsupported type sorting", func(t *testing.T) {
			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Sort: []slicer.SortField{
					{Field: "data", Desc: false},
				},
			}

			_, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			items := paginator.Items()
			// Should return all items but not sorted by the unsupported field
			if len(items) != 2 {
				t.Errorf("Expected 2 items for unsupported type sort, got %d", len(items))
			}
		})
	})

	t.Run("Type conversion edge cases", func(t *testing.T) {
		// Test different numeric types to improve coverage of toInt64, toFloat64, toString
		type NumericUser struct {
			ID      int     `json:"id"`
			Count   int64   `json:"count"`   // int64
			Rating  float32 `json:"rating"`  // float32
			Score   float64 `json:"score"`   // float64
			StrNum  string  `json:"str_num"` // string that can be parsed as number
		}

		numericUsers := []NumericUser{
			{ID: 1, Count: 100, Rating: 4.5, Score: 95.7, StrNum: "42"},
			{ID: 2, Count: 200, Rating: 3.8, Score: 88.2, StrNum: "invalid"},
		}

		numericFields := map[string]string{
			"id":      "id",
			"count":   "count",
			"rating":  "rating",
			"score":   "score",
			"str_num": "str_num",
		}

		paginator := slicer.NewSlicePaginator(numericUsers, numericFields)

		t.Run("int64 comparisons", func(t *testing.T) {
			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Comparisons: []slicer.ComparisonFilter{
					{Field: "count", Op: "lt", Value: "150"},
				},
			}

			_, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			items := paginator.Items()
			if len(items) != 1 {
				t.Errorf("Expected 1 user with count < 150, got %d", len(items))
			}
		})

		t.Run("Sort by int64", func(t *testing.T) {
			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Sort: []slicer.SortField{
					{Field: "count", Desc: true},
				},
			}

			_, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			items := paginator.Items()
			if len(items) != 2 {
				t.Errorf("Expected 2 users, got %d", len(items))
			}
			
			// Verify descending order
			if items[0].Count < items[1].Count {
				t.Error("Items not sorted by count descending")
			}
		})

		t.Run("String number parsing", func(t *testing.T) {
			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Comparisons: []slicer.ComparisonFilter{
					{Field: "str_num", Op: "eq", Value: "42"},
				},
			}

			_, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			items := paginator.Items()
			if len(items) != 1 {
				t.Errorf("Expected 1 user with str_num = '42', got %d", len(items))
			}
		})

		t.Run("String to number conversion edge cases", func(t *testing.T) {
			// Test different string numeric formats to improve toInt64 and toFloat64 coverage
			type StringNumUser struct {
				ID       int    `json:"id"`
				IntStr   string `json:"int_str"`   // Valid integer string
				FloatStr string `json:"float_str"` // Valid float string  
				BadStr   string `json:"bad_str"`   // Invalid string for number parsing
			}

			stringUsers := []StringNumUser{
				{ID: 1, IntStr: "123", FloatStr: "45.6", BadStr: "not_a_number"},
				{ID: 2, IntStr: "-456", FloatStr: "78.9", BadStr: "also_invalid"},
			}

			stringFields := map[string]string{
				"id":        "id",
				"int_str":   "int_str",
				"float_str": "float_str",
				"bad_str":   "bad_str",
			}

			stringPaginator := slicer.NewSlicePaginator(stringUsers, stringFields)

			// Test string parsing for integers
			optsInt := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Comparisons: []slicer.ComparisonFilter{
					{Field: "int_str", Op: "gt", Value: "0"},
				},
			}

			_, err := slicer.SlicePage(stringPaginator, optsInt)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			// Test string parsing for floats
			optsFloat := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Comparisons: []slicer.ComparisonFilter{
					{Field: "float_str", Op: "gte", Value: "40.0"},
				},
			}

			_, err = slicer.SlicePage(stringPaginator, optsFloat)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}

			// Test invalid string parsing (should trigger default cases)
			optsBad := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Comparisons: []slicer.ComparisonFilter{
					{Field: "bad_str", Op: "eq", Value: "not_a_number"},
				},
			}

			_, err = slicer.SlicePage(stringPaginator, optsBad)
			if err != nil {
				t.Fatal("SlicePage returned error:", err)
			}
		})
	})
}