package model

import (
	"time"
)

type RequestGetUserProfile struct {
	UserId string `json:"user_id" validate:"required"`
}

type RequestUpdateUserProfile struct {
	UserId    string     `validate:"required"`
	BirthDate *time.Time `json:"birth_date" validate:"required"`
	Nickname  string     `json:"nickname" validate:"required"`
	Biography string     `json:"biography" validate:"required"`
}

type UserProfileResponse struct {
	Id              string     `json:"id"`
	UserId          string     `json:"user_id"`
	BirthDate       *time.Time `json:"birth_date,omitempty"`
	Nickname        string     `json:"nickname"`
	Biography       string     `json:"biography"`
	ProfileUrl      string     `json:"profile_url"`
	ProfileCoverUrl string     `json:"profile_cover_url"`
	Similarity      string     `json:"similarity"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}
