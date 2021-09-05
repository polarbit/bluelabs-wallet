package db_test

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"

	"github.com/polarbit/bluelabs-wallet/config"
	"github.com/polarbit/bluelabs-wallet/db"
	"github.com/polarbit/bluelabs-wallet/wallet"
)

var runIt = flag.Bool("it", false, "Run integration test suite")
var appConfig *config.AppConfig
var walletRepo wallet.WalletRepository

func TestMain(m *testing.M) {

	appConfig = config.ReadConfig()
	walletRepo = db.NewRepository(appConfig)

	log.Println("Do stuff BEFORE the tests!")
	exitVal := m.Run()
	log.Println("Do stuff AFTER the tests!")

	os.Exit(exitVal)
}

func TestCreateWalletIt(t *testing.T) {
	if !*runIt {
		t.Skip("Skip it in short mode")
	}

	ctx := context.Background()
	model := &wallet.WalletModel{}
	walletRepo.CreateWallet(ctx, model)
}

func TestGetWalletIt(t *testing.T) {
	if !*runIt {
		t.Skip("Skip it in short mode")
	}
	log.Println("TestB running")
}
