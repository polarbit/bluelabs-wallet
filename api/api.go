package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/polarbit/bluelabs-wallet/config"
	"github.com/polarbit/bluelabs-wallet/db"
	"github.com/polarbit/bluelabs-wallet/service"
)

type (
	CreateWalletRequest struct {
		service.WalletModel
	}

	CreateWalletResponse struct {
		service.Wallet
	}

	CreateTransactionRequest struct {
		service.TransactionModel
	}

	CreateTransactionResponse struct {
		service.Transaction
	}

	GetTransactionResponse struct {
		service.Transaction
	}
)

func StartAPI() {
	e := echo.New()

	// init vallet handler
	h := func() *walletHandler {
		wc := config.ReadConfig()
		repo := db.NewRepository(wc.Db)
		service := service.NewWalletService(repo)
		validate := validator.New()
		return &walletHandler{s: service, v: validate}
	}()

	// Set validator
	e.Logger.SetLevel(log.DEBUG)
	e.Validator = &CustomEchoValidator{v: h.v}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/wallets", func(c echo.Context) error { return h.createWallet(c) })
	e.GET("/wallets/:id", func(c echo.Context) error { return h.getWallet(c) })
	e.POST("/wallets/:wid/transactions", func(c echo.Context) error { return h.createTransaction(c) })
	e.GET("/wallets/:wid/transactions/:id", func(c echo.Context) error { return h.getTransaction(c) })

	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
