package entity

import (
	"time"
)

type AccessToken struct {
	UserId    string
	Token     string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	ExpiresAt time.Time
}

type RefreshToken struct {
	UserId    string
	Token     string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	ExpiresAt time.Time
}

type EmployeeAccessToken struct {
	EmployeeUUID string
	Token        string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	ExpiresAt    time.Time
}
