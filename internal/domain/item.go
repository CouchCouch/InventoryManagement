package domain

type Item struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Notes         string    `json:"notes,omitempty"`
	DatePurchased string `json:"date_purchased,omitempty"`
}
