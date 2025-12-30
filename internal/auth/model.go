package auth

import "time"

type User struct {
	ID           string
	Phone        string
	PasswordHash string
	Status       string
	CreatedAt    time.Time
}

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
}
