package cmd

import (
	"fmt"

	"github.com/polarbit/bluelabs-wallet/worker"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(workerCmd)
}

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Actor worker",
	Long:  `Starts a worker node for the cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting worker ...")

		worker.StartWorker()
	},
}
