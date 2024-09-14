package cmd

import (
	"github.com/spf13/cobra"
	publicapi "go-wal/app/public_api"
)

func init() {
	rootCmd.AddCommand(publicCmd)
}

var publicCmd = &cobra.Command{
	Use:   "public-api",
	Short: "Serve public api application",
	Long: `A longer description that spans multiple lines and likely contains
			examples and usage of using your application. For example:

			Cobra is a CLI library for Go that empowers applications.
			This application is a tool to generate the needed files
			to quickly create a Cobra application.`,
	Run: func(_ *cobra.Command, _ []string) {
		publicapi.Launch()
	},
}
