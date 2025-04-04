package repository

import (
	"be-yourmoments/user-svc/internal/entity"
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type userImagePreparedStmt struct {
	FindByUserProfIdAndType *sqlx.Stmt
	FindByUserProfId        *sqlx.Stmt
}

func newUserImageStmt(db *sqlx.DB) (*userImagePreparedStmt, error) {
	FindByUserProfIdAndType, err := db.Preparex("SELECT * FROM user_images WHERE user_profile_id = $1 AND image_type = $2")
	if err != nil {
		return nil, err
	}

	FindByUserProfId, err := db.Preparex("SELECT * FROM user_images WHERE user_profile_id = $1")
	if err != nil {
		return nil, err
	}

	return &userImagePreparedStmt{
		FindByUserProfIdAndType: FindByUserProfIdAndType,
		FindByUserProfId:        FindByUserProfId,
	}, nil
}

type UserImageRepository interface {
	FindByUserProfIdAndType(ctx context.Context, userProfId, imageType string) (*entity.UserImage, error)
	FindByUserProfId(ctx context.Context, userProfId string) (*[]*entity.UserImage, error)
	Create(ctx context.Context, tx Querier, userImage *entity.UserImage) (*entity.UserImage, error)
	Update(ctx context.Context, tx Querier, userImage *entity.UserImage) (*entity.UserImage, error)
}

type userImageRepository struct {
	userImagePreparedStmt *userImagePreparedStmt
}

func NewUserImageRepository(db *sqlx.DB) (UserImageRepository, error) {
	userImagePreparedStmt, err := newUserImageStmt(db)
	if err != nil {
		log.Print("error initialize rest user image statement : ", err)
		return nil, err
	}

	return &userImageRepository{
		userImagePreparedStmt: userImagePreparedStmt,
	}, nil
}

func (r *userImageRepository) FindByUserProfIdAndType(ctx context.Context, userProfId, imageType string) (*entity.UserImage, error) {
	userImage := new(entity.UserImage)

	row := r.userImagePreparedStmt.FindByUserProfIdAndType.QueryRowxContext(ctx, userProfId, imageType)
	if err := row.StructScan(userImage); err != nil {

		return nil, err
	}

	return userImage, nil
}

func (r *userImageRepository) FindByUserProfId(ctx context.Context, userProfId string) (*[]*entity.UserImage, error) {
	userImages := make([]*entity.UserImage, 0)

	rows, err := r.userImagePreparedStmt.FindByUserProfId.QueryxContext(ctx, userProfId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		userImage := new(entity.UserImage)
		if err := rows.StructScan(userImage); err != nil {
			return nil, err
		}
		userImages = append(userImages, userImage)
	}

	return &userImages, nil
}

func (r *userImageRepository) Create(ctx context.Context, tx Querier, userImage *entity.UserImage) (*entity.UserImage, error) {
	query := `INSERT INTO user_images 
	(id, user_profile_id, file_name, file_key, image_type, size, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := tx.ExecContext(ctx, query, userImage.Id, userImage.UserProfileId, userImage.FileName, userImage.FileKey, userImage.ImageType, userImage.Size, userImage.CreatedAt, userImage.UpdatedAt)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return userImage, nil
}

func (r *userImageRepository) Update(ctx context.Context, tx Querier, userImage *entity.UserImage) (*entity.UserImage, error) {
	query := `UPDATE  user_images SET file_name = $1, file_key = $2, image_type = $3, size = $4, updated_at = $5 WHERE 
	user_profile_id = $6 AND image_type = $7`

	_, err := tx.ExecContext(ctx, query, userImage.FileName, userImage.FileKey, userImage.ImageType, userImage.Size, userImage.UpdatedAt, userImage.UserProfileId, userImage.ImageType)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return userImage, nil
}
