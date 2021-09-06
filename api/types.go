package api

import "github.com/polarbit/bluelabs-wallet/service"

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
