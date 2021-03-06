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

func (m *mockRepository) GetWalletBalance(ctx context.Context, wid int) (float64, error) {

	args := m.Called(ctx, wid)
	return args.Get(0).(float64), args.Error(1)
}

func (m *mockRepository) CreateTransaction(ctx context.Context, wid int, t *Transaction) error {
	args := m.Called(ctx, wid, t)
	return args.Error(0)
}

func (m *mockRepository) GetLatestTransaction(ctx context.Context, wid int) (*Transaction, error) {
	args := m.Called(ctx, wid)
	return args.Get(0).(*Transaction), args.Error(1)
}
