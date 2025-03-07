package repository

import (
	"be-yourmoments/photo-svc/internal/entity"
	"context"
	"fmt"
)

type PhotoRepository interface {
	Create(ctx context.Context, db Querier, photo *entity.Photo) error
	UpdateProcessedUrl(ctx context.Context, db Querier, photo *entity.Photo) error
	UpdateClaimedPhoto(ctx context.Context, db Querier, photo *entity.Photo) error
	UpdatePhotoStatus(ctx context.Context, db Querier, photo *entity.Photo) error
}

type photoRepository struct {
}

func NewPhotoRepository() PhotoRepository {
	return &photoRepository{}
}

func (r *photoRepository) Create(ctx context.Context, db Querier, photo *entity.Photo) error {

	query := `INSERT INTO photos 
			  (id, user_id, size, url, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err := db.QueryRow(ctx, query, photo.Id, photo.UserId, photo.Size, photo.RawUrl, photo.CreatedAt, photo.UpdatedAt).Scan(&photo.Id)

	if err != nil {
		return fmt.Errorf("failed to insert photo: %w", err)
	}

	return nil
}

func (r *photoRepository) UpdateProcessedUrl(ctx context.Context, db Querier, photo *entity.Photo) error {
	query := `UPDATE photos 
	          SET preview_url = $1, preview_with_bounding_url = $2, updated_at = $3 
	          WHERE id = $4`

	_, err := db.Exec(ctx, query, photo.PreviewUrl, photo.PreviewWithBoundingUrl, photo.UpdatedAt, photo.Id)

	if err != nil {
		return err
	}

	return nil
}

func (r *photoRepository) UpdateClaimedPhoto(ctx context.Context, db Querier, photo *entity.Photo) error {
	query := `UPDATE photos 
	          SET owned_by_user_id = $1, status = $2, updated_at = $3 
	          WHERE id = $4`

	_, err := db.Exec(ctx, query, photo.OwnedByUserId, photo.Status, photo.UpdatedAt, photo.Id)

	if err != nil {
		return err
	}

	return nil
}

func (r *photoRepository) UpdatePhotoStatus(ctx context.Context, db Querier, photo *entity.Photo) error {
	query := `UPDATE photos 
	          SET status = $1, updated_at = $3 
	          WHERE id = $4`

	_, err := db.Exec(ctx, query, photo.Status, photo.UpdatedAt, photo.Id)

	if err != nil {
		return err
	}

	return nil
}
