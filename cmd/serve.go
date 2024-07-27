package cmd

import (
	"github.com/spf13/cobra"
  "github.com/pnptcn/ventilator/service"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return service.NewHTTPS().Up()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
