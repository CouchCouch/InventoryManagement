package db

import "fmt"

// FieldType represents the data type of a field
type FieldType string

const (
	FieldTypeString    FieldType = "string"
	FieldTypeInt       FieldType = "int"
	FieldTypeTimestamp FieldType = "timestamp"
	FieldTypeBool      FieldType = "bool"
	FieldTypeUUID      FieldType = "uuid"
)

// FieldMetadata describes a queryable field
type FieldMetadata struct {
	Name       string
	Type       FieldType
	Filterable bool
	Sortable   bool
}

// TableRegistry manages field metadata for a table
type TableRegistry struct {
	TableName string
	Fields    map[string]FieldMetadata
}

// ValidateField checks if a field exists and is of the right type for the operation
func (tr *TableRegistry) ValidateField(fieldName string, operation string) error {
	field, exists := tr.Fields[fieldName]
	if !exists {
		return fmt.Errorf("invalid field '%s' for table '%s'", fieldName, tr.TableName)
	}

	if operation == "filter" && !field.Filterable {
		return fmt.Errorf("field '%s' is not filterable", fieldName)
	}
	if operation == "sort" && !field.Sortable {
		return fmt.Errorf("field '%s' is not sortable", fieldName)
	}

	return nil
}

// GetFieldType returns the type of a field
func (tr *TableRegistry) GetFieldType(fieldName string) (FieldType, error) {
	field, exists := tr.Fields[fieldName]
	if !exists {
		return "", fmt.Errorf("field '%s' not found", fieldName)
	}
	return field.Type, nil
}

// ItemsRegistry defines fields for the items table
var ItemsRegistry = &TableRegistry{
	TableName: "items i",
	Fields: map[string]FieldMetadata{
		"i.id": {
			Name:       "i.id",
			Type:       FieldTypeString,
			Filterable: true,
			Sortable:   true,
		},
		"i.name": {
			Name:       "i.name",
			Type:       FieldTypeString,
			Filterable: true,
			Sortable:   true,
		},
		"t.name": {
			Name:       "t.name",
			Type:       FieldTypeString,
			Filterable: true,
			Sortable:   true,
		},
		"i.notes": {
			Name:       "i.notes",
			Type:       FieldTypeString,
			Filterable: true,
			Sortable:   false,
		},
		"i.date_purchased": {
			Name:       "i.date_purchased",
			Type:       FieldTypeTimestamp,
			Filterable: true,
			Sortable:   true,
		},
		"i.deleted": {
			Name:       "i.deleted",
			Type:       FieldTypeBool,
			Filterable: true,
			Sortable:   false,
		},
	},
}

// CheckoutsRegistry defines fields for the checkouts table
var CheckoutsRegistry = &TableRegistry{
	TableName: "checkouts c",
	Fields: map[string]FieldMetadata{
		"c.id": {
			Name:       "c.id",
			Type:       FieldTypeInt,
			Filterable: true,
			Sortable:   true,
		},
		"c.checkout_date": {
			Name:       "c.checkout_date",
			Type:       FieldTypeTimestamp,
			Filterable: true,
			Sortable:   true,
		},
		"c.notes": {
			Name:       "c.notes",
			Type:       FieldTypeString,
			Filterable: true,
			Sortable:   false,
		},
		"c.personal": {
			Name:       "c.personal",
			Type:       FieldTypeBool,
			Filterable: true,
			Sortable:   true,
		},
		"u.name": {
			Name:       "u.name",
			Type:       FieldTypeString,
			Filterable: true,
			Sortable:   true,
		},
		"u.email": {
			Name:       "u.email",
			Type:       FieldTypeString,
			Filterable: true,
			Sortable:   true,
		},
		"u.id": {
			Name:       "u.id",
			Type:       FieldTypeUUID,
			Filterable: true,
			Sortable:   false,
		},
		"a.name": {
			Name:       "a.name",
			Type:       FieldTypeString,
			Filterable: true,
			Sortable:   true,
		},
		"a.email": {
			Name:       "a.email",
			Type:       FieldTypeString,
			Filterable: true,
			Sortable:   true,
		},
	},
}

