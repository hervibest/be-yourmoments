package model

import (
	"time"
)

type MinioFileResponse struct {
	ChecksumCRC32  string
	ChecksumCRC32C string
	ChecksumSHA1   string
	ChecksumSHA256 string
	ETag           string
	Expiration     time.Time
	URL            string
	FileKey        string
	Filename       string
	Mimetype       string
	Size           int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
