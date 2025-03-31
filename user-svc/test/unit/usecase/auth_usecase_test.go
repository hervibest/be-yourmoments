package usecase

import (
	"be-yourmoments/user-svc/internal/entity"
	mockadapter "be-yourmoments/user-svc/internal/mocks/adapter"
	mockdb "be-yourmoments/user-svc/internal/mocks/db"
	mockrepository "be-yourmoments/user-svc/internal/mocks/repository"
	"be-yourmoments/user-svc/internal/model"
	"be-yourmoments/user-svc/internal/usecase"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	// Import package internal (sesuaikan path-nya)
)

func TestRegisterByPhoneNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	// Buat mock untuk dependency
	mockUserRepo := mockrepository.NewMockUserRepository(ctrl)
	mockUserProfileRepo := mockrepository.NewMockUserProfileRepository(ctrl)
	mockEmailVerificationRepo := mockrepository.NewMockEmailVerificationRepository(ctrl)
	mockResetPasswordRepo := mockrepository.NewMockResetPasswordRepository(ctrl)
	mockDB := mockdb.NewMockBeginTx(ctrl)       // misal DB interface memiliki method BeginTxx(ctx, opts)
	mockTx := mockdb.NewMockTransactionTx(ctrl) // misal Tx interface dengan Commit() dan Rollback()

	mockCacheAdapter := mockadapter.NewMockCacheAdapter(ctrl)
	mockEmailAdapter := mockadapter.NewMockEmailAdapter(ctrl)
	mockGoogleTokenAdapter := mockadapter.NewMockGoogleTokenAdapter(ctrl)
	mockJwtAdapter := mockadapter.NewMockJWTAdapter(ctrl)
	mockSecurityAdapter := mockadapter.NewMockSecurityAdapter(ctrl)

	authUC := usecase.NewAuthUseCase(mockDB, mockUserRepo, mockUserProfileRepo, mockEmailVerificationRepo, mockResetPasswordRepo, mockGoogleTokenAdapter,
		mockEmailAdapter, mockJwtAdapter, mockSecurityAdapter, mockCacheAdapter)
	// Data request testing
	now := time.Now()
	req := &model.RegisterByPhoneRequest{
		PhoneNumber: "08123456789",
		Username:    "testuser",
		Password:    "secret",
		BirthDate:   &now,
	}

	t.Run("Phone number already taken", func(t *testing.T) {
		// Ekspektasi: jumlah data berdasarkan phone number > 0
		mockUserRepo.EXPECT().CountByPhoneNumber(ctx, req.PhoneNumber).Return(1, nil)

		resp, err := authUC.RegisterByPhoneNumber(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)

		// Pastikan error-nya sesuai dengan fiber.NewError(http.StatusBadRequest, "phone number has already taken")
		fiberErr, ok := err.(*fiber.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, fiberErr.Code)
		assert.Equal(t, "phone number has already taken", fiberErr.Message)
	})

	t.Run("Username already taken", func(t *testing.T) {
		// Nomor telepon tidak terdaftar
		mockUserRepo.EXPECT().CountByPhoneNumber(ctx, req.PhoneNumber).Return(0, nil)
		// Ekspektasi: username sudah ada
		mockUserRepo.EXPECT().CountByUsername(ctx, req.Username).Return(1, nil)

		resp, err := authUC.RegisterByPhoneNumber(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)

		fiberErr, ok := err.(*fiber.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, fiberErr.Code)
		assert.Equal(t, "username has already taken", fiberErr.Message)
	})

	t.Run("Successful registration", func(t *testing.T) {
		// Nomor telepon dan username belum terdaftar
		mockUserRepo.EXPECT().CountByPhoneNumber(ctx, req.PhoneNumber).Return(0, nil)
		mockUserRepo.EXPECT().CountByUsername(ctx, req.Username).Return(0, nil)

		// Mulai transaksi
		mockDB.EXPECT().BeginTxx(ctx, gomock.Any()).Return(mockTx, nil)

		// Ekspektasi pembuatan user; gunakan matcher Any() untuk parameter user
		var createdUser *entity.User
		mockUserRepo.EXPECT().CreateByPhoneNumber(ctx, mockTx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, tx interface{}, user *entity.User) (*entity.User, error) {
				createdUser = user
				// Simulasikan kembalian user yang telah dibuat
				return user, nil
			},
		)

		// Ekspektasi pembuatan user profile
		mockUserProfileRepo.EXPECT().Create(ctx, mockTx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, tx interface{}, profile *entity.UserProfile) (*entity.UserProfile, error) {
				// Misal set id profile dan pastikan UserId sesuai dengan user yang dibuat
				profile.Id = "profile-123"
				profile.UserId = createdUser.Id
				return profile, nil
			},
		)

		// Ekspektasi commit transaksi
		mockTx.EXPECT().Commit().Return(nil)

		resp, err := authUC.RegisterByPhoneNumber(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		// Verifikasi field yang diharapkan, misalnya username
		assert.Equal(t, req.Username, resp.Username)
		// Asersi tambahan dapat disesuaikan dengan hasil konversi dari converter.UserToResponse
	})

	t.Run("Transaction commit error", func(t *testing.T) {
		// Skenario: commit transaksi gagal

		// Nomor telepon dan username belum terdaftar
		mockUserRepo.EXPECT().CountByPhoneNumber(ctx, req.PhoneNumber).Return(0, nil)
		mockUserRepo.EXPECT().CountByUsername(ctx, req.Username).Return(0, nil)

		// Mulai transaksi
		mockDB.EXPECT().BeginTxx(ctx, gomock.Any()).Return(mockTx, nil)

		// Ekspektasi pembuatan user
		mockUserRepo.EXPECT().CreateByPhoneNumber(ctx, mockTx, gomock.Any()).Return(&entity.User{
			Id:          "user-123",
			Username:    req.Username,
			Password:    sql.NullString{String: "hashed"},
			PhoneNumber: sql.NullString{String: req.PhoneNumber},
			CreatedAt:   &time.Time{},
			UpdatedAt:   &time.Time{},
		}, nil)

		// Ekspektasi pembuatan user profile
		mockUserProfileRepo.EXPECT().Create(ctx, mockTx, gomock.Any()).Return(&entity.UserProfile{
			Id: "profile-123",
		}, nil)

		// Simulasikan error pada commit transaksi
		commitErr := errors.New("commit failed")
		mockTx.EXPECT().Commit().Return(commitErr)

		resp, err := authUC.RegisterByPhoneNumber(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, commitErr, err)
	})
}
