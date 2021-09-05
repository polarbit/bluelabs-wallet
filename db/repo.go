package db

import (
	"context"

	"github.com/polarbit/bluelabs-wallet/config"
	"github.com/polarbit/bluelabs-wallet/wallet"
)

type repository struct {
	config *config.AppConfig
}

func NewRepository(c *config.AppConfig) wallet.WalletRepository {
	return &repository{config: c}
}

func (r *repository) CreateWallet(ctx context.Context, w *wallet.WalletModel) error {
	return nil
}

func (r *repository) GetWallet(ctx context.Context, wid string) (*wallet.Wallet, error) {
	return nil, nil
}

func (r *repository) GetWalletBalance(ctx context.Context, wid string) (float64, error) {
	return 0, nil
}

func (r *repository) SetWalletBalance(ctx context.Context, wid string, amount float64) error {
	return nil
}
