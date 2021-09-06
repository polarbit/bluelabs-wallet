package service

import "time"

type (
	Wallet struct {
		ID         int               `json:"id"`
		Labels     map[string]string `json:"labels"`
		ExternalID string            `json:"externalid"`
		Created    time.Time         `json:"created"`
	}

	WalletModel struct {
		Labels     map[string]string `json:"labels" validate:"max=10"`
		ExternalID string            `json:"externalId" validate:"required,max=50"`
	}

	TransactionModel struct {
		Amount      float64           `json:"amount" validate:"required"`
		Description string            `json:"description" validate:"required,max=100"`
		Labels      map[string]string `json:"labels" validate:"max=10"`
		Fingerprint string            `json:"fingerprint" validate:"required,max=50"`
	}

	Transaction struct {
		ID          string            `json:"id"`
		RefNo       int               `json:"refno"`
		Amount      float64           `json:"amount"`
		Description string            `json:"description"`
		Labels      map[string]string `json:"labels"`
		Fingerprint string            `json:"fingerprint"`
		Created     time.Time         `json:"created"`
		OldBalance  float64           `json:"oldbalance"`
		NewBalance  float64           `json:"newbalance"`
	}
)
