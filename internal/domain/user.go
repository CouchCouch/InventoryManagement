package domain

import "github.com/google/uuid"

type User struct {
	ID    uuid.UUID `json:"id,omitempty"`
	Name  string    `json:"name,omitempty"`
	Email string    `json:"email"`
}

type Admin struct {
	User     User   `json:"user"`
	Role     string `json:"role,omitempty"`
	Password string `json:"password"`
}
