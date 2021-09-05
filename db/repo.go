package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/polarbit/bluelabs-wallet/controller"
)

type repository struct {
	url string
}

func NewRepository(url string) controller.Repository {
	parseUrl(url)
	return &repository{url: url}
}

func (r *repository) CreateWallet(ctx context.Context, w *controller.Wallet) error {
	conn, err := pgx.Connect(context.Background(), r.url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	stmt := `
	insert into wallets 
	(id, externalid, labels, created) 
	values ($1, $2, $3)`

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transactin failed: %v", err)
	}
	_, err = tx.Exec(ctx, stmt, stmt, w.ID, w.ExternalID, w.Labels, w.Created)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("insert wallet failed: %v", err)
	}

	_, err = tx.Exec(ctx,
		`insert into wallet_balances (wid, amount) values ($1, $2)`, w.ID, 0.)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("insert wallet balance failed: %v", err)
	}

	tx.Commit(ctx)

	return nil
}

func (r *repository) GetWallet(ctx context.Context, wid string) (*controller.Wallet, error) {
	return nil, nil
}

func (r *repository) GetWalletBalance(ctx context.Context, wid string) (float64, error) {
	return 0, nil
}
