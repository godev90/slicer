package slicer_test

import (
	"testing"

	"github.com/godev90/slicer"
)

// TestSlicePaginatorCore tests core slice paginator functionality
func TestSlicePaginatorCore(t *testing.T) {
	// Test data - using simple structure to avoid conflicts
	type User struct {
		ID         int     `json:"id"`
		Name       string  `json:"name"`
		Email      string  `json:"email"`
		Age        int     `json:"age"`
		Salary     int     `json:"salary"`
		Department string  `json:"department"`
		Role       string  `json:"role"`
		Status     string  `json:"status"`
		City       string  `json:"city"`
		Country    string  `json:"country"`
		Rating     float64 `json:"rating"`
		Score      int     `json:"score"`
	}

	testUsers := []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 25, Salary: 50000, Department: "engineering", Role: "junior", Status: "active", City: "NYC", Country: "USA", Rating: 4.2, Score: 85},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 30, Salary: 75000, Department: "marketing", Role: "senior", Status: "active", City: "LA", Country: "USA", Rating: 4.8, Score: 92},
		{ID: 3, Name: "Bob Johnson", Email: "bob@example.com", Age: 35, Salary: 80000, Department: "engineering", Role: "senior", Status: "inactive", City: "Chicago", Country: "USA", Rating: 3.9, Score: 78},
		{ID: 4, Name: "Alice Brown", Email: "alice@example.com", Age: 28, Salary: 65000, Department: "design", Role: "mid", Status: "active", City: "Seattle", Country: "USA", Rating: 4.6, Score: 88},
		{ID: 5, Name: "Charlie Wilson", Email: "charlie@example.com", Age: 40, Salary: 95000, Department: "engineering", Role: "lead", Status: "active", City: "Austin", Country: "USA", Rating: 4.9, Score: 95},
	}

	allowedFields := map[string]string{
		"id":         "id",
		"name":       "name",
		"email":      "email",
		"age":        "age",
		"salary":     "salary",
		"department": "department",
		"role":       "role",
		"status":     "status",
		"city":       "city",
		"country":    "country",
		"rating":     "rating",
		"score":      "score",
	}

	t.Run("NewSlicePaginator", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(testUsers, allowedFields)

		if paginator == nil {
			t.Fatal("NewSlicePaginator returned nil")
		}

		// Test Items() method
		items := paginator.Items()
		if items == nil {
			t.Error("Items() returned nil for new paginator")
		}

		// Test SetItems() method
		newItems := []User{{ID: 99, Name: "Test User"}}
		paginator.SetItems(newItems)
		retrievedItems := paginator.Items()
		if len(retrievedItems) != 1 || retrievedItems[0].ID != 99 {
			t.Error("SetItems() did not work correctly")
		}
	})

	t.Run("Basic pagination", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(testUsers, allowedFields)

		tests := []struct {
			name        string
			page        int
			limit       int
			expectedLen int
		}{
			{"First page with limit 3", 1, 3, 3},
			{"Second page with limit 3", 2, 3, 2},
			{"Large limit returns all items", 1, 10, 5},
			{"Page beyond available data", 5, 3, 0},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				opts := slicer.QueryOptions{
					Page:  tt.page,
					Limit: tt.limit,
				}

				result, err := slicer.SlicePage(paginator, opts)
				if err != nil {
					t.Fatalf("SlicePage returned error: %v", err)
				}

				resultUsers, ok := result.Items.([]User)
				if !ok {
					t.Fatalf("Expected []User, got %T", result.Items)
				}

				if len(resultUsers) != tt.expectedLen {
					t.Errorf("Expected %d items, got %d", tt.expectedLen, len(resultUsers))
				}

				if result.Total != int64(len(testUsers)) {
					t.Errorf("Expected total %d, got %d", len(testUsers), result.Total)
				}
			})
		}
	})

	t.Run("Comparison filters", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(testUsers, allowedFields)

		tests := []struct {
			name        string
			comparisons []slicer.ComparisonFilter
			expectedLen int
		}{
			{
				name: "Age greater than or equal to 30",
				comparisons: []slicer.ComparisonFilter{
					{Field: "age", Op: slicer.GTE, Value: "30"},
				},
				expectedLen: 3, // Jane, Bob, Charlie
			},
			{
				name: "Salary less than 70000",
				comparisons: []slicer.ComparisonFilter{
					{Field: "salary", Op: slicer.LT, Value: "70000"},
				},
				expectedLen: 2, // John, Alice
			},
			{
				name: "Age equal to 25",
				comparisons: []slicer.ComparisonFilter{
					{Field: "age", Op: slicer.EQ, Value: "25"},
				},
				expectedLen: 1, // John
			},
			{
				name: "Rating greater than 4.5",
				comparisons: []slicer.ComparisonFilter{
					{Field: "rating", Op: slicer.GT, Value: "4.5"},
				},
				expectedLen: 3, // Jane, Alice, Charlie
			},
			{
				name: "Score less than or equal to 85",
				comparisons: []slicer.ComparisonFilter{
					{Field: "score", Op: slicer.LTE, Value: "85"},
				},
				expectedLen: 2, // John, Bob
			},
			{
				name: "Multiple filters: age >= 30 AND salary < 90000",
				comparisons: []slicer.ComparisonFilter{
					{Field: "age", Op: slicer.GTE, Value: "30"},
					{Field: "salary", Op: slicer.LT, Value: "90000"},
				},
				expectedLen: 2, // Jane, Bob
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				opts := slicer.QueryOptions{
					Page:        1,
					Limit:       10,
					Comparisons: tt.comparisons,
				}

				result, err := slicer.SlicePage(paginator, opts)
				if err != nil {
					t.Fatalf("SlicePage returned error: %v", err)
				}

				resultUsers, ok := result.Items.([]User)
				if !ok {
					t.Fatalf("Expected []User, got %T", result.Items)
				}

				if len(resultUsers) != tt.expectedLen {
					t.Errorf("Expected %d items, got %d", tt.expectedLen, len(resultUsers))
				}
			})
		}
	})

	t.Run("Custom filters", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(testUsers, allowedFields)

		tests := []struct {
			name        string
			filters     map[string]string
			expectedLen int
		}{
			{
				name:        "Filter by department",
				filters:     map[string]string{"department": "engineering"},
				expectedLen: 3, // John, Bob, Charlie
			},
			{
				name:        "Filter by status",
				filters:     map[string]string{"status": "active"},
				expectedLen: 4, // John, Jane, Alice, Charlie
			},
			{
				name:        "Filter by city",
				filters:     map[string]string{"city": "NYC"},
				expectedLen: 1, // John
			},
			{
				name: "Multiple filters",
				filters: map[string]string{
					"department": "engineering",
					"status":     "active",
				},
				expectedLen: 2, // John, Charlie
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				opts := slicer.QueryOptions{
					Page:    1,
					Limit:   10,
					Filters: tt.filters,
				}

				result, err := slicer.SlicePage(paginator, opts)
				if err != nil {
					t.Fatalf("SlicePage returned error: %v", err)
				}

				resultUsers, ok := result.Items.([]User)
				if !ok {
					t.Fatalf("Expected []User, got %T", result.Items)
				}

				if len(resultUsers) != tt.expectedLen {
					t.Errorf("Expected %d items, got %d", tt.expectedLen, len(resultUsers))
				}
			})
		}
	})

	t.Run("Search functionality", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(testUsers, allowedFields)

		tests := []struct {
			name        string
			search      *slicer.SearchQuery
			expectedLen int
		}{
			{
				name: "Search by name",
				search: &slicer.SearchQuery{
					Fields:  []string{"name"},
					Keyword: "john",
				},
				expectedLen: 2, // John Doe, Bob Johnson (contains "john")
			},
			{
				name: "Search by email",
				search: &slicer.SearchQuery{
					Fields:  []string{"email"},
					Keyword: "example.com",
				},
				expectedLen: 5, // All users have example.com emails
			},
			{
				name: "Search in multiple fields",
				search: &slicer.SearchQuery{
					Fields:  []string{"name", "department"},
					Keyword: "engineering",
				},
				expectedLen: 3, // Bob, John, Charlie (engineering department)
			},
			{
				name: "Search with no matches",
				search: &slicer.SearchQuery{
					Fields:  []string{"name"},
					Keyword: "nonexistent",
				},
				expectedLen: 0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				opts := slicer.QueryOptions{
					Page:   1,
					Limit:  10,
					Search: tt.search,
				}

				result, err := slicer.SlicePage(paginator, opts)
				if err != nil {
					t.Fatalf("SlicePage returned error: %v", err)
				}

				resultUsers, ok := result.Items.([]User)
				if !ok {
					t.Fatalf("Expected []User, got %T", result.Items)
				}

				if len(resultUsers) != tt.expectedLen {
					t.Errorf("Expected %d items, got %d", tt.expectedLen, len(resultUsers))
				}
			})
		}
	})

	t.Run("Sorting functionality", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(testUsers, allowedFields)

		tests := []struct {
			name              string
			sort              []slicer.SortField
			expectedFirstName string
			expectedLastName  string
		}{
			{
				name: "Sort by name ascending",
				sort: []slicer.SortField{
					{Field: "name", Desc: false},
				},
				expectedFirstName: "Alice Brown",
				expectedLastName:  "John Doe",
			},
			{
				name: "Sort by age descending",
				sort: []slicer.SortField{
					{Field: "age", Desc: true},
				},
				expectedFirstName: "Charlie Wilson", // age 40
				expectedLastName:  "John Doe",       // age 25
			},
			{
				name: "Sort by salary ascending",
				sort: []slicer.SortField{
					{Field: "salary", Desc: false},
				},
				expectedFirstName: "John Doe",       // salary 50000
				expectedLastName:  "Charlie Wilson", // salary 95000
			},
			{
				name: "Multiple sort fields",
				sort: []slicer.SortField{
					{Field: "department", Desc: false},
					{Field: "age", Desc: true},
				},
				expectedFirstName: "Alice Brown", // design department
				expectedLastName:  "Jane Smith",  // marketing department
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				opts := slicer.QueryOptions{
					Page:  1,
					Limit: 10,
					Sort:  tt.sort,
				}

				result, err := slicer.SlicePage(paginator, opts)
				if err != nil {
					t.Fatalf("SlicePage returned error: %v", err)
				}

				resultUsers, ok := result.Items.([]User)
				if !ok {
					t.Fatalf("Expected []User, got %T", result.Items)
				}

				if len(resultUsers) < 2 {
					t.Fatalf("Need at least 2 users for sort testing, got %d", len(resultUsers))
				}

				if resultUsers[0].Name != tt.expectedFirstName {
					t.Errorf("Expected first user to be %s, got %s", tt.expectedFirstName, resultUsers[0].Name)
				}

				if resultUsers[len(resultUsers)-1].Name != tt.expectedLastName {
					t.Errorf("Expected last user to be %s, got %s", tt.expectedLastName, resultUsers[len(resultUsers)-1].Name)
				}
			})
		}
	})

	t.Run("Edge cases", func(t *testing.T) {
		t.Run("Empty slice", func(t *testing.T) {
			emptyPaginator := slicer.NewSlicePaginator([]User{}, allowedFields)

			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
			}

			result, err := slicer.SlicePage(emptyPaginator, opts)
			if err != nil {
				t.Fatalf("SlicePage returned error: %v", err)
			}

			resultUsers, ok := result.Items.([]User)
			if !ok {
				t.Fatalf("Expected []User, got %T", result.Items)
			}

			if len(resultUsers) != 0 {
				t.Errorf("Expected 0 items for empty slice, got %d", len(resultUsers))
			}

			if result.Total != 0 {
				t.Errorf("Expected total 0 for empty slice, got %d", result.Total)
			}
		})

		t.Run("Filter by non-existent field", func(t *testing.T) {
			paginator := slicer.NewSlicePaginator(testUsers, allowedFields)

			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Comparisons: []slicer.ComparisonFilter{
					{Field: "nonexistent_field", Op: slicer.EQ, Value: "value"},
				},
			}

			result, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatalf("SlicePage returned error: %v", err)
			}

			resultUsers, ok := result.Items.([]User)
			if !ok {
				t.Fatalf("Expected []User, got %T", result.Items)
			}

			// Should return all users since non-existent field filter is ignored
			if len(resultUsers) != len(testUsers) {
				t.Errorf("Expected %d items, got %d", len(testUsers), len(resultUsers))
			}
		})

		t.Run("Invalid comparison values", func(t *testing.T) {
			paginator := slicer.NewSlicePaginator(testUsers, allowedFields)

			opts := slicer.QueryOptions{
				Page:  1,
				Limit: 10,
				Comparisons: []slicer.ComparisonFilter{
					{Field: "age", Op: slicer.GTE, Value: "invalid_number"},
				},
			}

			result, err := slicer.SlicePage(paginator, opts)
			if err != nil {
				t.Fatalf("SlicePage returned error: %v", err)
			}

			resultUsers, ok := result.Items.([]User)
			if !ok {
				t.Fatalf("Expected []User, got %T", result.Items)
			}

			// Should handle gracefully and return appropriate results
			// Exact behavior depends on implementation, but shouldn't crash
			t.Logf("Handled invalid comparison value gracefully, returned %d items", len(resultUsers))
		})
	})
}
