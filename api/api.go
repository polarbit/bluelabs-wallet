package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/polarbit/bluelabs-wallet/worker"
)

type (
	CreateWalletReq struct {
		Labels map[string]string `json:"labels"`
	}

	CreateWalletRes struct {
		Wallet WalletRep `json:"wallet"`
	}

	WalletRep struct {
		ID      string            `json:"id"`
		Balance float64           `json:"balance"`
		Labels  map[string]string `json:"labels"`
		Created time.Time         `json:"created"`
	}

	CreateTransactionReq struct {
		Amount      float64           `json:"amount"`
		Description string            `json:"description"`
		ExternalID  string            `json:"externalId"`
		Labels      map[string]string `json:"labels"`
	}

	CreateTransactionRes struct {
		Transaction *TransactionRep `json:"transaction"`
		Wallet      *WalletRep      `json:"wallet"`
	}

	TransactionRep struct {
		ID          string            `json:"id"`
		Amount      float64           `json:"amount"`
		Description string            `json:"description"`
		ExternalID  string            `json:"externalId"`
		Labels      map[string]string `json:"labels"`
		Created     time.Time         `json:"created"`
	}
)

var (
	wallets      = map[string]*WalletRep{}
	seq          = 1
	transactions = map[string]*TransactionRep{}
)

//----------
// Handlers
//----------

func createWallet(c echo.Context) error {
	req := &CreateWalletReq{}
	if err := c.Bind(req); err != nil {
		return err
	}
	id := strconv.Itoa(seq)
	wallets[id] = &WalletRep{ID: id, Labels: req.Labels, Created: time.Now().UTC()}
	seq++
	return c.JSON(http.StatusCreated, CreateWalletRes{Wallet: *wallets[id]})
}

func getWallet(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, wallets[id])
}

func createTransaction(c echo.Context) error {
	wid := c.Param("wid")
	req := &CreateTransactionReq{}
	if err := c.Bind(req); err != nil {
		return err
	}

	t := &TransactionRep{ID: strconv.Itoa(seq),
		Amount:      req.Amount,
		Description: req.Description,
		ExternalID:  req.ExternalID,
		Created:     time.Now().UTC(),
		Labels:      req.Labels}

	t.Labels["WID"] = wid

	transactions[t.ID] = t

	w := wallets[wid]

	return c.JSON(http.StatusOK, CreateTransactionRes{Transaction: t, Wallet: w})
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

	// Subscribe to signal to finish interaction
	finish := make(chan os.Signal, 1)
	finito := make(chan bool, 1)
	signal.Notify(finish, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		worker.StartClient(finito)
	}()

	go func() {
		<-finish
		e.Server.Shutdown(context.Background())
		finito <- true
	}()

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
