package service

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
)

type Service interface {
	CreateWallet(ctx context.Context, m *WalletModel) (*Wallet, error)
	GetWallet(ctx context.Context, wid int) (*Wallet, error)
	GetWalletBalance(ctx context.Context, wid string) (float64, error)
}

type Repository interface {
	CreateWallet(ctx context.Context, w *Wallet) error
	GetWallet(ctx context.Context, wid int) (*Wallet, error)
	GetWalletBalance(ctx context.Context, wid int) (float64, error)
	CreateTransaction(ctx context.Context, wid int, t *Transaction) error
}

var (
	ErrWalletNotFound                        = errors.New("wallet not found")
	ErrWalletAlreadyExists                   = errors.New("wallet already exists with same external id")
	ErrTransactionFailedButRetriable         = errors.New("transaction failed due to consistency but retriable")
	ErrDatabaseConnectionFailed              = errors.New("could not connect to the database")
	ErrTransactionAlreadyExistsByRefNo       = errors.New("a transaction already exists with same refno")
	ErrTransactionAlreadyExistsByFingerprint = errors.New("a transaction already exists with same fingerprint")
)

type walletService struct {
	r Repository
	l zerolog.Logger
}

func NewWalletService(r Repository, logger zerolog.Logger) Service {
	return &walletService{r: r, l: logger}
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
