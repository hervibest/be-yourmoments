package usecase

import (
	"be-yourmoments/user-svc/internal/adapter"
	"be-yourmoments/user-svc/internal/entity"
	"be-yourmoments/user-svc/internal/helper"
	"be-yourmoments/user-svc/internal/model"
	"be-yourmoments/user-svc/internal/model/converter"
	"be-yourmoments/user-svc/internal/repository"
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	AccessTokenRequest(ctx context.Context, refreshToken string) (*model.UserResponse, *model.TokenResponse, error)
	Current(ctx context.Context, email string) (*model.UserResponse, error)
	Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, *model.TokenResponse, error)
	Logout(ctx context.Context, request *model.LogoutUserRequest) (bool, error)
	RegisterByEmail(ctx context.Context, request *model.RegisterByEmailRequest) (*model.UserResponse, error)
	RegisterByGoogleSignIn(ctx context.Context, request *model.RegisterByGoogleRequest) (*model.UserResponse, *model.TokenResponse, error)
	RegisterByPhoneNumber(ctx context.Context, request *model.RegisterByPhoneRequest) (*model.UserResponse, error)
	RequestResetPassword(ctx context.Context, email string) error
	ResendEmailVerification(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, request *model.ResetPasswordUserRequest) error
	ValidateResetPassword(ctx context.Context, request *model.ValidateResetTokenRequest) (bool, error)
	Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.AuthResponse, error)
	VerifyEmail(ctx context.Context, request *model.VerifyEmailUserRequest) error
}

type authUseCase struct {
	db                    repository.BeginTx
	userRepository        repository.UserRepository
	userProfileRepository repository.UserProfileRepository
	emailVerificationRepo repository.EmailVerificationRepository
	resetPasswordRepo     repository.ResetPasswordRepository
	googleTokenAdapter    adapter.GoogleTokenAdapter
	emailAdapter          adapter.EmailAdapter
	securityAdapter       adapter.SecurityAdapter
	jwtAdapter            adapter.JWTAdapter
	cacheAdapter          adapter.CacheAdapter
}

func NewAuthUseCase(db repository.BeginTx, userRepository repository.UserRepository, userProfileRepository repository.UserProfileRepository,
	emailVerificationRepo repository.EmailVerificationRepository, resetPasswordRepo repository.ResetPasswordRepository,
	googleTokenAdapter adapter.GoogleTokenAdapter, emailAdapter adapter.EmailAdapter, jwtAdapter adapter.JWTAdapter,
	securityAdapter adapter.SecurityAdapter, cacheAdapter adapter.CacheAdapter) AuthUseCase {
	return &authUseCase{
		db:                    db,
		userRepository:        userRepository,
		userProfileRepository: userProfileRepository,
		emailVerificationRepo: emailVerificationRepo,
		resetPasswordRepo:     resetPasswordRepo,
		googleTokenAdapter:    googleTokenAdapter,
		emailAdapter:          emailAdapter,
		securityAdapter:       securityAdapter,
		jwtAdapter:            jwtAdapter,
		cacheAdapter:          cacheAdapter,
	}
}

