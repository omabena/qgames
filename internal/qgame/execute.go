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
	ReadLogGame(ctx context.Context, reader io.Reader, matchChan chan<- []parser.Match, done chan<- bool)
}

type transformerLog interface {
	TransformToGame(ctx context.Context, matches []parser.Match)
	GetGames() []transformer.Game
}

func NewQGames(ctx context.Context, cfg *config.Config, readreadLogGame readLogGame, transformerLog transformerLog) *QGames {
	return &QGames{
		Config:         cfg,
		readLogGame:    readreadLogGame,
		transformerLog: transformerLog,
	}
}

func (qg QGames) Execute(ctx context.Context, doneReports chan<- bool) {
	file, err := os.Open(qg.Config.LogFilePath)
	if err != nil {
		zap.L().Error("could not open file", zap.Error(err))
		return
	}
	defer file.Close()

	matchChan := make(chan []parser.Match)
	doneReadLog := make(chan bool)

	go qg.readLogGame.ReadLogGame(ctx, file, matchChan, doneReadLog)
loop:
	for {
		select {
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
	go reportGamesGroup(&wg, gameGroupChan)
	go reportPlayersRanking(&wg, playersRankingChan)
	go reportDeathMod(&wg, deathModChan)

	gameGroupChan <- qg.transformerLog.GetGames()
	playersRankingChan <- qg.transformerLog.GetGames()
	deathModChan <- qg.transformerLog.GetGames()
	wg.Wait()
	zap.L().Info("finish executiong")
	doneReports <- true
	zap.L().Info("after done reports channel")
}

func reportGamesGroup(wg *sync.WaitGroup, gamesChan <-chan []transformer.Game) {
	defer wg.Done()
	zap.L().Info("Games group reporting started")
	games := <-gamesChan
	out := report.Matches(games)
	fmt.Println(out)
	zap.L().Info("finish games group repot")
}

func reportPlayersRanking(wg *sync.WaitGroup, gamesChan <-chan []transformer.Game) {
	defer wg.Done()
	zap.L().Info("Games players ranking started")
	games := <-gamesChan
	out := report.PlayersRanking(games)
	fmt.Println(out)
	zap.L().Info("finish players ranking repot")
}
func reportDeathMod(wg *sync.WaitGroup, gamesChan <-chan []transformer.Game) {
	defer wg.Done()
	zap.L().Info("Games death mod")
	games := <-gamesChan
	out := report.PlayersRanking(games)
	fmt.Println(out)
	zap.L().Info("finish death mod report")
}
