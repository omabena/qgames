package cmd

import (
	"context"

	"github.com/omabena/qgames/internal/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "qgames-parser",
		Short: "Commands to process a qgame file",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "group-data",
		Short: "Parse group data by match",
		RunE:  parse,
	})
	return cmd
}

func parse(_ *cobra.Command, _ []string) error {
	zap.L().Info("parse group data")
	ctx := context.Background()
	cfg := config.NewConfig(ctx)
	zap.L().Sugar().Infof("config setup: %+v", cfg)
	zap.L().Info("log file path", zap.String("path", cfg.LogFilePath))

	return nil
}
