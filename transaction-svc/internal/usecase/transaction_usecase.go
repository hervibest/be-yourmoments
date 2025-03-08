package usecase

import (
	"be-yourmoments/transaction-svc/internal/adapter"
	"be-yourmoments/transaction-svc/internal/repository"
)

type TransactionUsecase interface {
}

type transactionUsecase struct {
	transactionRepo repository.TransactionRepository
	photoAdapter    adapter.PhotoAdapter
}

func NewTransactionUsecase(transactionRepo repository.TransactionRepository, photoAdapter adapter.PhotoAdapter) TransactionUsecase {
	return &transactionUsecase{
		transactionRepo: transactionRepo,
		photoAdapter:    photoAdapter,
	}
}

func (u *transactionUsecase) BuyPhoto() {

}
