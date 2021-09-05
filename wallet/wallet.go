package wallet

import (
	"context"
	"time"
)

type Wallet struct {
	ID         int
	Labels     map[string]string
	ExternalID string
	Created    time.Time
}

type WalletModel struct {
	Labels     map[string]string
	ExternalID string
}

type TransactionModel struct {
	Amount      float64
	Description string
	Labels      map[string]string
	Fingerprint string
}

type Transaction struct {
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

type WalletController interface {
	CreateWallet(ctx context.Context, w *WalletModel) error
	GetWallet(ctx context.Context, wid string) (*Wallet, error)
	GetWalletBalance(ctx context.Context, wid string) (float64, error)
}

type WalletRepository interface {
	CreateWallet(ctx context.Context, w *WalletModel) error
	GetWallet(ctx context.Context, wid string) (*Wallet, error)
	GetWalletBalance(ctx context.Context, wid string) (float64, error)
	SetWalletBalance(ctx context.Context, wid string, amount float64) error
}

type wc struct {
	r *WalletRepository
}

func NewWalletController(r *WalletRepository) WalletController {
	return &wc{r: r}
}

func (c *wc) CreateWallet(ctx context.Context, w *WalletModel) error {
	return nil
}

func (c *wc) GetWallet(ctx context.Context, wid string) (*Wallet, error) {
	return &Wallet{}, nil
}

func (c *wc) GetWalletBalance(ctx context.Context, wid string) (float64, error) {
	return 0., nil
}
