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
	"github.com/rs/zerolog/log"

	"github.com/polarbit/bluelabs-wallet/config"
	"github.com/polarbit/bluelabs-wallet/db"
	"github.com/polarbit/bluelabs-wallet/service"
)

func StartAPI() {
	config.Init()

	// init wallet handler
	h := func() *walletHandler {
		repo := db.NewRepository(config.Config.Db, log.Logger)
		service := service.NewWalletService(repo, log.Logger)
		validate := validator.New()
		return &walletHandler{s: service, v: validate}
	}()

	// TODO: Push validator and minamount logic to elsewhere
	// h.v.RegisterValidation("minamount", func(fl validator.FieldLevel) bool {
	// 	if fl.Field().Kind() != reflect.Float32 && fl.Field().Kind() != reflect.Float64 {
	// 		return false
	// 	}
	// 	return math.Abs(fl.Field().Float()) >= 1.0
	// })

	e := echo.New()
	// e.Logger = lecho.From(log.Logger)                      // Set zerlogger as echo logger
	e.Validator = &CustomEchoValidator{v: h.v} // Set validator
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/wallets", func(c echo.Context) error { return h.createWallet(c) })
	e.GET("/wallets/:id", func(c echo.Context) error { return h.getWallet(c) })
	e.GET("/wallets/:id/balance", func(c echo.Context) error { return h.getWalletBalance(c) })
	e.POST("/wallets/:wid/transactions", func(c echo.Context) error { return h.createTransaction(c) })
	e.GET("/wallets/:wid/transactions/latest", func(c echo.Context) error { return h.getLatestTransaction(c) })

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
