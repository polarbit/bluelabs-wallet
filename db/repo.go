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

// TODO: Use connection pooling; or add to your notes.

func (r *repository) CreateWallet(ctx context.Context, w *controller.Wallet) error {
	conn, err := pgx.Connect(context.Background(), r.url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	stmt := `
	insert into wallets 
	(externalid, labels, created) 
	values ($1, $2, $3) 
	returning id`

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transactin failed: %v", err)
	}

	// insert wallet
	err = tx.QueryRow(ctx, stmt, w.ExternalID, w.Labels, w.Created).Scan(&w.ID)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("insert wallet failed: %v", err)
	}

	// insert balance
	_, err = tx.Exec(ctx,
		`insert into wallet_balances (wid, amount) values ($1, $2)`, w.ID, 0.)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("insert wallet balance failed: %v", err)
	}

	tx.Commit(ctx)

	return nil
}

func (r *repository) GetWallet(ctx context.Context, wid int32) (*controller.Wallet, error) {
	conn, err := pgx.Connect(context.Background(), r.url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	w := controller.Wallet{}
	stmt := `select id, externalid, labels, created from wallets where id = $1`
	err = conn.QueryRow(ctx, stmt, wid).Scan(&w.ID, &w.ExternalID, &w.Labels, &w.Created)
	if err != nil {
		return nil, fmt.Errorf("read wallet failed: %v", err)
	}

	return &w, nil
}

func (r *repository) GetWalletBalance(ctx context.Context, wid string) (b float64, err error) {
	conn, err := pgx.Connect(context.Background(), r.url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	stmt := `select (id, externalid, labels, created) from wallets where id = $1`
	err = conn.QueryRow(ctx, stmt, wid).Scan(&b)
	if err != nil {
		return 0., fmt.Errorf("read wallet balance failed: %v", err)
	}

	return b, err
}
