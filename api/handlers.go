package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/polarbit/bluelabs-wallet/service"
	"github.com/rs/zerolog"
)

type walletHandler struct {
	s service.Service
	v *validator.Validate
	l zerolog.Logger
}

func (h *walletHandler) createWallet(c echo.Context) error {
	req := &CreateWalletRequest{}

	// bind
	if err := c.Bind(req); err != nil {
		h.l.Debug().Err(err).Msg("bind error")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// validate
	if err := c.Validate(req); err != nil {
		h.l.Debug().Err(err).Msg("validate error")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// handle
	w, err := h.s.CreateWallet(c.Request().Context(), &req.WalletModel)
	if err != nil {
		if errors.Is(err, service.ErrWalletAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}

		h.l.Debug().Err(err).Msg("")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, CreateWalletResponse{*w})
}

func (h *walletHandler) getWallet(c echo.Context) error {
	// validate
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.l.Debug().Err(err).Msg("route error")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// handle
	w, err := h.s.GetWallet(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrWalletNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		h.l.Debug().Err(err).Msg("")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, CreateWalletResponse{*w})
}

func (h *walletHandler) createTransaction(c echo.Context) error {
	// route
	wid, err := strconv.Atoi(c.Param("wid"))
	if err != nil {
		h.l.Debug().Err(err).Msg("route error")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// bind
	req := &CreateTransactionRequest{}
	if err := c.Bind(req); err != nil {
		h.l.Debug().Err(err).Msg("bind error")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// validate
	if err := c.Validate(req); err != nil {
		h.l.Debug().Err(err).Msg("validate error")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// handle
	tr, err := h.s.CreateTransaction(c.Request().Context(), wid, &req.TransactionModel)
	if err != nil {
		if errors.Is(err, service.ErrWalletNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		if errors.Is(err, service.ErrNotEnoughWalletBalance) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if errors.Is(err, service.ErrTransactionAlreadyExistsByFingerprint) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		if errors.Is(err, service.ErrTransactionAlreadyExistsByRefNo) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		if errors.Is(err, service.ErrTransactionConsistency) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}

		h.l.Debug().Err(err).Msg("")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, CreateTransactionResponse{*tr})
}

func (h *walletHandler) getLatestTransaction(c echo.Context) error {
	// validate
	id, err := strconv.Atoi(c.Param("wid"))
	if err != nil {
		h.l.Debug().Err(err).Msg("route error")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// handle
	tr, err := h.s.GetLatestTransaction(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTransactionNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		h.l.Debug().Err(err).Msg("")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, GetTransactionResponse{*tr})
}

func (h *walletHandler) getWalletBalance(c echo.Context) error {
	// validate
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.l.Debug().Err(err).Msg("route error")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// handle
	b, err := h.s.GetWalletBalance(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrWalletNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		h.l.Debug().Err(err).Msg("")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, b)
}
