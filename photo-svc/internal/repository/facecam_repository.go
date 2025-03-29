package repository

import (
	"be-yourmoments/photo-svc/internal/entity"
	"fmt"
	"log"
)

type FacecamRepository interface {
	Create(tx Querier, facecam *entity.Facecam) (*entity.Facecam, error)
	UpdatedProcessedFacecam(tx Querier, facecam *entity.Facecam) error
}

type facecamRepository struct {
}

func NewFacecamRepository() FacecamRepository {
	return &facecamRepository{}
}

func (r *facecamRepository) Create(tx Querier, facecam *entity.Facecam) (*entity.Facecam, error) {
	query := `INSERT INTO facecams 
			  (id, user_id, file_name, file_key, title, size, checksum, url, original_at, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := tx.Exec(query, facecam.Id, facecam.UserId, facecam.FileName, facecam.FileKey, facecam.Title, facecam.Size,
		facecam.Checksum, facecam.Url, facecam.OriginalAt, facecam.CreatedAt, facecam.UpdatedAt)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert facecam: %w", err)
	}

	return facecam, nil
}

func (r *facecamRepository) UpdatedProcessedFacecam(tx Querier, facecam *entity.Facecam) error {
	log.Println("Updated accesed")
	query := `UPDATE facecams 
			  SET is_processed = $1, updated_at = $2 
			  WHERE user_id = $3`

	_, err := tx.Exec(query, facecam.IsProcessed, facecam.UpdatedAt, facecam.UserId)
	if err != nil {
		return fmt.Errorf("failed to update facecam: %w", err)
	}
	return nil
}

// func (r *facecamRepository) Update(ctx context.Context, db Querier, req *model.RequestUpdatePhoto) error {

// 	query := `UPDATE INTO user_simillars (id, user_id, size, url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING (id)`

// 	err := db.Exec(ctx, query)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
