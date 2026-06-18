package db

import (
	"strings"
	"testing"
)

func TestQueryBuilderBasicFilter(t *testing.T) {
	tests := []struct {
		name         string
		fn           func(*QueryBuilder)
		expectedLike string
		expectedVal  interface{}
	}{
		{
			name: "simple equality filter",
			fn: func(qb *QueryBuilder) {
				qb.Filter("id", OpEqual, "test-123")
			},
			expectedLike: "id = $1",
			expectedVal:  "test-123",
		},
		{
			name: "like filter",
			fn: func(qb *QueryBuilder) {
				qb.Filter("name", OpLike, "%test%")
			},
			expectedLike: "name LIKE $1",
			expectedVal:  "%test%",
		},
		{
			name: "greater than filter",
			fn: func(qb *QueryBuilder) {
				qb.Filter("age", OpGreater, 18)
			},
			expectedLike: "age > $1",
			expectedVal:  18,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder("items", "id, name")
			tt.fn(qb)
			query, params := qb.Build()

			if !strings.Contains(query, tt.expectedLike) {
				t.Errorf("Expected query to contain '%s', got: %s", tt.expectedLike, query)
			}
			if len(params) != 1 {
				t.Errorf("Expected 1 parameter, got %d", len(params))
			}
			if params[0] != tt.expectedVal {
				t.Errorf("Expected param '%v', got '%v'", tt.expectedVal, params[0])
			}
		})
	}
}

func TestQueryBuilderMultipleFilters(t *testing.T) {
	qb := NewQueryBuilder("items", "id, name")
	qb.Filter("deleted", OpEqual, false)
	qb.Filter("type_id", OpGreater, 0)

	query, params := qb.Build()

	if !strings.Contains(query, "deleted = $1") {
		t.Errorf("Expected first filter in query: %s", query)
	}
	if !strings.Contains(query, "type_id > $2") {
		t.Errorf("Expected second filter in query: %s", query)
	}
	if len(params) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(params))
	}
}

func TestQueryBuilderInOperator(t *testing.T) {
	qb := NewQueryBuilder("items", "id, name")
	ids := []interface{}{"id1", "id2", "id3"}
	qb.Filter("id", OpIn, ids)

	query, params := qb.Build()

	if !strings.Contains(query, "id IN ($1, $2, $3)") {
		t.Errorf("Expected IN clause in query: %s", query)
	}
	if len(params) != 3 {
		t.Errorf("Expected 3 parameters for IN clause, got %d", len(params))
	}
	for i, id := range ids {
		if params[i] != id {
			t.Errorf("Parameter %d mismatch: expected '%v', got '%v'", i, id, params[i])
		}
	}
}

func TestQueryBuilderSorting(t *testing.T) {
	qb := NewQueryBuilder("items", "id, name")
	qb.Sort("date_purchased", Desc)
	qb.Sort("name", Asc)

	query, _ := qb.Build()

	if !strings.Contains(query, "ORDER BY date_purchased DESC, name ASC") {
		t.Errorf("Expected ORDER BY clause in query: %s", query)
	}
}

func TestQueryBuilderLimitOffset(t *testing.T) {
	qb := NewQueryBuilder("items", "id, name")
	qb.Limit(10)
	qb.Offset(20)

	query, _ := qb.Build()

	if !strings.Contains(query, "LIMIT 10") {
		t.Errorf("Expected LIMIT 10 in query: %s", query)
	}
	if !strings.Contains(query, "OFFSET 20") {
		t.Errorf("Expected OFFSET 20 in query: %s", query)
	}
}

func TestQueryBuilderJoins(t *testing.T) {
	qb := NewQueryBuilder("items", "i.id, i.name, t.name")
	qb.AddJoin("LEFT JOIN item_types t ON i.item_type_id = t.id")

	query, _ := qb.Build()

	if !strings.Contains(query, "LEFT JOIN item_types t ON i.item_type_id = t.id") {
		t.Errorf("Expected JOIN clause in query: %s", query)
	}
}

func TestQueryBuilderBaseWhere(t *testing.T) {
	qb := NewQueryBuilder("items", "id, name")
	qb.SetBaseWhere("deleted = false")
	qb.Filter("name", OpLike, "%test%")

	query, _ := qb.Build()

	if !strings.Contains(query, "deleted = false") {
		t.Errorf("Expected base WHERE in query: %s", query)
	}
	if !strings.Contains(query, "name LIKE $1") {
		t.Errorf("Expected filter in query: %s", query)
	}
	// Both conditions should be combined with AND
	if !strings.Contains(query, "WHERE deleted = false AND name LIKE $1") {
		t.Errorf("Expected base WHERE and filter combined with AND: %s", query)
	}
}

func TestQueryBuilderIsNull(t *testing.T) {
	qb := NewQueryBuilder("items", "id, name")
	qb.Filter("deleted", OpIsNull, nil)

	query, _ := qb.Build()

	if !strings.Contains(query, "deleted IS NULL") {
		t.Errorf("Expected IS NULL in query: %s", query)
	}
}

