package config

import (
	"be-yourmoments/photo-svc/internal/helper/logger"
	"be-yourmoments/photo-svc/internal/helper/utils"
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	MinioClient     *minio.Client
	minioBucketName string
	enpoint         string
	Logs            *logger.Log
}

func NewMinio() *Minio {
	logger := logger.New("minio")
	ctx := context.Background()
	var minioClient *minio.Client
	var err error

	minioHost := utils.GetEnv("MINIO_HOST")
	minioPort := utils.GetEnv("MINIO_PORT")
	minioRootUser := utils.GetEnv("MINIO_ROOT_USER")
	minioRootPassword := utils.GetEnv("MINIO_ROOT_PASSWORD")
	minioTicketsBucket := utils.GetEnv("MINIO_TICKETS_BUCKET")
	minioLocation := utils.GetEnv("MINIO_LOCATION")
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

	return &Minio{
		MinioClient:     minioClient,
		minioBucketName: minioTicketsBucket,
		enpoint:         endpoint,
		Logs:            logger,
	}
}

func (m *Minio) GetBucketName() string {
	BucketName := m.minioBucketName
	return BucketName
}

func (m *Minio) GetEndpoint() string {
	Endpoint := m.enpoint
	return Endpoint
}
