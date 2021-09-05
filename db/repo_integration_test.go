package db_test

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/polarbit/bluelabs-wallet/config"
	"github.com/polarbit/bluelabs-wallet/controller"
	"github.com/polarbit/bluelabs-wallet/db"
	"github.com/stretchr/testify/assert"
)

var enabled = flag.Bool("integration", false, "run integration tests")
var c *config.AppConfig
var r controller.Repository

type testContext struct {
	ctx context.Context
	w   *controller.Wallet
}

func TestRepositoryIntegration(t *testing.T) {
	if !*enabled {
		t.Skip("Skip repository integration tests")
	}

	c = config.ReadConfig()
	r = db.NewRepository(c.Db)
	tc := &testContext{ctx: context.Background()}

	t.Run("CreateWallet", func(t *testing.T) {
		testCreateWallet(tc, t)
	})

	t.Run("GetWallet", func(t *testing.T) {
		testGetWallet(tc, t)
	})
}

func testCreateWallet(tc *testContext, t *testing.T) {
	tc.w = &controller.Wallet{
		ExternalID: uuid.NewString(),
		Labels:     map[string]string{"Source": "IntegrationTest"},
		Created:    time.Now().UTC().Truncate(time.Microsecond),
	}
	err := r.CreateWallet(tc.ctx, tc.w)

	assert.Nil(t, err)
	assert.Greater(t, tc.w.ID, int32(0))
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
