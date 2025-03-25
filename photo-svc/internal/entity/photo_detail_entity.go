package entity

import (
	"be-yourmoments/photo-svc/internal/enum"
	"time"
)

type PhotoDetail struct {
	Id              string               `db:"id"`
	PhotoId         string               `db:"photo_id"`
	Size            int64                `db:"size"`
	Type            string               `db:"type"`
	Checksum        string               `db:"checksum"`
	Width           int8                 `db:"width"`
	Height          int8                 `db:"height"`
	Url             string               `db:"url"`
	YourMomentsType enum.YourMomentsType `db:"your_moments_type"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
