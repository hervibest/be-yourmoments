package usecase

import (
	"be-yourmoments/photo-svc/internal/adapter"
	"be-yourmoments/photo-svc/internal/entity"
	"be-yourmoments/photo-svc/internal/enum"
	"be-yourmoments/photo-svc/internal/pb"
	"be-yourmoments/photo-svc/internal/repository"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type PhotoUsecase interface {
	CreatePhoto(ctx context.Context, request *pb.CreatePhotoRequest) error
	UpdatePhotoDetail(ctx context.Context, request *pb.UpdatePhotoDetailRequest) error
	// UpdateProcessedPhoto(ctx context.Context, req *model.RequestUpdateProcessedPhoto) (error, error)
}

type photoUsecase struct {
	db              *sqlx.DB
	photoRepo       repository.PhotoRepository
	photoDetailRepo repository.PhotoDetailRepository
	userSimilarRepo repository.UserSimilarRepository
	aiAdapter       adapter.AiAdapter
	uploadAdapter   adapter.UploadAdapter
}

func NewPhotoUsecase(db *sqlx.DB, photoRepo repository.PhotoRepository,
	photoDetailRepo repository.PhotoDetailRepository,
	userSimilarRepo repository.UserSimilarRepository,
	aiAdapter adapter.AiAdapter, uploadAdapter adapter.UploadAdapter) PhotoUsecase {
	return &photoUsecase{
		db:              db,
		photoRepo:       photoRepo,
		photoDetailRepo: photoDetailRepo,
		userSimilarRepo: userSimilarRepo,
		aiAdapter:       aiAdapter,
		uploadAdapter:   uploadAdapter,
	}
}

func (u *photoUsecase) CreatePhoto(ctx context.Context, request *pb.CreatePhotoRequest) error {
	tx, err := u.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	newPhoto := &entity.Photo{
		Id:            request.GetPhoto().GetId(),
		CreatorId:     "test-create-photo-case-2",
		Title:         request.GetPhoto().GetTitle(),
		CollectionUrl: request.GetPhoto().GetCollectionUrl(),
		Price:         request.GetPhoto().GetPrice(),
		PriceStr:      request.GetPhoto().GetPriceStr(),

		OriginalAt: request.GetPhoto().GetOriginalAt().AsTime(),
		CreatedAt:  request.GetPhoto().GetCreatedAt().AsTime(),
		UpdatedAt:  request.GetPhoto().GetUpdatedAt().AsTime(),
	}

	newPhoto, err = u.photoRepo.Create(tx, newPhoto)
	if err != nil {
		return err
	}

	newPhotoDetail := &entity.PhotoDetail{
		Id:              ulid.Make().String(),
		PhotoId:         newPhoto.Id,
		FileName:        request.GetPhoto().GetDetail().GetFileName(),
		FileKey:         request.GetPhoto().GetDetail().GetFileKey(),
		Size:            request.GetPhoto().GetDetail().GetSize(),
		Type:            request.GetPhoto().GetDetail().GetType(),
		Checksum:        request.GetPhoto().GetDetail().GetChecksum(),
		Width:           request.GetPhoto().GetDetail().GetWidth(),
		Height:          request.GetPhoto().GetDetail().GetHeight(),
		Url:             request.GetPhoto().GetDetail().GetUrl(),
		YourMomentsType: enum.YourMomentsType(request.GetPhoto().GetDetail().GetYourMomentsType()),
		CreatedAt:       request.GetPhoto().GetDetail().GetCreatedAt().AsTime(),
		UpdatedAt:       request.GetPhoto().GetDetail().GetUpdatedAt().AsTime(),
	}

	newPhotoDetail, err = u.photoDetailRepo.Create(tx, newPhotoDetail)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

func (u *photoUsecase) UpdatePhotoDetail(ctx context.Context, request *pb.UpdatePhotoDetailRequest) error {
	tx, err := u.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

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

	newPhotoDetail, err = u.photoDetailRepo.Create(tx, newPhotoDetail)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

// func (u *photoUsecase) UpdateProcessedPhoto(ctx context.Context, req *model.RequestUpdateProcessedPhoto) (error, error) {

// 	tx, err := u.db.Begin()
// 	if err != nil {
// 		return err, err
// 	}

// 	updatePhoto := &entity.Photo{
// 		Id:                     req.Id,
// 		PreviewUrl:             req.PreviewUrl,
// 		PreviewWithBoundingUrl: req.PreviewWithBoundingUrl,
// 		UpdatedAt:              time.Now(),
// 	}

// 	err = u.photoRepo.UpdateProcessedUrl(ctx, tx, updatePhoto)
// 	if err != nil {
// 		return err, err
// 	}

// 	err = u.userSimilarRepo.UpdateUsersForPhoto(ctx, tx, req.Id, req.UserId)
// 	if err != nil {
// 		return err, err
// 	}

// 	if err := tx.Commit(ctx); err != nil {
// 		return err, err
// 	}

// 	return nil, nil

// }

// func (u *photoUsecase) ClaimPhoto(ctx context.Context, req *model.RequestClaimPhoto) (error, error) {

// 	tx, err := u.db.Begin(ctx)
// 	if err != nil {
// 		return err, err
// 	}

// 	updatePhoto := &entity.Photo{
// 		Id:            req.Id,
// 		OwnedByUserId: req.UserId,
// 		Status:        "Claimed",
// 		UpdatedAt:     time.Now(),
// 	}

// 	err = u.photoRepo.UpdateClaimedPhoto(ctx, tx, updatePhoto)
// 	if err != nil {
// 		return err, err
// 	}

// 	// err = u.userSimilarRepo.UpdateUsersForPhoto(ctx, tx, req.Id, req.UserId)
// 	// if err != nil {
// 	// 	return err, err
// 	// }

// 	if err := tx.Commit(ctx); err != nil {
// 		return err, err
// 	}

// 	// Process photo service will be executed asyncronously by goroutine

// 	return nil, nil

// }

// func (u *photoUsecase) CancelClaimPhoto(ctx context.Context, req *model.RequestClaimPhoto) (error, error) {

// 	tx, err := u.db.Begin(ctx)
// 	if err != nil {
// 		return err, err
// 	}

// 	updatePhoto := &entity.Photo{
// 		Id:            req.Id,
// 		OwnedByUserId: "",
// 		Status:        "Unclaimed",
// 		UpdatedAt:     time.Now(),
// 	}

// 	err = u.photoRepo.UpdateClaimedPhoto(ctx, tx, updatePhoto)
// 	if err != nil {
// 		return err, err
// 	}

// 	// err = u.userSimilarRepo.UpdateUsersForPhoto(ctx, tx, req.Id, req.UserId)
// 	// if err != nil {
// 	// 	return err, err
// 	// }

// 	if err := tx.Commit(ctx); err != nil {
// 		return err, err
// 	}

// 	// Process photo service will be executed asyncronously by goroutine

// 	return nil, nil

// }

// func (u *photoUsecase) UpdateBuyyedPhoto(ctx context.Context, req *model.RequestClaimPhoto) (error, error) {

// 	tx, err := u.db.Begin(ctx)
// 	if err != nil {
// 		return err, err
// 	}

// 	updatePhoto := &entity.Photo{
// 		Id:        req.Id,
// 		Status:    "Owned",
// 		UpdatedAt: time.Now(),
// 	}

// 	err = u.photoRepo.UpdatePhotoStatus(ctx, tx, updatePhoto)
// 	if err != nil {
// 		return err, err
// 	}

// 	err = u.userSimilarRepo.DeleteSimilarUsers(ctx, tx, req.Id)
// 	if err != nil {
// 		return err, err
// 	}

// 	if err := tx.Commit(ctx); err != nil {
// 		return err, err
// 	}

// 	// Process photo service will be executed asyncronously by goroutine

// 	return nil, nil

// }
