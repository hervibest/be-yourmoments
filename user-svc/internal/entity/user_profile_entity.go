package entity

import (
	"database/sql"
	"time"
)

type UserProfile struct {
	Id              string         `db:"id"`
	UserId          string         `db:"user_id"`
	BirthDate       *time.Time     `db:"birth_date"`
	Nickname        string         `db:"nickname"`
	Biography       sql.NullString `db:"biography"`
	ProfileUrl      sql.NullString `db:"profile_url"`
	ProfileCoverUrl sql.NullString `db:"profile_cover_url"`
	Similarity      sql.NullString `db:"similarity"`
	CreatedAt       *time.Time     `db:"created_at"`
	UpdatedAt       *time.Time     `db:"updated_at"`
}
