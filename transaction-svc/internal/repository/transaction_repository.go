package repository

import (
	"be-yourmoments/transaction-svc/internal/entity"
	"context"
	"fmt"
)

type TransactionRepository interface {
	Create(ctx context.Context, db Querier, transaction *entity.Transaction) error
	// UpdateProcessedUrl(ctx context.Context, db Querier, transaction *entity.Transaction) error
	// UpdateClaimedPhoto(ctx context.Context, db Querier, transaction *entity.Transaction) error
	// UpdatePhotoStatus(ctx context.Context, db Querier, transaction *entity.Transaction) error
}

type transactionRepository struct {
}

func NewTransactionpository() TransactionRepository {
	return &transactionRepository{}
}

func (r *transactionRepository) Create(ctx context.Context, db Querier, transaction *entity.Transaction) error {

	query := `INSERT INTO transactions 
			  (id, user_id, photo_id, snap_token, external_callback_response, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := db.QueryRow(ctx, query, transaction.Id, transaction.UserId, transaction.PhotoId,
		transaction.SnapToken, transaction.ExternalCallbackResponse, transaction.PaidAt,
		transaction.CreatedAt, transaction.UpdatedAt).Scan(&transaction.Id)

	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	return nil
}

// func (r *transactionRepository) UpdateProcessedUrl(ctx context.Context, db Querier, transaction *entity.Transaction) error {
// 	query := `UPDATE photos
// 	          SET preview_url = $1, preview_with_bounding_url = $2, updated_at = $3
// 	          WHERE id = $4`

// 	_, err := db.Exec(ctx, query, transaction.PreviewUrl, transaction.PreviewWithBoundingUrl, transaction.UpdatedAt, transaction.Id)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *transactionRepository) UpdateClaimedPhoto(ctx context.Context, db Querier, transaction *entity.Transaction) error {
// 	query := `UPDATE photos
// 	          SET owned_by_user_id = $1, status = $2, updated_at = $3
// 	          WHERE id = $4`

// 	_, err := db.Exec(ctx, query, transaction.OwnedByUserId, transaction.Status, transaction.UpdatedAt, transaction.Id)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *transactionRepository) UpdatePhotoStatus(ctx context.Context, db Querier, transaction *entity.Transaction) error {
// 	query := `UPDATE photos
// 	          SET status = $1, updated_at = $3
// 	          WHERE id = $4`

// 	_, err := db.Exec(ctx, query, transaction.Status, transaction.UpdatedAt, transaction.Id)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
