package usecase

import (
	"be-yourmoments/upload-svc/internal/adapter"
	"be-yourmoments/upload-svc/internal/entity"
	"be-yourmoments/upload-svc/internal/enum"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/textproto"
	"os"
	"time"

	"github.com/gofiber/fiber"
	"github.com/oklog/ulid/v2"
)

type PhotoUsecase interface {
	UploadPhoto(ctx context.Context, file *multipart.FileHeader) error
	// UpdateProcessedPhoto(ctx context.Context, req *model.RequestUpdateProcessedPhoto) (error, error)
}

type photoUsecase struct {
	aiAdapter       adapter.AiAdapter
	photoAdapter    adapter.PhotoAdapter
	uploadAdapter   adapter.UploadAdapter
	compressAdapter adapter.CompressAdapter
}

func NewPhotoUsecase(aiAdapter adapter.AiAdapter, photoAdapter adapter.PhotoAdapter,
	uploadAdapter adapter.UploadAdapter,
	compressAdapter adapter.CompressAdapter) PhotoUsecase {
	return &photoUsecase{
		aiAdapter:       aiAdapter,
		photoAdapter:    photoAdapter,
		uploadAdapter:   uploadAdapter,
		compressAdapter: compressAdapter,
	}
}

func (u *photoUsecase) UploadPhoto(ctx context.Context, file *multipart.FileHeader) error {
	uploadFile, err := file.Open()
	if err != nil {
		log.Print("parse file error" + err.Error())
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	upload, err := u.uploadAdapter.UploadFile(ctx, file, uploadFile, "photo")
	if err != nil {
		return err
	}

	newPhoto := &entity.Photo{
		Id:            ulid.Make().String(),
		CreatorId:     "test-create-photo-case-2",
		Title:         upload.Filename,
		CollectionUrl: upload.URL,
		Price:         133,
		PriceStr:      "12312",

		OriginalAt: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// imgConfig, _, err := image.DecodeConfig(uploadFile)
	// if err != nil {
	// 	log.Print("image decode error:", err)
	// 	return fiber.NewError(fiber.StatusBadRequest, "Not a valid image")
	// }

	// WHAT TO DO JPG TYPE AND CHECKSUM
	newPhotoDetail := &entity.PhotoDetail{
		Id:       ulid.Make().String(),
		PhotoId:  newPhoto.Id,
		Size:     upload.Size,
		Type:     "JPG",
		Checksum: "",
		// Width:           int8(imgConfig.Width),
		// Height:          int8(imgConfig.Height),
		Width:           33,
		Height:          33,
		Url:             upload.URL,
		YourMomentsType: enum.YourMomentTypeCollection,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := u.photoAdapter.CreatePhoto(ctx, newPhoto, newPhotoDetail); err != nil {
		log.Printf("Error creating photo: %v", err)
		return err
	}

	go func() {
		defer uploadFile.Close()
		_, filePath, err := u.compressAdapter.CompressImage(uploadFile, "photo")
		if err != nil {
			log.Printf("Error compressing images: %v", err)
			return
		}

		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Error opening file: %v", err)
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			log.Printf("Error stating file: %v", err)
			return
		}

		mimeHeader := make(textproto.MIMEHeader)
		mimeHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileInfo.Name()))

		fileHeader := &multipart.FileHeader{
			Filename: fileInfo.Name(),
			Header:   mimeHeader,
			Size:     fileInfo.Size(),
		}

		uploadPath := "photo"

		response, err := u.uploadAdapter.UploadFile(ctx, fileHeader, file, uploadPath)
		if err != nil {
			log.Printf("Error uploading file: %v", err)
			return
		}

		log.Printf("File berhasil diupload: %+v", response)

		if err := os.Remove(filePath); err != nil {
			log.Printf("Gagal menghapus file: %v", err)
		} else {
			log.Printf("File sementara berhasil dihapus: %s", filePath)
		}

		u.aiAdapter.ProcessPhoto(ctx, newPhoto.Id, response.URL)
	}()

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