// CheckoutItemsRegistry defines fields for the checkout_items table
var CheckoutItemsRegistry = &TableRegistry{
	TableName: "checkout_items ci",
	Fields: map[string]FieldMetadata{
		"ci.checkout_id": {
			Name:       "ci.checkout_id",
			Type:       FieldTypeInt,
			Filterable: true,
			Sortable:   true,
		},
		"ci.item_id": {
			Name:       "ci.item_id",
			Type:       FieldTypeString,
			Filterable: true,
			Sortable:   true,
		},
		"ci.return_date": {
			Name:       "ci.return_date",
			Type:       FieldTypeTimestamp,
			Filterable: true,
			Sortable:   false,
		},
		"i.name": {
			Name:       "i.name",
			Type:       FieldTypeString,
			Filterable: true,
			Sortable:   true,
		},
		"t.name": {
			Name:       "t.name",
			Type:       FieldTypeString,
			Filterable: true,
			Sortable:   true,
		},
		"i.notes": {
			Name:       "i.notes",
			Type:       FieldTypeString,
			Filterable: false,
			Sortable:   false,
		},
		"i.deleted": {
			Name:       "i.deleted",
			Type:       FieldTypeBool,
			Filterable: true,
			Sortable:   false,
		},
	},
}

// SafeQueryBuilder wraps QueryBuilder with field validation
type SafeQueryBuilder struct {
	builder  *QueryBuilder
	registry *TableRegistry
}

// NewSafeQueryBuilder creates a validated query builder
func NewSafeQueryBuilder(registry *TableRegistry, selectCols string) *SafeQueryBuilder {
	return &SafeQueryBuilder{
		builder:  NewQueryBuilder(registry.TableName, selectCols),
		registry: registry,
	}
}

// Filter adds a validated filter
func (sqb *SafeQueryBuilder) Filter(field string, op FilterOperator, value interface{}) (*SafeQueryBuilder, error) {
	if err := sqb.registry.ValidateField(field, "filter"); err != nil {
		return sqb, err
	}
	sqb.builder.Filter(field, op, value)
	return sqb, nil
}

// OrFilter adds a validated filter with OR logic
func (sqb *SafeQueryBuilder) OrFilter(field string, op FilterOperator, value interface{}) (*SafeQueryBuilder, error) {
	if err := sqb.registry.ValidateField(field, "filter"); err != nil {
		return sqb, err
	}
	sqb.builder.OrFilter(field, op, value)
	return sqb, nil
}

// Sort adds a validated sort field
func (sqb *SafeQueryBuilder) Sort(field string, direction SortDirection) (*SafeQueryBuilder, error) {
	if err := sqb.registry.ValidateField(field, "sort"); err != nil {
		return sqb, err
	}
	sqb.builder.Sort(field, direction)
	return sqb, nil
}

// Limit sets the LIMIT clause
func (sqb *SafeQueryBuilder) Limit(limit int) *SafeQueryBuilder {
	sqb.builder.Limit(limit)
	return sqb
}

// Offset sets the OFFSET clause
func (sqb *SafeQueryBuilder) Offset(offset int) *SafeQueryBuilder {
	sqb.builder.Offset(offset)
	return sqb
}

// AddJoin adds a JOIN clause
func (sqb *SafeQueryBuilder) AddJoin(joinClause string) *SafeQueryBuilder {
	sqb.builder.AddJoin(joinClause)
	return sqb
}

// SetBaseWhere sets a base WHERE clause
func (sqb *SafeQueryBuilder) SetBaseWhere(where string) *SafeQueryBuilder {
	sqb.builder.SetBaseWhere(where)
	return sqb
}

// Build returns the SQL and parameters
func (sqb *SafeQueryBuilder) Build() (string, []interface{}) {
	return sqb.builder.Build()
}
