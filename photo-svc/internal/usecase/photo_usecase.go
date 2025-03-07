package usecase

import (
	"be-yourmoments/photo-svc/internal/adapter"
	"be-yourmoments/photo-svc/internal/entity"
	"be-yourmoments/photo-svc/internal/model"
	"be-yourmoments/photo-svc/internal/repository"
	"context"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PhotoUsecase interface {
	UploadPhoto(ctx context.Context, file *multipart.FileHeader) (error, error)
	UpdateProcessedPhoto(ctx context.Context, req *model.RequestUpdateProcessedPhoto) (error, error)
}

type photoUsecase struct {
	db               *pgxpool.Pool
	photoRepo        repository.PhotoRepository
	photoProfileRepo repository.PhotoProfileRepository
	userSimilarRepo  repository.UserSimilarRepository
	aiAdapter        adapter.AiAdapter
	uploadAdapter    adapter.Minio
}

func NewPhotoUsecase(db *pgxpool.Pool, photoRepo repository.PhotoRepository,
	photoProfileRepo repository.PhotoProfileRepository,
	userSimilarRepo repository.UserSimilarRepository,
	aiAdapter adapter.AiAdapter, uploadAdpter adapter.Minio) PhotoUsecase {
	return &photoUsecase{
		db:               db,
		photoRepo:        photoRepo,
		photoProfileRepo: photoProfileRepo,
		userSimilarRepo:  userSimilarRepo,
		aiAdapter:        aiAdapter,
		uploadAdapter:    uploadAdpter,
	}
}

func (u *photoUsecase) UploadPhoto(ctx context.Context, file *multipart.FileHeader) (error, error) {

	upload, err := u.uploadAdapter.UploadFile(ctx, file, "photo")
	if err != nil {
		return nil, err
	}

	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err, err
	}

	newPhoto := &entity.Photo{
		Id:        uuid.NewString(),
		UserId:    "",
		Size:      upload.Size,
		RawUrl:    upload.URL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = u.photoRepo.Create(ctx, tx, newPhoto)
	if err != nil {
		return err, err
	}

	if err := tx.Commit(ctx); err != nil {
		return err, err
	}

	// Process photo service will be executed asyncronously by goroutine
	go func() {
		u.aiAdapter.ProcessPhoto(ctx, newPhoto.Id, upload.URL)
	}()

	return nil, nil

}

func (u *photoUsecase) UpdateProcessedPhoto(ctx context.Context, req *model.RequestUpdateProcessedPhoto) (error, error) {

	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err, err
	}

	updatePhoto := &entity.Photo{
		Id:                     req.Id,
		PreviewUrl:             req.PreviewUrl,
		PreviewWithBoundingUrl: req.PreviewWithBoundingUrl,
		UpdatedAt:              time.Now(),
	}

	err = u.photoRepo.UpdateProcessedUrl(ctx, tx, updatePhoto)
	if err != nil {
		return err, err
	}

	err = u.userSimilarRepo.UpdateUsersForPhoto(ctx, tx, req.Id, req.UserId)
	if err != nil {
		return err, err
	}

	if err := tx.Commit(ctx); err != nil {
		return err, err
	}

	return nil, nil

}

func (u *photoUsecase) ClaimPhoto(ctx context.Context, req *model.RequestClaimPhoto) (error, error) {

	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err, err
	}

	updatePhoto := &entity.Photo{
		Id:            req.Id,
		OwnedByUserId: req.UserId,
		Status:        "Claimed",
		UpdatedAt:     time.Now(),
	}

	err = u.photoRepo.UpdateClaimedPhoto(ctx, tx, updatePhoto)
	if err != nil {
		return err, err
	}

	// err = u.userSimilarRepo.UpdateUsersForPhoto(ctx, tx, req.Id, req.UserId)
	// if err != nil {
	// 	return err, err
	// }

	if err := tx.Commit(ctx); err != nil {
		return err, err
	}

	// Process photo service will be executed asyncronously by goroutine

	return nil, nil

}

func (u *photoUsecase) CancelClaimPhoto(ctx context.Context, req *model.RequestClaimPhoto) (error, error) {

	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err, err
	}

	updatePhoto := &entity.Photo{
		Id:            req.Id,
		OwnedByUserId: "",
		Status:        "Unclaimed",
		UpdatedAt:     time.Now(),
	}

	err = u.photoRepo.UpdateClaimedPhoto(ctx, tx, updatePhoto)
	if err != nil {
		return err, err
	}

	// err = u.userSimilarRepo.UpdateUsersForPhoto(ctx, tx, req.Id, req.UserId)
	// if err != nil {
	// 	return err, err
	// }

	if err := tx.Commit(ctx); err != nil {
		return err, err
	}

	// Process photo service will be executed asyncronously by goroutine

	return nil, nil

}

func (u *photoUsecase) UpdateBuyyedPhoto(ctx context.Context, req *model.RequestClaimPhoto) (error, error) {

	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err, err
	}

	updatePhoto := &entity.Photo{
		Id:        req.Id,
		Status:    "Owned",
		UpdatedAt: time.Now(),
	}

	err = u.photoRepo.UpdatePhotoStatus(ctx, tx, updatePhoto)
	if err != nil {
		return err, err
	}

	err = u.userSimilarRepo.DeleteSimilarUsers(ctx, tx, req.Id)
	if err != nil {
		return err, err
	}

	if err := tx.Commit(ctx); err != nil {
		return err, err
	}

	// Process photo service will be executed asyncronously by goroutine

	return nil, nil

}
