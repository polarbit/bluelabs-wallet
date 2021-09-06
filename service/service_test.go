//go:build !integration
// +build !integration

package service

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mok = &mockRepository{}
var svc Service

func TestMain(m *testing.M) {
	mok.On("GetWallet", mock.Anything, 10).Return(&Wallet{ID: 10}, nil)
	mok.On("GetWallet", mock.Anything, 11).Return((*Wallet)(nil), ErrWalletNotFound)

	mok.On("CreateWallet", mock.Anything,
		mock.MatchedBy(func(w *Wallet) bool { w.ID = 99; return w.ExternalID == "99" })).
		Return(nil)
	mok.On("CreateWallet", mock.Anything,
		mock.MatchedBy(func(w *Wallet) bool { return w.ExternalID == "100" })).
		Return(ErrWalletAlreadyExists)

	logger := log.Logger
	svc = NewWalletService(mok, logger)

	exitcode := m.Run()
	os.Exit(exitcode)
}

func TestGetWallet(t *testing.T) {
	t.Run("WalletFound", func(t *testing.T) {
		w, err := svc.GetWallet(context.Background(), 10)
		assert.NoError(t, err)
		assert.NotNil(t, w)
		assert.Equal(t, 10, w.ID)
	})

	t.Run("WalletNotFound", func(t *testing.T) {
		w, err := svc.GetWallet(context.Background(), 11)
		assert.Nil(t, w)
		assert.NoError(t, err)
		assert.True(t, errors.Is(err, ErrWalletNotFound))
	})
}

func TestCreateWallet(t *testing.T) {
	t.Run("WalletCreated", func(t *testing.T) {
		w, err := svc.CreateWallet(context.Background(), &WalletModel{
			ExternalID: "99",
		})
		assert.NoError(t, err)
		assert.NotNil(t, w)
		assert.Equal(t, "99", w.ExternalID)
		assert.Equal(t, 99, w.ID)
	})

	t.Run("WalletAlreadyExists", func(t *testing.T) {
		w, err := svc.CreateWallet(context.Background(), &WalletModel{
			ExternalID: "100",
		})
		assert.ErrorIs(t, err, ErrWalletAlreadyExists)
		assert.Nil(t, w)
	})
}
