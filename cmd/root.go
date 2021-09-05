package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var dbname string

func init() {
	// rootCmd.PersistentFlags().StringVar(&dbname, "database", "", `Overrides database name in config (db.deneme) or environment variable (BL_DB_DENEME)`)
	// viper.BindPFlag("db.database", rootCmd.PersistentFlags().Lookup("database"))
}

var rootCmd = &cobra.Command{
	Use:   "bl-wallet",
	Short: "BlueLabs wallet management system",
	Long: `This is customer wallet management system for the BlueLabs betting platform.
				  Intended to be used only internal services and apps.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
