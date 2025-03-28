package repository

import (
	"be-yourmoments/photo-svc/internal/entity"
	"fmt"
	"log"
	"strings"
	"time"
)

type UserSimilarRepository interface {
	InsertOrUpdate(tx Querier, photoId string, userSimilarPhotos *[]*entity.UserSimilarPhoto) error
	// UpdateUsersForPhoto(ctx context.Context, db Querier, photoId string, userIds []string) error
	// GetSimilarPhotosByUser(ctx context.Context, db Querier, userId string) (*UserSimilarPhotosResponse, error)
	// DeleteSimilarUsers(ctx context.Context, db Querier, photoId string) error
}

type userSimilarRepository struct {
}

func NewUserSimilarRepository() UserSimilarRepository {
	return &userSimilarRepository{}
}

func (r *userSimilarRepository) InsertOrUpdate(tx Querier, photoId string, userSimilarPhotos *[]*entity.UserSimilarPhoto) error {
	now := time.Now()
	if len(*userSimilarPhotos) == 0 {
		log.Println("ALL DELETED BECAUSE OF ZERO")
		if _, err := tx.Exec("DELETE FROM user_similar_photos WHERE photo_id = $1", photoId); err != nil {
			return err
		}
		return nil
	}

	placeholders := make([]string, len(*userSimilarPhotos))
	deleteArgs := make([]interface{}, 0, len(*userSimilarPhotos)+1)
	deleteArgs = append(deleteArgs, photoId)
	for i, userSimilarPhoto := range *userSimilarPhotos {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		deleteArgs = append(deleteArgs, userSimilarPhoto.UserId)
	}

	deleteQuery := "DELETE FROM user_similar_photos WHERE photo_id = $1 AND user_id NOT IN (" + strings.Join(placeholders, ", ") + ")"
	if _, err := tx.Exec(deleteQuery, deleteArgs...); err != nil {
		log.Println("Error at delete query:", err)
		return err
	}

	insertValues := make([]string, 0, len(*userSimilarPhotos))
	insertArgs := make([]interface{}, 0, len(*userSimilarPhotos)*4)
	placeholderCounter := 1
	for _, userSimilarPhoto := range *userSimilarPhotos {
		// Misalnya, baris pertama: ($1, $2, $3, $4), baris kedua: ($5, $6, $7, $8), dst.
		insertValues = append(insertValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", placeholderCounter, placeholderCounter+1, placeholderCounter+2, placeholderCounter+3, placeholderCounter+4))
		insertArgs = append(insertArgs, photoId, userSimilarPhoto.UserId, userSimilarPhoto.Similarity, now, now)
		placeholderCounter += 5
	}

	insertQuery := "INSERT INTO user_similar_photos (photo_id, user_id, similarity, created_at, updated_at) VALUES " +
		strings.Join(insertValues, ", ") +
		" ON CONFLICT (photo_id, user_id) DO UPDATE SET updated_at = EXCLUDED.updated_at"

	if _, err := tx.Exec(insertQuery, insertArgs...); err != nil {
		log.Println("Error at insert query:", err)
		return err
	}

	return nil
}

// func (r *userSimilarRepository) UpdateUsersForPhoto(ctx context.Context, db Querier, photoId string, userIds []string) error {
// 	// Jika userIds kosong, hapus semua user_similar dengan photoId ini
// 	if len(userIds) == 0 {
// 		if err := r.DeleteSimilarUsers(ctx, db, photoId); err != nil {
// 			return err
// 		}
// 	}

// 	// 1. Hapus user_id yang tidak ada dalam daftar baru
// 	queryDelete := `
// 		DELETE FROM user_similars
// 		WHERE photo_id = $1 AND user_id NOT IN (` + generatePlaceholders(len(userIds), 2) + `)
// 	`
// 	args := append([]interface{}{photoId}, convertToInterface(userIds)...)

// 	_, err := db.Exec(ctx, queryDelete, args...)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete outdated user_similars: %w", err)
// 	}

// 	// 2. Insert user_id baru jika belum ada
// 	queryInsert := `
// 		INSERT INTO user_similars (photo_id, user_id, created_at, updated_at)
// 		VALUES ` + generateInsertValues(len(userIds)) + `
// 		ON CONFLICT DO NOTHING
// 	`
// 	args = append([]interface{}{photoId}, convertToInterface(userIds)...)

// 	_, err = db.Exec(ctx, queryInsert, args...)
// 	if err != nil {
// 		return fmt.Errorf("failed to insert new user_similars: %w", err)
// 	}

// 	return nil
// }

// func (r *userSimilarRepository) DeleteSimilarUsers(ctx context.Context, db Querier, photoId string) error {
// 	query := `DELETE FROM user_similar_photos WHERE photo_id = ?`
// 	_, err := db.ExecContext(ctx, query, photoId)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete user_similars: %w", err)
// 	}
// 	return nil
// }

// func (r *userSimilarRepository) DeleteSimilarExceptHerself(ctx context.Context, db Querier, photoId string) error {
// 	query := `DELETE FROM user_similars WHERE photo_id = $1 AND NOT owner_id`
// 	_, err := db.Exec(ctx, query, photoId)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete user_similars: %w", err)
// 	}
// 	return nil
// }

// type UserSimilarPhotosResponse struct {
// 	UserID string         `json:"user_id"`
// 	Photos []PhotoPreview `json:"photos"`
// }

// type PhotoPreview struct {
// 	ID         string `json:"id"`
// 	PreviewUrl string `json:"preview_url"`
// }

// func (r *userSimilarRepository) GetSimilarPhotosByUser(ctx context.Context, db Querier, userId string) (*UserSimilarPhotosResponse, error) {
// 	/*  !!! TODO !!!
// 	Query salah, pastikan logic wishlish dihandle

// 	*/

// 	query := `
// 		SELECT p.id,
// 		       CASE
// 		         WHEN p.owned_by_user_id IS NULL THEN p.preview_with_bounding_url
// 		         WHEN p.owned_by_user_id = $1 THEN p.preview_url
// 		         ELSE NULL
// 		       END AS preview_url
// 		FROM photos p
// 		INNER JOIN user_similar_photos usp ON p.id = usp.photo_id
// 		WHERE usp.user_id = $1
// 	`

// 	rows, err := db.Query(ctx, query, userId)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get similar photos: %w", err)
// 	}
// 	defer rows.Close()

// 	var photos []PhotoPreview
// 	for rows.Next() {
// 		var photo PhotoPreview
// 		if err := rows.Scan(&photo.ID, &photo.PreviewUrl); err != nil {
// 			return nil, fmt.Errorf("failed to scan photo: %w", err)
// 		}
// 		photos = append(photos, photo)
// 	}

// 	return &UserSimilarPhotosResponse{
// 		UserID: userId,
// 		Photos: photos,
// 	}, nil
// }
