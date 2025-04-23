package application

import (
	"fmt"
	"tic-tac-toe/internal/domain/game"
	"time"
)

// MatchmakingService manages pairing players for two-player games.
type MatchmakingService struct {
	gameRepo game.GameRepository
	waiting  chan string
}

func NewMatchmakingService(gameRepo game.GameRepository) *MatchmakingService {
	return &MatchmakingService{
		gameRepo: gameRepo,
		waiting:  make(chan string, 100),
	}
}

func (s *MatchmakingService) JoinTwoPlayerGame(username string) (string, error) {
	select {
	case opponent := <-s.waiting:
		gameID := fmt.Sprintf("game-%d", time.Now().UnixNano())
		g := game.NewGame(gameID, []string{username, opponent}, false)
		if err := s.gameRepo.Save(g); err != nil {
			return "", err
		}
		return gameID, nil
	default:
		s.waiting <- username
		return "", nil // Waiting for opponent
	}
}
