package slicer_test

import (
	"testing"
	"time"

	"github.com/godev90/slicer"
)

type AdvancedUser struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Age        int       `json:"age"`
	Salary     float64   `json:"salary"`
	CreatedAt  time.Time `json:"created_at"`
	IsActive   bool      `json:"is_active"`
	Department string    `json:"department"`
}

func TestSlicePaginatorAdvancedFiltering(t *testing.T) {
	baseTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	users := []AdvancedUser{
		{ID: 1, Name: "Alice", Email: "alice@example.com", Age: 25, Salary: 50000.50, CreatedAt: baseTime, IsActive: true, Department: "Engineering"},
		{ID: 2, Name: "Bob", Email: "bob@example.com", Age: 30, Salary: 60000.75, CreatedAt: baseTime.Add(24 * time.Hour), IsActive: false, Department: "Marketing"},
		{ID: 3, Name: "Charlie", Email: "charlie@example.com", Age: 35, Salary: 70000.25, CreatedAt: baseTime.Add(48 * time.Hour), IsActive: true, Department: "Engineering"},
		{ID: 4, Name: "Diana", Email: "diana@example.com", Age: 28, Salary: 55000.00, CreatedAt: baseTime.Add(72 * time.Hour), IsActive: true, Department: "Sales"},
		{ID: 5, Name: "Eve", Email: "eve@example.com", Age: 32, Salary: 65000.80, CreatedAt: baseTime.Add(96 * time.Hour), IsActive: false, Department: "Engineering"},
	}

	allowedFields := map[string]string{
		"id":         "id",
		"name":       "name",
		"email":      "email",
		"age":        "age",
		"salary":     "salary",
		"created_at": "created_at",
		"is_active":  "is_active",
		"department": "department",
	}

	t.Run("String comparisons", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		tests := []struct {
			name     string
			op       slicer.ComparisonOp
			field    string
			value    string
			expected int
		}{
			{"EQ string", slicer.EQ, "name", "Alice", 1},
			{"GT string", slicer.GT, "name", "Bob", 3},       // Charlie, Diana, Eve
			{"GTE string", slicer.GTE, "name", "Bob", 4},     // Bob, Charlie, Diana, Eve
			{"LT string", slicer.LT, "name", "Charlie", 2},   // Alice, Bob
			{"LTE string", slicer.LTE, "name", "Charlie", 3}, // Alice, Bob, Charlie
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				opts := slicer.QueryOptions{
					Page:  1,
					Limit: 50,
					Comparisons: []slicer.ComparisonFilter{
						{Field: tt.field, Op: tt.op, Value: tt.value},
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

				if len(resultUsers) != tt.expected {
					t.Errorf("Expected %d users, got %d", tt.expected, len(resultUsers))
				}
			})
		}
	})

	t.Run("Integer comparisons", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		tests := []struct {
			name     string
			op       slicer.ComparisonOp
			field    string
			value    string
			expected int
		}{
			{"EQ int", slicer.EQ, "age", "30", 1},   // Bob
			{"GT int", slicer.GT, "age", "30", 2},   // Charlie (35), Eve (32)
			{"GTE int", slicer.GTE, "age", "30", 3}, // Bob (30), Charlie (35), Eve (32)
			{"LT int", slicer.LT, "age", "30", 2},   // Alice (25), Diana (28)
			{"LTE int", slicer.LTE, "age", "30", 3}, // Alice (25), Bob (30), Diana (28)
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				opts := slicer.QueryOptions{
					Page:  1,
					Limit: 50,
					Comparisons: []slicer.ComparisonFilter{
						{Field: tt.field, Op: tt.op, Value: tt.value},
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

				if len(resultUsers) != tt.expected {
					t.Errorf("Expected %d users, got %d", tt.expected, len(resultUsers))
				}
			})
		}
	})

	t.Run("Float comparisons", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		tests := []struct {
			name     string
			op       slicer.ComparisonOp
			field    string
			value    string
			expected int
		}{
			{"EQ float", slicer.EQ, "salary", "60000.75", 1},   // Bob
			{"GT float", slicer.GT, "salary", "60000", 3},      // Bob (60000.75), Charlie (70000.25), Eve (65000.80)
			{"GTE float", slicer.GTE, "salary", "60000.75", 3}, // Bob (60000.75), Charlie (70000.25), Eve (65000.80)
			{"LT float", slicer.LT, "salary", "60000", 2},      // Alice (50000.50), Diana (55000.00)
			{"LTE float", slicer.LTE, "salary", "60000.75", 3}, // Alice (50000.50), Bob (60000.75), Diana (55000.00)
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				opts := slicer.QueryOptions{
					Page:  1,
					Limit: 50,
					Comparisons: []slicer.ComparisonFilter{
						{Field: tt.field, Op: tt.op, Value: tt.value},
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

				if len(resultUsers) != tt.expected {
					t.Errorf("Expected %d users, got %d", tt.expected, len(resultUsers))
				}
			})
		}
	})

	t.Run("Time comparisons", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		tests := []struct {
			name     string
			op       slicer.ComparisonOp
			field    string
			value    string
			expected int
		}{
			{"EQ time", slicer.EQ, "created_at", "2023-01-01T00:00:00Z", 1},   // Alice
			{"GT time", slicer.GT, "created_at", "2023-01-02T00:00:00Z", 3},   // Charlie, Diana, Eve
			{"GTE time", slicer.GTE, "created_at", "2023-01-02T00:00:00Z", 4}, // Bob, Charlie, Diana, Eve
			{"LT time", slicer.LT, "created_at", "2023-01-02T00:00:00Z", 1},   // Alice
			{"LTE time", slicer.LTE, "created_at", "2023-01-02T00:00:00Z", 2}, // Alice, Bob
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				opts := slicer.QueryOptions{
					Page:  1,
					Limit: 50,
					Comparisons: []slicer.ComparisonFilter{
						{Field: tt.field, Op: tt.op, Value: tt.value},
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

				if len(resultUsers) != tt.expected {
					t.Errorf("Expected %d users, got %d", tt.expected, len(resultUsers))
				}
			})
		}
	})
}

