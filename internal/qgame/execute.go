package qgame

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"go.uber.org/zap"

	"github.com/omabena/qgames/internal/config"
	"github.com/omabena/qgames/internal/parser"
	"github.com/omabena/qgames/internal/report"
	"github.com/omabena/qgames/internal/transformer"
)

type QGames struct {
	Config         *config.Config
	readLogGame    readLogGame
	transformerLog transformerLog
}

type readLogGame interface {
	ReadLogGame(ctx context.Context, reader io.Reader, matchChan chan<- []parser.Match, done chan<- struct{})
}

type transformerLog interface {
	TransformToGame(ctx context.Context, matches []parser.Match)
	GetGames() []transformer.Game
}

func NewQGames(cfg *config.Config, readreadLogGame readLogGame, transformerLog transformerLog) *QGames {
	return &QGames{
		Config:         cfg,
		readLogGame:    readreadLogGame,
		transformerLog: transformerLog,
	}
}

func (qg QGames) Execute(ctx context.Context, doneReports chan<- struct{}) error {
	defer func() {
		doneReports <- struct{}{}
	}()
	file, err := os.Open(qg.Config.LogFilePath)
	if err != nil {
		zap.L().Error("could not open file", zap.Error(err))
		return err
	}
	defer file.Close()

	matchChan := make(chan []parser.Match)
	doneReadLog := make(chan struct{})

	go qg.readLogGame.ReadLogGame(ctx, file, matchChan, doneReadLog)
loop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-doneReadLog:
			break loop
		case entries := <-matchChan:
			qg.transformerLog.TransformToGame(ctx, entries)
		}
	}

	gameGroupChan := make(chan []transformer.Game)
	playersRankingChan := make(chan []transformer.Game)
	deathModChan := make(chan []transformer.Game)

	var wg sync.WaitGroup
	wg.Add(3)

	go reportGamesGroup(ctx, &wg, gameGroupChan)
	go reportPlayersRanking(ctx, &wg, playersRankingChan)
	go reportDeathMod(ctx, &wg, deathModChan)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case gameGroupChan <- qg.transformerLog.GetGames():
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case playersRankingChan <- qg.transformerLog.GetGames():
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case deathModChan <- qg.transformerLog.GetGames():
	}

	wg.Wait()
	return nil
}

func reportGamesGroup(ctx context.Context, wg *sync.WaitGroup, gamesChan <-chan []transformer.Game) {
	defer wg.Done()
	zap.L().Info("Games group reporting started")
	select {
	case <-ctx.Done():
		return
	case games := <-gamesChan:
		gamesReport := report.Matches(games)
		fmt.Println(gamesReport)
		zap.L().Info("finish games group repot")
	}
}

func reportPlayersRanking(ctx context.Context, wg *sync.WaitGroup, gamesChan <-chan []transformer.Game) {
	defer wg.Done()
	zap.L().Info("Games players ranking started")
	select {
	case <-ctx.Done():
		return
	case games := <-gamesChan:
		out := report.PlayersRanking(games)
		fmt.Println(out)
		zap.L().Info("finish players ranking repot")
	}
}
func reportDeathMod(ctx context.Context, wg *sync.WaitGroup, gamesChan <-chan []transformer.Game) {
	defer wg.Done()
	zap.L().Info("Games death mod")
	select {
	case <-ctx.Done():
		return
	case games := <-gamesChan:
		deathModReport := report.DeathMod(games)
		fmt.Println(deathModReport)
		zap.L().Info("finish death mod report")
	}
}
