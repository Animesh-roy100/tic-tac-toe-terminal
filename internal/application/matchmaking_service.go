package application

import (
	"fmt"
	"sync"
	"tic-tac-toe/internal/domain/game"
	"time"
)

// MatchmakingService manages pairing players for two-player games.
type MatchmakingService struct {
	gameRepo game.GameRepository
	waiting  []string
	mu       sync.Mutex
}

func NewMatchmakingService(gameRepo game.GameRepository) *MatchmakingService {
	return &MatchmakingService{
		gameRepo: gameRepo,
		waiting:  make([]string, 0),
	}
}

func (s *MatchmakingService) JoinTwoPlayerGame(username string) (string, error) {
	s.mu.Lock()
	if len(s.waiting) > 0 {
		opponent := s.waiting[0]
		if opponent == "" {
			s.waiting = s.waiting[1:]
			s.mu.Unlock()
			return s.JoinTwoPlayerGame(username) // Retry
		}
		s.waiting = s.waiting[1:]
		s.mu.Unlock()

		gameID := fmt.Sprintf("game-%d", time.Now().UnixNano())
		g := game.NewGame(gameID, []string{opponent, username}, false)
		if err := s.gameRepo.Save(g); err != nil {
			return "", err
		}
		return gameID, nil
	} else {
		s.waiting = append(s.waiting, username)
		s.mu.Unlock()
		return "", nil // Waiting for opponent
	}
}

func (s *MatchmakingService) AddToWaiting(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, u := range s.waiting {
		if u == username {
			return // Already in waiting list
		}
	}
	s.waiting = append(s.waiting, username)
}

func (s *MatchmakingService) RemoveFromWaiting(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, u := range s.waiting {
		if u == username {
			s.waiting = append(s.waiting[:i], s.waiting[i+1:]...)
			break
		}
	}
}
