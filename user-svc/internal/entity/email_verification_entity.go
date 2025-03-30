package entity

import "time"

type EmailVerification struct {
	Email     string     `db:"email"`
	Token     string     `db:"token"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
