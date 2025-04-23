package handler

import "net"

type Player struct {
	Conn     net.Conn
	Username string
	GameID   string
}

func NewPlayer(conn net.Conn) *Player {
	return &Player{Conn: conn}
}
