package repository

import (
	"be-yourmoments/photo-svc/internal/entity"
	"context"
	"fmt"
)

type PhotoProfileRepository interface {
	Create(ctx context.Context, db Querier, photo *entity.PhotoProfile) error
	// Update(ctx context.Context, db Querier, req *model.RequestUpdatePhoto) error
}

type photoProfileRepository struct {
}

func NewPhotoProfileRepository() PhotoProfileRepository {
	return &photoProfileRepository{}
}

func (r *photoProfileRepository) Create(ctx context.Context, db Querier, photo *entity.PhotoProfile) error {

	query := `INSERT INTO photo_profiles (id, user_id, size, url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING (id)`

	err := db.QueryRow(ctx, query, photo.Id, photo.UserId, photo.Size, photo.RawUrl, photo.CreatedAt, photo.UpdatedAt).Scan(&photo.Id)
	if err != nil {
		return fmt.Errorf("failed to insert photo: %w", err)
	}

	return nil
}

// func (r *photoProfileRepository) Update(ctx context.Context, db Querier, req *model.RequestUpdatePhoto) error {

// 	query := `UPDATE INTO user_simillars (id, user_id, size, url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING (id)`

// 	err := db.Exec(ctx, query)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
