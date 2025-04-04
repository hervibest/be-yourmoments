package entity

import (
	"database/sql"
	"time"
)

type SocialMedia struct {
	Id          string         `db:"id"`
	Name        string         `db:"name"`
	BaseUrl     sql.NullString `db:"base_url"`
	LogoUrl     sql.NullString `db:"logo_url"`
	Description sql.NullString `db:"description"`
	IsActive    bool           `db:"is_active"`
	CreatedAt   *time.Time     `db:"created_at"`
	UpdatedAt   *time.Time     `db:"updated_at"`
}
