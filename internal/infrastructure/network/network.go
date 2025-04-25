package network

import (
	"bufio"
	"log"
	"net"
	"strings"
	"sync"
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
	mu          sync.Mutex // for thread safety
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
			log.Printf("Error reading from %s: %v", player.Username, err)
			s.mu.Lock()
			delete(s.players, player.Username)
			s.mu.Unlock()
			return
		}
		username = strings.TrimSpace(username)
		if _, err := s.userRepo.FindByUsername(username); err == nil {
			types.SendMessage(player, "Username already taken. Please choose another one:")
		} else {
			player.Username = username
			u := user.NewUser(username)
			s.userRepo.Save(u)
			s.mu.Lock()
			s.players[username] = player
			s.mu.Unlock()
			types.SendMessage(player, "Welcome, "+username)
			types.SendMessage(player, "Commands: join <two-player|ai>, move <1-9>, leaderboard, exit")
			break
		}
	}

	// Command loop
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading from %s: %v", player.Username, err)
			s.mu.Lock()
			delete(s.players, player.Username)
			s.mu.Unlock()
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
			if err.Error() == "exit requested" {
				return // clean exit
			}
			types.SendMessage(player, "Error: "+err.Error())
		}
	}
}

func (s *TCPServer) AddPlayerToGame(gameID string, player *types.Player) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gamePlayers[gameID] = append(s.gamePlayers[gameID], player)
}

func (s *TCPServer) BroadcastToGame(gameID string, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, player := range s.gamePlayers[gameID] {
		types.SendMessage(player, message)
	}
}

func (s *TCPServer) GetPlayer(username string) *types.Player {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.players[username]
}

func (s *TCPServer) GetPlayers() map[string]*types.Player {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.players
}

func (s *TCPServer) ExitPlayer(player *types.Player) {
	s.mu.Lock()
	log.Printf("Exiting player: %s", player.Username)
	delete(s.players, player.Username)

	var remainingPlayers []*types.Player
	if player.GameID != "" {
		// Player was in a game
		gameID := player.GameID
		log.Printf("%s was in game %s", player.Username, gameID)
		gamePlayers := s.gamePlayers[gameID]
		for i, p := range gamePlayers {
			if p.Username == player.Username {
				s.gamePlayers[gameID] = append(gamePlayers[:i], gamePlayers[i+1:]...)
				break
			}
		}
		remainingPlayers := s.gamePlayers[gameID]
		log.Printf("Remaining players in game %s: %v", gameID, remainingPlayers[0])

		if len(remainingPlayers) == 1 {
			remainingPlayer := remainingPlayers[0]
			log.Printf("Notifying %s and moving to waiting queue", remainingPlayer.Username)
			types.SendMessage(remainingPlayer, "Your opponent has left. Waiting for a new opponent...")
			s.matchmaking.AddToWaiting(remainingPlayer.Username)
			remainingPlayer.GameID = ""
			s.gameService.DeleteGame(gameID)
		}
		log.Printf("Deleting game %s from gamePlayers", gameID)
		delete(s.gamePlayers, gameID)
	}
	s.mu.Unlock()

	if remainingPlayers != nil && len(remainingPlayers) > 0 {
		log.Printf("Broadcasting to game %s: %s has left", player.GameID, player.Username)
		for _, p := range remainingPlayers {
			types.SendMessage(p, player.Username+" has left the game.")
		}
	}
}

func (s *TCPServer) EndGame(gameID string, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if players, ok := s.gamePlayers[gameID]; ok {
		for _, p := range players {
			types.SendMessage(p, message)
			p.GameID = ""
		}
		delete(s.gamePlayers, gameID)
	}
	s.gameService.DeleteGame(gameID)
}
