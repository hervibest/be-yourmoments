package grpc

import (
	"be-yourmoments/photo-svc/internal/model/converter"
	"be-yourmoments/photo-svc/internal/pb"
	"be-yourmoments/photo-svc/internal/usecase"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PhotoGRPCHandler struct {
	usecase usecase.PhotoUsecase
	pb.UnimplementedPhotoServiceServer
}

func NewPhotoGRPCHandler(server *grpc.Server, usecase usecase.PhotoUsecase) {
	handler := &PhotoGRPCHandler{
		usecase: usecase,
	}

	pb.RegisterPhotoServiceServer(server, handler)
}

func (h *PhotoGRPCHandler) UpdatePhotographerPhoto(ctx context.Context,
	pbReq *pb.UpdatePhotographerPhotoRequest) (
	*pb.UpdatePhotographerPhotoResponse, error) {

	req := converter.GrpcToCreateRequest(pbReq)
	h.usecase.UpdatePhoto(ctx, req)

	return nil, nil
}
func (h *PhotoGRPCHandler) UpdateFaceRecogPhoto(ctx context.Context, req *pb.UpdateFaceRecogPhotoRequest) (*pb.UpdateFaceRecogPhotoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateFaceRecogPhoto not implemented")
}