func (u *authUseCase) RegisterByPhoneNumber(ctx context.Context, request *model.RegisterByPhoneRequest) (*model.UserResponse, error) {
	countByNumberTotal, err := u.userRepository.CountByPhoneNumber(ctx, request.PhoneNumber)
	if err != nil {
		return nil, err
	}

	if countByNumberTotal > 0 {
		return nil, fiber.NewError(http.StatusBadRequest, "phone number has already taken")
	}

	countByUsernameTotal, err := u.userRepository.CountByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}

	if countByUsernameTotal > 0 {
		return nil, fiber.NewError(http.StatusBadRequest, "username has already taken")
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

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	now := time.Now()
	user := &entity.User{
		Id:       ulid.Make().String(),
		Username: request.Username,
		Password: sql.NullString{
			String: string(hashedPassword),
		},
		PhoneNumber: sql.NullString{
			String: request.PhoneNumber,
		},
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	user, err = u.userRepository.CreateByPhoneNumber(ctx, tx, user)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	userProfile := &entity.UserProfile{
		Id:        ulid.Make().String(),
		UserId:    user.Id,
		BirthDate: request.BirthDate,
		Nickname:  helper.GenerateNickname(),
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	_, err = u.userProfileRepository.Create(ctx, tx, userProfile)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//WHAT TO DO WHATS APP BUSINESS VERIF OTP

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return nil, err
	}

	return converter.UserToResponse(user), nil
}

// TODO EFICIENT QUERY ISSUE (COUNT COUNT AND COUNT)
func (u *authUseCase) RegisterByGoogleSignIn(ctx context.Context, request *model.RegisterByGoogleRequest) (*model.UserResponse, *model.TokenResponse, error) {
	claims, err := u.googleTokenAdapter.ValidateGoogleToken(ctx, request.Token)
	if err != nil {
		fiber.NewError(http.StatusBadGateway, err.Error())
	}

	countByNotGoogleTotal, err := u.userRepository.CountByEmailNotGoogle(ctx, claims.Email)
	if err != nil {
		log.Println("eror count by email not google")
		return nil, nil, err
	}

	if countByNotGoogleTotal > 0 {
		return nil, nil, fiber.NewError(http.StatusBadRequest, "email has already takens")
	}

	countByGoogleTotal, err := u.userRepository.CountByEmailGoogleId(ctx, claims.Email, claims.GoogleId)
	if err != nil {
		log.Println("error count by email google id")
		return nil, nil, err
	}

	if countByGoogleTotal > 0 {
		user, err := u.userRepository.FindByEmail(ctx, claims.Email)
		if err != nil && errors.Is(sql.ErrNoRows, err) {
			log.Println(err)
			return nil, nil, fiber.NewError(fiber.StatusBadRequest, "invalid email")
		}

		token, err := u.generateToken(ctx, user.Id)
		if err != nil {
			log.Println(err)
			return nil, nil, err
		}

		return converter.UserToResponse(user), token, nil
	} else {
		tx, err := u.db.BeginTxx(ctx, nil)
		if err != nil {
			return nil, nil, err
		}

		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()

		now := time.Now()
		user := &entity.User{
			Id: ulid.Make().String(),
			Email: sql.NullString{
				String: claims.Email,
			},
			Username: claims.Username,
			GoogleId: sql.NullString{
				String: claims.Email,
			},
			CreatedAt: &now,
			UpdatedAt: &now,
		}

		user, err = u.userRepository.CreateByGoogleSignIn(ctx, tx, user)
		if err != nil {
			log.Println(err)
			return nil, nil, err
		}

		userProfile := &entity.UserProfile{
			Id:     ulid.Make().String(),
			UserId: user.Id,
			// BirthDate: request.BirthDate,
			Nickname:   helper.GenerateNickname(),
			ProfileUrl: claims.ProfilePictureUrl,
			CreatedAt:  &now,
			UpdatedAt:  &now,
		}

		_, err = u.userProfileRepository.CreateWithProfileUrl(ctx, tx, userProfile)
		if err != nil {
			log.Println(err)
			return nil, nil, err
		}

		if err := tx.Commit(); err != nil {
			log.Println(err)
			return nil, nil, err
		}

		token, err := u.generateToken(ctx, user.Id)
		if err != nil {
			log.Println(err)
			return nil, nil, err
		}

		return converter.UserToResponse(user), token, nil
	}
}

func (u *authUseCase) RegisterByEmail(ctx context.Context, request *model.RegisterByEmailRequest) (*model.UserResponse, error) {
	countByNumberTotal, err := u.userRepository.CountByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if countByNumberTotal > 0 {
		return nil, fiber.NewError(http.StatusBadRequest, "email has already takens")
	}

	countByUsernameTotal, err := u.userRepository.CountByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}

	if countByUsernameTotal > 0 {
		return nil, fiber.NewError(http.StatusBadRequest, "username has already takens")
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

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	now := time.Now()
	user := &entity.User{
		Id:       ulid.Make().String(),
		Username: request.Username,
		Email: sql.NullString{
			String: request.Email,
		},
		Password: sql.NullString{
			String: string(hashedPassword),
		},
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	user, err = u.userRepository.CreateByEmail(ctx, tx, user)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	userProfile := &entity.UserProfile{
		Id:        ulid.Make().String(),
		UserId:    user.Id,
		BirthDate: request.BirthDate,
		Nickname:  helper.GenerateNickname(),
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	_, err = u.userProfileRepository.Create(ctx, tx, userProfile)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// token := uuid.NewString()

	// emailVerification := &entity.EmailVerification{
	// 	Email:     request.Email,
	// 	Token:     token,
	// 	CreatedAt: &now,
	// 	UpdatedAt: &now,
	// }

	// emailVerification, err = u.emailVerificationRepo.Insert(ctx, tx, emailVerification)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return nil, err
	}

	if err := u.requestEmailVerification(ctx, request.Email, true); err != nil {
		log.Println(err)
		return nil, err
	}

	return converter.UserToResponse(user), nil
}

func (u *authUseCase) requestEmailVerification(ctx context.Context, email string, newUser bool) error {
	token := uuid.NewString()
	now := time.Now()
	emailVerification := &entity.EmailVerification{
		Email:     email,
		Token:     token,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if newUser {
		_, err = u.emailVerificationRepo.Insert(ctx, tx, emailVerification)
		if err != nil {
			return err
		}
	} else {
		_, err = u.emailVerificationRepo.Update(ctx, tx, emailVerification)
		if err != nil {
			return err
		}
	}

	encryptedToken, err := u.securityAdapter.Encrypt(token)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}

	if err := u.emailAdapter.SendEmail(email, encryptedToken, "new email verification"); err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (u *authUseCase) ResendEmailVerification(ctx context.Context, email string) error {
	user, err := u.userRepository.FindByEmail(ctx, email)
	if err != nil && errors.Is(sql.ErrNoRows, err) {
		log.Println(err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid email")
	}

	if user.HasVerifiedEmail() {
		return fiber.NewError(fiber.StatusBadRequest, "user already verified")
	}

	_, err = u.emailVerificationRepo.FindByEmail(ctx, email)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "resend request is not valid")
	}

	if err := u.requestEmailVerification(ctx, email, false); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "resend request is not valid")
	}

	return nil
}

