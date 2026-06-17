package db

import (
	"fmt"
	"strings"
)

// FilterOperator represents comparison operators for filters
type FilterOperator string

const (
	OpEqual     FilterOperator = "="
	OpNotEqual  FilterOperator = "!="
	OpGreater   FilterOperator = ">"
	OpLess      FilterOperator = "<"
	OpGreaterEq FilterOperator = ">="
	OpLessEq    FilterOperator = "<="
	OpLike      FilterOperator = "LIKE"
	OpNotLike   FilterOperator = "NOT LIKE"
	OpIn        FilterOperator = "IN"
	OpNotIn     FilterOperator = "NOT IN"
	OpIsNull    FilterOperator = "IS NULL"
	OpIsNotNull FilterOperator = "IS NOT NULL"
)

// LogicOperator represents how to combine conditions
type LogicOperator string

const (
	OpAnd LogicOperator = "AND"
	OpOr  LogicOperator = "OR"
)

// SortDirection represents sort order
type SortDirection string

const (
	Asc  SortDirection = "ASC"
	Desc SortDirection = "DESC"
)

// FilterCondition represents a single filter condition
type FilterCondition struct {
	Field    string
	Operator FilterOperator
	Value    any
}

// FilterGroup represents a group of conditions with a logic operator
type FilterGroup struct {
	Conditions []FilterCondition
	Groups     []FilterGroup
	LogicOp    LogicOperator
	IsAnd      bool // true for AND, false for OR
}

// SortField represents a sort specification
type SortField struct {
	Field     string
	Direction SortDirection
}

// QueryBuilder provides a safe, fluent API for building queries
type QueryBuilder struct {
	tableName   string
	selectCols  string
	joinClauses []string
	baseWhere   string
	filterGroup *FilterGroup
	sortFields  []SortField
	limit       int
	offset      int
	params      []any
	paramCount  int
}

// NewQueryBuilder creates a new query builder for a table
func NewQueryBuilder(tableName string, selectCols string) *QueryBuilder {
	return &QueryBuilder{
		tableName:  tableName,
		selectCols: selectCols,
		filterGroup: &FilterGroup{
			Conditions: []FilterCondition{},
			Groups:     []FilterGroup{},
			IsAnd:      true,
		},
		sortFields: []SortField{},
		limit:      -1,
		offset:     0,
		params:     []any{},
		paramCount: 0,
	}
}

// Filter adds a filter condition to the query (AND logic by default at top level)
func (qb *QueryBuilder) Filter(field string, op FilterOperator, value any) *QueryBuilder {
	qb.filterGroup.Conditions = append(qb.filterGroup.Conditions, FilterCondition{
		Field:    field,
		Operator: op,
		Value:    value,
	})
	return qb
}

// OrFilter adds a filter condition with OR logic
func (qb *QueryBuilder) OrFilter(field string, op FilterOperator, value any) *QueryBuilder {
	// If this is the first filter with OR logic, convert to OR group
	if len(qb.filterGroup.Conditions) > 0 && qb.filterGroup.IsAnd {
		// Move existing conditions to a sub-group
		mainGroup := FilterGroup{
			Conditions: qb.filterGroup.Conditions,
			Groups:     qb.filterGroup.Groups,
			IsAnd:      true,
		}
		qb.filterGroup = &FilterGroup{
			Conditions: []FilterCondition{},
			Groups:     []FilterGroup{mainGroup},
			IsAnd:      false,
		}
	} else if len(qb.filterGroup.Conditions) == 0 && len(qb.filterGroup.Groups) == 0 {
		qb.filterGroup.IsAnd = false
	}

	qb.filterGroup.Conditions = append(qb.filterGroup.Conditions, FilterCondition{
		Field:    field,
		Operator: op,
		Value:    value,
	})
	return qb
}

// Sort adds a sort field to the query
func (qb *QueryBuilder) Sort(field string, direction SortDirection) *QueryBuilder {
	qb.sortFields = append(qb.sortFields, SortField{
		Field:     field,
		Direction: direction,
	})
	return qb
}

// Limit sets the LIMIT clause
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

// Offset sets the OFFSET clause
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

