package entity

import "time"

type Photo struct {
	Id                string
	UserId            string
	Size              int64
	RawUrl            string
	PreviewIsUserUrl  string
	PreviewNotUserUrl string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
