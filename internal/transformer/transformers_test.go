package transformer

import (
	"context"
	"os"
	"testing"

	"github.com/omabena/qgames/internal/parser"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatches(t *testing.T) {
	logFile := "fixtures/singlematch.log"
	file, err := os.Open(logFile)
	require.NoError(t, err)
	defer file.Close()

	ctx := context.Background()
	matchChan := make(chan []parser.Match)
	doneChan := make(chan struct{})
	readLog := parser.New()
	go readLog.ReadLogGame(ctx, file, matchChan, doneChan)
	entries := <-matchChan
	require.NotEmpty(t, entries)

	transformer := New()
	transformer.TransformToGame(ctx, entries)
	require.NoError(t, err)

	assert.Equal(t, 1, len(transformer.Games))
	assert.Equal(t, "game_1", transformer.Games[0].Name)
	assert.Equal(t, 15, transformer.Games[0].TotalKills)
	assert.Equal(t, 2, transformer.Games[0].Kills["Isgalamido"])
	assert.Equal(t, -8, transformer.Games[0].Scores["Isgalamido"])
	assert.Equal(t, -2, transformer.Games[0].Scores["Mocinha"])
	assert.Equal(t, -2, transformer.Games[0].Scores["Zeh"])
	assert.Equal(t, -1, transformer.Games[0].Scores["Dono da Bola"])
	assert.Equal(t, 9, transformer.Games[0].Mods["MOD_TRIGGER_HURT"])
	assert.Equal(t, 3, transformer.Games[0].Mods["MOD_ROCKET_SPLASH"])
	assert.Equal(t, 2, transformer.Games[0].Mods["MOD_FALLING"])
	assert.Equal(t, 1, transformer.Games[0].Mods["MOD_ROCKET"])
}

func TestMultiplesMatchesTransformation(t *testing.T) {
	logFile := "fixtures/multiplematches.log"
	file, err := os.Open(logFile)
	require.NoError(t, err)
	defer file.Close()

	ctx := context.Background()
	matchChan := make(chan []parser.Match)
	doneChan := make(chan struct{})
	readLog := parser.New()

	go readLog.ReadLogGame(ctx, file, matchChan, doneChan)

	transformer := New()
gameTransformer:
	for {
		select {
		case <-doneChan:
			break gameTransformer
		case entries := <-matchChan:
			require.NotEmpty(t, entries)
			transformer.TransformToGame(ctx, entries)
			require.NoError(t, err)
		}
	}
	assert.Equal(t, 2, len(transformer.Games))
	assert.Equal(t, "game_1", transformer.Games[0].Name)
	assert.Equal(t, "game_2", transformer.Games[1].Name)
}
