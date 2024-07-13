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
	assert.NotEmpty(t, report)
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
	assert.NotEmpty(t, report)
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
	assert.NotEmpty(t, report)
}
