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
	facecamUseCase          usecase.FacecamUseCase
	userSimilarPhotoUseCase usecase.UserSimilarUsecase
	pb.UnimplementedPhotoServiceServer
}

func NewPhotoGRPCHandler(server *grpc.Server, photoUseCase usecase.PhotoUsecase,
	facecamUseCase usecase.FacecamUseCase, userSimilarPhotoUseCase usecase.UserSimilarUsecase) {
	handler := &PhotoGRPCHandler{
		photoUseCase:            photoUseCase,
		facecamUseCase:          facecamUseCase,
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

func (h *PhotoGRPCHandler) CreateUserSimilar(ctx context.Context, pbReq *pb.CreateUserSimilarPhotoRequest) (
	*pb.CreateUserSimilarPhotoResponse, error) {
	log.Println("----  CreatePhoto user similar Requets via GRPC in photo-svc ------")
	if err := h.userSimilarPhotoUseCase.CreateUserSimilar(context.Background(), pbReq); err != nil {
		log.Println("error hapening in create user similar ", err.Error())
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

func (h *PhotoGRPCHandler) CreateFacecam(ctx context.Context, pbReq *pb.CreateFacecamRequest) (
	*pb.CreateFacecamResponse, error) {
	log.Println("----  Create facecam Requets via GRPC in photo-svc ------")
	if err := h.facecamUseCase.CreateFacecam(context.Background(), pbReq); err != nil {
		return &pb.CreateFacecamResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	return &pb.CreateFacecamResponse{
		Status: http.StatusCreated,
	}, nil
}

func (h *PhotoGRPCHandler) CreateUserSimilarFacecam(ctx context.Context, pbReq *pb.CreateUserSimilarFacecamRequest) (
	*pb.CreateUserSimilarFacecamResponse, error) {
	log.Println("----  CreatePhoto user similar Requets via GRPC in photo-svc ------")
	if err := h.userSimilarPhotoUseCase.CreateUserFacecam(context.Background(), pbReq); err != nil {
		log.Println("error hapening in create user similar ", err.Error())
		return &pb.CreateUserSimilarFacecamResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	return &pb.CreateUserSimilarFacecamResponse{
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
