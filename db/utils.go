package db

import (
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
