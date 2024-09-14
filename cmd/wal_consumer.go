package cmd

import (
	"github.com/spf13/cobra"
	"go-wal/app/consumer/wal"
)

func init() {
	rootCmd.AddCommand(walConsumer)
}

var walConsumer = &cobra.Command{
	Use:   "wal_consumer",
	Short: "Consume WAL",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		wal.ConsumerLaunch()
	},
}
