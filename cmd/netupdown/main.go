package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/moli-xia/netupdown/internal/app"
	"github.com/moli-xia/netupdown/internal/config"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	root := &cobra.Command{Use: "netupdown", Short: "Personal app publishing and software distribution"}
	var configPath string
	serve := &cobra.Command{
		Use:   "serve",
		Short: "start the HTTP server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load(configPath)
			if err != nil {
				return err
			}
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			return app.Serve(ctx, cfg, version)
		},
	}
	serve.Flags().StringVarP(&configPath, "config", "c", "config.yaml", "configuration file")
	root.AddCommand(serve)

	var adminConfig, username string
	var passwordStdin bool
	adminCreate := &cobra.Command{
		Use:   "create",
		Short: "create an administrator",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load(adminConfig)
			if err != nil {
				return err
			}
			return app.CreateAdmin(cmd.Context(), cfg, username, passwordStdin)
		},
	}
	adminCreate.Flags().StringVarP(&adminConfig, "config", "c", "config.yaml", "configuration file")
	adminCreate.Flags().StringVar(&username, "username", "", "administrator username")
	adminCreate.Flags().BoolVar(&passwordStdin, "password-stdin", false, "read password from stdin")
	_ = adminCreate.MarkFlagRequired("username")
	admin := &cobra.Command{Use: "admin", Short: "administrator commands"}
	admin.AddCommand(adminCreate)
	root.AddCommand(admin)

	root.AddCommand(&cobra.Command{Use: "version", Run: func(*cobra.Command, []string) { fmt.Println(version) }})
	if err := root.ExecuteContext(context.Background()); err != nil {
		slog.Error("command failed", "error", err)
		os.Exit(1)
	}
}
