package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

type ConnInfo struct {
	Host     string
	Username string
	Password string
	Database string
}

func (c *ConnInfo) getUrl() string {
	var url string
	if c.Database == "" {
		url = fmt.Sprintf("postgresql://%s:%s@%s", c.Username, c.Password, c.Host)
	} else {
		url = fmt.Sprintf("postgresql://%s:%s@%s/%s", c.Username, c.Password, c.Host, c.Database)
	}

	_, err := pgx.ParseConfig(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid connection parameters: %v", c)
		os.Exit(1)
	}

	return url
}

func InitDb(c *ConnInfo) {
	var c2 = *c
	c2.Database = ""

	conn, err := pgx.Connect(context.Background(), c2.getUrl())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	stmt := fmt.Sprintf("create database %s", c.Database)
	_, err = conn.Exec(context.Background(), stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("database is created")

	createSchema(c)
}

func DropDb(c *ConnInfo) {
	var c2 = *c
	c2.Database = ""

	conn, err := pgx.Connect(context.Background(), c2.getUrl())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	stmt := fmt.Sprintf("drop database %s", c.Database)
	t, err := conn.Exec(context.Background(), stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("t: %v", t)
}

func createSchema(c *ConnInfo) {
	conn, err := pgx.Connect(context.Background(), c.getUrl())
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

	fmt.Println("database schema is created")
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
