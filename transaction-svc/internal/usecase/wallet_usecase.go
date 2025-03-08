package usecase

type WalletUsecase interface{}
type walletUsecase struct{}

func NewWalletUsecase() WalletUsecase {
	return &walletUsecase{}
}
