package cmd

import (
	"errors"

	"github.com/polarbit/bluelabs-wallet/config"
	"github.com/polarbit/bluelabs-wallet/db"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dbCmd)
}

var dbCmd = func() *cobra.Command {
	var initdb bool
	var dropdb bool

	c := &cobra.Command{
		Use:   "db",
		Short: "Initializes or drops wallet database. See help for parameters.",
		Long:  `See bl-wallet db --help`,
		RunE: func(cmd *cobra.Command, args []string) error {

			config := config.ReadConfig()

			if initdb && dropdb {
				return errors.New("only one of --initdb or --dropdb  parameters should be provided")
			} else if initdb {
				db.InitDb(config.Db)
			} else if dropdb {
				db.DropDb(config.Db)
			} else {
				return errors.New("either one of --initdb or --dropdb  parameters should be provided")
			}

			return nil
		},
	}

	c.Flags().BoolVar(&initdb, "initdb", false, "If given, a wallet database is initialized")
	c.Flags().BoolVar(&dropdb, "dropdb", false, "If given, a wallet database is dropped")

	return c
}()
