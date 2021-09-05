package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

func parseUrl(url string) *pgx.ConnConfig {
	config, err := pgx.ParseConfig(url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid database url: %v\n", err)
		os.Exit(1)
	}

	return config
}

func InitDb(url string) {
	config := parseUrl(url)
	db := config.Database
	config.Database = ""

	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	stmt := fmt.Sprintf("create database %s", db)
	_, err = conn.Exec(context.Background(), stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database is created")

	createSchema(url)
}

func DropDb(url string) {
	config := parseUrl(url)
	db := config.Database
	config.Database = ""

	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	stmt := fmt.Sprintf("Drop database %s", db)
	_, err = conn.Exec(context.Background(), stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database is dropped")
}

func createSchema(url string) {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), schemaSql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database schema is created")
}

const schemaSql = `
CREATE TABLE IF NOT EXISTS wallets (
	id			serial		PRIMARY KEY,
	labels		jsonb		not null,
	created 	timestamp	not null,
	externalid	varchar(50)	not null unique
);

CREATE TABLE IF NOT EXISTS wallet_transactions (
	id			uuid			PRIMARY KEY,
	wid			integer			not null,
	refno		integer			not null,
	amount		numeric(10,2)	not null,
	description	varchar(100)	not null,
	labels		jsonb			not null,
	fingerprint	varchar(50)		not null unique,
	created 	timestamp		not null,
	old_balance	numeric(10,2)	not null,
	new_balance numeric(10,2)	not null
);

CREATE TABLE IF NOT EXISTS wallet_balances (
	wid			integer			PRIMARY KEY,
	amount		numeric(10,2)	not null
);

CREATE UNIQUE INDEX ix_wid_refno ON wallet_transactions (wid, refno);
`
