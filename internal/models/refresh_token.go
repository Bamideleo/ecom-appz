package models

import "time"

type RefreshToken struct {
	ID        int
	UserId    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}