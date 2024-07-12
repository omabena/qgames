package parser

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
)

func ParseLogMatch(ctx context.Context, reader io.Reader) ([]Match, error) {
	var logEntries []Match
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		entry, err := parseLine(line)
		if err != nil {
			continue
		}
		logEntries = append(logEntries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return logEntries, nil
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
	dataKill := strings.Fields(parts[1])
	if len(dataKill) < 5 {
		return Kill{}, fmt.Errorf("invalid kill line: %s", rest)
	}
	return Kill{Timestamp: timestamp, Killer: dataKill[0], Killed: dataKill[2], Mod: dataKill[4]}, nil
}

func parseKeyValuePairs(s string) map[string]string {
	result := make(map[string]string)
	pairs := strings.Split(s, "\\")
	for i := 0; i < len(pairs)-1; i += 2 {
		result[pairs[i]] = pairs[i+1]
	}
	return result
}
