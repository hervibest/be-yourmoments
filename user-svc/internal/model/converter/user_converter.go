package converter

import (
	"be-yourmoments/user-svc/internal/entity"
	"be-yourmoments/user-svc/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		Id:                    user.Id,
		Username:              user.Username,
		Email:                 user.Email.String,
		EmailVerifiedAt:       user.EmailVerifiedAt,
		PhoneNumber:           user.PhoneNumber.String,
		PhoneNumberVerifiedAt: user.PhoneNumberVerifiedAt,
		GoogleId:              user.GoogleId.String,
		CreatedAt:             user.CreatedAt,
		UpdatedAt:             user.UpdatedAt,
	}
}
