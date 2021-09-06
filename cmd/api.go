package cmd

import (
	"fmt"

	"github.com/polarbit/bluelabs-wallet/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "REST API for wallet management",
	Long:  `This api enables to create wallets and transactions and query wallet and transactions`,
	Run: func(cmd *cobra.Command, args []string) {
		api.StartAPI()
		fmt.Println("wallet api is started")
	},
}
