package domain

import "errors"

const (
	// ErrCodeItemNotFound is returned when an item is not found in the inventory.
	ErrCodeItemNotFound = "item_not_found"
	// ErrCodeItemAlreadyExists is returned when an item already exists in the inventory.
	ErrCodeItemAlreadyExists = "item_already_exists"
	// ErrCodeInvalidItemID is returned when an item ID is invalid.
	ErrCodeInvalidItemID = "invalid_item_id"
	// ErrCodeInvalidItemType is returned when an item type is invalid.
	ErrCodeInvalidItemType = "invalid_item_type"
	// ErrCodeWrongSchema is returned when the database schema is incorrect.
	ErrCodeWrongSchema = "wrong_schema"
	// ErrCodeWrongPassword is returned when the login password is wrong
	ErrCodeWrongPassword = "wrong_password"
	// ErrCodeUserNotFound is returned when a user is not found in the database.
	ErrCodeUserNotFound = "user_not_found"
	// ErrCodeUserAlreadyExists is returned when a user already exists in the database.
	ErrCodeUserAlreadyExists = "user_already_exists"
	// ErrCodeInvalidSortField
	ErrCodeInvalidSortField = "invalid_sort_field"
	// ErrCodeInvalidFilterField
	ErrCodeInvalidFilterField = "invalid_filter_field"
	// ErrCodeInvalidHash is returned when a stored hash doesn't mathch the right format
	ErrCodeInvalidHash = "invalid_password_hash"
)

var (
	ErrItemNotFound       = errors.New(ErrCodeItemNotFound)
	ErrItemAlreadyExists  = errors.New(ErrCodeItemAlreadyExists)
	ErrInvalidItemID      = errors.New(ErrCodeInvalidItemID)
	ErrInvalidItemType    = errors.New(ErrCodeInvalidItemType)
	ErrWrongSchema        = errors.New(ErrCodeWrongSchema)
	ErrWrongPassword      = errors.New(ErrCodeWrongPassword)
	ErrUserNotFound       = errors.New(ErrCodeUserNotFound)
	ErrUserAlreadyExists  = errors.New(ErrCodeUserAlreadyExists)
	ErrInvalidSortField   = errors.New(ErrCodeInvalidSortField)
	ErrInvalidFilterField = errors.New(ErrCodeInvalidFilterField)
	ErrInvalidHash        = errors.New(ErrCodeInvalidHash)
)
