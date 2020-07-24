package server

import (
	"context"
	"fmt"

	"github.com/angelbarrera92/registry-auth-server/internal/configs"
	"github.com/angelbarrera92/registry-auth-server/internal/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var configPath string

func NewTokenServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "run registry token auth server",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error
			configPath, err = cmd.Flags().GetString("config")
			if err != nil || configPath == "" {
				logrus.WithField("Stage", "Load Server Config").Infoln("Config is not specialfied, will use default config file")
				return
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			cfg := configs.NewConfigs(configPath)
			fmt.Println(cfg)
			s := server.NewRegistryAuthServer(cfg)
			if s == nil {
				return
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if err := s.Run(ctx); err != nil {
				logrus.Panicf("RunServerFailed: %s", err.Error())
			}
			<-ctx.Done()
		},
	}

	initFlags(cmd)
	return cmd
}

func initFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&configPath, "config", "c", "", "--config/-c")
}
