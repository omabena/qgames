package report

import (
	"fmt"
	"github.com/omabena/qgames/internal/transformer"
	"go.uber.org/zap"
	"sort"
	"strings"
)

func Matches(g []transformer.Game) string {
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
	sb := strings.Builder{}

	_, _ = sb.WriteString("***************\n")
	_, _ = sb.WriteString("***************\n")
	_, _ = sb.WriteString("Games report\n")
	for _, game := range games.Games {
		_, _ = sb.WriteString(fmt.Sprintf("Game: %s\n", game.Name))
		_, _ = sb.WriteString(fmt.Sprintf("Total Kills: %d\n", game.TotalKills))
		_, _ = sb.WriteString("Players: [")
		_, _ = sb.WriteString(strings.Join(game.Players, ", "))
		_, _ = sb.WriteString("]\n")

		_, _ = sb.WriteString("Kills {\n")
		for player, kill := range game.Kills {
			_, _ = sb.WriteString(fmt.Sprintf("\t%s: %d\n", player, kill))
		}
		_, _ = sb.WriteString("}\n")
	}

	return sb.String()
}

func PlayersRanking(games []transformer.Game) string {
	zap.L().Info("Players Ranking report executing")
	rankings := make(map[string]int)
	for _, game := range games {
		for player, score := range game.Scores {
			rankings[player] += score
		}
	}

	type kv struct {
		Key   string
		Value int
	}
	var kvSlice []kv
	for k, v := range rankings {
		kvSlice = append(kvSlice, kv{k, v})
	}

	sort.Slice(kvSlice, func(i, j int) bool {
		return kvSlice[i].Value > kvSlice[j].Value
	})

	var sortedPlayers []string
	for _, kv := range kvSlice {
		sortedPlayers = append(sortedPlayers, kv.Key)
	}

	sb := strings.Builder{}
	_, _ = sb.WriteString("***************\n")
	_, _ = sb.WriteString("***************\n")
	_, _ = sb.WriteString("Players Ranking\n")
	for _, player := range sortedPlayers {
		_, _ = sb.WriteString(fmt.Sprintf("%s: %d\n", player, rankings[player]))
	}
	return sb.String()
}

func DeathMod(games []transformer.Game) string {
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

	sb := strings.Builder{}
	_, _ = sb.WriteString("***************\n")
	_, _ = sb.WriteString("***************\n")
	_, _ = sb.WriteString("Death Mode report\n")
	for _, gameDeath := range report.DeathMode {
		game := fmt.Sprintf("Game: %s\n", gameDeath.Game)
		_, _ = sb.WriteString(game)

		for mod, count := range gameDeath.Count {
			modCount := fmt.Sprintf("%s, Count: %d\n", mod, count)
			_, _ = sb.WriteString(modCount)
		}
	}
	return sb.String()
}
