package repository

import (
	"be-yourmoments/user-svc/internal/entity"
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type emailVerifPreparedStmt struct {
	findByEmail         *sqlx.Stmt
	findByEmailAndToken *sqlx.Stmt
}

func newEmailVerifStmt(db *sqlx.DB) (*emailVerifPreparedStmt, error) {
	findByEmailStmt, err := db.Preparex("SELECT * FROM email_verifications WHERE email = $1")
	if err != nil {
		return nil, err
	}

	findByEmailAndTokenStmt, err := db.Preparex("SELECT * FROM email_verifications WHERE email = $1 AND token = $2")
	if err != nil {
		return nil, err
	}

	return &emailVerifPreparedStmt{
		findByEmail:         findByEmailStmt,
		findByEmailAndToken: findByEmailAndTokenStmt,
	}, nil
}

type EmailVerificationRepository interface {
	Insert(ctx context.Context, tx Querier, emailVerification *entity.EmailVerification) (*entity.EmailVerification, error)
	Update(ctx context.Context, tx Querier, emailVerification *entity.EmailVerification) (*entity.EmailVerification, error)
	FindByEmail(ctx context.Context, email string) (*entity.EmailVerification, error)
	FindByEmailAndToken(ctx context.Context, email, token string) (*entity.EmailVerification, error)
	Delete(ctx context.Context, tx Querier, emailVerification *entity.EmailVerification) error
}

type emailVerificationRepository struct {
	emailVerifStmt *emailVerifPreparedStmt
}

func NewEmailVerificationRepository(db *sqlx.DB) (EmailVerificationRepository, error) {
	emailVerifStmt, err := newEmailVerifStmt(db)
	if err != nil {
		log.Print("error initialize email verification statement : ", err)
		return nil, err
	}
	return &emailVerificationRepository{
		emailVerifStmt: emailVerifStmt,
	}, nil
}

func (r *emailVerificationRepository) Insert(ctx context.Context, tx Querier, emailVerification *entity.EmailVerification) (*entity.EmailVerification, error) {
	query := ` INSERT INTO email_verifications  (email, token, created_at, updated_at)  
	VALUES ($1, $2, $3, $4) `

	_, err := tx.ExecContext(ctx, query, emailVerification.Email, emailVerification.Token, emailVerification.CreatedAt, emailVerification.UpdatedAt)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert emailVerification: %w", err)
	}

	return emailVerification, nil
}

func (r *emailVerificationRepository) Update(ctx context.Context, tx Querier, emailVerification *entity.EmailVerification) (*entity.EmailVerification, error) {
	query := `UPDATE email_verifications  set token = $1, updated_at = $2 WHERE email = $3`
	_, err := tx.ExecContext(ctx, query, emailVerification.Token, emailVerification.UpdatedAt, emailVerification.Email)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert emailVerification: %w", err)
	}

	return emailVerification, nil
}

func (r *emailVerificationRepository) FindByEmail(ctx context.Context, email string) (*entity.EmailVerification, error) {
	emailVerification := new(entity.EmailVerification)

	row := r.emailVerifStmt.findByEmail.QueryRowxContext(ctx, email)
	if err := row.StructScan(emailVerification); err != nil {

		return nil, err
	}

	return emailVerification, nil
}

func (r *emailVerificationRepository) FindByEmailAndToken(ctx context.Context, email, token string) (*entity.EmailVerification, error) {
	emailVerification := new(entity.EmailVerification)

	row := r.emailVerifStmt.findByEmailAndToken.QueryRowxContext(ctx, email, token)
	if err := row.StructScan(emailVerification); err != nil {

		return nil, err
	}

	return emailVerification, nil
}

func (r *emailVerificationRepository) Delete(ctx context.Context, tx Querier, emailVerification *entity.EmailVerification) error {
	query := `DELETE FROM email_verifications WHERE email = $1 AND token = $2`
	_, err := tx.ExecContext(ctx, query, emailVerification.Email, emailVerification.Token)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to delete emailVerification: %w", err)
	}

	return nil
}
