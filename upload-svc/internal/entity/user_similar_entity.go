package entity

import "time"

type UserSimilar struct {
	Id        string
	UserId    string
	PhotoId   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
