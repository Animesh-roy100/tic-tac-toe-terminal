package types

import "net"

type Player struct {
	Conn     net.Conn
	Username string
	GameID   string
}

func NewPlayer(conn net.Conn) *Player {
	return &Player{Conn: conn}
}

type Server interface {
	AddPlayerToGame(gameID string, player *Player)
	BroadcastToGame(gameID string, message string)
	GetPlayer(username string) *Player
	GetPlayers() map[string]*Player
	ExitPlayer(player *Player)
	EndGame(gameID string, message string)
}
