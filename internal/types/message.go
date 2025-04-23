package types

func SendMessage(p *Player, msg string) {
	p.Conn.Write([]byte(msg + "\n"))
}