// AddJoin adds a JOIN clause
func (qb *QueryBuilder) AddJoin(joinClause string) *QueryBuilder {
	qb.joinClauses = append(qb.joinClauses, joinClause)
	return qb
}

// SetBaseWhere sets a base WHERE clause that's always applied
func (qb *QueryBuilder) SetBaseWhere(where string) *QueryBuilder {
	qb.baseWhere = where
	return qb
}

// Build constructs the final SQL query and returns the query string and parameters
func (qb *QueryBuilder) Build() (string, []any) {
	qb.params = []any{}
	qb.paramCount = 0

	var query strings.Builder
	query.WriteString("SELECT ")
	query.WriteString(qb.selectCols)
	query.WriteString(" FROM ")
	query.WriteString(qb.tableName)

	// Add JOIN clauses
	for _, join := range qb.joinClauses {
		query.WriteString(" ")
		query.WriteString(join)
	}

	// Build WHERE clause
	whereParts := []string{}
	if qb.baseWhere != "" {
		whereParts = append(whereParts, qb.baseWhere)
	}

	if len(qb.filterGroup.Conditions) > 0 || len(qb.filterGroup.Groups) > 0 {
		filterClause := qb.buildFilterGroup(qb.filterGroup)
		if filterClause != "" {
			whereParts = append(whereParts, filterClause)
		}
	}

	if len(whereParts) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(whereParts, " AND "))
	}

	// Build ORDER BY clause
	if len(qb.sortFields) > 0 {
		query.WriteString(" ORDER BY ")
		sortParts := []string{}
		for _, sf := range qb.sortFields {
			sortParts = append(sortParts, fmt.Sprintf("%s %s", sf.Field, sf.Direction))
		}
		query.WriteString(strings.Join(sortParts, ", "))
	}

	// Add LIMIT clause
	if qb.limit > 0 {
		fmt.Fprintf(&query, " LIMIT %d", qb.limit)
	}

	// Add OFFSET clause
	if qb.offset > 0 {
		fmt.Fprintf(&query, " OFFSET %d", qb.offset)
	}

	query.WriteString(";")

	return query.String(), qb.params
}

// buildFilterGroup recursively builds filter conditions
func (qb *QueryBuilder) buildFilterGroup(group *FilterGroup) string {
	var parts []string

	// Add individual conditions from this group
	for _, cond := range group.Conditions {
		condStr := qb.buildCondition(cond)
		if condStr != "" {
			parts = append(parts, condStr)
		}
	}

	// Add sub-groups
	for _, subGroup := range group.Groups {
		subStr := qb.buildFilterGroup(&subGroup)
		if subStr != "" {
			parts = append(parts, fmt.Sprintf("(%s)", subStr))
		}
	}

	if len(parts) == 0 {
		return ""
	}

	logicOp := "AND"
	if !group.IsAnd {
		logicOp = "OR"
	}

	return strings.Join(parts, fmt.Sprintf(" %s ", logicOp))
}

// buildCondition builds a single filter condition
func (qb *QueryBuilder) buildCondition(cond FilterCondition) string {
	switch cond.Operator {
	case OpIsNull:
		return fmt.Sprintf("%s IS NULL", cond.Field)
	case OpIsNotNull:
		return fmt.Sprintf("%s IS NOT NULL", cond.Field)
	case OpIn, OpNotIn:
		// For IN/NOT IN, value should be a slice
		if vals, ok := cond.Value.([]any); ok {
			placeholders := []string{}
			for _, val := range vals {
				qb.paramCount++
				placeholders = append(placeholders, fmt.Sprintf("$%d", qb.paramCount))
				qb.params = append(qb.params, val)
			}
			if len(placeholders) == 0 {
				return ""
			}
			return fmt.Sprintf("%s %s (%s)", cond.Field, cond.Operator, strings.Join(placeholders, ", "))
		}
		return ""
	default:
		// For standard operators, use a parameter
		qb.paramCount++
		qb.params = append(qb.params, cond.Value)
		return fmt.Sprintf("%s %s $%d", cond.Field, cond.Operator, qb.paramCount)
	}
}

// GetParams returns the accumulated parameters
func (qb *QueryBuilder) GetParams() []any {
	return qb.params
}
