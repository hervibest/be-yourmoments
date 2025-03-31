package repository

import (
	"be-yourmoments/user-svc/internal/entity"
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type userPreparedStmt struct {
	findById            *sqlx.Stmt
	findByEmail         *sqlx.Stmt
	findByMultipleParam *sqlx.Stmt

	countByEmail          *sqlx.Stmt
	countByUsername       *sqlx.Stmt
	countByPhoneNumber    *sqlx.Stmt
	countByEmailGoogleId  *sqlx.Stmt
	countByEmailNotGoogle *sqlx.Stmt
}

func newUserPreparedStmt(db *sqlx.DB) (*userPreparedStmt, error) {
	//WHAT TO DO QUERY JANGAN ASAL * DOANG
	findByIdStmt, err := db.Preparex("SELECT * FROM users WHERE id = $1")
	if err != nil {
		return nil, err
	}

	findByEmailStmt, err := db.Preparex("SELECT * FROM users WHERE email = $1")
	if err != nil {
		return nil, err
	}

	findByMultipleParamStmt, err := db.Preparex(`SELECT * FROM users WHERE email = $1 
	OR username = $1 OR phone_number = $1 AND google_id IS NULL`)
	if err != nil {
		return nil, err
	}

	countByEmailStmt, err := db.Preparex("SELECT COUNT(*) FROM users WHERE email = $1")
	if err != nil {
		return nil, err
	}

	countByEmailNotGoogleStmt, err := db.Preparex("SELECT COUNT(*) FROM users WHERE email = $1 AND google_id IS NULL")
	if err != nil {
		return nil, err
	}

	countByUsernameStmt, err := db.Preparex("SELECT COUNT(*) FROM users WHERE username = $1")
	if err != nil {
		return nil, err
	}

	countByPhoneNumberStmt, err := db.Preparex("SELECT COUNT(*) FROM users WHERE phone_number  = $1")
	if err != nil {
		return nil, err
	}

	countByEmailGoogleIdStmt, err := db.Preparex("SELECT COUNT(*) FROM users WHERE email = $1 AND google_id = $2")
	if err != nil {
		return nil, err
	}

	return &userPreparedStmt{
		findById:              findByIdStmt,
		findByEmail:           findByEmailStmt,
		findByMultipleParam:   findByMultipleParamStmt,
		countByEmail:          countByEmailStmt,
		countByUsername:       countByUsernameStmt,
		countByPhoneNumber:    countByPhoneNumberStmt,
		countByEmailGoogleId:  countByEmailGoogleIdStmt,
		countByEmailNotGoogle: countByEmailNotGoogleStmt,
	}, nil
}

type UserRepository interface {
	Close() error

	CreateByPhoneNumber(ctx context.Context, tx Querier, user *entity.User) (*entity.User, error)
	CreateByGoogleSignIn(ctx context.Context, tx Querier, user *entity.User) (*entity.User, error)
	CreateByEmail(ctx context.Context, tx Querier, user *entity.User) (*entity.User, error)

	CountByEmail(ctx context.Context, email string) (int, error)
	CountByUsername(ctx context.Context, email string) (int, error)
	CountByPhoneNumber(ctx context.Context, email string) (int, error)
	CountByEmailGoogleId(ctx context.Context, email, googleId string) (int, error)
	CountByEmailNotGoogle(ctx context.Context, email string) (int, error)

	FindById(ctx context.Context, userId string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByMultipleParam(ctx context.Context, multipleParam string) (*entity.User, error)

	UpdateEmailVerifiedAt(ctx context.Context, tx Querier, user *entity.User) (*entity.User, error)
	UpdatePassword(ctx context.Context, tx Querier, user *entity.User) (*entity.User, error)
}

type userRepository struct {
	userPreparedStmt *userPreparedStmt
}

func NewUserRepository(db *sqlx.DB) (UserRepository, error) {

	userPreparedStmt, err := newUserPreparedStmt(db)
	if err != nil {
		log.Print("error initialize user statement : ", err)
		return nil, err
	}

	return &userRepository{
		userPreparedStmt: userPreparedStmt,
	}, nil
}

func (r *userRepository) Close() error {
	if err := r.userPreparedStmt.findById.Close(); err != nil {
		log.Print(err)
		return err
	}

	if err := r.userPreparedStmt.findByEmail.Close(); err != nil {
		log.Print(err)
		return err
	}

	if err := r.userPreparedStmt.countByEmail.Close(); err != nil {
		log.Print(err)
		return err
	}

	if err := r.userPreparedStmt.countByPhoneNumber.Close(); err != nil {
		log.Print(err)
		return err
	}

	if err := r.userPreparedStmt.countByUsername.Close(); err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (r *userRepository) CreateByPhoneNumber(ctx context.Context, tx Querier, user *entity.User) (*entity.User, error) {

	query := `INSERT INTO users 
	(id, username, password, phone_number, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := tx.ExecContext(ctx, query, user.Id, user.Username, user.Password.String, user.PhoneNumber.String, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return user, nil
}

func (r *userRepository) CreateByGoogleSignIn(ctx context.Context, tx Querier, user *entity.User) (*entity.User, error) {
	query := `INSERT INTO users 
	(id, email, username, google_id, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := tx.ExecContext(ctx, query, user.Id, user.Email.String, user.Username, user.GoogleId.String, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return user, nil
}

func (r *userRepository) CreateByEmail(ctx context.Context, tx Querier, user *entity.User) (*entity.User, error) {
	query := `INSERT INTO users 
	(id, username, email, password, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := tx.ExecContext(ctx, query, user.Id, user.Username, user.Email.String, user.Password.String, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return user, nil
}

func (r *userRepository) CountByEmail(ctx context.Context, email string) (int, error) {
	var total int

	row := r.userPreparedStmt.countByEmail.QueryRowxContext(ctx, email)
	if err := row.Scan(&total); err != nil {

		return 0, err
	}

	return total, nil
}

func (r *userRepository) CountByUsername(ctx context.Context, username string) (int, error) {
	var total int

	row := r.userPreparedStmt.countByUsername.QueryRowxContext(ctx, username)
	if err := row.Scan(&total); err != nil {

		return 0, err
	}

	return total, nil
}

func (r *userRepository) CountByPhoneNumber(ctx context.Context, phoneNumber string) (int, error) {
	var total int

	row := r.userPreparedStmt.countByPhoneNumber.QueryRowxContext(ctx, phoneNumber)
	if err := row.Scan(&total); err != nil {

		return 0, err
	}

	return total, nil
}

func (r *userRepository) CountByEmailGoogleId(ctx context.Context, email, googleId string) (int, error) {
	var total int

	row := r.userPreparedStmt.countByEmailGoogleId.QueryRowxContext(ctx, email, googleId)
	if err := row.Scan(&total); err != nil {

		return 0, err
	}

	return total, nil
}

func (r *userRepository) CountByEmailNotGoogle(ctx context.Context, email string) (int, error) {
	var total int

	row := r.userPreparedStmt.countByEmailNotGoogle.QueryRowxContext(ctx, email)
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (r *userRepository) FindById(ctx context.Context, userId string) (*entity.User, error) {
	user := new(entity.User)

	row := r.userPreparedStmt.findById.QueryRowxContext(ctx, userId)
	if err := row.StructScan(user); err != nil {

		return nil, err
	}

	return user, nil
}

// STRUCT SCAN
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := new(entity.User)

	row := r.userPreparedStmt.findByEmail.QueryRowxContext(ctx, email)
	if err := row.StructScan(user); err != nil {

		return nil, err
	}

	return user, nil
}

func (r *userRepository) FindByMultipleParam(ctx context.Context, multipleParam string) (*entity.User, error) {
	user := new(entity.User)

	row := r.userPreparedStmt.findByMultipleParam.QueryRowxContext(ctx, multipleParam)
	if err := row.StructScan(user); err != nil {

		return nil, err
	}

	return user, nil
}

func (r *userRepository) UpdateEmailVerifiedAt(ctx context.Context, tx Querier, user *entity.User) (*entity.User, error) {
	query := `UPDATE users set email_verified_at = $1, updated_at = $2 WHERE email = $3`

	_, err := tx.ExecContext(ctx, query, user.EmailVerifiedAt, user.UpdatedAt, user.Email.String)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return user, nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, tx Querier, user *entity.User) (*entity.User, error) {
	query := `UPDATE users set password = $1, updated_at = $2 WHERE email = $3`

	_, err := tx.ExecContext(ctx, query, user.Password.String, user.UpdatedAt, user.Email.String)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return user, nil
}
