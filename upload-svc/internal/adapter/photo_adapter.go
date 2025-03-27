package adapter

import (
	"be-yourmoments/upload-svc/internal/entity"
	discovery "be-yourmoments/upload-svc/internal/helper"
	"be-yourmoments/upload-svc/internal/pb"
	"context"
	"errors"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type PhotoAdapter interface {
	CreatePhoto(ctx context.Context, photo *entity.Photo, photoDetail *entity.PhotoDetail) error
	UpdatePhotoDetail(ctx context.Context, photoDetail *entity.PhotoDetail) error
}

type photoAdapter struct {
	client pb.PhotoServiceClient
}

func NewPhotoAdapter(ctx context.Context, registry discovery.Registry) (PhotoAdapter, error) {
	conn, err := discovery.ServiceConnection(ctx, "photo-svc-grpc", registry)
	if err != nil {
		return nil, err
	}
	client := pb.NewPhotoServiceClient(conn)

	return &photoAdapter{
		client: client,
	}, nil
}

func (a *photoAdapter) CreatePhoto(ctx context.Context, photo *entity.Photo, photoDetail *entity.PhotoDetail) error {

	photoDetailpb := &pb.PhotoDetail{
		Id:              photoDetail.Id,
		PhotoId:         photoDetail.PhotoId,
		FileName:        photoDetail.FileName,
		FileKey:         photoDetail.FileKey,
		Size:            photoDetail.Size,
		Type:            photoDetail.Type,
		Checksum:        photoDetail.Checksum,
		Width:           int32(photoDetail.Width),
		Height:          int32(photoDetail.Height),
		Url:             photoDetail.Url,
		YourMomentsType: string(photoDetail.YourMomentsType),
		CreatedAt: &timestamppb.Timestamp{
			Seconds: photoDetail.CreatedAt.Unix(),
			Nanos:   int32(photoDetail.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: photoDetail.UpdatedAt.Unix(),
			Nanos:   int32(photoDetail.UpdatedAt.Nanosecond()),
		},
	}

	photoPb := &pb.Photo{
		Id:            photo.Id,
		CreatorId:     photo.CreatorId,
		Title:         photo.Title,
		CollectionUrl: photo.CollectionUrl,
		Price:         int32(photo.Price),
		PriceStr:      photo.PriceStr,

		OriginalAt: &timestamppb.Timestamp{
			Seconds: photo.OriginalAt.Unix(),
			Nanos:   int32(photo.OriginalAt.Nanosecond()),
		},
		CreatedAt: &timestamppb.Timestamp{
			Seconds: photo.CreatedAt.Unix(),
			Nanos:   int32(photo.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: photo.UpdatedAt.Unix(),
			Nanos:   int32(photo.UpdatedAt.Nanosecond()),
		},

		Detail: photoDetailpb,
	}

	pbRequest := &pb.CreatePhotoRequest{
		Photo: photoPb,
	}

	res, err := a.client.CreatePhoto(context.Background(), pbRequest)
	if res.Status >= 400 || res.Error != "" || err != nil {
		return errors.New(res.Error)
	}

	return nil
}

func (a *photoAdapter) UpdatePhotoDetail(ctx context.Context, photoDetail *entity.PhotoDetail) error {

	photoDetailpb := &pb.PhotoDetail{
		Id:              photoDetail.Id,
		PhotoId:         photoDetail.PhotoId,
		FileName:        photoDetail.FileName,
		FileKey:         photoDetail.FileKey,
		Size:            photoDetail.Size,
		Type:            photoDetail.Type,
		Checksum:        photoDetail.Checksum,
		Width:           int32(photoDetail.Width),
		Height:          int32(photoDetail.Height),
		Url:             photoDetail.Url,
		YourMomentsType: string(photoDetail.YourMomentsType),
		CreatedAt: &timestamppb.Timestamp{
			Seconds: photoDetail.CreatedAt.Unix(),
			Nanos:   int32(photoDetail.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: photoDetail.UpdatedAt.Unix(),
			Nanos:   int32(photoDetail.UpdatedAt.Nanosecond()),
		},
	}
	pbRequest := &pb.UpdatePhotoDetailRequest{
		PhotoDetail: photoDetailpb,
	}

	res, err := a.client.UpdatePhotoDetail(context.Background(), pbRequest)
	if res.Status >= 400 || res.Error != "" || err != nil {
		return errors.New(res.Error)
	}

	return nil
}
