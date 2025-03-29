package usecase

import (
	"be-yourmoments/upload-svc/internal/adapter"
	"be-yourmoments/upload-svc/internal/entity"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/textproto"
	"os"
	"time"

	"github.com/gofiber/fiber"
	"github.com/oklog/ulid/v2"
)

type FacecamUseCase interface {
	UploadFacecam(ctx context.Context, file *multipart.FileHeader, userId string) error
	// UpdateProcessedPhoto(ctx context.Context, req *model.RequestUpdateProcessedPhoto) (error, error)
}

type facecamUseCase struct {
	aiAdapter       adapter.AiAdapter
	photoAdapter    adapter.PhotoAdapter
	storageAdapter  adapter.StorageAdapter
	compressAdapter adapter.CompressAdapter
}

func NewFacecamUseCase(aiAdapter adapter.AiAdapter, photoAdapter adapter.PhotoAdapter,
	storageAdapter adapter.StorageAdapter, compressAdapter adapter.CompressAdapter) FacecamUseCase {
	return &facecamUseCase{
		aiAdapter:       aiAdapter,
		photoAdapter:    photoAdapter,
		storageAdapter:  storageAdapter,
		compressAdapter: compressAdapter,
	}
}

func (u *facecamUseCase) UploadFacecam(ctx context.Context, file *multipart.FileHeader, userId string) error {
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

	checksum := fmt.Sprintf("%x", sha256.Sum256(data))

	go func() {
		_, filePath, err := u.compressAdapter.CompressImage(file, wrappedReader, "facecam")
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

		uploadPath := "facecam/compressed"
		compressedPhoto, err := u.storageAdapter.UploadFile(ctx, fileHeader, fileComp, uploadPath)
		if err != nil {
			log.Printf("Error uploading file: %v", err)
			return
		}

		newFacecam := &entity.Facecam{
			Id:         ulid.Make().String(),
			UserId:     userId,
			FileName:   compressedPhoto.Filename,
			FileKey:    compressedPhoto.FileKey,
			Title:      compressedPhoto.Filename,
			Size:       compressedPhoto.Size,
			Checksum:   checksum,
			Url:        compressedPhoto.URL,
			OriginalAt: time.Now(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := u.photoAdapter.CreateFacecam(ctx, newFacecam); err != nil {
			log.Printf("Error creating facecam: %v", err)
			return
		}

		if err := os.Remove(filePath); err != nil {
			log.Printf("Gagal menghapus file: %v", err)
		} else {
			log.Printf("File sementara berhasil dihapus: %s", filePath)
		}

		u.aiAdapter.ProcessFacecam(ctx, userId, compressedPhoto.URL)
	}()

	return nil

}
