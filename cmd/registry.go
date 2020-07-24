package cmd

import (
	"github.com/angelbarrera92/registry-auth-server/cmd/server"
	"github.com/spf13/cobra"
)

func Execute() {
	app := newRootCmd()
	app.Execute()
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "registry",
	}
	cmd.AddCommand(server.NewTokenServerCommand())
	return cmd
}
