package repository

import (
	"be-yourmoments/user-svc/internal/entity"
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type resetPasswordPreparedStmt struct {
	findByEmail         *sqlx.Stmt
	findByEmailAndToken *sqlx.Stmt
	countByEmail        *sqlx.Stmt
}

func newResetPasswordStmt(db *sqlx.DB) (*resetPasswordPreparedStmt, error) {
	findByEmailStmt, err := db.Preparex("SELECT * FROM reset_passwords WHERE email = $1")
	if err != nil {
		return nil, err
	}

	findByEmailAndTokenStmt, err := db.Preparex("SELECT * FROM reset_passwords WHERE email = $1 AND token = $2")
	if err != nil {
		return nil, err
	}

	countByEmailStmt, err := db.Preparex("SELECT COUNT(*) FROM reset_passwords WHERE email = $1")
	if err != nil {
		return nil, err
	}

	return &resetPasswordPreparedStmt{
		findByEmail:         findByEmailStmt,
		findByEmailAndToken: findByEmailAndTokenStmt,
		countByEmail:        countByEmailStmt,
	}, nil
}

type ResetPasswordRepository interface {
	Insert(ctx context.Context, tx Querier, resetPassword *entity.ResetPassword) (*entity.ResetPassword, error)
	Update(ctx context.Context, tx Querier, resetPassword *entity.ResetPassword) (*entity.ResetPassword, error)
	FindByEmail(ctx context.Context, email string) (*entity.ResetPassword, error)
	FindByEmailAndToken(ctx context.Context, email, token string) (*entity.ResetPassword, error)
	Delete(ctx context.Context, tx Querier, resetPassword *entity.ResetPassword) error

	CountByEmail(ctx context.Context, email string) (int, error)
}

type resetPasswordRepository struct {
	resetPasswordStmt *resetPasswordPreparedStmt
}

func NewResetPasswordRepository(db *sqlx.DB) (ResetPasswordRepository, error) {
	resetPasswordStmt, err := newResetPasswordStmt(db)
	if err != nil {
		log.Print("error initialize rest password statement : ", err)
		return nil, err
	}
	return &resetPasswordRepository{
		resetPasswordStmt: resetPasswordStmt,
	}, nil
}

func (r *resetPasswordRepository) Insert(ctx context.Context, tx Querier, resetPassword *entity.ResetPassword) (*entity.ResetPassword, error) {
	query := ` INSERT INTO reset_passwords (email, token, created_at, updated_at)  
	VALUES ($1, $2, $3, $4) `

	_, err := tx.ExecContext(ctx, query, resetPassword.Email, resetPassword.Token, resetPassword.CreatedAt, resetPassword.UpdatedAt)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert resetPassword: %w", err)
	}

	return resetPassword, nil
}

func (r *resetPasswordRepository) Update(ctx context.Context, tx Querier, resetPassword *entity.ResetPassword) (*entity.ResetPassword, error) {
	query := `UPDATE reset_passwords set token = $1, updated_at = $2 WHERE email = $3`
	_, err := tx.ExecContext(ctx, query, resetPassword.Token, resetPassword.UpdatedAt, resetPassword.Email)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert resetPassword: %w", err)
	}

	return resetPassword, nil
}

func (r *resetPasswordRepository) FindByEmail(ctx context.Context, email string) (*entity.ResetPassword, error) {
	resetPassword := new(entity.ResetPassword)

	row := r.resetPasswordStmt.findByEmail.QueryRowxContext(ctx, email)
	if err := row.StructScan(resetPassword); err != nil {

		return nil, err
	}

	return resetPassword, nil
}

func (r *resetPasswordRepository) FindByEmailAndToken(ctx context.Context, email, token string) (*entity.ResetPassword, error) {
	resetPassword := new(entity.ResetPassword)

	row := r.resetPasswordStmt.findByEmailAndToken.QueryRowxContext(ctx, email, token)
	if err := row.StructScan(resetPassword); err != nil {

		return nil, err
	}

	return resetPassword, nil
}

func (r *resetPasswordRepository) Delete(ctx context.Context, tx Querier, resetPassword *entity.ResetPassword) error {
	query := `DELETE FROM reset_passwords WHERE email = $1 AND token = $2`
	_, err := tx.ExecContext(ctx, query, resetPassword.Email, resetPassword.Token)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to delete resetPassword: %w", err)
	}

	return nil
}

func (r *resetPasswordRepository) CountByEmail(ctx context.Context, email string) (int, error) {
	var total int

	row := r.resetPasswordStmt.countByEmail.QueryRowxContext(ctx, email)
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}
