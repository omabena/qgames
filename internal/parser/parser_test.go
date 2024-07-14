package parser

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseInitGame(t *testing.T) {
	ctx := context.Background()
	matchLog := strings.NewReader(`
  0:00 InitGame: \sv_floodProtect\1\sv_maxPing\0\sv_minPing\0\sv_maxRate\10000\sv_minRate\0\sv_hostname\Code Miner Server\g_gametype\0\sv_privateClients\2\sv_maxclients\16\sv_allowDownload\0\dmflags\0\fraglimit\20\timelimit\15\g_maxGameClients\0\capturelimit\8\version\ioq3 1.36 linux-x86_64 Apr 12 2009\protocol\68\mapname\q3dm17\gamename\baseq3\g_needpass\0
 20:37 ShutdownGame:
`)
	matchChan := make(chan []Match)
	doneChan := make(chan struct{})
	readLog := New()
	go readLog.ReadLogGame(ctx, matchLog, matchChan, doneChan)
	entries := <-matchChan
	<-doneChan
	require.Empty(t, entries)
}

func TestParseClientUserinfoChanged(t *testing.T) {
	ctx := context.Background()
	matchLog := strings.NewReader(`
  0:00 InitGame: \sv_floodProtect\1\sv_maxPing\0\sv_minPing\0\sv_maxRate\10000\sv_minRate\0\sv_hostname\Code Miner Server\g_gametype\0\sv_privateClients\2\sv_maxclients\16\sv_allowDownload\0\dmflags\0\fraglimit\20\timelimit\15\g_maxGameClients\0\capturelimit\8\version\ioq3 1.36 linux-x86_64 Apr 12 2009\protocol\68\mapname\q3dm17\gamename\baseq3\g_needpass\0
 20:34 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\xian/default\hmodel\xian/default\g_redteam\\g_blueteam\\c1\4\c2\5\hc\100\w\0\l\0\tt\0\tl\0
 20:37 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0
 21:51 ClientUserinfoChanged: 3 n\Dono da Bola\t\0\model\sarge/krusade\hmodel\sarge/krusade\g_redteam\\g_blueteam\\c1\5\c2\5\hc\95\w\0\l\0\tt\0\tl\0
 21:53 ClientUserinfoChanged: 3 n\Mocinha\t\0\model\sarge\hmodel\sarge\g_redteam\\g_blueteam\\c1\4\c2\5\hc\95\w\0\l\0\tt\0\tl\0
 20:37 ShutdownGame:
`)
	matchChan := make(chan []Match)
	doneChan := make(chan struct{})
	readLog := New()
	go readLog.ReadLogGame(ctx, matchLog, matchChan, doneChan)

	entries := <-matchChan
	require.NotEmpty(t, entries)
	assert.IsType(t, entries[0], ClientUserinfoChanged{})
	assert.Equal(t, "Isgalamido", entries[0].(ClientUserinfoChanged).Player)

	assert.IsType(t, entries[1], ClientUserinfoChanged{})
	assert.Equal(t, "Isgalamido", entries[1].(ClientUserinfoChanged).Player)

	assert.IsType(t, entries[2], ClientUserinfoChanged{})
	assert.Equal(t, "Dono da Bola", entries[2].(ClientUserinfoChanged).Player)

	assert.IsType(t, entries[3], ClientUserinfoChanged{})
	assert.Equal(t, "Mocinha", entries[3].(ClientUserinfoChanged).Player)
}

func TestParseLogFromFile(t *testing.T) {
	logFile := "test.log"
	file, err := os.Open(logFile)
	require.NoError(t, err)
	defer file.Close()

	ctx := context.Background()

	matchChan := make(chan []Match)
	doneChan := make(chan struct{})
	readLog := New()
	go readLog.ReadLogGame(ctx, file, matchChan, doneChan)
	entries := <-matchChan

	require.NotEmpty(t, entries)
	assert.IsType(t, entries[0], ClientUserinfoChanged{})
	assert.IsType(t, entries[1], ClientUserinfoChanged{})
}

func TestParseKill(t *testing.T) {
	ctx := context.Background()
	matchLog := strings.NewReader(`
  0:00 InitGame: \sv_floodProtect\1\sv_maxPing\0\sv_minPing\0\sv_maxRate\10000\sv_minRate\0\sv_hostname\Code Miner Server\g_gametype\0\sv_privateClients\2\sv_maxclients\16\sv_allowDownload\0\dmflags\0\fraglimit\20\timelimit\15\g_maxGameClients\0\capturelimit\8\version\ioq3 1.36 linux-x86_64 Apr 12 2009\protocol\68\mapname\q3dm17\gamename\baseq3\g_needpass\0
 20:54 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT
 22:06 Kill: 2 3 7: Isgalamido killed Mocinha by MOD_ROCKET_SPLASH
  1:41 Kill: 1022 2 19: <world> killed Dono da Bola by MOD_FALLING
  1:42 Kill: 1022 2 19: Dono da Bola killed Arcius malus by MOD_FALLING
 20:37 ShutdownGame:
`)
	matchChan := make(chan []Match)
	doneChan := make(chan struct{})
	readLog := New()
	go readLog.ReadLogGame(ctx, matchLog, matchChan, doneChan)
	entries := <-matchChan

	require.NotEmpty(t, entries)
	assert.IsType(t, entries[0], Kill{})
	assert.Equal(t, entries[0].(Kill).Killer, "<world>")
	assert.Equal(t, entries[0].(Kill).Killed, "Isgalamido")
	assert.Equal(t, entries[0].(Kill).Mod, "MOD_TRIGGER_HURT")

	assert.Equal(t, entries[1].(Kill).Killer, "Isgalamido")
	assert.Equal(t, entries[1].(Kill).Killed, "Mocinha")
	assert.Equal(t, entries[1].(Kill).Mod, "MOD_ROCKET_SPLASH")

	assert.Equal(t, entries[2].(Kill).Killer, "<world>")
	assert.Equal(t, entries[2].(Kill).Killed, "Dono da Bola")
	assert.Equal(t, entries[2].(Kill).Mod, "MOD_FALLING")

	assert.Equal(t, entries[3].(Kill).Killer, "Dono da Bola")
	assert.Equal(t, entries[3].(Kill).Killed, "Arcius malus")
	assert.Equal(t, entries[3].(Kill).Mod, "MOD_FALLING")
}

