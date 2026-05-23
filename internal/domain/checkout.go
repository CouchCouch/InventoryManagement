package domain

import (
	"time"
)

type Checkout struct {
	ID           int
	User         User
	Items        []CheckoutItem
	CheckoutDate time.Time
	CreatedBy    User
	Notes        string
}

type CheckoutItem struct {
	Item       Item      `json:"item"`
	ReturnDate time.Time `json:"return_date"`
}

type CreateCheckoutRequest struct {
	UserEmail    string    `json:"user_email"`
	Items        []string  `json:"items"`
	CheckoutDate string    `json:"checkout_date"`
	CreatedBy    string    `json:"created_by"`
	Notes        string    `json:"notes,omitempty"`
}

type CheckoutResponse struct {
	ID           int            `json:"id"`
	User         UserResponse   `json:"user"`
	Items        []CheckoutItem `json:"items"`
	CheckoutDate time.Time      `json:"checkout_date"`
	CreatedBy    UserResponse   `json:"created_by"`
	Notes        string         `json:"notes,omitempty"`
}
