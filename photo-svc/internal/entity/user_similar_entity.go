package entity

import (
	"be-yourmoments/photo-svc/internal/enum"
	"time"
)

type UserSimilarPhoto struct {
	Id         string                   `db:"id"`
	PhotoId    string                   `db:"photo_id"`
	UserId     string                   `db:"user_id"`
	Similarity enum.SimilarityLevelEnum `db:"similarity"`
	IsWishlist bool                     `db:"is_wishlist"`
	IsResend   bool                     `db:"is_resend"`
	IsCart     bool                     `db:"is_cart"`
	IsFavorite bool                     `db:"is_favorite"`
	CreatedAt  time.Time                `db:"created_at"`
	UpdatedAt  time.Time                `db:"updated_at"`
}
