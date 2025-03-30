package repository

import (
	"be-yourmoments/user-svc/internal/entity"
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type userProfilePreparedStmt struct {
	findById     *sqlx.Stmt
	findByUserId *sqlx.Stmt
}

func newUserProfilePreraredStmt(db *sqlx.DB) (*userProfilePreparedStmt, error) {
	findByUserIdStmt, err := db.Preparex("SELECT * FROM user_profiles WHERE user_id = $1")
	if err != nil {
		return nil, err
	}

	findByIdStmt, err := db.Preparex("SELECT * FROM user_profiles WHERE id = $1")
	if err != nil {
		return nil, err
	}

	return &userProfilePreparedStmt{
		findById:     findByIdStmt,
		findByUserId: findByUserIdStmt,
	}, nil
}

type UserProfileRepository interface {
	Close() error
	CreateWithProfileUrl(ctx context.Context, tx Querier, userProfile *entity.UserProfile) (*entity.UserProfile, error)
	Create(ctx context.Context, tx Querier, userProfile *entity.UserProfile) (*entity.UserProfile, error)
}

type userProfileRepository struct {
	userProfilePreparedStmt *userProfilePreparedStmt
}

func NewUserProfileRepository(db *sqlx.DB) (UserProfileRepository, error) {

	userProfilePreparedStmt, err := newUserProfilePreraredStmt(db)
	if err != nil {
		log.Print("error initialize userProfile category statement : ", err)
		return nil, err
	}

	return &userProfileRepository{
		userProfilePreparedStmt: userProfilePreparedStmt,
	}, nil
}

func (r *userProfileRepository) Close() error {
	if err := r.userProfilePreparedStmt.findById.Close(); err != nil {
		log.Print(err)
		return err
	}

	if err := r.userProfilePreparedStmt.findByUserId.Close(); err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (r *userProfileRepository) CreateWithProfileUrl(ctx context.Context, tx Querier, userProfile *entity.UserProfile) (*entity.UserProfile, error) {
	query := `INSERT INTO user_profiles 
	(id, user_id, nickname, profile_url, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := tx.ExecContext(ctx, query, userProfile.Id, userProfile.UserId, userProfile.Nickname, userProfile.ProfileUrl, userProfile.CreatedAt, userProfile.UpdatedAt)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert user profile: %w", err)
	}

	return userProfile, nil
}

func (r *userProfileRepository) Create(ctx context.Context, tx Querier, userProfile *entity.UserProfile) (*entity.UserProfile, error) {
	query := `INSERT INTO user_profiles 
	(id, user_id, birth_date, nickname, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := tx.ExecContext(ctx, query, userProfile.Id, userProfile.UserId, userProfile.BirthDate, userProfile.Nickname, userProfile.CreatedAt, userProfile.UpdatedAt)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert user profile: %w", err)
	}

	return userProfile, nil
}