func TestQueryBuilderIsNotNull(t *testing.T) {
	qb := NewQueryBuilder("items", "id, name")
	qb.Filter("deleted", OpIsNotNull, nil)

	query, _ := qb.Build()

	if !strings.Contains(query, "deleted IS NOT NULL") {
		t.Errorf("Expected IS NOT NULL in query: %s", query)
	}
}

func TestQueryBuilderNotIn(t *testing.T) {
	qb := NewQueryBuilder("items", "id, name")
	ids := []interface{}{"deleted1", "deleted2"}
	qb.Filter("id", OpNotIn, ids)

	query, params := qb.Build()

	if !strings.Contains(query, "id NOT IN ($1, $2)") {
		t.Errorf("Expected NOT IN clause in query: %s", query)
	}
	if len(params) != 2 {
		t.Errorf("Expected 2 parameters for NOT IN clause, got %d", len(params))
	}
}

func TestSafeQueryBuilderFieldValidation(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		operation     string
		shouldError   bool
		errorContains string
	}{
		{
			name:        "valid filter field",
			field:       "i.name",
			operation:   "filter",
			shouldError: false,
		},
		{
			name:        "valid sort field",
			field:       "i.date_purchased",
			operation:   "sort",
			shouldError: false,
		},
		{
			name:          "invalid field",
			field:         "i.invalid_field",
			operation:     "filter",
			shouldError:   true,
			errorContains: "invalid field",
		},
		{
			name:        "non-filterable field",
			field:       "i.notes",
			operation:   "filter",
			shouldError: false, // notes is filterable in ItemsRegistry
		},
		{
			name:          "non-sortable field",
			field:         "i.notes",
			operation:     "sort",
			shouldError:   true,
			errorContains: "not sortable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ItemsRegistry.ValidateField(tt.field, tt.operation)
			if tt.shouldError && err == nil {
				t.Errorf("Expected error but got nil")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Did not expect error but got: %v", err)
			}
			if tt.shouldError && tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
				t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
			}
		})
	}
}

func TestSafeQueryBuilderFilter(t *testing.T) {
	sqb := NewSafeQueryBuilder(ItemsRegistry, "i.id, i.name")
	sqb2, err := sqb.Filter("i.name", OpLike, "%test%")
	if err != nil {
		t.Errorf("Unexpected error on valid filter: %v", err)
	}
	if sqb2 != sqb {
		t.Errorf("Expected same builder instance returned")
	}

	query, _ := sqb2.Build()
	if !strings.Contains(query, "i.name LIKE $1") {
		t.Errorf("Expected filter in query: %s", query)
	}
}

func TestSafeQueryBuilderInvalidFilter(t *testing.T) {
	sqb := NewSafeQueryBuilder(ItemsRegistry, "i.id, i.name")
	_, err := sqb.Filter("invalid.field", OpEqual, "value")
	if err == nil {
		t.Errorf("Expected error for invalid field")
	}
	if !strings.Contains(err.Error(), "invalid field") {
		t.Errorf("Expected error to mention invalid field: %v", err)
	}
}

func TestSafeQueryBuilderSort(t *testing.T) {
	sqb := NewSafeQueryBuilder(ItemsRegistry, "i.id, i.name")
	sqb2, err := sqb.Sort("i.date_purchased", Desc)
	if err != nil {
		t.Errorf("Unexpected error on valid sort: %v", err)
	}
	if sqb2 != sqb {
		t.Errorf("Expected same builder instance returned")
	}

	query, _ := sqb2.Build()
	if !strings.Contains(query, "ORDER BY i.date_purchased DESC") {
		t.Errorf("Expected sort in query: %s", query)
	}
}

func TestSafeQueryBuilderInvalidSort(t *testing.T) {
	sqb := NewSafeQueryBuilder(ItemsRegistry, "i.id, i.name")
	_, err := sqb.Sort("i.notes", Asc) // notes is not sortable
	if err == nil {
		t.Errorf("Expected error for non-sortable field")
	}
	if !strings.Contains(err.Error(), "not sortable") {
		t.Errorf("Expected error to mention not sortable: %v", err)
	}
}

