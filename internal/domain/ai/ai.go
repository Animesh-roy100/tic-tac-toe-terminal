package ai

import (
	"math/rand"
	"tic-tac-toe/internal/domain/game"
	"time"
)

func AIMove(g *game.Game, aiUsername string) int {
	aiSymbol := "O"
	playerSymbol := "X"

	// Try to win
	if winMove := getWinningMove(g, aiSymbol); winMove != -1 {
		return winMove
	}

	// Block player's win
	if blockMove := getWinningMove(g, playerSymbol); blockMove != -1 {
		return blockMove
	}

	// Choose random empty cell
	emptyCells := []int{}
	for i, cell := range g.Board {
		if cell == " " {
			emptyCells = append(emptyCells, i)
		}
	}
	if len(emptyCells) > 0 {
		rand.Seed(time.Now().UnixNano())
		return emptyCells[rand.Intn(len(emptyCells))]
	}
	return -1 // Should not happen
}

func getWinningMove(g *game.Game, symbol string) int {
	for i := 0; i < 9; i++ {
		if g.Board[i] == " " {
			g.Board[i] = symbol
			if g.CheckWin(symbol) {
				g.Board[i] = " " // Undo
				return i
			}
			g.Board[i] = " " // Undo
		}
	}
	return -1
}
