package cmd

import (
	"github.com/spf13/cobra"
	"go-wal/app/consumer/wal"
)

func init() {
	rootCmd.AddCommand(walCaptureCmd)
}

var walCaptureCmd = &cobra.Command{
	Use:   "wal_capture",
	Short: "Capture WAL",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Do something
		wal.CaptureListen()
	},
}
