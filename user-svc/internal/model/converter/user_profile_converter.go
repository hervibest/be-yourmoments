package converter

import (
	"be-yourmoments/user-svc/internal/entity"
	"be-yourmoments/user-svc/internal/model"
)

func UserProfileToResponse(userProfile *entity.UserProfile, profileUrl, coverUrl string) *model.UserProfileResponse {
	return &model.UserProfileResponse{
		Id:              userProfile.Id,
		UserId:          userProfile.UserId,
		BirthDate:       userProfile.BirthDate,
		Nickname:        userProfile.Nickname,
		Biography:       userProfile.Biography.String,
		ProfileUrl:      profileUrl,
		ProfileCoverUrl: coverUrl,
		Similarity:      userProfile.Similarity.String,
		CreatedAt:       userProfile.CreatedAt,
		UpdatedAt:       userProfile.UpdatedAt,
	}
}
