package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/omabena/qgames/internal/config"
	"github.com/omabena/qgames/internal/qgame"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/omabena/qgames/internal/parser"
	"github.com/omabena/qgames/internal/transformer"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "qgames-parser",
		Short: "Commands to process a qgame file",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "Parse group data by match",
		RunE:  run,
	})
	return cmd
}

func run(_ *cobra.Command, _ []string) error {
	zap.L().Info("parse group data")
	ctx := cancelationSignal()
	cfg := config.NewConfig(ctx)
	zap.L().Sugar().Infof("config setup: %+v", cfg)
	zap.L().Info("log file path", zap.String("path", cfg.LogFilePath))

	parser := parser.New()
	transformer := transformer.New()

	doneReports := make(chan struct{})
	qgames := qgame.NewQGames(&cfg, parser, transformer)
	go qgames.Execute(ctx, doneReports)
	<-doneReports
	fmt.Println("\ncompleted processing qgame log")
	fmt.Println("\nwaiting for signal to exit ... press ctrl+c to exit")

	<-ctx.Done()
	return nil
}

func cancelationSignal() context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	go func() {
		<-c
		cancel()
	}()
	return ctx
}
