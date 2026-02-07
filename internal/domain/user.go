package domain

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID    `json:"id"`
	Name  string `json:"name"`
	Email     string `json:"email"`
}

type Admin struct {
	User     User `json:"user"`
	Role     string `json:"role"`
	Password string `json:"password"`
}