func TestQueryBuilderComplexQuery(t *testing.T) {
	qb := NewQueryBuilder("items", "i.id, i.name, t.name, i.notes")
	qb.AddJoin("LEFT JOIN item_types t ON i.item_type_id = t.id")
	qb.SetBaseWhere("i.deleted = false")
	qb.Filter("t.name", OpLike, "%Electronics%")
	qb.Filter("i.date_purchased", OpGreaterEq, "2024-01-01")
	qb.Sort("i.date_purchased", Desc)
	qb.Sort("i.name", Asc)
	qb.Limit(50)
	qb.Offset(10)

	query, params := qb.Build()

	// Verify all components are present
	if !strings.Contains(query, "LEFT JOIN item_types t") {
		t.Errorf("JOIN missing: %s", query)
	}
	if !strings.Contains(query, "i.deleted = false") {
		t.Errorf("Base WHERE missing: %s", query)
	}
	if !strings.Contains(query, "t.name LIKE") {
		t.Errorf("First filter missing: %s", query)
	}
	if !strings.Contains(query, "i.date_purchased >= $2") {
		t.Errorf("Second filter missing: %s", query)
	}
	if !strings.Contains(query, "ORDER BY i.date_purchased DESC, i.name ASC") {
		t.Errorf("ORDER BY missing or wrong: %s", query)
	}
	if !strings.Contains(query, "LIMIT 50") {
		t.Errorf("LIMIT missing: %s", query)
	}
	if !strings.Contains(query, "OFFSET 10") {
		t.Errorf("OFFSET missing: %s", query)
	}

	// Verify parameter count and values
	if len(params) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(params))
	}
	if params[0] != "%Electronics%" {
		t.Errorf("First param should be '%%Electronics%%', got '%v'", params[0])
	}
	if params[1] != "2024-01-01" {
		t.Errorf("Second param should be '2024-01-01', got '%v'", params[1])
	}
}

func TestQueryBuilderSQLInjectionPrevention(t *testing.T) {
	// Try to inject SQL through a filter value
	maliciousInput := "' OR '1'='1"
	qb := NewQueryBuilder("items", "id, name")
	qb.Filter("name", OpEqual, maliciousInput)

	query, params := qb.Build()

	// The SQL should have a parameter placeholder, not the malicious input
	if strings.Contains(query, "OR") && !strings.Contains(query, "$1") {
		t.Errorf("Potential SQL injection detected in query: %s", query)
	}

	// The parameter should be the malicious input unchanged, but safely parameterized
	if len(params) != 1 || params[0] != maliciousInput {
		t.Errorf("Parameter not properly preserved: %v", params)
	}
}

func TestOrFilterAfterFilter(t *testing.T) {
	qb := NewQueryBuilder("items", "id, name")
	qb.Filter("deleted", OpEqual, false)
	qb.OrFilter("type_id", OpEqual, 5)

	query, params := qb.Build()

	if !strings.Contains(query, "deleted = $1 AND (type_id = $2)") {
		t.Errorf("Expected AND with OR sub-group, got: %s", query)
	}
	if len(params) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(params))
	}
}

func TestOrFilterMultiple(t *testing.T) {
	qb := NewQueryBuilder("items", "id, name")
	qb.Filter("deleted", OpEqual, false)
	qb.OrFilter("type_id", OpEqual, 5)
	qb.OrFilter("type_id", OpEqual, 6)

	query, params := qb.Build()

	if !strings.Contains(query, "deleted = $1 AND (type_id = $2 OR type_id = $3)") {
		t.Errorf("Expected multi-condition OR sub-group, got: %s", query)
	}
	if len(params) != 3 {
		t.Errorf("Expected 3 parameters, got %d", len(params))
	}
}

func TestFilterClosesOrGroup(t *testing.T) {
	qb := NewQueryBuilder("items", "id, name")
	qb.OrFilter("type_id", OpEqual, 5)
	qb.Filter("deleted", OpEqual, false)

	query, params := qb.Build()

	// Filter conditions appear before sub-groups; AND is commutative so the order is cosmetic.
	if !strings.Contains(query, "deleted = $1 AND (type_id = $2)") {
		t.Errorf("Expected AND with OR sub-group, got: %s", query)
	}
	if len(params) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(params))
	}
}

func TestInterleavedFilterOrFilter(t *testing.T) {
	qb := NewQueryBuilder("items", "id, name")
	qb.Filter("a", OpEqual, 1)
	qb.OrFilter("b", OpEqual, 2)
	qb.OrFilter("c", OpEqual, 3)
	qb.Filter("d", OpEqual, 4)

	query, params := qb.Build()

	// Filter conditions appear before sub-groups; AND is commutative so ordering is cosmetic.
	if !strings.Contains(query, "a = $1 AND d = $2 AND (b = $3 OR c = $4)") {
		t.Errorf("Expected interleaved AND/OR, got: %s", query)
	}
	if len(params) != 4 {
		t.Errorf("Expected 4 parameters, got %d", len(params))
	}
}

func TestOrFilterOnly(t *testing.T) {
	qb := NewQueryBuilder("items", "id, name")
	qb.OrFilter("a", OpEqual, 1)
	qb.OrFilter("b", OpEqual, 2)

	query, params := qb.Build()

	if !strings.Contains(query, "(a = $1 OR b = $2)") {
		t.Errorf("Expected OR-only sub-group, got: %s", query)
	}
	if len(params) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(params))
	}
}

func TestSafeQueryBuilderOrFilter(t *testing.T) {
	sqb := NewSafeQueryBuilder(ItemsRegistry, "i.id, i.name")
	_, err := sqb.Filter("i.name", OpLike, "%test%")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	_, err = sqb.OrFilter("i.id", OpEqual, "123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	query, _ := sqb.Build()
	if !strings.Contains(query, "i.name LIKE $1 AND (i.id = $2)") {
		t.Errorf("Expected SafeQueryBuilder OR sub-group, got: %s", query)
	}
}
