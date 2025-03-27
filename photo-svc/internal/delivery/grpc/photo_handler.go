package grpc

import (
	"be-yourmoments/photo-svc/internal/pb"
	"be-yourmoments/photo-svc/internal/usecase"
	"context"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PhotoGRPCHandler struct {
	photoUseCase            usecase.PhotoUsecase
	userSimilarPhotoUseCase usecase.UserSimilarUsecase
	pb.UnimplementedPhotoServiceServer
}

func NewPhotoGRPCHandler(server *grpc.Server, photoUseCase usecase.PhotoUsecase,
	userSimilarPhotoUseCase usecase.UserSimilarUsecase) {
	handler := &PhotoGRPCHandler{
		photoUseCase:            photoUseCase,
		userSimilarPhotoUseCase: userSimilarPhotoUseCase,
	}

	pb.RegisterPhotoServiceServer(server, handler)
}

func (h *PhotoGRPCHandler) CreatePhoto(ctx context.Context, pbReq *pb.CreatePhotoRequest) (
	*pb.CreatePhotoResponse, error) {
	log.Println("----  CreatePhoto Requets via GRPC in photo-svc ------")
	if err := h.photoUseCase.CreatePhoto(context.Background(), pbReq); err != nil {
		return &pb.CreatePhotoResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	return &pb.CreatePhotoResponse{
		Status: http.StatusCreated,
	}, nil
}

func (h *PhotoGRPCHandler) CreateUserSimilarPhoto(ctx context.Context, pbReq *pb.CreateUserSimilarPhotoRequest) (
	*pb.CreateUserSimilarPhotoResponse, error) {
	log.Println("----  CreatePhoto Requets via GRPC in photo-svc ------")
	if err := h.userSimilarPhotoUseCase.CreateUserSimilarPhoto(context.Background(), pbReq); err != nil {
		return &pb.CreateUserSimilarPhotoResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	return &pb.CreateUserSimilarPhotoResponse{
		Status: http.StatusCreated,
	}, nil
}

func (h *PhotoGRPCHandler) UpdatePhotoDetail(ctx context.Context, pbReq *pb.UpdatePhotoDetailRequest) (
	*pb.UpdatePhotoDetailResponse, error) {
	log.Println("----  UpdatePhoto Requets via GRPC in photo-svc ------")
	if err := h.photoUseCase.UpdatePhotoDetail(context.Background(), pbReq); err != nil {
		return &pb.UpdatePhotoDetailResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	return &pb.UpdatePhotoDetailResponse{
		Status: http.StatusCreated,
	}, nil
}

func (h *PhotoGRPCHandler) UpdatePhotographerPhoto(ctx context.Context,
	pbReq *pb.UpdatePhotographerPhotoRequest) (
	*pb.UpdatePhotographerPhotoResponse, error) {

	// req := converter.GrpcToCreateRequest(pbReq)
	// h.usecase.UpdatePhoto(ctx, req)

	return nil, nil
}
func (h *PhotoGRPCHandler) UpdateFaceRecogPhoto(ctx context.Context, req *pb.UpdateFaceRecogPhotoRequest) (*pb.UpdateFaceRecogPhotoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateFaceRecogPhoto not implemented")
}
