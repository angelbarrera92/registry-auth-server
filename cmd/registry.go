package cmd

import (
	"fmt"

	"github.com/angelbarrera92/registry-auth-server/cmd/server"
	"github.com/spf13/cobra"
)

func Execute() {
	app := newRootCmd()
	app.AddCommand(versionCmd)
	app.Execute()
}

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "registry",
	}
	cmd.AddCommand(server.NewTokenServerCommand())
	return cmd
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version information",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s\ncommit: %s\ndate: %s\n", version, commit, date)
	},
}
