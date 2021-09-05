package cmd

import (
	"github.com/polarbit/bluelabs-wallet/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Dumps current application configuration",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		config.Dump()
	},
}
