package entity

import "time"

type Facecam struct {
	Id        string `db:"id"`
	CreatorId string `db:"creator_id"`
	Title     string `db:"title"`
	Size      int64  `db:"size"`
	Checksum  int    `db:"checksum"`
	Url       string `db:"url"`

	OriginalAt time.Time `db:"original_at"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
