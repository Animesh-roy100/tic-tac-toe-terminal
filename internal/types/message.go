package types

func SendMessage(player *Player, message string) {
	if player.Conn != nil {
		player.Conn.Write([]byte(message + "\n"))
	}
}
