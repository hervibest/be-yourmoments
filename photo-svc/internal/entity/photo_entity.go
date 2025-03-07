package entity

import "time"

type Photo struct {
	Id                     string
	UserId                 string
	OwnedByUserId          string
	Status                 string
	Size                   int64
	RawUrl                 string
	PreviewUrl             string
	PreviewWithBoundingUrl string

	CreatedAt time.Time
	UpdatedAt time.Time
}
