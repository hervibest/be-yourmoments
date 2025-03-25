package usecase

import (
	"be-yourmoments/photo-svc/internal/adapter"
	"be-yourmoments/photo-svc/internal/entity"
	"be-yourmoments/photo-svc/internal/repository"
	"context"
	"mime/multipart"
	"time"

	"github.com/emicklei/go-restful/v3/log"
	"github.com/gofiber/fiber"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type FacecamUseCase interface {
	UploadFacecam(ctx context.Context, file *multipart.FileHeader) error
	// UpdateProcessedPhoto(ctx context.Context, req *model.RequestUpdateProcessedPhoto) (error, error)
}

type facecamUseCase struct {
	db              *sqlx.DB
	facecamRepo     repository.FacecamRepository
	userSimilarRepo repository.UserSimilarRepository
	aiAdapter       adapter.AiAdapter
	uploadAdapter   adapter.UploadAdapter
}

func NewFacecamUseCase(db *sqlx.DB, facecamRepo repository.FacecamRepository,
	userSimilarRepo repository.UserSimilarRepository, aiAdapter adapter.AiAdapter,
	uploadAdapter adapter.UploadAdapter) FacecamUseCase {
	return &facecamUseCase{
		db:              db,
		facecamRepo:     facecamRepo,
		userSimilarRepo: userSimilarRepo,
		aiAdapter:       aiAdapter,
		uploadAdapter:   uploadAdapter,
	}
}

func (u *facecamUseCase) UploadFacecam(ctx context.Context, file *multipart.FileHeader) error {
	uploadFile, err := file.Open()
	if err != nil {
		log.Print("parse file error" + err.Error())
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	upload, err := u.uploadAdapter.UploadFile(ctx, file, uploadFile, "facecam")
	if err != nil {
		return err
	}

	tx, err := u.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	newPhoto := &entity.Facecam{
		Id:        ulid.Make().String(),
		CreatorId: "test-create-facecam-case",
		Title:     upload.Filename,
		Size:      upload.Size,
		Url:       upload.URL,

		OriginalAt: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	newPhoto, err = u.facecamRepo.Create(tx, newPhoto)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// go func() {
	// 	defer uploadFile.Close()
	// 	_, filePath, err := u.compressAdapter.CompressImage(uploadFile, "facecam")
	// 	if err != nil {
	// 		log.Printf("Error compressing image: %v", err)
	// 		return
	// 	}

	// 	file, err := os.Open(filePath)
	// 	if err != nil {
	// 		log.Printf("Error opening file: %v", err)
	// 		return
	// 	}
	// 	defer file.Close()

	// 	fileInfo, err := file.Stat()
	// 	if err != nil {
	// 		log.Printf("Error stating file: %v", err)
	// 		return
	// 	}

	// 	mimeHeader := make(textproto.MIMEHeader)
	// 	mimeHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileInfo.Name()))

	// 	fileHeader := &multipart.FileHeader{
	// 		Filename: fileInfo.Name(),
	// 		Header:   mimeHeader,
	// 		Size:     fileInfo.Size(),
	// 	}

	// 	uploadPath := "photo"

	// 	response, err := u.uploadAdapter.UploadFile(ctx, fileHeader, file, uploadPath)
	// 	if err != nil {
	// 		log.Printf("Error uploading file: %v", err)
	// 		return
	// 	}

	// 	log.Printf("File berhasil diupload: %+v", response)
	// 	u.aiAdapter.ProcessPhoto(ctx, newPhoto.Id, response.URL)
	// }()

	return nil

}
