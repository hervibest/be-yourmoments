package adapter

import (
	"be-yourmoments/internal/model"
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Minio interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, path string) (*model.MinioFileResponse, error)
	DeleteFile(ctx context.Context, fileName string) (bool, error)
}

type minioImpl struct {
	MinioClient     *minio.Client
	minioBucketName string
	enpoint         string
	log             *logrus.Logger
}

func NewMinio(viper *viper.Viper, log *logrus.Logger) Minio {
	ctx := context.Background()
	var minioClient *minio.Client
	var err error

	minioHost := viper.GetString("minio.minio_host")
	minioPort := viper.GetString("minio.minio_port")
	minioRootUser := viper.GetString("minio.minio_root_user")
	minioRootPassword := viper.GetString("minio.minio_root_password")
	minioTicketsBucket := viper.GetString("minio.minio_tickets_bucket")
	minioLocation := viper.GetString("minio.minio_location")
	endpoint := minioHost + ":" + minioPort

	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioRootUser, minioRootPassword, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	err = minioClient.MakeBucket(ctx, minioTicketsBucket, minio.MakeBucketOptions{Region: minioLocation})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, minioTicketsBucket)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", minioTicketsBucket)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", minioTicketsBucket)
	}

	log.Printf("Successfully connected %s\n", minioTicketsBucket)

	return &minioImpl{
		MinioClient:     minioClient,
		minioBucketName: minioTicketsBucket,
		log:             log,
	}
}

func (m *minioImpl) getBucketName() string {
	BucketName := m.minioBucketName
	return BucketName
}

func (m *minioImpl) getEndpoint() string {
	Endpoint := m.enpoint
	return Endpoint
}

func (m *minioImpl) UploadFile(ctx context.Context, file *multipart.FileHeader, path string) (*model.MinioFileResponse, error) {
	uploadFile, err := file.Open()
	if err != nil {
		m.log.Warnf("parse file error" + err.Error())
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	defer uploadFile.Close()

	fileKey := path + string(RandomNumber(31)) + "_" + file.Filename
	contentType := file.Header.Get("Content-Type")

	s3PutObjectOutput, err := m.MinioClient.PutObject(ctx, m.getBucketName(), fileKey, uploadFile, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		m.log.Warnf("failed to upload file to S3" + err.Error())

		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	fileResponse := new(model.MinioFileResponse)
	fileResponse.ChecksumCRC32 = s3PutObjectOutput.ChecksumCRC32
	fileResponse.ChecksumCRC32C = s3PutObjectOutput.ChecksumCRC32C
	fileResponse.ChecksumSHA1 = s3PutObjectOutput.ChecksumSHA1
	fileResponse.ChecksumSHA256 = s3PutObjectOutput.ChecksumSHA256
	fileResponse.ETag = s3PutObjectOutput.ETag
	fileResponse.Expiration = s3PutObjectOutput.Expiration

	fileURL, err := m.MinioClient.PresignedGetObject(ctx, m.getBucketName(), fileKey, 1*time.Hour, nil)
	if err != nil {
		m.log.Warnf("failed to generate presigned URL:" + err.Error())

		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	fileResponse.URL = fileURL.String()
	fileResponse.Filename = fileKey
	fileResponse.Mimetype = contentType
	fileResponse.Size = file.Size

	return fileResponse, nil
}

func (m *minioImpl) DeleteFile(ctx context.Context, fileName string) (bool, error) {

	err := m.MinioClient.RemoveObject(ctx, m.getBucketName(), fileName, minio.RemoveObjectOptions{ForceDelete: true})
	if err != nil {
		return false, fmt.Errorf("failed to delete file: %w", err)
	}

	return true, nil
}
