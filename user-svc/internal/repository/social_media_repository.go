package repository

import (
	"be-yourmoments/user-svc/internal/entity"
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type socialMediaPreparedStmt struct {
	findAll     *sqlx.Stmt
	findById    *sqlx.Stmt
	findByName  *sqlx.Stmt
	countByName *sqlx.Stmt
}

func newSocialMediaStmt(db *sqlx.DB) (*socialMediaPreparedStmt, error) {
	findAllStmt, err := db.Preparex("SELECT * FROM social_medias")
	if err != nil {
		return nil, err
	}

	findByIdStmt, err := db.Preparex("SELECT * FROM social_medias WHERE id=$1")
	if err != nil {
		return nil, err
	}

	findByNameStmt, err := db.Preparex("SELECT * FROM social_medias WHERE name=$1")
	if err != nil {
		return nil, err
	}

	countByNameStmt, err := db.Preparex("SELECT COUNT(*) FROM social_medias WHERE name=$1")
	if err != nil {
		return nil, err
	}

	return &socialMediaPreparedStmt{
		findAll:     findAllStmt,
		findById:    findByIdStmt,
		findByName:  findByNameStmt,
		countByName: countByNameStmt,
	}, nil
}

type SocialMediaRepository interface {
	Insert(ctx context.Context, tx Querier, socialMedia *entity.SocialMedia) (*entity.SocialMedia, error)
	Update(ctx context.Context, tx Querier, socialMedia *entity.SocialMedia) (*entity.SocialMedia, error)
	FindAll(ctx context.Context) (*[]*entity.SocialMedia, error)
	FindByName(ctx context.Context, name string) (*entity.SocialMedia, error)
	Delete(ctx context.Context, tx Querier, name string) error
}

type socialMediaRepository struct {
	socialMediaStmt *socialMediaPreparedStmt
}

func NewSocialMediaRepository(db *sqlx.DB) (SocialMediaRepository, error) {
	socialMediaStmt, err := newSocialMediaStmt(db)
	if err != nil {
		log.Print("error initialize social media  statement : ", err)
		return nil, err
	}
	return &socialMediaRepository{
		socialMediaStmt: socialMediaStmt,
	}, nil
}

func (r *socialMediaRepository) Insert(ctx context.Context, tx Querier, socialMedia *entity.SocialMedia) (*entity.SocialMedia, error) {
	query := ` INSERT INTO social_medias  (id, name, base_url, logo_url, description, is_active, created_at, updated_at)  
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) `

	_, err := tx.ExecContext(ctx, query, socialMedia.Id, socialMedia.Name, socialMedia.BaseUrl, socialMedia.LogoUrl,
		socialMedia.Description, socialMedia.IsActive, socialMedia.CreatedAt, socialMedia.UpdatedAt)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to insert social media: %w", err)
	}

	return socialMedia, nil
}

func (r *socialMediaRepository) Update(ctx context.Context, tx Querier, socialMedia *entity.SocialMedia) (*entity.SocialMedia, error) {
	query := `UPDATE social_medias  SET base_url = $1, logo_url = $2, description = $3, is_active = $4, updated_at = $5 WHERE id = $6`

	_, err := tx.ExecContext(ctx, query, socialMedia.BaseUrl, socialMedia.LogoUrl, socialMedia.Description, socialMedia.IsActive,
		socialMedia.UpdatedAt, socialMedia.Id)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to update social media: %w", err)
	}

	return socialMedia, nil
}

func (r *socialMediaRepository) FindAll(ctx context.Context) (*[]*entity.SocialMedia, error) {
	socialMedias := make([]*entity.SocialMedia, 0)

	rows, err := r.socialMediaStmt.findAll.QueryxContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		socialMedia := new(entity.SocialMedia)
		if err := rows.StructScan(socialMedia); err != nil {
			return nil, err
		}
		socialMedias = append(socialMedias, socialMedia)
	}

	return &socialMedias, nil
}

func (r *socialMediaRepository) FindByName(ctx context.Context, name string) (*entity.SocialMedia, error) {
	socialMedia := new(entity.SocialMedia)

	row := r.socialMediaStmt.findByName.QueryRowxContext(ctx, name)
	if err := row.StructScan(socialMedia); err != nil {
		return nil, err
	}

	return socialMedia, nil
}

func (r *socialMediaRepository) Delete(ctx context.Context, tx Querier, name string) error {
	query := `DELETE FROM social_medias WHERE name = $1`
	_, err := tx.ExecContext(ctx, query, name)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to delete socialMedia: %w", err)
	}

	return nil
}
