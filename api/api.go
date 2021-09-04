package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	CreateWalletRequest struct {
		ID     string            `json:"id"`
		Labels map[string]string `json:"labels"`
	}

	CreateWalletResponse struct {
		Wallet WalletRepresentation `json:"wallet"`
	}

	WalletRepresentation struct {
		ID      string            `json:"id"`
		Balance float64           `json:"balance"`
		Labels  map[string]string `json:"labels"`
		Created time.Time         `json:"created"`
	}

	CreateTransactionRequest struct {
		Amount      float64           `json:"amount"`
		Description string            `json:"description"`
		Labels      map[string]string `json:"labels"`
		Fingerprint string            `json:"fingerprint"`
	}

	CreateTransactionResponse struct {
		Transaction *TransactionRepresentation `json:"transaction"`
	}

	TransactionRepresentation struct {
		ID          string            `json:"id"`
		RefNo       int32             `json:"refNo"`
		Amount      float64           `json:"amount"`
		Description string            `json:"description"`
		Labels      map[string]string `json:"labels"`
		Fingerprint string            `json:"fingerprint"`
		Created     time.Time         `json:"created"`
		OldBalance  float64           `json:"oldBalance"`
		NewBalance  float64           `json:"newBalance"`
	}
)

var (
	wallets      = map[string]*WalletRepresentation{}
	seq          = 1
	transactions = map[string]*TransactionRepresentation{}
)

//----------
// Handlers
//----------

func createWallet(c echo.Context) error {
	req := &CreateWalletRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}
	id := strconv.Itoa(seq)
	wallets[id] = &WalletRepresentation{ID: id, Labels: req.Labels, Created: time.Now().UTC()}
	seq++
	return c.JSON(http.StatusCreated, CreateWalletResponse{Wallet: *wallets[id]})
}

func getWallet(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, wallets[id])
}

func createTransaction(c echo.Context) error {
	wid := c.Param("wid")
	req := &CreateTransactionRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	t := &TransactionRepresentation{ID: strconv.Itoa(seq),
		Amount:      req.Amount,
		Description: req.Description,
		Fingerprint: req.Fingerprint,
		Created:     time.Now().UTC(),
		Labels:      req.Labels}

	t.Labels["WID"] = wid

	transactions[t.ID] = t

	w := wallets[wid]

	return c.JSON(http.StatusOK, CreateTransactionResponse{Transaction: t, Wallet: w})
}

func getTransaction(c echo.Context) error {
	wid := c.Param("wid")
	id := c.Param("id")

	t := transactions[id]

	if t == nil {
		return c.String(http.StatusNotFound, "")
	}

	if t.Labels["WID"] != wid {
		return c.String(http.StatusBadRequest, "")
	}

	return c.JSON(http.StatusOK, t)
}

func StartAPI() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/wallets", createWallet)
	e.GET("/wallets/:id", getWallet)

	e.POST("wallets/:wid/transactions", createTransaction)
	e.GET("wallets/:wid/transactions/:id", getTransaction)

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
