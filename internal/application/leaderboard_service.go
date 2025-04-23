package application

import (
	"fmt"
	"strings"
	"tic-tac-toe/internal/domain/user"
)

type LeaderboardService struct {
	userRepo user.UserRepository
}

func NewLeaderboardService(userRepo user.UserRepository) *LeaderboardService {
	return &LeaderboardService{userRepo: userRepo}
}

func (s *LeaderboardService) GetLeaderboard() (string, error) {
	users, err := s.userRepo.All()
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.WriteString("Leaderboard:\n")
	for _, u := range users {
		sb.WriteString(fmt.Sprintf("%s: %d points, %d win streak\n", u.Username, u.Score, u.WinStreak))
	}
	return sb.String(), nil
}
