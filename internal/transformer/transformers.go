package transformer

import (
	"context"
	"fmt"
	"github.com/omabena/qgames/internal/parser"
	"go.uber.org/zap"
)

type Transformer struct {
	Games []Game
}

func New() *Transformer {
	return &Transformer{
		Games: []Game{},
	}
}

func (t *Transformer) TransformToGame(ctx context.Context, match []parser.Match) {
	game := &Game{
		Name:   fmt.Sprintf("game_%d", len(t.Games)+1),
		Kills:  make(map[string]int),
		Scores: make(map[string]int),
		Mods:   make(map[string]int),
	}
	for _, action := range match {
		switch v := action.(type) {
		case parser.InitGame:
			zap.L().Info("InitGame", zap.Any("action", v))
		case parser.ClientUserinfoChanged:
			if err := t.transfromClientUpdate(game, v); err != nil {
				zap.L().Error("error transforming client update", zap.Error(err))
			}
		case parser.Kill:
			if err := t.transformKill(game, v); err != nil {
				zap.L().Error("error transforming kill", zap.Error(err))
			}
		default:
			zap.L().Debug("ignore not process", zap.Any("action", v))
		}
	}
	t.Games = append(t.Games, *game)
}

func (t *Transformer) GetGames() []Game {
	return t.Games
}

func (t *Transformer) transformKill(game *Game, kill parser.Kill) error {
	zap.L().Info("Kill", zap.Any("action", kill))
	game.TotalKills++
	if kill.Killer != "<world>" {
		game.Kills[kill.Killer]++
		game.Scores[kill.Killer]++
	}
	game.Scores[kill.Killed]--
	game.Mods[kill.Mod]++
	return nil
}

func (t *Transformer) transfromClientUpdate(game *Game, client parser.ClientUserinfoChanged) error {
	zap.L().Info("ClientUserinfoChanged", zap.Any("action", client))
	if _, ok := game.Scores[client.Player]; !ok {
		game.Scores[client.Player] = 0
	}
	return nil
}
