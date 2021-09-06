//go:build integration
// +build integration

package db

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/polarbit/bluelabs-wallet/config"
	"github.com/polarbit/bluelabs-wallet/service"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var r service.Repository

type testContext struct {
	ctx context.Context
	w   *service.Wallet
	t   *service.Transaction
}

func TestRepositoryIntegration(t *testing.T) {
	config.Init()
	r = NewRepository(config.Config.Db, log.Logger)
	tc := &testContext{ctx: context.Background()}

	t.Run("CreateWalletOk", func(t *testing.T) {
		testCreateWalletOk(tc, t)
	})

	t.Run("GetWalletOk", func(t *testing.T) {
		testGetWalletOk(tc, t)
	})

	t.Run("GetWalletBalanceZeroOk", func(t *testing.T) {
		testGetWalletBalanceZeroOk(tc, t)
	})

	t.Run("GetWalletBalanceReturnsNotFound", func(t *testing.T) {
		testGetWalletBalanceReturnsNotFound(tc, t)
	})

	t.Run("CreateTransactionOk", func(t *testing.T) {
		testCreateTransactionOk(tc, t)
	})

	t.Run("CreateTransactionFailsByRefno", func(t *testing.T) {
		testCreateTransactionFailsByRefNo(tc, t)
	})

	t.Run("CreateTransactionFailsByFingerprint", func(t *testing.T) {
		testCreateTransactionFailsByFingerprint(tc, t)
	})

	t.Run("CreateTransactionFailsByConcurrency", func(t *testing.T) {
		testCreateTransactionFailsByConcurrency(tc, t)
	})

	t.Run("GetWalletBalanceOk", func(t *testing.T) {
		testGetWalletBalanceOk(tc, t)
	})

	t.Run("GetLatestTransactionOk", func(t *testing.T) {
		testGetLatestTransactionOk(tc, t)
	})
}

func testCreateWalletOk(tc *testContext, t *testing.T) {
	tc.w = &service.Wallet{
		ExternalID: uuid.NewString(),
		Labels:     map[string]string{"Source": "IntegrationTest"},
		Created:    time.Now().UTC().Truncate(time.Microsecond),
	}
	err := r.CreateWallet(tc.ctx, tc.w)

	assert.Nil(t, err)
	assert.Greater(t, tc.w.ID, 0)
}

func testGetWalletOk(tc *testContext, t *testing.T) {
	w, err := r.GetWallet(tc.ctx, tc.w.ID)
	assert.Nil(t, err)
	assert.NotNil(t, w)
	assert.Equal(t, tc.w.ID, w.ID)
	assert.Equal(t, tc.w.ExternalID, w.ExternalID)
	assert.Equal(t, tc.w.Created, w.Created)
	assert.Contains(t, w.Labels, "Source")
	assert.Equal(t, tc.w.Labels["Source"], w.Labels["Source"])
}

func testGetWalletBalanceZeroOk(tc *testContext, t *testing.T) {
	b, err := r.GetWalletBalance(tc.ctx, tc.w.ID)
	assert.NoError(t, err)
	assert.Equal(t, 0., b)
}

func testGetWalletBalanceReturnsNotFound(tc *testContext, t *testing.T) {
	_, err := r.GetWalletBalance(tc.ctx, -1)
	assert.ErrorIs(t, err, service.ErrWalletNotFound)
}

func testCreateTransactionOk(tc *testContext, t *testing.T) {
	tc.t = &service.Transaction{
		ID:          uuid.NewString(),
		RefNo:       1,
		Amount:      24.75,
		Description: "test transaction",
		Labels:      map[string]string{"test": "true"},
		Fingerprint: uuid.NewString(),
		Created:     time.Now().UTC().Truncate(time.Millisecond),
		OldBalance:  0.,
		NewBalance:  24.75,
	}
	err := r.CreateTransaction(tc.ctx, tc.w.ID, tc.t)
	assert.NoError(t, err)
}

func testCreateTransactionFailsByRefNo(tc *testContext, t *testing.T) {
	tr := &service.Transaction{
		ID:          uuid.NewString(),
		RefNo:       1,
		Amount:      30,
		Description: "test duplicate refno",
		Labels:      map[string]string{"test": "true"},
		Fingerprint: uuid.NewString(),
		Created:     time.Now().UTC().Truncate(time.Millisecond),
		OldBalance:  24.75,
		NewBalance:  54.75,
	}
	err := r.CreateTransaction(tc.ctx, tc.w.ID, tr)
	assert.ErrorIs(t, err, service.ErrTransactionAlreadyExistsByRefNo)
}

func testCreateTransactionFailsByFingerprint(tc *testContext, t *testing.T) {
	tr := &service.Transaction{
		ID:          uuid.NewString(),
		RefNo:       2,
		Amount:      30,
		Description: "test duplicate fingerprint",
		Labels:      map[string]string{"test": "true"},
		Fingerprint: tc.t.Fingerprint,
		Created:     time.Now().UTC().Truncate(time.Millisecond),
		OldBalance:  24.75,
		NewBalance:  54.75,
	}
	err := r.CreateTransaction(tc.ctx, tc.w.ID, tr)
	assert.ErrorIs(t, err, service.ErrTransactionAlreadyExistsByFingerprint)
}

func testCreateTransactionFailsByConcurrency(tc *testContext, t *testing.T) {
	tr := &service.Transaction{
		ID:          uuid.NewString(),
		RefNo:       2,
		Amount:      30.,
		Description: "test duplicate fingerprint",
		Labels:      map[string]string{"test": "true"},
		Fingerprint: uuid.NewString(),
		Created:     time.Now().UTC().Truncate(time.Millisecond),
		OldBalance:  0.,
		NewBalance:  30.,
	}
	err := r.CreateTransaction(tc.ctx, tc.w.ID, tr)
	assert.ErrorIs(t, err, service.ErrTransactionConsistency)
}

func testGetWalletBalanceOk(tc *testContext, t *testing.T) {
	b, err := r.GetWalletBalance(tc.ctx, tc.w.ID)
	assert.NoError(t, err)
	assert.Equal(t, 24.75, b)
}

func testGetLatestTransactionOk(tc *testContext, t *testing.T) {
	lt, err := r.GetLatestTransaction(tc.ctx, tc.w.ID)
	assert.NoError(t, err)
	assert.Equal(t, tc.t.ID, lt.ID)
	assert.Equal(t, tc.t.RefNo, lt.RefNo)
	assert.Equal(t, tc.t.Amount, lt.Amount)
	assert.Equal(t, tc.t.Description, lt.Description)
	assert.Equal(t, tc.t.Fingerprint, lt.Fingerprint)
	assert.Equal(t, tc.t.Labels, lt.Labels)
	assert.Equal(t, tc.t.OldBalance, lt.OldBalance)
	assert.Equal(t, tc.t.NewBalance, lt.NewBalance)
}
