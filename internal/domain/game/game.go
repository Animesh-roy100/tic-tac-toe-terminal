package game

import (
	"errors"
	"log"
	"strings"
)

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
		Board:       [9]string{" ", " ", " ", " ", " ", " ", " ", " ", " "},
		Players:     players,
		CurrentTurn: players[0],
		IsAIGame:    isAIGame,
	}
}

func (g *Game) MakeMove(player string, position int) error {
	log.Printf("MakeMove: player=%s, position=%d, CurrentTurn=%s, IsAIGame=%v", player, position, g.CurrentTurn, g.IsAIGame)
	if g.Winner != "" || g.IsDraw {
		log.Println("MakeMove: game is already over")
		return errors.New("game is already over")
	}
	// Allow moves in AI games regardless of CurrentTurn, as GameService handles AI turns
	if !g.IsAIGame && g.CurrentTurn != player {
		log.Println("MakeMove: not your turn")
		return errors.New("not your turn")
	}
	if position < 0 || position > 8 {
		log.Println("MakeMove: invalid position (out of bounds)")
		return errors.New("invalid position")
	}
	if g.Board[position] != " " {
		log.Printf("MakeMove: cell %d already taken (value: %s)", position, g.Board[position])
		return errors.New("cell already taken")
	}
	symbol := "X"
	if g.Players[0] != player {
		symbol = "O"
	}
	g.Board[position] = symbol
	log.Printf("MakeMove: placed %s at position %d", symbol, position)
	if g.CheckWin(symbol) {
		g.Winner = player
		log.Printf("MakeMove: %s wins", player)
	} else if g.CheckDraw() {
		g.IsDraw = true
		log.Println("MakeMove: game is a draw")
	} else {
		// Update CurrentTurn only for two-player games
		if !g.IsAIGame {
			g.CurrentTurn = g.Players[0]
			if g.CurrentTurn == player {
				g.CurrentTurn = g.Players[1]
			}
		}
		log.Printf("MakeMove: updated CurrentTurn to %s", g.CurrentTurn)
	}
	return nil
}

func (g *Game) CheckWin(symbol string) bool {
	winningCombos := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // Rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // Columns
		{0, 4, 8}, {2, 4, 6}, // Diagonals
	}
	for _, combo := range winningCombos {
		if g.Board[combo[0]] == symbol && g.Board[combo[1]] == symbol && g.Board[combo[2]] == symbol {
			return true
		}
	}
	return false
}

func (g *Game) CheckDraw() bool {
	for _, cell := range g.Board {
		if cell == " " {
			return false
		}
	}
	return true
}

// DisplayString formats the board for terminal output.
func (g *Game) DisplayString() string {
	var sb strings.Builder
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			index := row*3 + col
			sb.WriteString(g.Board[index])
			if col < 2 {
				sb.WriteString(" | ")
			}
		}
		sb.WriteString("\n")
		if row < 2 {
			sb.WriteString("-----------\n")
		}
	}
	return sb.String()
}
