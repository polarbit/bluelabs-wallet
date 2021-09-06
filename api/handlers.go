package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/polarbit/bluelabs-wallet/service"
)

type walletHandler struct {
	s service.Service
	v *validator.Validate
}

func (h *walletHandler) createWallet(c echo.Context) error {
	req := &CreateWalletRequest{}

	// bind
	if err := c.Bind(req); err != nil {
		c.Logger().Debug("bind: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	// validate
	if err := c.Validate(req); err != nil {
		c.Logger().Debug("validate: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	// handle
	w, err := h.s.CreateWallet(c.Request().Context(), &req.WalletModel)
	if err != nil {
		if errors.Is(err, service.ErrWalletAlreadyExists) {
			return c.String(http.StatusConflict, err.Error())
		}

		c.Logger().Debug("handle: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, CreateWalletResponse{*w})
}

func (h *walletHandler) getWallet(c echo.Context) error {
	// validate
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Debug("route: ", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	// handle
	w, err := h.s.GetWallet(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrWalletNotFound) {
			return c.String(http.StatusNotFound, err.Error())
		}

		c.Logger().Debug("handle: ", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, CreateWalletResponse{*w})
}

func (h *walletHandler) createTransaction(c echo.Context) error {

	return c.JSON(http.StatusOK, CreateTransactionResponse{})
}

func (h *walletHandler) getTransaction(c echo.Context) error {

	return c.JSON(http.StatusOK, GetTransactionResponse{})
}
