//go:build !integration
// +build !integration

package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	exitcode := m.Run()
	os.Exit(exitcode)
}

func TestGetWallet(t *testing.T) {
	t.Run("WalletFound", func(t *testing.T) {
		var mok = &mockRepository{}
		mok.On("GetWallet", mock.Anything, 10).Return(&Wallet{ID: 10}, nil)
		svc := NewWalletService(mok, log.Logger)

		w, err := svc.GetWallet(context.Background(), 10)
		assert.NoError(t, err)
		assert.NotNil(t, w)
		assert.Equal(t, 10, w.ID)
	})

	t.Run("WalletNotFound", func(t *testing.T) {
		var mok = &mockRepository{}
		mok.On("GetWallet", mock.Anything, 11).Return((*Wallet)(nil), ErrWalletNotFound)
		svc := NewWalletService(mok, log.Logger)

		w, err := svc.GetWallet(context.Background(), 11)
		assert.Nil(t, w)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrWalletNotFound)
	})
}

func TestCreateWallet(t *testing.T) {
	t.Run("WalletCreated", func(t *testing.T) {
		var mok = &mockRepository{}
		mok.On("CreateWallet", mock.Anything,
			mock.MatchedBy(func(w *Wallet) bool { w.ID = 99; return w.ExternalID == "99" })).
			Return(nil)
		svc := NewWalletService(mok, log.Logger)

		w, err := svc.CreateWallet(context.Background(), &WalletModel{
			ExternalID: "99",
		})
		assert.NoError(t, err)
		assert.NotNil(t, w)
		assert.Equal(t, "99", w.ExternalID)
		assert.Equal(t, 99, w.ID)
	})

	t.Run("WalletAlreadyExists", func(t *testing.T) {
		var mok = &mockRepository{}
		mok.On("CreateWallet", mock.Anything,
			mock.MatchedBy(func(w *Wallet) bool { return w.ExternalID == "100" })).
			Return(ErrWalletAlreadyExists)
		svc := NewWalletService(mok, log.Logger)

		w, err := svc.CreateWallet(context.Background(), &WalletModel{
			ExternalID: "100",
		})
		assert.ErrorIs(t, err, ErrWalletAlreadyExists)
		assert.Nil(t, w)
	})
}

func TestGetWalletBalance(t *testing.T) {
	t.Run("WalletNotFound", func(t *testing.T) {
		var mok = &mockRepository{}
		mok.On("GetWalletBalance", mock.Anything, mock.Anything).Return(0., ErrWalletNotFound)
		svc := NewWalletService(mok, log.Logger)

		b, err := svc.GetWalletBalance(context.Background(), 10)
		assert.ErrorIs(t, err, ErrWalletNotFound)
		assert.Equal(t, 0., b)
	})

	t.Run("ReturnSomeBalance", func(t *testing.T) {
		var mok = &mockRepository{}
		mok.On("GetWalletBalance", mock.Anything, mock.Anything).Return(99., nil)
		svc := NewWalletService(mok, log.Logger)

		b, err := svc.GetWalletBalance(context.Background(), 10)
		assert.NoError(t, err)
		assert.Equal(t, 99., b)
	})
}

