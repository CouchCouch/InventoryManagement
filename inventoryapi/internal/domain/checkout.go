package domain

import "time"

type Checkout struct {
	ID           int            `json:"id"`
	User         User           `json:"user"`
	Items        []CheckoutItem `json:"items"`
	CheckoutDate time.Time      `json:"checkout_date"`
	CreatedBy    int            `json:"created_by"`
	Notes        string         `json:"notes,omitempty"`
}

type CheckoutItem struct {
	Item       Item      `json:"item"`
	ReturnDate time.Time `json:"return_date"`
}
