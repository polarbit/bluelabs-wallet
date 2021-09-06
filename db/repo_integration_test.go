//go:build integration
// +build integration

package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/polarbit/bluelabs-wallet/config"
	"github.com/polarbit/bluelabs-wallet/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var c *config.AppConfig
var r service.Repository

type testContext struct {
	ctx context.Context
	w   *service.Wallet
	t   *service.Transaction
}

func TestRepositoryIntegration(t *testing.T) {
	c = config.ReadConfig()

	var logger zerolog.Logger
	if level, err := zerolog.ParseLevel(c.LogLevel); err != nil {
		fmt.Printf("level %v\n", level)
		logger = log.Logger.Level(zerolog.DebugLevel)
	} else {
		logger = log.Logger.Level(level)
	}

	r = NewRepository(c.Db, logger)
	tc := &testContext{ctx: context.Background()}

	t.Run("CreateWalletOk", func(t *testing.T) {
		testCreateWalletOk(tc, t)
	})

	t.Run("GetWalletOk", func(t *testing.T) {
		testGetWalletOk(tc, t)
	})

	t.Run("GetWalletBalanceOk", func(t *testing.T) {
		testGetWalletBalanceOk(tc, t)
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

func testGetWalletBalanceOk(tc *testContext, t *testing.T) {
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
	assert.ErrorIs(t, err, service.ErrTransactionFailedButRetriable)
}
