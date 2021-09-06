package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
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
	l := s.l.With().Int("wid", wid).Str("fingerprint", m.Fingerprint).Logger()

	_, err := s.r.GetWallet(ctx, wid)
	if err != nil {
		l.Info().Err(err).Send()
		return nil, err
	}

	lt, err := s.r.GetLatestTransaction(ctx, wid)
	if err != nil && !errors.Is(err, ErrTransactionNotFound) {
		l.Info().Err(err).Send()
		return nil, err
	}

	tr := &Transaction{
		ID:          uuid.NewString(),
		Amount:      m.Amount,
		Description: m.Description,
		Labels:      m.Labels,
		Fingerprint: m.Fingerprint,
		Created:     time.Now().UTC().Truncate(time.Millisecond),
	}

	if lt == nil {
		tr.OldBalance = 0.
		tr.NewBalance = tr.Amount
		tr.RefNo = 1
	} else {
		tr.OldBalance = lt.NewBalance
		tr.NewBalance = lt.NewBalance + tr.Amount
		tr.RefNo = lt.RefNo + 1
	}

	if tr.NewBalance < 0. {
		return nil, ErrNotEnoughWalletBalance
	}

	err = s.r.CreateTransaction(ctx, wid, tr)
	if err != nil {
		l.Info().Err(err).Send()
		return nil, err
	}

	return tr, nil
}
