package handler

import (
	"errors"
	"strconv"
	"tic-tac-toe/internal/application"
	"tic-tac-toe/internal/types"
	"time"
)

var ErrExit = errors.New("exit requested")

type CommandHandler func(player *types.Player, args []string, gameService *application.GameService, leaderboard *application.LeaderboardService, matchmaking *application.MatchmakingService, server types.Server) error

var handlers = map[string]CommandHandler{
	"join":        JoinGameHandler,
	"move":        MakeMoveHandler,
	"leaderboard": LeaderboardHandler,
	"exit":        ExitHandler,
}

func HandleCommand(player *types.Player, command string, args []string, gameService *application.GameService, leaderboard *application.LeaderboardService, matchmaking *application.MatchmakingService, server types.Server) error {
	handler, ok := handlers[command]
	if !ok {
		return errors.New("unknown command")
	}
	return handler(player, args, gameService, leaderboard, matchmaking, server)
}

func JoinGameHandler(player *types.Player, args []string, gameService *application.GameService, _ *application.LeaderboardService, matchmaking *application.MatchmakingService, server types.Server) error {
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
			types.SendMessage(player, "Waiting for an opponent...")
			return nil
		}

		// Game has been created, add both players to the game
		g, err := gameService.FindGameByID(gameID)
		if err != nil {
			return err
		}

		for _, username := range g.Players {
			p := server.GetPlayer(username)
			if p != nil {
				p.GameID = gameID
				server.AddPlayerToGame(gameID, p)
			}
		}

		server.BroadcastToGame(gameID, "Game started. "+gameService.GetCurrentTurn(gameID)+"'s turn.")
		server.BroadcastToGame(gameID, gameService.GetBoard(gameID))
	} else if mode == "ai" {
		gameID, err = gameService.StartAIGame(player.Username)
		if err != nil {
			return err
		}
		player.GameID = gameID
		server.AddPlayerToGame(gameID, player)
		types.SendMessage(player, "Game started. Your turn.")
		types.SendMessage(player, gameService.GetBoard(gameID))
	} else {
		return errors.New("invalid mode")
	}
	return nil
}

func MakeMoveHandler(player *types.Player, args []string, gameService *application.GameService, _ *application.LeaderboardService, _ *application.MatchmakingService, server types.Server) error {
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

	if gameService.IsAIGame(player.GameID) {
		types.SendMessage(player, message)
	} else {
		server.BroadcastToGame(player.GameID, message)
	}

	// Notify next player if game continues
	g, err := gameService.FindGameByID(player.GameID)
	if err == nil && (g.Winner != "" || g.IsDraw) {
		server.EndGame(player.GameID, "Game has ended. You can start a new game.")
	} else if err == nil && !g.IsDraw {
		// Notify next player if game continues
		currentTurn := g.CurrentTurn
		if currentTurn == "" {
			return errors.New("error: current turn not set")
		}
		if gameService.IsAIGame(player.GameID) {
			types.SendMessage(player, "Your turn.")
		} else {
			server.BroadcastToGame(player.GameID, currentTurn+"'s turn.")
		}
	}

	return nil
}

func LeaderboardHandler(player *types.Player, args []string, _ *application.GameService, leaderboard *application.LeaderboardService, _ *application.MatchmakingService, server types.Server) error {
	leaderboardStr, err := leaderboard.GetLeaderboard()
	if err != nil {
		return err
	}
	types.SendMessage(player, leaderboardStr)
	return nil
}

func ExitHandler(player *types.Player, args []string, _ *application.GameService, _ *application.LeaderboardService, _ *application.MatchmakingService, server types.Server) error {
	server.ExitPlayer(player)

	types.SendMessage(player, "Goodbye!")

	// Close the connection
	// if player.Conn != nil {
	// 	player.Conn.Close()
	// }
	time.Sleep(100 * time.Millisecond)
	return ErrExit
}
