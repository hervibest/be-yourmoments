package usecase

import (
	"be-yourmoments/user-svc/internal/adapter"
	"be-yourmoments/user-svc/internal/entity"
	"be-yourmoments/user-svc/internal/enum"
	"be-yourmoments/user-svc/internal/model"
	"be-yourmoments/user-svc/internal/model/converter"
	"be-yourmoments/user-svc/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

type UserUseCase interface {
	GetUserProfile(ctx context.Context, userId string) (*model.UserProfileResponse, error)
	UpdateUserProfile(ctx context.Context, request *model.RequestUpdateUserProfile) (*model.UserProfileResponse, error)
	UpdateUserProfileImage(ctx context.Context, file *multipart.FileHeader, userProfId string) (bool, error)
	UpdateUserCoverImage(ctx context.Context, file *multipart.FileHeader, userProfId string) (bool, error)
}

type userUseCase struct {
	db                    repository.BeginTx
	userRepository        repository.UserRepository
	userProfileRepository repository.UserProfileRepository
	userImageRepository   repository.UserImageRepository
	uploadAdapter         adapter.UploadAdapter
}

func NewUserUseCase(db repository.BeginTx, userRepository repository.UserRepository, userProfileRepository repository.UserProfileRepository,
	userImageRepository repository.UserImageRepository, uploadAdapter adapter.UploadAdapter) UserUseCase {
	return &userUseCase{
		db:                    db,
		userRepository:        userRepository,
		userProfileRepository: userProfileRepository,
		userImageRepository:   userImageRepository,
		uploadAdapter:         uploadAdapter,
	}
}

func (u *userUseCase) GetUserProfile(ctx context.Context, userId string) (*model.UserProfileResponse, error) {
	log.Print(userId)
	userProfile, err := u.userProfileRepository.FindByUserId(ctx, userId)
	if err != nil && errors.Is(sql.ErrNoRows, err) {
		return nil, fiber.ErrNotFound
	}

	userImages, err := u.userImageRepository.FindByUserProfId(ctx, userProfile.Id)
	if err != nil && errors.Is(sql.ErrNoRows, err) {
		return nil, fiber.ErrNotFound
	}

	profileUrl, coverUrl, err := u.getUserImageUrl(ctx, userImages)
	if err != nil {
		return nil, fiber.ErrNotFound
	}

	return converter.UserProfileToResponse(userProfile, profileUrl, coverUrl), nil
}

func (u *userUseCase) getUserImageUrl(ctx context.Context, userImages *[]*entity.UserImage) (string, string, error) {
	var (
		profileUrl string
		coverUrl   string
		err        error
	)

	for _, userImage := range *userImages {
		if userImage.ImageType == enum.ImageTypeProfile {
			profileUrl, err = u.uploadAdapter.GetPresignedUrl(ctx, userImage.FileName, userImage.FileKey)
			if err != nil {
			}
		} else if userImage.ImageType == enum.ImageTypeCover {
			coverUrl, err = u.uploadAdapter.GetPresignedUrl(ctx, userImage.FileName, userImage.FileKey)
			if err != nil {
			}
		}
	}

	return profileUrl, coverUrl, nil
}

func (u *userUseCase) UpdateUserProfile(ctx context.Context, request *model.RequestUpdateUserProfile) (*model.UserProfileResponse, error) {
	now := time.Now()
	userProfile := &entity.UserProfile{
		UserId:    request.UserId,
		BirthDate: request.BirthDate,
		Nickname:  request.Nickname,
		Biography: sql.NullString{
			String: request.Biography,
		},
		UpdatedAt: &now,
	}

	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	userProfile, err = u.userProfileRepository.Update(ctx, tx, userProfile)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return nil, err
	}

	userImages, err := u.userImageRepository.FindByUserProfId(ctx, userProfile.Id)
	if err != nil && errors.Is(sql.ErrNoRows, err) {
		return nil, fiber.ErrNotFound
	}

	profileUrl, coverUrl, err := u.getUserImageUrl(ctx, userImages)
	if err != nil {
		return nil, fiber.ErrNotFound
	}

	return converter.UserProfileToResponse(userProfile, profileUrl, coverUrl), nil
}

