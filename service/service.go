package service

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

type Service interface {
	CreateWallet(ctx context.Context, m *WalletModel) (*Wallet, error)
	GetWallet(ctx context.Context, wid int) (*Wallet, error)
	GetWalletBalance(ctx context.Context, wid int) (float64, error)
	CreateTransaction(ctx context.Context, wid int, m *TransactionModel) (*Transaction, error)
}

type Repository interface {
	CreateWallet(ctx context.Context, w *Wallet) error
	GetWallet(ctx context.Context, wid int) (*Wallet, error)
	GetWalletBalance(ctx context.Context, wid int) (float64, error)
	CreateTransaction(ctx context.Context, wid int, t *Transaction) error
	GetLatestTransaction(ctx context.Context, wid int) (*Transaction, error)
}

type walletService struct {
	r Repository
	l zerolog.Logger
}

func NewWalletService(r Repository, logger zerolog.Logger) Service {
	return &walletService{r: r, l: logger}
}

func (s *walletService) CreateWallet(ctx context.Context, m *WalletModel) (*Wallet, error) {
	w := &Wallet{
		ExternalID: m.ExternalID,
		Labels:     m.Labels,
		Created:    time.Now().UTC().Truncate(time.Millisecond),
	}

	if err := s.r.CreateWallet(ctx, w); err != nil {
		return nil, err
	}

	return w, nil
}

func (s *walletService) GetWallet(ctx context.Context, wid int) (*Wallet, error) {
	w, err := s.r.GetWallet(ctx, wid)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *walletService) GetWalletBalance(ctx context.Context, wid int) (float64, error) {
	return s.r.GetWalletBalance(ctx, wid)
}

func (s *walletService) CreateTransaction(ctx context.Context, wid int, m *TransactionModel) (*Transaction, error) {
	w, err := s.r.GetWallet(ctx, wid)
	if err != nil {
		s.l.Info().Int("wid", wid).Err(err).Send()
		return nil, err
	}
	if w == nil {
		s.l.Error().Int("wid", wid).Msg("unexpected situation, no error no wallet")
		panic("unexpected situation")
	}

	lt, err := s.r.GetLatestTransaction(ctx, wid)

	return lt, nil
}
