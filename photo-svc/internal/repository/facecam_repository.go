package repository

import (
	"be-yourmoments/photo-svc/internal/entity"
	"fmt"
	"log"
)

type FacecamRepository interface {
	Create(tx Querier, facecam *entity.Facecam) (*entity.Facecam, error)
}

type facecamRepository struct {
}

func NewFacecamRepository() FacecamRepository {
	return &facecamRepository{}
}

func (r *facecamRepository) Create(tx Querier, facecam *entity.Facecam) (*entity.Facecam, error) {
	query := `INSERT INTO photos 
			  (id, creator_id, title, size, url, original_at, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := tx.Exec(query, facecam.Id, facecam.CreatorId, facecam.Title, facecam.Size,
		facecam.Url, facecam.OriginalAt, facecam.CreatedAt, facecam.UpdatedAt)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert facecam: %w", err)
	}

	return facecam, nil
}

// func (r *facecamRepository) Update(ctx context.Context, db Querier, req *model.RequestUpdatePhoto) error {

// 	query := `UPDATE INTO user_simillars (id, user_id, size, url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING (id)`

// 	err := db.Exec(ctx, query)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
