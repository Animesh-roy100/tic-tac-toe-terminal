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
	"tic-tac-toe/internal/types"
)

type TCPServer struct {
	listener    net.Listener
	userRepo    user.UserRepository
	gameService *application.GameService
	matchmaking *application.MatchmakingService
	leaderboard *application.LeaderboardService
	players     map[string]*types.Player
	gamePlayers map[string][]*types.Player
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
		players:     make(map[string]*types.Player),
		gamePlayers: make(map[string][]*types.Player),
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

	// Get Username from User
	player := types.NewPlayer(conn)
	types.SendMessage(player, "Welcome to Tic Tac Toe!\nEnter username: ")

	for {
		username, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading username: %v", err)
			return
		}
		username = strings.TrimSpace(username)
		if _, err := s.userRepo.FindByUsername(username); err == nil {
			types.SendMessage(player, "Username already taken. Please choose another one:")
		} else {
			player.Username = username
			u := user.NewUser(username)
			s.userRepo.Save(u)
			s.players[username] = player
			types.SendMessage(player, "Welcome, "+username)
			types.SendMessage(player, "Commands: join <two-player|ai>, move <1-9>, leaderboard")
			break
		}
	}

	// Command loop
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading from %s: %v", player.Username, err)
			delete(s.players, player.Username)
			return
		}
		message = strings.TrimSpace(message)
		parts := strings.Split(message, " ")
		if len(parts) == 0 {
			continue
		}
		command := parts[0]
		args := parts[1:]
		if err := handler.HandleCommand(player, command, args, s.gameService, s.leaderboard, s.matchmaking, s); err != nil {
			types.SendMessage(player, "Error: "+err.Error())
		}
	}
}

func (s *TCPServer) AddPlayerToGame(gameID string, player *types.Player) {
	s.gamePlayers[gameID] = append(s.gamePlayers[gameID], player)
}

func (s *TCPServer) BroadcastToGame(gameID string, message string) {
	for _, player := range s.gamePlayers[gameID] {
		types.SendMessage(player, message)
	}
}
