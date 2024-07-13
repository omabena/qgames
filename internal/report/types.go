package report

type GamesReport struct {
	Games []Match
}

type Match struct {
	Name       string
	Players    []string
	Scores     map[string]int
	Kills      map[string]int
	TotalKills int
}

type PlayerReport struct {
	Ranking map[string]int
}

type DeathModeReport struct {
	DeathMode []DeathMode
}

type DeathMode struct {
	Game      string
	Count map[string]int
}
