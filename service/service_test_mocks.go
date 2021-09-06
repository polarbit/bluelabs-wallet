//go:build !integration
// +build !integration

package service

import (
	"context"

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

func (m *mockRepository) GetWalletBalance(ctx context.Context, wid int) (b float64, err error) {

	args := m.Called(ctx, wid)
	return args.Get(0).(float64), args.Error(1)
}
