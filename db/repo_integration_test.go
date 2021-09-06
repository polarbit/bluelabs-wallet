//go:build integration
// +build integration

package db

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/polarbit/bluelabs-wallet/config"
	"github.com/polarbit/bluelabs-wallet/service"
	"github.com/stretchr/testify/assert"
)

var c *config.AppConfig
var r service.Repository

type testContext struct {
	ctx context.Context
	w   *service.Wallet
}

func TestRepositoryIntegration(t *testing.T) {
	c = config.ReadConfig()
	r = NewRepository(c.Db)
	tc := &testContext{ctx: context.Background()}

	t.Run("CreateWallet", func(t *testing.T) {
		testCreateWallet(tc, t)
	})

	t.Run("GetWallet", func(t *testing.T) {
		testGetWallet(tc, t)
	})

	t.Run("GetWalletBalance", func(t *testing.T) {
		testGetWalletBalance(tc, t)
	})

	t.Run("GetWalletBalanceReturnsNotFound", func(t *testing.T) {
		testGetWalletBalanceReturnsNotFound(tc, t)
	})
}

func testCreateWallet(tc *testContext, t *testing.T) {
	tc.w = &service.Wallet{
		ExternalID: uuid.NewString(),
		Labels:     map[string]string{"Source": "IntegrationTest"},
		Created:    time.Now().UTC().Truncate(time.Microsecond),
	}
	err := r.CreateWallet(tc.ctx, tc.w)

	assert.Nil(t, err)
	assert.Greater(t, tc.w.ID, 0)
}

func testGetWallet(tc *testContext, t *testing.T) {
	w, err := r.GetWallet(tc.ctx, tc.w.ID)
	assert.Nil(t, err)
	assert.NotNil(t, w)
	assert.Equal(t, tc.w.ID, w.ID)
	assert.Equal(t, tc.w.ExternalID, w.ExternalID)
	assert.Equal(t, tc.w.Created, w.Created)
	assert.Contains(t, w.Labels, "Source")
	assert.Equal(t, tc.w.Labels["Source"], w.Labels["Source"])
}

func testGetWalletBalance(tc *testContext, t *testing.T) {
	b, err := r.GetWalletBalance(tc.ctx, tc.w.ID)
	assert.Nil(t, err)
	assert.Equal(t, 0., b)
}

func testGetWalletBalanceReturnsNotFound(tc *testContext, t *testing.T) {
	_, err := r.GetWalletBalance(tc.ctx, -1)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, service.ErrWalletNotFound))
}
