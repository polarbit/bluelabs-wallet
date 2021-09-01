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
	Long:  `This is API to manage wallets`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting API ...")

		api.StartAPI()
	},
}
