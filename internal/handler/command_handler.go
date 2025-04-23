package handler

import (
	"errors"
	"strconv"
	"tic-tac-toe/internal/application"
)

type CommandHandler func(player *Player, args []string, gameService *application.GameService, leaderboard *application.LeaderboardService, matchmaking *application.MatchmakingService) error

var handlers = map[string]CommandHandler{
	"join":        JoinGameHandler,
	"move":        MakeMoveHandler,
	"leaderboard": LeaderboardHandler,
}

func HandleCommand(player *Player, command string, args []string, gameService *application.GameService, leaderboard *application.LeaderboardService, matchmaking *application.MatchmakingService) error {
	handler, ok := handlers[command]
	if !ok {
		return errors.New("unknown command")
	}
	return handler(player, args, gameService, leaderboard, matchmaking)
}

func JoinGameHandler(player *Player, args []string, gameService *application.GameService, _ *application.LeaderboardService, matchmaking *application.MatchmakingService) error {
	if len(args) < 1 {
		return errors.New("mode required: two-player or ai")
	}
	mode := args[0]
	var gameID string
	var err error
	if mode == "two-player" {
		gameID, err = matchmaking.JoinTwoPlayerGame(player.Username)
		if err != nil {
			return err
		}
		if gameID == "" {
			SendMessage(player, "Waiting for an opponent...")
			return nil
		}
	} else if mode == "ai" {
		gameID, err = gameService.StartAIGame(player.Username)
		if err != nil {
			return err
		}
	} else {
		return errors.New("invalid mode")
	}
	player.GameID = gameID
	SendMessage(player, "Game started. Your turn.")
	return nil
}

func MakeMoveHandler(player *Player, args []string, gameService *application.GameService, _ *application.LeaderboardService, _ *application.MatchmakingService) error {
	if len(args) < 1 {
		return errors.New("position required")
	}
	position, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid position")
	}
	if player.GameID == "" {
		return errors.New("not in a game")
	}
	board, result, bonusMsg, err := gameService.MakeMove(player.GameID, player.Username, position-1)
	if err != nil {
		return err
	}
	message := "Board:\n" + board
	if result != "" {
		message += "\n" + result
	}
	if bonusMsg != "" {
		message += "\n" + bonusMsg
	}
	SendMessage(player, message)
	return nil
}

func LeaderboardHandler(player *Player, args []string, _ *application.GameService, leaderboard *application.LeaderboardService, _ *application.MatchmakingService) error {
	leaderboardStr, err := leaderboard.GetLeaderboard()
	if err != nil {
		return err
	}
	SendMessage(player, leaderboardStr)
	return nil
}

func SendMessage(p *Player, msg string) {
	p.Conn.Write([]byte(msg + "\n"))
}
