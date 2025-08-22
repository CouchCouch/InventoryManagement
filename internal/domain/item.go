package domain

type Item struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Identifiers string `json:"identifiers"`
}
