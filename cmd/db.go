package cmd

import (
	"errors"

	"github.com/polarbit/bluelabs-wallet/db"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dbCmd)
}

var dbCmd = func() *cobra.Command {
	var address string
	var username string
	var password string
	var initdb string
	var dropdb string

	c := &cobra.Command{
		Use:   "db",
		Short: "Initializes or drops wallet database. See help for parameters.",
		Long:  `See bl-wallet db --help`,
		RunE: func(cmd *cobra.Command, args []string) error {

			c := &db.ConnInfo{Host: address, Username: username, Password: password}

			if initdb != "" && dropdb != "" {
				return errors.New("only one of --initdb or --dropdb  parameters should be provided")
			} else if initdb != "" {
				c.Database = initdb
				db.InitDb(c)
			} else if dropdb != "" {
				c.Database = dropdb
				db.DropDb(c)
			} else {
				return errors.New("either one of --initdb or --dropdb  parameters should be provided")
			}

			return nil
		},
	}

	c.Flags().StringVarP(&address, "addr", "a", "localhost:5432", "Database address")
	c.Flags().StringVarP(&username, "user", "u", "postgres", "Database username")
	c.Flags().StringVarP(&password, "pass", "p", "1234", "Database password")
	c.Flags().StringVar(&initdb, "initdb", "", "--initdb parameter should be provided with database name to initalize a wallet database")
	c.Flags().StringVar(&dropdb, "dropdb", "", "--dropdb parameter should be provided with database name to drop a wallet database")

	return c
}()
