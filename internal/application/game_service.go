package application

import (
	"errors"
	"fmt"
	"log"
	"tic-tac-toe/internal/domain/ai"
	"tic-tac-toe/internal/domain/game"
	"tic-tac-toe/internal/domain/user"
	"time"
)

// GameService manages game-related operations.
type GameService struct {
	gameRepo game.GameRepository
	userRepo user.UserRepository
}

func NewGameService(gameRepo game.GameRepository, userRepo user.UserRepository) *GameService {
	return &GameService{
		gameRepo: gameRepo,
		userRepo: userRepo,
	}
}

func (s *GameService) StartAIGame(username string) (string, error) {
	aiUsername := "AI"
	gameID := fmt.Sprintf("game-%d", time.Now().UnixNano())
	g := game.NewGame(gameID, []string{username, aiUsername}, true)
	if err := s.gameRepo.Save(g); err != nil {
		return "", err
	}
	log.Printf("Started AI game: gameID=%s, player=%s", gameID, username)
	return gameID, nil
}

func (s *GameService) MakeMove(gameID, username string, position int) (string, string, string, error) {
	g, err := s.gameRepo.FindByID(gameID)
	if err != nil {
		log.Printf("MakeMove: gameID=%s not found", gameID)
		return "", "", "", err
	}
	log.Printf("MakeMove: gameID=%s, username=%s, position=%d, board=%v", gameID, username, position, g.Board)
	if err := g.MakeMove(username, position); err != nil {
		log.Printf("MakeMove: invalid move for gameID=%s, username=%s: %v", gameID, username, err)
		return "", "", "", err
	}
	result := ""
	bonusMsg := ""
	if g.Winner != "" {
		winner, err := s.userRepo.FindByUsername(g.Winner)
		if err != nil {
			log.Printf("MakeMove: winner %s not found", g.Winner)
			return "", "", "", err
		}
		isAIGame := g.IsAIGame && g.Winner != "AI"
		bonusMsg = winner.WinGame(isAIGame)
		result = fmt.Sprintf("%s wins!", g.Winner)
		if !g.IsAIGame {
			loserUsername := g.Players[0]
			if loserUsername == g.Winner {
				loserUsername = g.Players[1]
			}
			loser, err := s.userRepo.FindByUsername(loserUsername)
			if err != nil {
				log.Printf("MakeMove: loser %s not found", loserUsername)
				return "", "", "", err
			}
			loser.LoseGame()
			s.userRepo.Save(loser)
		}
	} else if g.IsDraw {
		result = "It's a draw!"
		for _, player := range g.Players {
			if player != "AI" {
				u, err := s.userRepo.FindByUsername(player)
				if err != nil {
					log.Printf("MakeMove: player %s not found", player)
					return "", "", "", err
				}
				u.DrawGame()
				s.userRepo.Save(u)
			}
		}
	} else if g.IsAIGame && g.CurrentTurn == "AI" {
		aiMove := ai.AIMove(g, "AI")
		if aiMove == -1 {
			log.Printf("MakeMove: AI failed to make a move for gameID=%s", gameID)
			return "", "", "", errors.New("AI failed to make a move")
		}
		if err := g.MakeMove("AI", aiMove); err != nil {
			log.Printf("MakeMove: AI move failed for gameID=%s: %v", gameID, err)
			return "", "", "", err
		}
		result = fmt.Sprintf("AI chooses position %d", aiMove+1)
		if g.Winner == "AI" {
			result = "AI wins!"
			u, err := s.userRepo.FindByUsername(username)
			if err != nil {
				log.Printf("MakeMove: player %s not found", username)
				return "", "", "", err
			}
			u.LoseGame()
			s.userRepo.Save(u)
		} else if g.IsDraw {
			result = "It's a draw!"
			u, err := s.userRepo.FindByUsername(username)
			if err != nil {
				log.Printf("MakeMove: player %s not found", username)
				return "", "", "", err
			}
			u.DrawGame()
			s.userRepo.Save(u)
		}
	}
	if err := s.gameRepo.Save(g); err != nil {
		log.Printf("MakeMove: failed to save gameID=%s: %v", gameID, err)
		return "", "", "", err
	}
	log.Printf("MakeMove: gameID=%s, result=%s, bonusMsg=%s", gameID, result, bonusMsg)
	return g.DisplayString(), result, bonusMsg, nil
}
