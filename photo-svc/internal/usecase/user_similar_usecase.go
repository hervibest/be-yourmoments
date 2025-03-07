package usecase

type UserSimilarUsecase interface {
}

type userSimilarUsecase struct {
}

func NewUserSimilarUsecase() UserSimilarUsecase {
	return &userSimilarUsecase{}
}

func (u *userSimilarUsecase) UserGetSimilarPhoto() {

}
