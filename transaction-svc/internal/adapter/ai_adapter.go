package adapter

import (
	discovery "be-yourmoments/transaction-svc/internal/helper"
	"be-yourmoments/transaction-svc/internal/pb"
	"context"
	"log"
)

type PhotoAdapter interface {
	ProcessPhoto(ctx context.Context, fileId, fileUrl string) error
}

type photoAdapter struct {
	client pb.PhotoServiceClient
}

func NewPhotoAdapter(ctx context.Context, registry discovery.Registry) (PhotoAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "transaction-svc-grpc", registry)
	if err != nil {
		return nil, err
	}

	client := pb.NewPhotoServiceClient(conn)

	return &photoAdapter{
		client: client,
	}, nil
}

func (a *photoAdapter) ProcessPhoto(ctx context.Context, fileId, fileUrl string) error {
	processPhotoRequest := &pb.ProcessPhotoRequest{
		Id:  fileId,
		Url: fileUrl,
	}

	res, err := a.client.ProcessPhoto(ctx, processPhotoRequest)
	if err != nil {
		return err
	}

	log.Println(res)
	return nil
}
