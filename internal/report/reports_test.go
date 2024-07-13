package report

import (
	"testing"

	"github.com/omabena/qgames/internal/transformer"
	"github.com/stretchr/testify/assert"
)

func TestMatches(t *testing.T) {
	games := []transformer.Game{
		{
			Name: "game1",
			Scores: map[string]int{
				"player1": 1,
				"player2": 2,
			},
			Kills: map[string]int{
				"player1": 1,
				"player2": 2,
			},
			TotalKills: 3,
		},
		{
			Name: "game1",
			Scores: map[string]int{
				"player1": 1,
				"player2": 2,
			},
			Kills: map[string]int{
				"player1": 1,
				"player2": 2,
			},
			TotalKills: 3,
		},
	}

	report := Matches(games)
	assert.NotEmpty(t, report.Games)
	assert.Equal(t, 2, len(report.Games))
	assert.Equal(t, 2, len(report.Games[0].Players))
	assert.Equal(t, 2, len(report.Games[1].Players))
	assert.Equal(t, 1, report.Games[0].Kills["player1"])
	assert.Equal(t, 2, report.Games[0].Kills["player2"])
	assert.Equal(t, 1, report.Games[1].Kills["player1"])
	assert.Equal(t, 2, report.Games[1].Kills["player2"])
}

func TestPlayersRanking(t *testing.T) {
	games := []transformer.Game{
		{
			Name: "game1",
			Scores: map[string]int{
				"player1": 1,
				"player2": 2,
			},
			Kills: map[string]int{
				"player1": 1,
				"player2": 2,
			},
			TotalKills: 3,
		},
		{
			Name: "game1",
			Scores: map[string]int{
				"player1": 1,
				"player2": 2,
			},
			Kills: map[string]int{
				"player1": 1,
				"player2": 2,
			},
			TotalKills: 3,
		},
	}

	report := PlayersRanking(games)
	assert.NotEmpty(t, report.Ranking)
	assert.Equal(t, 2, len(report.Ranking))
	assert.Equal(t, 2, report.Ranking["player1"])
	assert.Equal(t, 4, report.Ranking["player2"])
}

func TestDeathMod(t *testing.T) {
	games := []transformer.Game{
		{
			Name: "game1",
			Scores: map[string]int{
				"player1": 1,
				"player2": 2,
			},
			Kills: map[string]int{
				"player1": 1,
				"player2": 2,
			},
			TotalKills: 3,
			Mods: map[string]int{
				"MOD1": 1,
				"MOD2": 2,
			},
		},
		{
			Name: "game1",
			Scores: map[string]int{
				"player1": 1,
				"player2": 2,
			},
			Kills: map[string]int{
				"player1": 1,
				"player2": 2,
			},
			TotalKills: 3,
			Mods: map[string]int{
				"MOD1": 1,
				"MOD2": 2,
			},
		},
	}

	report := DeathMod(games)
	assert.NotEmpty(t, report.DeathMode)
	assert.Equal(t, 2, len(report.DeathMode))
	assert.Equal(t, 1, report.DeathMode[0].Count["MOD1"])
	assert.Equal(t, 2, report.DeathMode[0].Count["MOD2"])
}
