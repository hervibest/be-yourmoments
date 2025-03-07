package usecase

import (
	"be-yourmoments/photo-svc/internal/entity"
	"be-yourmoments/photo-svc/internal/repository"
	"context"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type PhotoProfileUsecase interface {
}

type photoProfileUsecase struct {
	photoProfileRepository repository.PhotoProfileRepository
}

func NewPhotoProfileUsecase() PhotoProfileUsecase {
	return &photoProfileUsecase{}
}

func (u *photoUsecase) UploadPhotoProfile(ctx context.Context, file *multipart.FileHeader) (error, error) {

	upload, err := u.uploadAdapter.UploadFile(ctx, file, "photo-profile")
	if err != nil {
		return nil, err
	}

	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err, err
	}

	newPhotoProfile := &entity.PhotoProfile{
		Id:        uuid.NewString(),
		UserId:    "",
		Size:      upload.Size,
		RawUrl:    upload.URL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = u.photoProfileRepo.Create(ctx, tx, newPhotoProfile)
	if err != nil {
		return err, err
	}

	if err := tx.Commit(ctx); err != nil {
		return err, err
	}

	go func() {
		u.aiAdapter.ProcessPhoto(ctx, newPhotoProfile.Id, upload.URL)
	}()

	return nil, nil

}
