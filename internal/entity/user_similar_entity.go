package entity

import "time"

type UserSimilarEntity struct {
	Id        string
	UserId    string
	PhotoId   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
