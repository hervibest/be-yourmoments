package repository

import (
	"context"
	"fmt"
	"strings"
)

type UserSimilarRepository interface {
	UpdateUsersForPhoto(ctx context.Context, db Querier, photoId string, userIds []string) error
	GetSimilarPhotosByUser(ctx context.Context, db Querier, userId string) (*UserSimilarPhotosResponse, error)
	DeleteSimilarUsers(ctx context.Context, db Querier, photoId string) error
}

type userSimilarRepository struct {
}

func NewUserSimilarRepository() UserSimilarRepository {
	return &userSimilarRepository{}
}

func (r *userSimilarRepository) UpdateUsersForPhoto(ctx context.Context, db Querier, photoId string, userIds []string) error {
	// Jika userIds kosong, hapus semua user_similar dengan photoId ini
	if len(userIds) == 0 {
		if err := r.DeleteSimilarUsers(ctx, db, photoId); err != nil {
			return err
		}
	}

	// 1. Hapus user_id yang tidak ada dalam daftar baru
	queryDelete := `
		DELETE FROM user_similars 
		WHERE photo_id = $1 AND user_id NOT IN (` + generatePlaceholders(len(userIds), 2) + `)
	`
	args := append([]interface{}{photoId}, convertToInterface(userIds)...)

	_, err := db.Exec(ctx, queryDelete, args...)
	if err != nil {
		return fmt.Errorf("failed to delete outdated user_similars: %w", err)
	}

	// 2. Insert user_id baru jika belum ada
	queryInsert := `
		INSERT INTO user_similars (photo_id, user_id, created_at, updated_at) 
		VALUES ` + generateInsertValues(len(userIds)) + `
		ON CONFLICT DO NOTHING
	`
	args = append([]interface{}{photoId}, convertToInterface(userIds)...)

	_, err = db.Exec(ctx, queryInsert, args...)
	if err != nil {
		return fmt.Errorf("failed to insert new user_similars: %w", err)
	}

	return nil
}

func (r *userSimilarRepository) DeleteSimilarUsers(ctx context.Context, db Querier, photoId string) error {
	query := `DELETE FROM user_similars WHERE photo_id = $1`
	_, err := db.Exec(ctx, query, photoId)
	if err != nil {
		return fmt.Errorf("failed to delete user_similars: %w", err)
	}
	return nil
}

func (r *userSimilarRepository) DeleteSimilarExceptHerself(ctx context.Context, db Querier, photoId string) error {
	query := `DELETE FROM user_similars WHERE photo_id = $1 AND NOT owner_id`
	_, err := db.Exec(ctx, query, photoId)
	if err != nil {
		return fmt.Errorf("failed to delete user_similars: %w", err)
	}
	return nil
}

type UserSimilarPhotosResponse struct {
	UserID string         `json:"user_id"`
	Photos []PhotoPreview `json:"photos"`
}

type PhotoPreview struct {
	ID         string `json:"id"`
	PreviewUrl string `json:"preview_url"`
}

func (r *userSimilarRepository) GetSimilarPhotosByUser(ctx context.Context, db Querier, userId string) (*UserSimilarPhotosResponse, error) {
	query := `
		SELECT p.id, 
		       CASE 
		         WHEN p.owned_by_user_id IS NULL THEN p.preview_with_bounding_url 
		         WHEN p.owned_by_user_id = $1 THEN p.preview_url 
		         ELSE NULL 
		       END AS preview_url
		FROM photos p
		INNER JOIN user_similar_photos usp ON p.id = usp.photo_id
		WHERE usp.user_id = $1
	`

	rows, err := db.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get similar photos: %w", err)
	}
	defer rows.Close()

	var photos []PhotoPreview
	for rows.Next() {
		var photo PhotoPreview
		if err := rows.Scan(&photo.ID, &photo.PreviewUrl); err != nil {
			return nil, fmt.Errorf("failed to scan photo: %w", err)
		}
		photos = append(photos, photo)
	}

	return &UserSimilarPhotosResponse{
		UserID: userId,
		Photos: photos,
	}, nil
}

// Helper untuk membuat placeholder dinamis ($2, $3, $4, ...)
func generatePlaceholders(count, start int) string {
	placeholders := []string{}
	for i := 0; i < count; i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%d", start+i))
	}
	return strings.Join(placeholders, ", ")
}

// Helper untuk membuat values dinamis untuk batch INSERT
func generateInsertValues(count int) string {
	values := []string{}
	for i := 0; i < count; i++ {
		n := i + 2 // karena $1 untuk photo_id
		values = append(values, fmt.Sprintf("($1, $%d, NOW(), NOW())", n))
	}
	return strings.Join(values, ", ")
}

// Helper untuk konversi []string ke []interface{}
func convertToInterface(slice []string) []interface{} {
	interfaceSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}
