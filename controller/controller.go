package controller

import (
	"context"
	"time"
)

type (
	Wallet struct {
		ID         int32
		Labels     map[string]string
		ExternalID string
		Created    time.Time
	}

	WalletModel struct {
		Labels     map[string]string
		ExternalID string
	}

	TransactionModel struct {
		Amount      float64
		Description string
		Labels      map[string]string
		Fingerprint string
	}

	Transaction struct {
		ID          string
		RefNo       int32
		Amount      float64
		Description string
		Labels      map[string]string
		Fingerprint string
		Created     time.Time
		OldBalance  float64
		NewBalance  float64
	}
)

type Controller interface {
	CreateWallet(ctx context.Context, w *WalletModel) error
	GetWallet(ctx context.Context, wid int32) (*Wallet, error)
	GetWalletBalance(ctx context.Context, wid string) (float64, error)
}

type Repository interface {
	CreateWallet(ctx context.Context, w *Wallet) error
	GetWallet(ctx context.Context, wid int32) (*Wallet, error)
	GetWalletBalance(ctx context.Context, wid string) (float64, error)
}

type wc struct {
	r *Repository
}

func NewWalletController(r *Repository) Controller {
	return &wc{r: r}
}

func (c *wc) CreateWallet(ctx context.Context, w *WalletModel) error {
	return nil
}

func (c *wc) GetWallet(ctx context.Context, wid int32) (*Wallet, error) {
	return &Wallet{}, nil
}

func (c *wc) GetWalletBalance(ctx context.Context, wid string) (float64, error) {
	return 0., nil
}
