package db

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/polarbit/bluelabs-wallet/service"
	"github.com/rs/zerolog"
)

const (
	errTextWalletAlreadyExists                   = `duplicate key value violates unique constraint "wallets_externalid_key"`
	errTextRowNotFound                           = `no rows in result set`
	errTextTransactionAlreadyExistsByRefno       = `duplicate key value violates unique constraint "ix_wid_refno"`
	errTextTransactionAlreadyExistsByFingerprint = `duplicate key value violates unique constraint "wallet_transactions_fingerprint_key"`
)

type repository struct {
	url string
	l   zerolog.Logger
}

func NewRepository(url string, logger zerolog.Logger) service.Repository {
	parseUrl(url)
	return &repository{url: url, l: logger}
}

func (r *repository) CreateWallet(ctx context.Context, w *service.Wallet) error {
	conn, err := pgx.Connect(context.Background(), r.url)
	if err != nil {
		r.l.Error().Err(err).Send()
		return service.NewDbError(err)
	}
	defer conn.Close(context.Background())

	stmt := `
	insert into wallets 
	(externalid, labels, created) 
	values ($1, $2, $3) 
	returning id`

	tx, err := conn.Begin(ctx)
	if err != nil {
		r.l.Error().Err(err).Send()
		return service.NewDbError(err)
	}

	// insert wallet
	err = tx.QueryRow(ctx, stmt, w.ExternalID, w.Labels, w.Created).Scan(&w.ID)
	if err != nil {
		tx.Rollback(ctx)
		if strings.Contains(err.Error(), errTextWalletAlreadyExists) {
			return service.ErrWalletAlreadyExists
		}
		r.l.Error().Err(err).Send()
		return service.NewDbError(err)
	}

	// insert balance
	_, err = tx.Exec(ctx,
		`insert into wallet_balances (wid, amount) values ($1, $2)`, w.ID, 0.)
	if err != nil {
		tx.Rollback(ctx)
		r.l.Error().Err(err).Send()
		return service.NewDbError(err)
	}

	tx.Commit(ctx)

	return nil
}

func (r *repository) GetWallet(ctx context.Context, wid int) (*service.Wallet, error) {
	conn, err := pgx.Connect(context.Background(), r.url)
	if err != nil {
		r.l.Error().Err(err).Send()
		return nil, service.NewDbError(err)
	}
	defer conn.Close(context.Background())

	stmt := `select id, externalid, labels, created from wallets where id = $1`
	rows, err := conn.Query(ctx, stmt, wid)

	if err != nil {
		r.l.Error().Err(err).Send()
		return nil, service.NewDbError(err)
	}
	defer rows.Close()

	w := service.Wallet{}
	if ok := rows.Next(); !ok {
		return nil, service.ErrWalletNotFound
	}

	err = rows.Scan(&w.ID, &w.ExternalID, &w.Labels, &w.Created)
	if err != nil {
		r.l.Error().Err(err).Send()
		return nil, service.NewDbError(err)
	}

	return &w, nil
}

func (r *repository) GetWalletBalance(ctx context.Context, wid int) (b float64, err error) {
	conn, err := pgx.Connect(context.Background(), r.url)
	if err != nil {
		r.l.Error().Err(err).Send()
		return 0., service.NewDbError(err)
	}
	defer conn.Close(context.Background())

	stmt := `select amount from wallet_balances where wid = $1`
	err = conn.QueryRow(ctx, stmt, wid).Scan(&b)
	if err != nil {
		if strings.Contains(err.Error(), errTextRowNotFound) {
			return 0., service.ErrWalletNotFound
		}
		r.l.Error().Err(err).Send()
		return 0., service.NewDbError(err)
	}

	return b, err
}

func (r *repository) CreateTransaction(ctx context.Context, wid int, t *service.Transaction) error {
	conn, err := pgx.Connect(context.Background(), r.url)
	if err != nil {
		r.l.Error().Err(err).Send()
		return service.NewDbError(err)
	}
	defer conn.Close(context.Background())

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		r.l.Error().Err(err).Send()
		return service.NewDbError(err)
	}

	// insert transaction
	stmt := `insert into wallet_transactions 
	(id, wid, refno, amount, description, labels, fingerprint, old_balance, new_balance, created) 
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = tx.Exec(ctx, stmt, t.ID, wid, t.RefNo, t.Amount, t.Description, t.Labels,
		t.Fingerprint, t.OldBalance, t.NewBalance, t.Created)
	if err != nil {
		tx.Rollback(ctx)
		r.l.Error().Err(err).Msg("insert transaction failed")
		if strings.Contains(err.Error(), errTextTransactionAlreadyExistsByRefno) {
			return service.ErrTransactionAlreadyExistsByRefNo
		}
		if strings.Contains(err.Error(), errTextTransactionAlreadyExistsByFingerprint) {
			return service.ErrTransactionAlreadyExistsByFingerprint
		}
		r.l.Error().Err(err).Send()
		return service.NewDbError(err)
	}

	// update balance
	stmt = `update wallet_balances set amount = $2 where wid=$1 and amount = $3`
	ctag, err := tx.Exec(ctx, stmt, wid, t.NewBalance, t.OldBalance)
	if err != nil {
		tx.Rollback(ctx)
		r.l.Error().Err(err).Send()
		return service.NewDbError(err)
	}
	if ctag.RowsAffected() != 1 {
		tx.Rollback(ctx)
		r.l.Error().Err(err).Send()
		return service.ErrTransactionConsistency
	}

	tx.Commit(ctx)

	return nil
}

func (r *repository) GetLatestTransaction(ctx context.Context, wid int) (*service.Transaction, error) {
	conn, err := pgx.Connect(context.Background(), r.url)
	if err != nil {
		r.l.Error().Err(err).Send()
		return nil, service.NewDbError(err)
	}
	defer conn.Close(context.Background())

	stmt := `select id, refno, amount, description, labels, fingerprint, created, old_balance, new_balance
	from wallet_transactions where wid = $1 
	order by wid, refno desc
	limit 1`
	rows, err := conn.Query(ctx, stmt, wid)

	if err != nil {
		r.l.Error().Err(err).Send()
		return nil, service.NewDbError(err)
	}
	defer rows.Close()

	t := service.Transaction{}
	if ok := rows.Next(); !ok {
		return nil, service.ErrTransactionNotFound
	}

	err = rows.Scan(&t.ID, &t.RefNo, &t.Amount, &t.Description, &t.Labels,
		&t.Fingerprint, &t.Created, &t.OldBalance, &t.NewBalance)
	if err != nil {
		r.l.Error().Err(err).Send()
		return nil, service.NewDbError(err)
	}

	return &t, nil
}
