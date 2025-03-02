package repository

import (
	"be-yourmoments/internal/entity"
	"context"
	"fmt"
)

type PhotoRepository interface {
	Create(ctx context.Context, db Querier, photo *entity.Photo) error
	Update(ctx context.Context, db Querier, photoId string) error
}

type photoRepository struct {
}

func NewPhotoRepository() PhotoRepository {
	return &photoRepository{}
}

func (r *photoRepository) Create(ctx context.Context, db Querier, photo *entity.Photo) error {

	query := `INSERT INTO photos (id, user_id, size, url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING (id)`

	err := db.QueryRow(ctx, query, photo.Id, photo.UserId, photo.Size, photo.RawUrl, photo.CreatedAt, photo.UpdatedAt).Scan(&photo.Id)
	if err != nil {
		return fmt.Errorf("failed to insert photo: %w", err)
	}

	return nil
}

func (r *photoRepository) Update(ctx context.Context, db Querier, photoId string) error {

	query := `UPDATE INTO photos (id, user_id, size, url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING (id)`

	err := db.QueryRow(ctx, query, ).Scan(&)
	if err != nil {
		return err
	}

	return nil
}