func TestCreateTransaction(t *testing.T) {
	t.Run("WalletNotExists", func(t *testing.T) {
		var mok = &mockRepository{}
		mok.On("GetWallet", mock.Anything, 10).Return((*Wallet)(nil), ErrWalletNotFound)
		svc := NewWalletService(mok, log.Logger)

		_, err := svc.CreateTransaction(context.Background(), 10, &TransactionModel{})
		assert.ErrorIs(t, err, ErrWalletNotFound)
	})

	t.Run("FirstTransaction", func(t *testing.T) {
		var mok = &mockRepository{}
		mok.On("GetWallet", mock.Anything, 10).Return(&Wallet{}, nil)
		mok.On("GetLatestTransaction", mock.Anything, 10).Return((*Transaction)(nil), nil)
		mok.On("CreateTransaction", mock.Anything, 10, mock.Anything).Return(nil)
		svc := NewWalletService(mok, log.Logger)

		m := &TransactionModel{Amount: 25, Description: "desc",
			Fingerprint: uuid.NewString(), Labels: map[string]string{}}

		tr, err := svc.CreateTransaction(context.Background(), 10, m)
		assert.NoError(t, err)
		assert.NotNil(t, tr)
		assert.NotEmpty(t, tr.ID)
		assert.Equal(t, 1, tr.RefNo)
		assert.Equal(t, m.Amount, tr.Amount)
		assert.Equal(t, 0., tr.OldBalance)
		assert.Equal(t, 25., tr.NewBalance)
		assert.Equal(t, m.Fingerprint, tr.Fingerprint)
		assert.Equal(t, m.Description, tr.Description)
		assert.Equal(t, m.Labels, tr.Labels)
		assert.Greater(t, tr.Created.UnixMilli(), time.Now().Add(-3*time.Second).UnixMilli())
	})

	t.Run("SecondTransaction", func(t *testing.T) {
		var mok = &mockRepository{}
		mok.On("GetWallet", mock.Anything, 10).Return(&Wallet{}, nil)
		mok.On("GetLatestTransaction", mock.Anything, mock.Anything).Return(
			&Transaction{RefNo: 1, NewBalance: 99}, nil)
		mok.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		svc := NewWalletService(mok, log.Logger)

		m := &TransactionModel{Amount: -25.}

		tr, err := svc.CreateTransaction(context.Background(), 10, m)
		assert.NoError(t, err)
		assert.NotNil(t, tr)
		assert.NotEmpty(t, tr.ID)
		assert.Equal(t, 2, tr.RefNo)
		assert.Equal(t, m.Amount, tr.Amount)
		assert.Equal(t, 99., tr.OldBalance)
		assert.Equal(t, 99.+m.Amount, tr.NewBalance)
	})

	t.Run("NotEnoughBalance", func(t *testing.T) {
		var mok = &mockRepository{}
		mok.On("GetWallet", mock.Anything, 10).Return(&Wallet{}, nil)
		mok.On("GetLatestTransaction", mock.Anything, mock.Anything).Return(
			&Transaction{RefNo: 1, NewBalance: 1}, nil)
		mok.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		svc := NewWalletService(mok, log.Logger)

		m := &TransactionModel{Amount: -2}

		tr, err := svc.CreateTransaction(context.Background(), 10, m)
		assert.ErrorIs(t, err, ErrNotEnoughWalletBalance)
		assert.Nil(t, tr)
	})

	t.Run("ConsistencyError", func(t *testing.T) {
		var mok = &mockRepository{}
		mok.On("GetWallet", mock.Anything, 10).Return(&Wallet{}, nil)
		mok.On("GetLatestTransaction", mock.Anything, mock.Anything).Return(
			&Transaction{RefNo: 1, NewBalance: 99.}, nil)
		mok.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).
			Return(ErrTransactionConsistency)
		svc := NewWalletService(mok, log.Logger)

		m := &TransactionModel{Amount: -2}

		tr, err := svc.CreateTransaction(context.Background(), 10, m)
		assert.ErrorIs(t, err, ErrTransactionConsistency)
		assert.Nil(t, tr)
	})
}

func TestGetLatestTransaction(t *testing.T) {
	t.Run("TransactionNotFound", func(t *testing.T) {
		var mok = &mockRepository{}
		mok.On("GetLatestTransaction", mock.Anything, mock.Anything).Return((*Transaction)(nil), ErrTransactionNotFound)
		svc := NewWalletService(mok, log.Logger)

		_, err := svc.GetLatestTransaction(context.Background(), 10)
		assert.ErrorIs(t, err, ErrTransactionNotFound)
	})

	t.Run("FoundTransaction", func(t *testing.T) {
		var mok = &mockRepository{}
		mok.On("GetLatestTransaction", mock.Anything, mock.Anything).Return(
			&Transaction{RefNo: 1, NewBalance: 99}, nil)
		svc := NewWalletService(mok, log.Logger)

		tr, err := svc.GetLatestTransaction(context.Background(), 1)
		assert.NoError(t, err)
		assert.NotNil(t, tr)
		assert.Equal(t, 1, tr.RefNo)
	})
}
