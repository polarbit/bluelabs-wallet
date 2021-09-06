//go:build integration
// +build integration

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/polarbit/bluelabs-wallet/config"
	"github.com/polarbit/bluelabs-wallet/db"
	"github.com/polarbit/bluelabs-wallet/service"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var r service.Repository

type testContext struct {
	ctx context.Context
	h   *walletHandler
	w   *CreateWalletResponse
	t   *CreateTransactionResponse
}

func TestApiIntegration(t *testing.T) {
	config.Init()

	// init wallet handler
	h := func() *walletHandler {
		repo := db.NewRepository(config.Config.Db, log.Logger)
		service := service.NewWalletService(repo, log.Logger)
		validate := validator.New()
		return &walletHandler{s: service, v: validate}
	}()

	tc := &testContext{ctx: context.Background(), h: h}

	t.Run("CreateWalletOk", func(t *testing.T) {
		createWalletOk(tc, t)
	})

	t.Run("CreateTransactionOk", func(t *testing.T) {
		createTransactionOk(tc, t)
	})
}

func createWalletOk(tc *testContext, t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomEchoValidator{v: tc.h.v}
	model := service.WalletModel{ExternalID: uuid.NewString(), Labels: map[string]string{"somekey": uuid.NewString()}}
	data, err := json.Marshal(model)
	if err != nil {
		panic("invalid model")
	}
	req := httptest.NewRequest(http.MethodPost, "/wallets", bytes.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, tc.h.createWallet(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		resp := &CreateWalletResponse{}
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), resp)) {
			assert.Equal(t, model.ExternalID, resp.ExternalID)
			assert.Equal(t, model.Labels, resp.Labels)
			assert.Positive(t, resp.ID)
		}

		tc.w = resp
	}
}

func createTransactionOk(tc *testContext, t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomEchoValidator{v: tc.h.v}
	model := service.TransactionModel{Amount: 9, Description: uuid.NewString(),
		Fingerprint: uuid.NewString(), Labels: map[string]string{"somekey": uuid.NewString()}}
	data, err := json.Marshal(model)
	if err != nil {
		panic("invalid model")
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/wallets/:wid/transactions")
	c.SetParamNames("wid")
	c.SetParamValues(strconv.Itoa(tc.w.ID))

	// Assertions
	if assert.NoError(t, tc.h.createTransaction(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		resp := &CreateTransactionResponse{}
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), resp)) {
			assert.Equal(t, model.Fingerprint, resp.Fingerprint)
			assert.Equal(t, model.Labels, resp.Labels)
			assert.NotEmpty(t, resp.ID)
			assert.Equal(t, model.Amount, resp.Amount)
			assert.Equal(t, model.Description, resp.Description)
			assert.Equal(t, 0., resp.OldBalance)
			assert.Equal(t, model.Amount, resp.NewBalance)
		}

		tc.t = resp
	}
}
