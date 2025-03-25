package repository

import (
	"be-yourmoments/photo-svc/internal/entity"
	"fmt"
	"log"
)

type PhotoRepository interface {
	Create(tx Querier, photo *entity.Photo) (*entity.Photo, error)
	UpdateProcessedUrl(tx Querier, photo *entity.Photo) error
	// UpdateClaimedPhoto(ctx context.Context, db Querier, photo *entity.Photo) error
	// UpdatePhotoStatus(ctx context.Context, db Querier, photo *entity.Photo) error
}

type photoRepository struct {
}

func NewPhotoRepository() PhotoRepository {
	return &photoRepository{}
}

func (r *photoRepository) Create(tx Querier, photo *entity.Photo) (*entity.Photo, error) {
	query := `INSERT INTO photos 
			  (id, creator_id, title, collection_url, price, price_str, original_at, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := tx.Exec(query, photo.Id, photo.CreatorId, photo.Title, photo.CollectionUrl, photo.Price, photo.PriceStr,
		photo.OriginalAt, photo.CreatedAt, photo.UpdatedAt)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert photo: %w", err)
	}

	return photo, nil
}

func (r *photoRepository) UpdateProcessedUrl(tx Querier, photo *entity.Photo) error {
	query := `UPDATE photos 
			  SET is_this_you_url = $1, your_moments_url = $2, updated_at = $3 
			  WHERE id = $4`

	_, err := tx.Exec(query, photo.IsThisYouURL, photo.YourMomentsUrl, photo.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update photo: %w", err)
	}
	return nil
}

// func (r *photoRepository) UpdateClaimedPhoto(ctx context.Context, db Querier, photo *entity.Photo) error {
// 	query := `UPDATE photos
// 	          SET owned_by_user_id = :owned_by_user_id, updated_at = $3
// 	          WHERE id = $4`

// 	_, err := db.NamedExecContext(ctx, query, photo)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *photoRepository) UpdatePhotoStatus(ctx context.Context, db Querier, photo *entity.Photo) error {
// 	query := `UPDATE photos
// 	          SET status = $1, updated_at = $3
// 	          WHERE id = $4`

// 	_, err := db.Exec(ctx, query, photo.Status, photo.UpdatedAt, photo.Id)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
