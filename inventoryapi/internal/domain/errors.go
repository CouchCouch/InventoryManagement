package domain

const (
	// ErrCodeItemNotFound is returned when an item is not found in the inventory.
	ErrCodeItemNotFound = "item_not_found"
	// ErrCodeItemAlreadyExists is returned when an item already exists in the inventory.
	ErrCodeItemAlreadyExists = "item_already_exists"
	// ErrCodeInvalidItemID is returned when an item ID is invalid.
	ErrCodeInvalidItemID = "invalid_item_id"
	// ErrCodeInvalidItemType is returned when an item type is invalid.
	ErrCodeInvalidItemType = "invalid_item_type"
)
