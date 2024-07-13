package parser

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
)

type ReadLog struct {
	Match []Match
}

func New() *ReadLog {
	return &ReadLog{
		Match: nil,
	}
}

func (r *ReadLog) ReadLogGame(ctx context.Context, reader io.Reader, matchChan chan<- []Match, done chan<- bool) {
	defer func() {
		done <- true
	}()
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		r.processMatch(matchChan, line)
	}
	if err := scanner.Err(); err != nil {
		return
	}
}

func (r *ReadLog) processMatch(matchChan chan<- []Match, line string) {
	entry, err := parseLine(line)
	if err != nil {
		return
	}
	switch entry.(type) {
	case InitGame:
		if r.Match == nil {
			r.Match = []Match{}
		}
	case ShutdownGame:
		// create a match and send it to the channel
		newMatch := make([]Match, len(r.Match))
		copy(newMatch, r.Match)
		matchChan <- newMatch
		r.Match = nil
	default:
		r.Match = append(r.Match, entry)
	}
}

func parseLine(line string) (Match, error) {
	initLine := strings.TrimLeft(line, " ")
	parts := strings.SplitN(initLine, " ", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid line: %s", line)
	}
	timestamp, rest := parts[0], parts[1]
	switch {
	case strings.HasPrefix(rest, "InitGame:"):
		return parseInitGame(timestamp)
	case strings.HasPrefix(rest, "ClientUserinfoChanged:"):
		return parseClientUserinfoChanged(timestamp, rest)
	case strings.HasPrefix(rest, "ShutdownGame:"):
		return parseShutdownGame(timestamp, rest)
	case strings.HasPrefix(rest, "Kill:"):
		return parseKill(timestamp, rest)
	default:
		return nil, fmt.Errorf("unknown log entry: %s", line)
	}
}

func parseInitGame(timestamp string) (InitGame, error) {
	return InitGame{Timestamp: timestamp}, nil
}

func parseClientUserinfoChanged(timestamp, rest string) (ClientUserinfoChanged, error) {
	parts := strings.SplitN(strings.TrimPrefix(rest, "ClientUserinfoChanged: "), " ", 2)
	info := parseKeyValuePairs(parts[1])
	player, ok := info["n"]
	if !ok {
		return ClientUserinfoChanged{}, fmt.Errorf("player name not found in client info: %s", parts[1])
	}
	return ClientUserinfoChanged{Timestamp: timestamp, Player: player}, nil
}

func parseShutdownGame(timestamp, _ string) (ShutdownGame, error) {
	return ShutdownGame{Timestamp: timestamp}, nil
}

func parseKill(timestamp, rest string) (Kill, error) {
	parts := strings.SplitN(strings.TrimPrefix(rest, "Kill: "), ": ", 2)
	dataKill := strings.SplitN(parts[1], "killed", 2)
	if len(dataKill) < 2 {
		return Kill{}, fmt.Errorf("invalid kill line: %s", rest)
	}
	killer := strings.TrimSuffix(dataKill[0], " ")
	dataKilled := strings.SplitN(strings.TrimPrefix(dataKill[1], " "), "by", 2)
	if len(dataKill) < 2 {
		return Kill{}, fmt.Errorf("invalid kill line: %s", rest)
	}
	killed := strings.TrimSuffix(dataKilled[0], " ")
	mod := strings.TrimPrefix(dataKilled[1], " ")
	return Kill{Timestamp: timestamp, Killer: killer, Killed: killed, Mod: mod}, nil
}

func parseKeyValuePairs(s string) map[string]string {
	result := make(map[string]string)
	pairs := strings.Split(s, "\\")
	for i := 0; i < len(pairs)-1; i += 2 {
		result[pairs[i]] = pairs[i+1]
	}
	return result
}
