package main

import (
	"os"

	"github.com/omabena/qgames/cmd"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {
	zap.L().Info("Starting qgames")
	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
