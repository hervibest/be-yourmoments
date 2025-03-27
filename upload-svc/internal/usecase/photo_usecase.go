package usecase

import (
	"be-yourmoments/upload-svc/internal/adapter"
	"be-yourmoments/upload-svc/internal/entity"
	"be-yourmoments/upload-svc/internal/enum"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"image"
	"io"
	"log"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"
	"time"

	_ "image/jpeg"
	_ "image/png"

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
	storageAdapter  adapter.StorageAdapter
	compressAdapter adapter.CompressAdapter
}

func NewPhotoUsecase(aiAdapter adapter.AiAdapter, photoAdapter adapter.PhotoAdapter,
	storageAdapter adapter.StorageAdapter,
	compressAdapter adapter.CompressAdapter) PhotoUsecase {
	return &photoUsecase{
		aiAdapter:       aiAdapter,
		photoAdapter:    photoAdapter,
		storageAdapter:  storageAdapter,
		compressAdapter: compressAdapter,
	}
}

type nopReadSeekCloser struct {
	*bytes.Reader
}

func (n nopReadSeekCloser) Close() error {
	return nil
}

func (u *photoUsecase) UploadPhoto(ctx context.Context, file *multipart.FileHeader) error {
	uploadFile, err := file.Open()
	if err != nil {
		log.Print("parse file error: " + err.Error())
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	data, err := io.ReadAll(uploadFile)
	if err != nil {
		log.Print("failed to read file: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}
	uploadFile.Close()

	readerForUpload := bytes.NewReader(data)
	wrappedReader := nopReadSeekCloser{readerForUpload}

	upload, err := u.storageAdapter.UploadFile(ctx, file, wrappedReader, "photo")
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
		OriginalAt:    time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	readerForDecode := bytes.NewReader(data)
	imgConfig, format, err := image.DecodeConfig(readerForDecode)
	if err != nil {
		log.Print("image decode error:", err)
		return fiber.NewError(fiber.StatusBadRequest, "Not a valid images")
	}

	log.Println("Decoded image format:", format)
	log.Println("Decoded image format:", imgConfig.Width, imgConfig.Height)

	var imageType string
	if format == "jpeg" {
		imageType = "JPG"
	} else {
		imageType = strings.ToUpper(format)
	}

	checksum := fmt.Sprintf("%x", sha256.Sum256(data))

	newPhotoDetail := &entity.PhotoDetail{
		Id:              ulid.Make().String(),
		PhotoId:         newPhoto.Id,
		FileName:        upload.Filename,
		FileKey:         upload.FileKey,
		Size:            upload.Size,
		Type:            imageType,
		Checksum:        checksum,
		Width:           imgConfig.Width,  // disesuaikan tipe data jika perlu
		Height:          imgConfig.Height, // disesuaikan tipe data jika perlu
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
		_, filePath, err := u.compressAdapter.CompressImage(wrappedReader, "photo")
		if err != nil {
			log.Printf("Error compressing images: %v", err)
			return
		}

		fileComp, err := os.Open(filePath)
		if err != nil {
			log.Printf("Error opening file: %v", err)
			return
		}
		defer fileComp.Close()

		fileInfo, err := fileComp.Stat()
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
		compressedPhoto, err := u.storageAdapter.UploadFile(ctx, fileHeader, fileComp, uploadPath)
		if err != nil {
			log.Printf("Error uploading file: %v", err)
			return
		}

		compressedPhotoDetail := &entity.PhotoDetail{
			Id:              ulid.Make().String(),
			PhotoId:         newPhoto.Id,
			FileName:        compressedPhoto.Filename,
			FileKey:         compressedPhoto.FileKey,
			Size:            compressedPhoto.Size,
			Url:             compressedPhoto.URL,
			YourMomentsType: enum.YourMomentTypeCompressed,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		if err := u.photoAdapter.UpdatePhotoDetail(ctx, compressedPhotoDetail); err != nil {
			log.Printf("Error creating photo: %v", err)
			return
		}

		if err := os.Remove(filePath); err != nil {
			log.Printf("Gagal menghapus file: %v", err)
		} else {
			log.Printf("File sementara berhasil dihapus: %s", filePath)
		}

		// u.aiAdapter.ProcessPhoto(ctx, newPhoto.Id, response.URL)
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
