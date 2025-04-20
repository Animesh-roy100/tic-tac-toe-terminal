package game

type Game struct {
	ID          string
	Board       [9]string
	Players     []string
	CurrentTurn string
	Winner      string
	IsDraw      bool
	IsAIGame    bool
}

func NewGame(id string, players []string, isAIGame bool) *Game {
	return &Game{
		ID:          id,
		Board:       [9]string{"", "", "", "", "", "", "", "", ""},
		Players:     players,
		CurrentTurn: players[0],
		IsAIGame:    isAIGame,
	}
}

func (g *Game) MakeMove(player string, position int) error {
	return nil
}

func (g *Game) CheckWin(symbol string) bool {
	return false
}

func (g *Game) CheckDraw() bool {
	return false
}

func (g *Game) DisplayBoard() string {
	return ""
}
