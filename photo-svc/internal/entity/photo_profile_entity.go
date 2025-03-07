package entity

import "time"

type PhotoProfile struct {
	Id     string
	UserId string
	Size   int64
	RawUrl string

	CreatedAt time.Time
	UpdatedAt time.Time
}
