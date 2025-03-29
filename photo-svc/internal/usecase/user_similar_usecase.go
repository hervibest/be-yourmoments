package usecase

import (
	"be-yourmoments/photo-svc/internal/entity"
	"be-yourmoments/photo-svc/internal/enum"
	"be-yourmoments/photo-svc/internal/pb"
	"be-yourmoments/photo-svc/internal/repository"
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type UserSimilarUsecase interface {
	CreateUserSimilar(ctx context.Context, request *pb.CreateUserSimilarPhotoRequest) error
	CreateUserFacecam(ctx context.Context, request *pb.CreateUserSimilarFacecamRequest) error
}

type userSimilarUsecase struct {
	db              *sqlx.DB
	photoRepo       repository.PhotoRepository
	photoDetailRepo repository.PhotoDetailRepository
	facecamRepo     repository.FacecamRepository
	userSimilarRepo repository.UserSimilarRepository
}

func NewUserSimilarUsecase(db *sqlx.DB, photoRepo repository.PhotoRepository,
	photoDetailRepo repository.PhotoDetailRepository, facecamRepo repository.FacecamRepository,
	userSimilarRepo repository.UserSimilarRepository) UserSimilarUsecase {
	return &userSimilarUsecase{
		db:              db,
		photoRepo:       photoRepo,
		photoDetailRepo: photoDetailRepo,
		facecamRepo:     facecamRepo,
		userSimilarRepo: userSimilarRepo,
	}
}

func (u *userSimilarUsecase) CreateUserSimilar(ctx context.Context, request *pb.CreateUserSimilarPhotoRequest) error {
	log.Println("Create user similar ")

	tx, err := u.db.Beginx()
	if err != nil {
		log.Println(err)

		return err
	}

	defer func() {
		if err != nil {
			log.Println(err)

			tx.Rollback()
		}
	}()

	photo := &entity.Photo{
		Id:             request.GetPhotoDetail().PhotoId,
		IsThisYouURL:   "",
		YourMomentsUrl: request.GetPhotoDetail().Url,
		UpdatedAt:      time.Now(),
	}

	err = u.photoRepo.UpdateProcessedUrl(tx, photo)
	if err != nil {
		log.Println(err)

		return err
	}

	newPhotoDetail := &entity.PhotoDetail{
		Id:              ulid.Make().String(),
		PhotoId:         request.GetPhotoDetail().GetPhotoId(),
		FileName:        request.GetPhotoDetail().GetFileName(),
		FileKey:         request.GetPhotoDetail().GetFileKey(),
		Size:            request.GetPhotoDetail().GetSize(),
		Type:            "JPG",
		Checksum:        "1212",
		Height:          121,
		Width:           1212,
		Url:             request.GetPhotoDetail().GetUrl(),
		YourMomentsType: enum.YourMomentsType(request.GetPhotoDetail().GetYourMomentsType()),
		CreatedAt:       request.GetPhotoDetail().GetCreatedAt().AsTime(),
		UpdatedAt:       request.GetPhotoDetail().GetUpdatedAt().AsTime(),
	}

	_, err = u.photoDetailRepo.Create(tx, newPhotoDetail)
	if err != nil {
		log.Println(err)

		return err
	}

	userSimilarPhotos := make([]*entity.UserSimilarPhoto, 0, len(request.GetUserSimilarPhoto()))
	for _, userSimilarPhotoRequest := range request.GetUserSimilarPhoto() {
		log.Println("ALL DELETED BECAUSE OF ZERO")
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

	err = u.userSimilarRepo.InsertOrUpdate(tx, request.GetPhotoDetail().PhotoId, &userSimilarPhotos)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func (u *userSimilarUsecase) CreateUserFacecam(ctx context.Context, request *pb.CreateUserSimilarFacecamRequest) error {
	log.Println("Create user similar ")

	tx, err := u.db.Beginx()
	if err != nil {
		log.Println(err)

		return err
	}

	defer func() {
		if err != nil {
			log.Println(err)

			tx.Rollback()
		}
	}()

	facecam := &entity.Facecam{
		UserId:      request.GetFacecam().GetUserId(),
		IsProcessed: request.GetFacecam().GetIsProcessed(),
		UpdatedAt:   request.GetFacecam().GetUpdatedAt().AsTime(),
	}

	err = u.facecamRepo.UpdatedProcessedFacecam(tx, facecam)
	if err != nil {
		log.Println(err)

		return err
	}

	userSimilarPhotos := make([]*entity.UserSimilarPhoto, 0, len(request.GetUserSimilarPhoto()))
	for _, userSimilarPhotoRequest := range request.GetUserSimilarPhoto() {
		log.Println("UPDATE UserSimilarPhoto from facecams")
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

	err = u.userSimilarRepo.InserOrUpdateByUserId(tx, request.GetFacecam().UserId, &userSimilarPhotos)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}

	return nil

}