func TestSlicePaginatorAdvancedSorting(t *testing.T) {
	baseTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	users := []AdvancedUser{
		{ID: 3, Name: "Charlie", Age: 35, Salary: 70000.25, CreatedAt: baseTime.Add(48 * time.Hour)},
		{ID: 1, Name: "Alice", Age: 25, Salary: 50000.50, CreatedAt: baseTime},
		{ID: 5, Name: "Eve", Age: 32, Salary: 65000.80, CreatedAt: baseTime.Add(96 * time.Hour)},
		{ID: 2, Name: "Bob", Age: 30, Salary: 60000.75, CreatedAt: baseTime.Add(24 * time.Hour)},
		{ID: 4, Name: "Diana", Age: 28, Salary: 55000.00, CreatedAt: baseTime.Add(72 * time.Hour)},
	}

	allowedFields := map[string]string{
		"id":         "id",
		"name":       "name",
		"age":        "age",
		"salary":     "salary",
		"created_at": "created_at",
	}

	t.Run("String sorting", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		t.Run("Ascending", func(t *testing.T) {
			opts := slicer.QueryOptions{
					Page:  1,
					Limit: 50,
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

			expected := []string{"Alice", "Bob", "Charlie", "Diana", "Eve"}
			for i, user := range resultUsers {
				if user.Name != expected[i] {
					t.Errorf("Expected user %d to be %s, got %s", i, expected[i], user.Name)
				}
			}
		})

		t.Run("Descending", func(t *testing.T) {
			opts := slicer.QueryOptions{
					Page:  1,
					Limit: 50,
				Sort: []slicer.SortField{
					{Field: "name", Desc: true},
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

			expected := []string{"Eve", "Diana", "Charlie", "Bob", "Alice"}
			for i, user := range resultUsers {
				if user.Name != expected[i] {
					t.Errorf("Expected user %d to be %s, got %s", i, expected[i], user.Name)
				}
			}
		})
	})

	t.Run("Integer sorting", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		t.Run("Ascending by age", func(t *testing.T) {
			opts := slicer.QueryOptions{
					Page:  1,
					Limit: 50,
				Sort: []slicer.SortField{
					{Field: "age", Desc: false},
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

			expected := []int{25, 28, 30, 32, 35} // Alice, Diana, Bob, Eve, Charlie
			for i, user := range resultUsers {
				if user.Age != expected[i] {
					t.Errorf("Expected user %d age to be %d, got %d", i, expected[i], user.Age)
				}
			}
		})

		t.Run("Descending by age", func(t *testing.T) {
			opts := slicer.QueryOptions{
					Page:  1,
					Limit: 50,
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

			expected := []int{35, 32, 30, 28, 25} // Charlie, Eve, Bob, Diana, Alice
			for i, user := range resultUsers {
				if user.Age != expected[i] {
					t.Errorf("Expected user %d age to be %d, got %d", i, expected[i], user.Age)
				}
			}
		})
	})

	t.Run("Float sorting", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		t.Run("Ascending by salary", func(t *testing.T) {
			opts := slicer.QueryOptions{
					Page:  1,
					Limit: 50,
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

			// Expected order: Alice (50000.50), Diana (55000.00), Bob (60000.75), Eve (65000.80), Charlie (70000.25)
			expected := []float64{50000.50, 55000.00, 60000.75, 65000.80, 70000.25}
			for i, user := range resultUsers {
				if user.Salary != expected[i] {
					t.Errorf("Expected user %d salary to be %f, got %f", i, expected[i], user.Salary)
				}
			}
		})

		t.Run("Descending by salary", func(t *testing.T) {
			opts := slicer.QueryOptions{
					Page:  1,
					Limit: 50,
				Sort: []slicer.SortField{
					{Field: "salary", Desc: true},
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

			// Expected order: Charlie (70000.25), Eve (65000.80), Bob (60000.75), Diana (55000.00), Alice (50000.50)
			expected := []float64{70000.25, 65000.80, 60000.75, 55000.00, 50000.50}
			for i, user := range resultUsers {
				if user.Salary != expected[i] {
					t.Errorf("Expected user %d salary to be %f, got %f", i, expected[i], user.Salary)
				}
			}
		})
	})

	t.Run("Time sorting", func(t *testing.T) {
		paginator := slicer.NewSlicePaginator(users, allowedFields)

		t.Run("Ascending by created_at", func(t *testing.T) {
			opts := slicer.QueryOptions{
					Page:  1,
					Limit: 50,
				Sort: []slicer.SortField{
					{Field: "created_at", Desc: false},
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

			// Check that times are in ascending order
			for i := 0; i < len(resultUsers)-1; i++ {
				if resultUsers[i].CreatedAt.After(resultUsers[i+1].CreatedAt) {
					t.Errorf("Time sorting failed: user %d created at %v should be before user %d created at %v",
						i, resultUsers[i].CreatedAt, i+1, resultUsers[i+1].CreatedAt)
				}
			}
		})

		t.Run("Descending by created_at", func(t *testing.T) {
			opts := slicer.QueryOptions{
					Page:  1,
					Limit: 50,
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
	})

	t.Run("Multiple field sorting", func(t *testing.T) {
		// Add users with same age to test secondary sorting
		usersWithDuplicates := append(users, AdvancedUser{
			ID: 6, Name: "Frank", Age: 30, Salary: 45000.00, CreatedAt: baseTime.Add(120 * time.Hour),
		})

		paginator := slicer.NewSlicePaginator(usersWithDuplicates, allowedFields)

		opts := slicer.QueryOptions{
					Page:  1,
					Limit: 50,
			Sort: []slicer.SortField{
				{Field: "age", Desc: false},  // Primary sort by age ascending
				{Field: "name", Desc: false}, // Secondary sort by name ascending
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

		// Check that primary sort by age works
		for i := 0; i < len(resultUsers)-1; i++ {
			if resultUsers[i].Age > resultUsers[i+1].Age {
				t.Errorf("Primary sort failed: user %d age %d should not be greater than user %d age %d",
					i, resultUsers[i].Age, i+1, resultUsers[i+1].Age)
			}
		}

		// Check that users with age 30 are sorted by name (Bob, Frank)
		age30Users := []AdvancedUser{}
		for _, user := range resultUsers {
			if user.Age == 30 {
				age30Users = append(age30Users, user)
			}
		}

		if len(age30Users) == 2 {
			if age30Users[0].Name != "Bob" || age30Users[1].Name != "Frank" {
				t.Errorf("Secondary sort failed: expected Bob then Frank for age 30, got %s then %s",
					age30Users[0].Name, age30Users[1].Name)
			}
		}
	})
}

func TestSlicePaginatorErrorHandling(t *testing.T) {
	users := []AdvancedUser{
		{ID: 1, Name: "Alice", Age: 25},
	}

	allowedFields := map[string]string{
		"id":   "id",
		"name": "name",
		"age":  "age",
	}

	paginator := slicer.NewSlicePaginator(users, allowedFields)

	t.Run("Invalid comparison values", func(t *testing.T) {
		tests := []struct {
			name  string
			field string
			value string
		}{
			{"Invalid int comparison", "age", "not_a_number"},
			{"Invalid float comparison", "salary", "not_a_float"},
			{"Invalid time comparison", "created_at", "not_a_date"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				opts := slicer.QueryOptions{
					Page:  1,
					Limit: 50,
					Comparisons: []slicer.ComparisonFilter{
						{Field: tt.field, Op: slicer.EQ, Value: tt.value},
					},
				}

				result, err := slicer.SlicePage(paginator, opts)
				if err != nil {
					t.Fatalf("SlicePage returned error: %v", err)
				}

				// Should not error but should return no results or handle gracefully
				resultUsers, ok := result.Items.([]AdvancedUser)
				if !ok {
					t.Fatalf("Expected []AdvancedUser, got %T", result.Items)
				}

				// Invalid comparisons should result in no matches
				if len(resultUsers) > len(users) {
					t.Errorf("Expected <= %d users for invalid comparison, got %d", len(users), len(resultUsers))
				}
			})
		}
	})

	t.Run("Field lookup by JSON tag vs struct name", func(t *testing.T) {
		// Test that findFieldByColumn works with both JSON tags and struct field names
		opts := slicer.QueryOptions{
					Page:  1,
					Limit: 50,
			Sort: []slicer.SortField{
				{Field: "id", Desc: false}, // Should work via JSON tag
			},
			Comparisons: []slicer.ComparisonFilter{
				{Field: "name", Op: slicer.EQ, Value: "Alice"}, // Should work via JSON tag
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
			t.Errorf("Field lookup failed: expected 1 Alice, got %d users", len(resultUsers))
		}
	})
}