func (u *authUseCase) VerifyEmail(ctx context.Context, request *model.VerifyEmailUserRequest) error {
	decryptedToken, err := u.securityAdapter.Decrypt(request.Token)
	if err != nil {
		log.Printf("failed to decrypt : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid token")
	}

	emailVerification, err := u.emailVerificationRepo.FindByEmailAndToken(ctx, request.Email, decryptedToken)
	if err != nil {
		log.Printf("failed find email and token : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid token")
	}

	if time.Since(*emailVerification.UpdatedAt) > 15*time.Minute {
		return fiber.NewError(fiber.StatusUnauthorized, "verification token expired")
	}

	now := time.Now()
	user := &entity.User{
		Email: sql.NullString{
			String: request.Email,
		},
		EmailVerifiedAt: &now,
		UpdatedAt:       &now,
	}

	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = u.userRepository.UpdateEmailVerifiedAt(ctx, tx, user)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := u.emailVerificationRepo.Delete(ctx, tx, emailVerification); err != nil {
		log.Println(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (u *authUseCase) RequestResetPassword(ctx context.Context, email string) error {
	_, err := u.userRepository.FindByEmail(ctx, email)
	if err != nil && errors.Is(sql.ErrNoRows, err) {
		log.Println(err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid email")

	}

	countByEmailTotal, err := u.resetPasswordRepo.CountByEmail(ctx, email)
	if err != nil {
		return err
	}

	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	token := uuid.NewString()
	now := time.Now()
	resetPassword := &entity.ResetPassword{
		Email:     email,
		Token:     token,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	if countByEmailTotal > 0 {
		_, err := u.resetPasswordRepo.Update(ctx, tx, resetPassword)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		_, err := u.resetPasswordRepo.Insert(ctx, tx, resetPassword)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}

	encryptedToken, err := u.securityAdapter.Encrypt(token)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := u.emailAdapter.SendEmail(email, encryptedToken, "reset password"); err != nil {
		return err
	}

	return nil
}

func (u *authUseCase) ValidateResetPassword(ctx context.Context, request *model.ValidateResetTokenRequest) (bool, error) {
	decryptedToken, err := u.securityAdapter.Decrypt(request.Token)
	if err != nil {
		log.Println(err)
		return false, fiber.NewError(fiber.StatusBadRequest, "invalid reset password token")
	}

	_, err = u.resetPasswordRepo.FindByEmailAndToken(ctx, request.Email, decryptedToken)
	if err != nil && errors.Is(sql.ErrNoRows, err) {
		log.Println(err)
		return false, fiber.NewError(fiber.StatusBadRequest, "invalid email")
	}

	return true, nil
}

func (u *authUseCase) ResetPassword(ctx context.Context, request *model.ResetPasswordUserRequest) error {
	decryptedToken, err := u.securityAdapter.Decrypt(request.Token)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid token")
	}

	resetPassword, err := u.resetPasswordRepo.FindByEmailAndToken(ctx, request.Email, decryptedToken)
	if err != nil && errors.Is(sql.ErrNoRows, err) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid token")
	}

	if time.Since(*resetPassword.UpdatedAt) > 15*time.Minute {
		return fiber.NewError(fiber.StatusUnauthorized, "reset password token expired")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	now := time.Now()
	user := &entity.User{
		Email: sql.NullString{
			String: request.Email,
		},
		Password: sql.NullString{
			String: string(hashedPassword),
		},
		UpdatedAt: &now,
	}

	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = u.userRepository.UpdatePassword(ctx, tx, user)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := u.resetPasswordRepo.Delete(ctx, tx, resetPassword); err != nil {
		log.Println(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (u *authUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, *model.TokenResponse, error) {
	user, err := u.userRepository.FindByMultipleParam(ctx, request.MultipleParam)
	if err != nil && errors.Is(sql.ErrNoRows, err) {
		log.Println(err)
		return nil, nil, fiber.NewError(fiber.StatusBadRequest, "invalid email or password")
	}

	if !(user.HasEmail() && user.HasVerifiedEmail()) &&
		!(user.HasPhoneNumber() && user.HasVerifiedPhoneNumber()) {
		return nil, nil, fiber.NewError(fiber.StatusBadRequest, "email or phone number must be verified")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(request.Password)); err != nil {
		return nil, nil, fiber.NewError(fiber.StatusBadRequest, "invalid password")
	}

	token, err := u.generateToken(ctx, user.Id)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	return converter.UserToResponse(user), token, nil
}

func (u *authUseCase) generateToken(ctx context.Context, userId string) (*model.TokenResponse, error) {
	accessTokenDetail, err := u.jwtAdapter.GenerateAccessToken(userId)
	if err != nil {
		log.Printf("failed to generate access token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	refreshTokenDetail, err := u.jwtAdapter.GenerateRefreshToken(userId)
	if err != nil {
		log.Printf("failed to generate refresh token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := u.cacheAdapter.Set(ctx, refreshTokenDetail.Token, userId, time.Until(refreshTokenDetail.ExpiresAt)); err != nil {
		log.Printf("failed to save refresh token to cache : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	token := &model.TokenResponse{
		AccessToken:  accessTokenDetail.Token,
		RefreshToken: refreshTokenDetail.Token,
	}

	return token, nil
}

func (u *authUseCase) Current(ctx context.Context, email string) (*model.UserResponse, error) {
	user, err := u.userRepository.FindByEmail(ctx, email)
	if err != nil && errors.Is(sql.ErrNoRows, err) {
		return nil, fiber.ErrNotFound
	}

	return converter.UserToResponse(user), nil
}

func (u *authUseCase) Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.AuthResponse, error) {
	accessTokenDetail, err := u.jwtAdapter.VerifyAccessToken(request.Token)
	if err != nil {
		log.Printf("failed to verify access token : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	userId, err := u.cacheAdapter.Get(ctx, request.Token)
	if userId != "" {
		log.Println("access token found in Redis, user has already signed out")
		return nil, fiber.ErrUnauthorized
	}

	user, err := u.userRepository.FindById(ctx, accessTokenDetail.UserId)
	if err != nil && errors.Is(sql.ErrNoRows, err) {
		log.Printf("failed to verify access token : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	authResponse := &model.AuthResponse{
		UserId:      user.Id,
		Username:    user.Username,
		Email:       user.Email.String,
		PhoneNumber: user.PhoneNumber.String,
		Token:       request.Token,
		ExpiresAt:   accessTokenDetail.ExpiresAt,
	}

	return authResponse, nil
}

// TODO does ErrInternalServerError considered to be a best practice ?
func (u *authUseCase) Logout(ctx context.Context, request *model.LogoutUserRequest) (bool, error) {

	if err := u.cacheAdapter.Set(ctx, request.AccessToken, "revoked", time.Until(request.ExpiresAt)); err != nil {
		log.Printf("failed to save access token to cache : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := u.cacheAdapter.Del(ctx, request.RefreshToken); err != nil {
		log.Printf("failed to save refresh token to cache : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}

func (u *authUseCase) AccessTokenRequest(ctx context.Context, refreshToken string) (*model.UserResponse, *model.TokenResponse, error) {

	refreshTokenDetail, err := u.jwtAdapter.VerifyRefreshToken(refreshToken)
	if err != nil {
		log.Printf("failed to verify refresh token : %+v", err)
		return nil, nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")

	}

	userId, err := u.cacheAdapter.Get(ctx, refreshToken)
	if userId == "" {
		log.Printf("refresh token not found in Redis : %+v", err)
		return nil, nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	user, err := u.userRepository.FindById(ctx, refreshTokenDetail.UserId)
	if err != nil && errors.Is(sql.ErrNoRows, err) {
		log.Println(err)
		return nil, nil, fiber.NewError(fiber.StatusUnauthorized, "invalid refresh token")
	}

	//TODO !!! ACCESS TOKEN TO RESPONSE
	accessTokenDetail, err := u.jwtAdapter.GenerateAccessToken(user.Id)
	if err != nil {
		log.Printf("failed to generate access token : %+v", err)
		return nil, nil, fiber.ErrInternalServerError
	}

	tokenResponse := &model.TokenResponse{
		AccessToken: accessTokenDetail.Token,
	}

	//TODO time duration refresh token to be renewed
	if time.Until(refreshTokenDetail.ExpiresAt) < time.Hour*24 {
		newRefreshTokenDetail, err := u.jwtAdapter.GenerateRefreshToken(user.Id)
		if err != nil {
			log.Printf("failed to generate refresh token : %+v", err)
			return nil, nil, fiber.ErrInternalServerError
		}

		if err := u.cacheAdapter.Set(ctx, newRefreshTokenDetail.Token, user.Id, time.Until(refreshTokenDetail.ExpiresAt)); err != nil {
			log.Printf("failed to save refresh token to cache : %+v", err)
			return nil, nil, fiber.ErrInternalServerError
		}

		tokenResponse.RefreshToken = newRefreshTokenDetail.Token

	} else {
		tokenResponse.RefreshToken = refreshTokenDetail.Token

	}

	return converter.UserToResponse(user), tokenResponse, nil
}
