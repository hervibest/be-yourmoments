package entity

import "time"

type UserProfile struct {
	Id              string     `db:"id"`
	UserId          string     `db:"user_id"`
	BirthDate       *time.Time `db:"birth_date"`
	Nickname        string     `db:"nickname"`
	Biography       string     `db:"biography"`
	ProfileUrl      string     `db:"profile_url"`
	ProfileCoverUrl string     `db:"profile_cover_url"`
	Similarity      string     `db:"similarity"`
	CreatedAt       *time.Time `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
}
