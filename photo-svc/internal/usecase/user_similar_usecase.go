package usecase

import (
	"be-yourmoments/photo-svc/internal/entity"
	"be-yourmoments/photo-svc/internal/enum"
	"be-yourmoments/photo-svc/internal/pb"
	"be-yourmoments/photo-svc/internal/repository"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type UserSimilarUsecase interface {
	CreateUserSimilarPhoto(ctx context.Context, request *pb.CreateUserSimilarPhotoRequest) error
}

type userSimilarUsecase struct {
	db              *sqlx.DB
	photoRepo       repository.PhotoRepository
	photoDetailRepo repository.PhotoDetailRepository
	userSimilarRepo repository.UserSimilarRepository
}

func NewUserSimilarUsecase(db *sqlx.DB, photoRepo repository.PhotoRepository,
	photoDetailRepo repository.PhotoDetailRepository,
	userSimilarRepo repository.UserSimilarRepository) UserSimilarUsecase {
	return &userSimilarUsecase{
		db:              db,
		userSimilarRepo: userSimilarRepo,
	}
}

func (u *userSimilarUsecase) CreateUserSimilarPhoto(ctx context.Context, request *pb.CreateUserSimilarPhotoRequest) error {
	tx, err := u.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	photo := &entity.Photo{
		Id:             request.GetUserSimilarPhoto()[0].PhotoId,
		IsThisYouURL:   request.GetIsThisYouUrl(),
		YourMomentsUrl: request.GetYourMomentsUrl(),
		UpdatedAt:      time.Now(),
	}

	err = u.photoRepo.UpdateProcessedUrl(tx, photo)
	if err != nil {
		return err
	}

	userSimilarPhotos := make([]*entity.UserSimilarPhoto, 0, len(request.GetUserSimilarPhoto()))
	for _, userSimilarPhotoRequest := range request.GetUserSimilarPhoto() {
		userSimilarPhoto := &entity.UserSimilarPhoto{
			Id:         ulid.Make().String(),
			PhotoId:    userSimilarPhotoRequest.GetPhotoId(),
			UserId:     userSimilarPhotoRequest.GetUserId(),
			Similarity: enum.SimilarityLevelEnum(userSimilarPhotoRequest.GetSimilarity().String()),
			CreatedAt:  userSimilarPhotoRequest.GetCreatedAt().AsTime(),
			UpdatedAt:  userSimilarPhotoRequest.GetUpdatedAt().AsTime(),
		}
		userSimilarPhotos = append(userSimilarPhotos, userSimilarPhoto)
	}

	err = u.userSimilarRepo.InsertOrUpdate(tx, request.GetUserSimilarPhoto()[0].PhotoId, &userSimilarPhotos)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}
