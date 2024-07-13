package parser

type Match interface { }

type InitGame struct {
	Timestamp string
}

type ClientUserinfoChanged struct {
	Timestamp string
	Player    string
}

type ShutdownGame struct {
	Timestamp string
}

type Kill struct {
	Timestamp string
	Killer    string
	Killed    string
	Mod       string
}

type ClientDisconnect struct {
	Timestamp string
	ClientID  int
}
