package network

import (
	"bufio"
	"log"
	"net"
	"strings"
	"tic-tac-toe/internal/application"
	"tic-tac-toe/internal/domain/game"
	"tic-tac-toe/internal/domain/user"
	"tic-tac-toe/internal/handler"
)

type TCPServer struct {
	listener    net.Listener
	userRepo    user.UserRepository
	gameService *application.GameService
	leaderboard *application.LeaderboardService
	matchmaking *application.MatchmakingService
}

func NewTCPServer(addr string, userRepo user.UserRepository, gameRepo game.GameRepository) *TCPServer {
	gameService := application.NewGameService(gameRepo, userRepo)
	leaderboard := application.NewLeaderboardService(userRepo)
	matchmaking := application.NewMatchmakingService(gameRepo)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}
	return &TCPServer{
		listener:    listener,
		userRepo:    userRepo,
		gameService: gameService,
		leaderboard: leaderboard,
		matchmaking: matchmaking,
	}
}

func (s *TCPServer) Start() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handleClient(conn)
	}
}

func (s *TCPServer) handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Get username
	player := handler.NewPlayer(conn)
	handler.SendMessage(player, "Welcome to Tic Tac Toe!")
	handler.SendMessage(player, "Enter username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading username: %v", err)
		return
	}
	username = strings.TrimSpace(username)
	player.Username = username
	u := user.NewUser(username)
	s.userRepo.Save(u)
	handler.SendMessage(player, "Welcome, "+username)
	handler.SendMessage(player, "Commands: join <two-player|ai>, move <1-9>, leaderboard")

	// Command loop
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading from %s: %v", username, err)
			return
		}
		message = strings.TrimSpace(message)
		parts := strings.Split(message, " ")
		if len(parts) == 0 {
			continue
		}
		command := parts[0]
		args := parts[1:]
		if err := handler.HandleCommand(player, command, args, s.gameService, s.leaderboard, s.matchmaking); err != nil {
			handler.SendMessage(player, "Error: "+err.Error())
		}
	}
}
