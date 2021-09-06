//go:build !integration
// +build !integration

package service

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) CreateWallet(ctx context.Context, w *Wallet) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *mockRepository) GetWallet(ctx context.Context, wid int) (*Wallet, error) {
	args := m.Called(ctx, wid)
	return args.Get(0).(*Wallet), args.Error(1)
}

func (m *mockRepository) GetWalletBalance(ctx context.Context, wid string) (b float64, err error) {

	args := m.Called(ctx, wid)
	return args.Get(0).(float64), args.Error(1)
}

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

	svc = NewWalletService(mok)

	exitcode := m.Run()
	os.Exit(exitcode)
}

func TestGetWallet(t *testing.T) {
	t.Run("WalletFound", func(t *testing.T) {
		w, err := svc.GetWallet(context.Background(), 10)
		assert.Nil(t, err)
		assert.NotNil(t, w)
		assert.Equal(t, 10, w.ID)
	})

	t.Run("WalletNotFound", func(t *testing.T) {
		w, err := svc.GetWallet(context.Background(), 11)
		assert.Nil(t, w)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, ErrWalletNotFound))
	})
}

func TestCreateWallet(t *testing.T) {
	t.Run("WalletCreated", func(t *testing.T) {
		w, err := svc.CreateWallet(context.Background(), &WalletModel{
			ExternalID: "99",
		})
		assert.Nil(t, err)
		assert.NotNil(t, w)
		assert.Equal(t, "99", w.ExternalID)
		assert.Equal(t, 99, w.ID)
	})

	t.Run("WalletAlreadyExists", func(t *testing.T) {
		w, err := svc.CreateWallet(context.Background(), &WalletModel{
			ExternalID: "100",
		})
		assert.NotNil(t, err)
		assert.Nil(t, w)
		assert.True(t, errors.Is(err, ErrWalletAlreadyExists))
	})
}
