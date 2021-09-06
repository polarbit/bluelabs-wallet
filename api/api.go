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
	"github.com/ziflex/lecho/v2"

	"github.com/polarbit/bluelabs-wallet/config"
	"github.com/polarbit/bluelabs-wallet/db"
	"github.com/polarbit/bluelabs-wallet/service"
	"github.com/rs/zerolog/log"
)

func StartAPI() {

	logger := log.Logger

	// init wallet handler
	h := func() *walletHandler {
		wc := config.ReadConfig()
		repo := db.NewRepository(wc.Db, logger)
		service := service.NewWalletService(repo, logger)
		validate := validator.New()
		return &walletHandler{s: service, v: validate}
	}()

	e := echo.New()
	e.Logger = lecho.From(logger)                          // Set zerlogger as echo logger
	e.Validator = &CustomEchoValidator{v: validator.New()} // Set validator
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/wallets", func(c echo.Context) error { return h.createWallet(c) })
	e.GET("/wallets/:id", func(c echo.Context) error { return h.getWallet(c) })
	e.POST("/wallets/:wid/transactions", func(c echo.Context) error { return h.createTransaction(c) })
	e.GET("/wallets/:wid/transactions/:id", func(c echo.Context) error { return h.getTransaction(c) })

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
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
