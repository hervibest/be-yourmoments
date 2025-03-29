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
	CreatePhoto(ctx context.Context, photo *entity.Photo, facecam *entity.PhotoDetail) error
	UpdatePhotoDetail(ctx context.Context, facecam *entity.PhotoDetail) error
	CreateFacecam(ctx context.Context, facecam *entity.Facecam) error
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

func (a *photoAdapter) CreatePhoto(ctx context.Context, photo *entity.Photo, facecam *entity.PhotoDetail) error {

	facecampb := &pb.PhotoDetail{
		Id:              facecam.Id,
		PhotoId:         facecam.PhotoId,
		FileName:        facecam.FileName,
		FileKey:         facecam.FileKey,
		Size:            facecam.Size,
		Type:            facecam.Type,
		Checksum:        facecam.Checksum,
		Width:           int32(facecam.Width),
		Height:          int32(facecam.Height),
		Url:             facecam.Url,
		YourMomentsType: string(facecam.YourMomentsType),
		CreatedAt: &timestamppb.Timestamp{
			Seconds: facecam.CreatedAt.Unix(),
			Nanos:   int32(facecam.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: facecam.UpdatedAt.Unix(),
			Nanos:   int32(facecam.UpdatedAt.Nanosecond()),
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

		Detail: facecampb,
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

func (a *photoAdapter) UpdatePhotoDetail(ctx context.Context, facecam *entity.PhotoDetail) error {

	facecampb := &pb.PhotoDetail{
		Id:              facecam.Id,
		PhotoId:         facecam.PhotoId,
		FileName:        facecam.FileName,
		FileKey:         facecam.FileKey,
		Size:            facecam.Size,
		Type:            facecam.Type,
		Checksum:        facecam.Checksum,
		Width:           int32(facecam.Width),
		Height:          int32(facecam.Height),
		Url:             facecam.Url,
		YourMomentsType: string(facecam.YourMomentsType),
		CreatedAt: &timestamppb.Timestamp{
			Seconds: facecam.CreatedAt.Unix(),
			Nanos:   int32(facecam.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: facecam.UpdatedAt.Unix(),
			Nanos:   int32(facecam.UpdatedAt.Nanosecond()),
		},
	}
	pbRequest := &pb.UpdatePhotoDetailRequest{
		PhotoDetail: facecampb,
	}

	res, err := a.client.UpdatePhotoDetail(context.Background(), pbRequest)
	if res.Status >= 400 || res.Error != "" || err != nil {
		return errors.New(res.Error)
	}

	return nil
}

func (a *photoAdapter) CreateFacecam(ctx context.Context, facecam *entity.Facecam) error {

	facecamPb := &pb.Facecam{
		Id:       facecam.Id,
		UserId:   facecam.UserId,
		FileName: facecam.FileName,
		FileKey:  facecam.FileKey,
		Size:     facecam.Size,
		Checksum: facecam.Checksum,
		Url:      facecam.Url,
		CreatedAt: &timestamppb.Timestamp{
			Seconds: facecam.CreatedAt.Unix(),
			Nanos:   int32(facecam.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: facecam.UpdatedAt.Unix(),
			Nanos:   int32(facecam.UpdatedAt.Nanosecond()),
		},
	}

	pbRequest := &pb.CreateFacecamRequest{
		Facecam: facecamPb,
	}

	res, err := a.client.CreateFacecam(context.Background(), pbRequest)
	if res.Status >= 400 || res.Error != "" || err != nil {
		return errors.New(res.Error)
	}

	return nil
}
