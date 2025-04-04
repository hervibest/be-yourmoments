package usecase

import "context"

type AdminUseCase interface {
	CreateSocialMedia(ctx context.Context)
	GetSocialMedia(ctx context.Context)
	GetAllSocialMedia(ctx context.Context)
	UpdateSocialMedia(ctx context.Context)
	DeleteSocialMedia(ctx context.Context)
}

type adminUseCase struct{}

func NewAdminUseCase() AdminUseCase {
	return &adminUseCase{}
}

func (u *adminUseCase) CreateSocialMedia(ctx context.Context)
func (u *adminUseCase) GetSocialMedia(ctx context.Context)
func (u *adminUseCase) GetAllSocialMedia(ctx context.Context)
func (u *adminUseCase) UpdateSocialMedia(ctx context.Context)
func (u *adminUseCase) DeleteSocialMedia(ctx context.Context)
