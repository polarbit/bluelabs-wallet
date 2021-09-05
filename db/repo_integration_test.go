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
)

var enabled = flag.Bool("integration", false, "run integration tests")
var c *config.AppConfig
var r controller.Repository

type testContext struct {
	t   *testing.T
	ctx context.Context
	w   *controller.Wallet
}

func TestRepositoryIntegration(t *testing.T) {
	if !*enabled {
		t.Skip("Skip repository integration tests")
	}

	c = config.ReadConfig()
	r = db.NewRepository(c.Db)

	t.Run("CreateWallet", func(t *testing.T) {
		tc := &testContext{t: t, ctx: context.Background()}
		testCreateWallet(tc)
	})

	t.Run("GetWallet", func(t *testing.T) {
		tc := &testContext{t: t, ctx: context.Background()}
		testGetWallet(tc)
	})
}

func testCreateWallet(tc *testContext) {
	tc.w = &controller.Wallet{
		ExternalID: uuid.NewString(),
		Labels:     map[string]string{"Source": "IntegrationTest"},
		Created:    time.Now().UTC().Truncate(time.Microsecond),
	}
	err := r.CreateWallet(tc.ctx, tc.w)
	if err != nil {
		tc.t.Errorf("test failed %v", err)
	}
}

func testGetWallet(tc *testContext) {
	w, err := r.GetWallet(tc.ctx, tc.w.ID)
	if err != nil {
		tc.t.Errorf("test failed %v", err)
	}
	if w == nil {
		tc.t.Error("returned wallet is nil")
	}
}
