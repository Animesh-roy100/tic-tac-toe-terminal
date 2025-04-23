package application

import (
	"fmt"
	"sort"
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
	sort.Slice(users, func(i, j int) bool {
		return users[i].Score > users[j].Score
	})
	leaderboard := "Leaderboard:\n"
	for _, u := range users {
		leaderboard += fmt.Sprintf("%s: %d points, %d win streak\n", u.Username, u.Score, u.WinStreak)
	}
	return leaderboard, nil
}