func (u *userUseCase) UpdateUserProfileImage(ctx context.Context, file *multipart.FileHeader, userProfId string) (bool, error) {
	ok, err := u.updateUserImage(ctx, file, userProfId, enum.ImageTypeProfile)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return ok, err
}

func (u *userUseCase) UpdateUserCoverImage(ctx context.Context, file *multipart.FileHeader, userProfId string) (bool, error) {
	ok, err := u.updateUserImage(ctx, file, userProfId, enum.ImageTypeCover)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return ok, err
}

func (u *userUseCase) updateUserImage(ctx context.Context, file *multipart.FileHeader, userProfId string, imageType enum.ImageTypeEnum) (bool, error) {
	prevUserProfileImage, errRepo := u.userImageRepository.FindByUserProfIdAndType(ctx, userProfId, string(imageType))
	if errRepo != nil && !errors.Is(sql.ErrNoRows, errRepo) {
		log.Println(errRepo)
		return false, errRepo
	}

	uploadFile, err := file.Open()
	if err != nil {
		return false, err
	}
	defer uploadFile.Close()

	upload, err := u.uploadAdapter.UploadFile(ctx, file, uploadFile, fmt.Sprintf("user/profile/%s", userProfId))
	if err != nil {
		return false, err
	}

	now := time.Now()

	if errors.Is(sql.ErrNoRows, errRepo) {
		newUserProfileImage := &entity.UserImage{
			Id:            ulid.Make().String(),
			UserProfileId: userProfId,
			FileName:      upload.Filename,
			FileKey:       upload.FileKey,
			ImageType:     imageType,
			Size:          upload.Size,
			CreatedAt:     &now,
			UpdatedAt:     &now,
		}

		tx, err := u.db.BeginTxx(ctx, nil)
		if err != nil {
			return false, err
		}

		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()

		_, err = u.userImageRepository.Create(ctx, tx, newUserProfileImage)
		if err != nil {
			log.Println(err)
			return false, err
		}

		if err := tx.Commit(); err != nil {
			log.Println(err)
			return false, err
		}

	} else {
		updatedUserProfileImage := &entity.UserImage{
			UserProfileId: userProfId,
			FileName:      upload.Filename,
			FileKey:       upload.FileKey,
			ImageType:     imageType,
			Size:          upload.Size,
			CreatedAt:     &now,
			UpdatedAt:     &now,
		}

		tx, err := u.db.BeginTxx(ctx, nil)
		if err != nil {
			return false, err
		}

		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()

		_, err = u.userImageRepository.Update(ctx, tx, updatedUserProfileImage)
		if err != nil {
			log.Println(err)
			return false, err
		}

		_, err = u.uploadAdapter.DeleteFile(ctx, prevUserProfileImage.FileName)
		if err != nil {
			log.Println(err)
			return false, err
		}

		if err := tx.Commit(); err != nil {
			log.Println(err)
			return false, err
		}

	}

	return true, nil
}

// func (u *userUseCase) UpdateUserProfileCover(ctx context.Context, userId string) (*model.UserProfileResponse, error) {
// 	now := time.Now()
// 	userProfile := &entity.UserProfile{
// 		UserId: userId,
// 		ProfileCoverUrl: sql.NullString{
// 			String: request.ProfileCoverUrl,
// 		},

// 		UpdatedAt: &now,
// 	}

// 	tx, err := u.db.BeginTxx(ctx, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	userProfile, err = u.userProfileRepository.UpdateUserProfileCover(ctx, tx, userProfile)
// 	if err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}

// 	if err := tx.Commit(); err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}

// 	return converter.UserProfileToResponse(userProfile), nil
// }
