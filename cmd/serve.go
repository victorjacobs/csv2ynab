package cmd

import (
	"github.com/spf13/cobra"
	"github.com/victorjacobs/csv2ynab/web"
)

var bindAddr string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serves a WebDAV server to send files to",
	RunE: func(cmd *cobra.Command, args []string) error {
		web.Serve(bindAddr, cfg)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVar(&bindAddr, "bind", ":8080", "address to listen on")
}
