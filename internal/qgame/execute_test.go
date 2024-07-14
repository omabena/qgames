package qgame

import (
	"context"
	"testing"

	"github.com/omabena/qgames/internal/config"
	"github.com/omabena/qgames/internal/parser"
	"github.com/omabena/qgames/internal/transformer"
)

func TestQGame(t *testing.T) {
	ctx := context.Background()
	cfg := config.NewConfig(ctx)
	cfg.LogFilePath = "fixtures/qgames.log"

	doneReports := make(chan struct{})

	parser := parser.New()
	transformer := transformer.New()
	qgames := NewQGames(&cfg, parser, transformer)
	go qgames.Execute(ctx, doneReports)
	t.Log("executing ...")
	<-doneReports
	t.Log("completed processing qgame log")
}
