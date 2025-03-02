package adapter

import (
	"be-yourmoments/internal/adapter/pb"
	discovery "be-yourmoments/internal/helper"
	"context"
	"log"
)

type AiAdapter interface {
	ProcessPhoto(ctx context.Context, fileId, fileUrl string) error
}

type aiAdapter struct {
	client pb.AiServiceClient
}

func NewAiAdapter(ctx context.Context, registry discovery.Registry) (AiAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "ai-svc-grpc", registry)
	if err != nil {
		return nil, err
	}

	client := pb.NewAiServiceClient(conn)

	return &aiAdapter{
		client: client,
	}, nil
}

func (a *aiAdapter) ProcessPhoto(ctx context.Context, fileId, fileUrl string) error {
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