func TestParseKillMultipleGames(t *testing.T) {
	ctx := context.Background()
	matchLog := strings.NewReader(`
  0:00 ------------------------------------------------------------
  0:00 InitGame: \sv_floodProtect\1\sv_maxPing\0\sv_minPing\0\sv_maxRate\10000\sv_minRate\0\sv_hostname\Code Miner Server\g_gametype\0\sv_privateClients\2\sv_maxclients\16\sv_allowDownload\0\dmflags\0\fraglimit\20\timelimit\15\g_maxGameClients\0\capturelimit\8\version\ioq3 1.36 linux-x86_64 Apr 12 2009\protocol\68\mapname\q3dm17\gamename\baseq3\g_needpass\0
 20:54 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT
 22:06 Kill: 2 3 7: Isgalamido killed Mocinha by MOD_ROCKET_SPLASH
  1:41 Kill: 1022 2 19: <world> killed Dono da Bola by MOD_FALLING
  1:42 Kill: 1022 2 19: Dono da Bola killed Arcius malus by MOD_FALLING
 20:37 ShutdownGame:
 20:37 ------------------------------------------------------------
 20:37 ------------------------------------------------------------
  0:00 InitGame: \sv_floodProtect\1\sv_maxPing\0\sv_minPing\0\sv_maxRate\10000\sv_minRate\0\sv_hostname\Code Miner Server\g_gametype\0\sv_privateClients\2\sv_maxclients\16\sv_allowDownload\0\dmflags\0\fraglimit\20\timelimit\15\g_maxGameClients\0\capturelimit\8\version\ioq3 1.36 linux-x86_64 Apr 12 2009\protocol\68\mapname\q3dm17\gamename\baseq3\g_needpass\0
 20:54 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT
 22:06 Kill: 2 3 7: Isgalamido killed Mocinha by MOD_ROCKET_SPLASH
  1:41 Kill: 1022 2 19: <world> killed Dono da Bola by MOD_FALLING
  1:42 Kill: 1022 2 19: Dono da Bola killed Arcius malus by MOD_FALLING
 20:37 ShutdownGame:
`)
	matchChan := make(chan []Match)
	doneChan := make(chan struct{})
	readLog := New()
	go readLog.ReadLogGame(ctx, matchLog, matchChan, doneChan)
	entries := <-matchChan

	require.NotEmpty(t, entries)
	assert.IsType(t, entries[0], Kill{})
	assert.Equal(t, entries[0].(Kill).Killer, "<world>")
	assert.Equal(t, entries[0].(Kill).Killed, "Isgalamido")
	assert.Equal(t, entries[0].(Kill).Mod, "MOD_TRIGGER_HURT")

	assert.Equal(t, entries[1].(Kill).Killer, "Isgalamido")
	assert.Equal(t, entries[1].(Kill).Killed, "Mocinha")
	assert.Equal(t, entries[1].(Kill).Mod, "MOD_ROCKET_SPLASH")

	assert.Equal(t, entries[2].(Kill).Killer, "<world>")
	assert.Equal(t, entries[2].(Kill).Killed, "Dono da Bola")
	assert.Equal(t, entries[2].(Kill).Mod, "MOD_FALLING")

	assert.Equal(t, entries[3].(Kill).Killer, "Dono da Bola")
	assert.Equal(t, entries[3].(Kill).Killed, "Arcius malus")
	assert.Equal(t, entries[3].(Kill).Mod, "MOD_FALLING")

	entries = <-matchChan

	require.NotEmpty(t, entries)
	assert.IsType(t, entries[0], Kill{})
	assert.Equal(t, entries[0].(Kill).Killer, "<world>")
	assert.Equal(t, entries[0].(Kill).Killed, "Isgalamido")
	assert.Equal(t, entries[0].(Kill).Mod, "MOD_TRIGGER_HURT")

	assert.Equal(t, entries[1].(Kill).Killer, "Isgalamido")
	assert.Equal(t, entries[1].(Kill).Killed, "Mocinha")
	assert.Equal(t, entries[1].(Kill).Mod, "MOD_ROCKET_SPLASH")

	assert.Equal(t, entries[2].(Kill).Killer, "<world>")
	assert.Equal(t, entries[2].(Kill).Killed, "Dono da Bola")
	assert.Equal(t, entries[2].(Kill).Mod, "MOD_FALLING")

	assert.Equal(t, entries[3].(Kill).Killer, "Dono da Bola")
	assert.Equal(t, entries[3].(Kill).Killed, "Arcius malus")
	assert.Equal(t, entries[3].(Kill).Mod, "MOD_FALLING")
}
