package converter

import (
	"be-yourmoments/photo-svc/internal/model"
	"be-yourmoments/photo-svc/internal/pb"
)

func GrpcToCreateRequest(req *pb.UpdatePhotographerPhotoRequest) *model.RequestUpdateProcessedPhoto {

	userId := make([]string, len(req.UserId))
	for _, value := range req.UserId {
		userId = append(userId, value)
	}

	return &model.RequestUpdateProcessedPhoto{
		Id:     req.GetId(),
		UserId: userId,
	}

}
