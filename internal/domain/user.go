package domain

import "github.com/google/uuid"

type User struct {
	ID    uuid.UUID
	Name  string
	Email string
}

type Admin struct {
	User     User
	Role     string
	Password string
}

type UserResponse struct {
	ID    uuid.UUID `json:"id,omitempty"`
	Name  string    `json:"name,omitempty"`
	Email string    `json:"email"`
}

type AdminResponse struct {
	User UserResponse `json:"user"`
	Role string       `json:"role,omitempty"`
}

type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
