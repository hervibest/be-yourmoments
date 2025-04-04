package entity

import (
	"be-yourmoments/user-svc/internal/enum"
	"database/sql"
	"time"
)

type UserImage struct {
	Id            string             `db:"id"`
	UserProfileId string             `db:"user_profile_id"`
	FileName      string             `db:"file_name"`
	FileKey       string             `db:"file_key"`
	ImageType     enum.ImageTypeEnum `db:"image_type"`
	Size          int64              `db:"size"`
	Checksum      sql.NullString     `db:"checksum"`
	Url           sql.NullString     `db:"url"`
	CreatedAt     *time.Time         `db:"created_at"`
	UpdatedAt     *time.Time         `db:"updated_at"`
}
