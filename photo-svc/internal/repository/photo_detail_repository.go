package repository

import (
	"be-yourmoments/photo-svc/internal/entity"
	"fmt"
	"log"
)

type PhotoDetailRepository interface {
	Create(tx Querier, photoDetail *entity.PhotoDetail) (*entity.PhotoDetail, error)
}

type photoDetailRepository struct {
}

func NewPhotoDetailRepository() PhotoDetailRepository {
	return &photoDetailRepository{}
}

func (r *photoDetailRepository) Create(tx Querier, photoDetail *entity.PhotoDetail) (*entity.PhotoDetail, error) {
	log.Println("Created accessed")
	query := `INSERT INTO photo_details 
			  (id, photo_id, file_name, file_key, size, type, checksum, width, height, url, your_moments_type, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := tx.Exec(query, photoDetail.Id, photoDetail.PhotoId, photoDetail.FileName, photoDetail.FileKey, photoDetail.Size, photoDetail.Type,
		photoDetail.Checksum, photoDetail.Width, photoDetail.Height, photoDetail.Url, photoDetail.YourMomentsType, photoDetail.CreatedAt, photoDetail.UpdatedAt)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert photoDetail: %w", err)
	}

	return photoDetail, nil
}

// func (r *photoDetailRepository) Update(ctx context.Context, db Querier, req *model.RequestUpdatePhoto) error {

// 	query := `UPDATE INTO user_simillars (id, user_id, size, url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING (id)`

// 	err := db.Exec(ctx, query)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
