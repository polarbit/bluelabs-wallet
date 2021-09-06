package service

import (
	"context"
	"errors"
	"time"
)

type Service interface {
	CreateWallet(ctx context.Context, m *WalletModel) (*Wallet, error)
	GetWallet(ctx context.Context, wid int) (*Wallet, error)
	GetWalletBalance(ctx context.Context, wid string) (float64, error)
}

type Repository interface {
	CreateWallet(ctx context.Context, w *Wallet) error
	GetWallet(ctx context.Context, wid int) (*Wallet, error)
	GetWalletBalance(ctx context.Context, wid string) (float64, error)
}

var (
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrWalletAlreadyExists = errors.New("Wallet already exists with same external id")
)

type walletService struct {
	r Repository
}

func NewWalletService(r Repository) Service {
	return &walletService{r: r}
}

func (c *walletService) CreateWallet(ctx context.Context, m *WalletModel) (*Wallet, error) {
	w := &Wallet{
		ExternalID: m.ExternalID,
		Labels:     m.Labels,
		Created:    time.Now().UTC().Truncate(time.Millisecond),
	}

	if err := c.r.CreateWallet(ctx, w); err != nil {
		return nil, err
	}

	return w, nil
}

func (c *walletService) GetWallet(ctx context.Context, wid int) (*Wallet, error) {
	w, err := c.r.GetWallet(ctx, wid)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (c *walletService) GetWalletBalance(ctx context.Context, wid string) (float64, error) {
	return 0., nil
}
