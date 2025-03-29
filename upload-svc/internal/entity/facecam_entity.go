package entity

import "time"

type Facecam struct {
	Id       string `db:"id"`
	UserId   string `db:"creator_id"`
	FileName string `db:"file_name"`
	FileKey  string `db:"file_key"`
	Title    string `db:"title"`
	Size     int64  `db:"size"`
	Checksum string `db:"checksum"`
	Url      string `db:"url"`

	OriginalAt time.Time `db:"original_at"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
