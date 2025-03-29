package adapter

import (
	discovery "be-yourmoments/upload-svc/internal/helper"
	"be-yourmoments/upload-svc/internal/pb"
	"context"
	"log"
)

type AiAdapter interface {
	ProcessPhoto(ctx context.Context, fileId, fileUrl string) error
	ProcessFacecam(ctx context.Context, userId, fileUrl string) error
}

type aiAdapter struct {
	client pb.AiServiceClient
}

func NewAiAdapter(ctx context.Context, registry discovery.Registry) (AiAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "grpc-ai-service", registry)
	if err != nil {
		return nil, err
	}

	log.Print("successfuly connected to grpc-ai-service")
	client := pb.NewAiServiceClient(conn)

	return &aiAdapter{
		client: client,
	}, nil
}

func (a *aiAdapter) ProcessPhoto(ctx context.Context, userId, fileUrl string) error {
	processPhotoRequest := &pb.ProcessPhotoRequest{
		Id:  userId,
		Url: fileUrl,
	}

	res, err := a.client.ProcessPhoto(ctx, processPhotoRequest)
	if err != nil {
		return err
	}

	log.Println(res)
	return nil
}

func (a *aiAdapter) ProcessFacecam(ctx context.Context, fileId, fileUrl string) error {
	log.Println("REQUESTED PROCESS FACECAM VIA GRPC TO AI SERVER")
	processPhotoRequest := &pb.ProcessFacecamRequest{
		Id:  fileId,
		Url: fileUrl,
	}

	res, err := a.client.ProcessFacecam(ctx, processPhotoRequest)
	if err != nil {
		return err
	}

	log.Println(res)
	return nil
}
