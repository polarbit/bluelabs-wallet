package service

import "time"

type (
	Wallet struct {
		ID         int               `json:"id"`
		Labels     map[string]string `json:"labels"`
		ExternalID string            `json:"externalId"`
		Created    time.Time         `json:"created"`
	}

	WalletModel struct {
		Labels     map[string]string `json:"labels" validate:"max=1"`
		ExternalID string            `json:"externalId" validate:"required,max=50"`
	}

	TransactionModel struct {
		Amount      float64
		Description string
		Labels      map[string]string
		Fingerprint string
	}

	Transaction struct {
		ID          string
		RefNo       int
		Amount      float64
		Description string
		Labels      map[string]string
		Fingerprint string
		Created     time.Time
		OldBalance  float64
		NewBalance  float64
	}
)
