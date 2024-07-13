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
	entries, error := parser.ParseLogMatch(ctx, file)
	require.NoError(t, error)
	require.NotEmpty(t, entries)

	transformer := New()
	err = transformer.Matches(ctx, entries)
	require.NoError(t, err)

	assert.Equal(t, 1, len(transformer.Games))
	assert.Equal(t, "game_1", transformer.Games[0].Name)
	assert.Equal(t, 15, transformer.Games[0].TotalKills)
	assert.Equal(t, 4, transformer.Games[0].Kills["Isgalamido"])
	assert.Equal(t, -6, transformer.Games[0].Scores["Isgalamido"])
	assert.Equal(t, -2, transformer.Games[0].Scores["Mocinha"])
	assert.Equal(t, -2, transformer.Games[0].Scores["Zeh"])
	assert.Equal(t, -1, transformer.Games[0].Scores["Dono da Bola"])
}
