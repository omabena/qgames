package transformer

type Game struct {
	Name       string
	Scores     map[string]int
	Kills      map[string]int
	TotalKills int
	Mods       map[string]int
}
