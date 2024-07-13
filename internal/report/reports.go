package report

import (
	"go.uber.org/zap"
	"github.com/omabena/qgames/internal/transformer"
)

func Matches(g []transformer.Game) GamesReport {
	zap.L().Info("Matches Ranking report executing")
	games := GamesReport{
		Games: []Match{},
	}
	for _, game := range g {
		players := []string{}
		for player := range game.Scores {
			players = append(players, player)
		}
		games.Games = append(games.Games, Match{
			Name:       game.Name,
			Players:    players,
			Kills:      game.Kills,
			TotalKills: game.TotalKills,
		})
	}
	return games
}

func PlayersRanking(games []transformer.Game) PlayerReport {
	zap.L().Info("Players Ranking report executing")
	rankings := make(map[string]int)
	for _, game := range games {
		for player, score := range game.Scores {
			rankings[player] += score
		}
	}
	return PlayerReport{
		Ranking: rankings,
	}
}

func DeathMod(games []transformer.Game) DeathModeReport {
	zap.L().Info("DeathMod report executing")
	report := DeathModeReport{
		DeathMode: []DeathMode{},
	}
	for _, game := range games {
		count := make(map[string]int)
		for mod, value := range game.Mods {
			count[mod] = value
		}
		deathMode := DeathMode{
			Game:  game.Name,
			Count: count,
		}

		report.DeathMode = append(report.DeathMode, deathMode)
	}

	return report
}
