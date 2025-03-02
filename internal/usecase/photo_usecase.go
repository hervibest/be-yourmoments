package usecase

import (
	"be-yourmoments/internal/adapter"
	"be-yourmoments/internal/entity"
	"be-yourmoments/internal/repository"
	"context"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PhotoUsecase interface {
	UploadPhoto(ctx context.Context, file *multipart.FileHeader) (error, error)
	UpdatePhoto(ctx context.Context, photoId string) (error, error)
}

type photoUsecase struct {
	db            *pgxpool.Pool
	photoRepo     repository.PhotoRepository
	aiAdapter     adapter.AiAdapter
	uploadAdapter adapter.Minio
}

func NewPhotoUsecase(db *pgxpool.Pool, photoRepo repository.PhotoRepository,
	aiAdapter adapter.AiAdapter, uploadAdpter adapter.Minio) PhotoUsecase {
	return &photoUsecase{
		db:            db,
		aiAdapter:     aiAdapter,
		photoRepo:     photoRepo,
		uploadAdapter: uploadAdpter,
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
	go u.aiAdapter.ProcessPhoto(ctx, newPhoto.Id, upload.URL)

	return nil, nil

}

func (u *photoUsecase) UpdatePhoto(ctx context.Context, photoId string) (error, error) {

	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err, err
	}

	err = u.photoRepo.Update(ctx, tx, photoId)
	if err != nil {
		return err, err
	}

	if err := tx.Commit(ctx); err != nil {
		return err, err
	}

	// Process photo service will be executed asyncronously by goroutine

	return nil, nil

}
