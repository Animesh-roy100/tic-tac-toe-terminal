package main

import (
	"log"
	"tic-tac-toe/internal/infrastructure/network"
	"tic-tac-toe/internal/infrastructure/repository"
)

func main() {

	userRepo := repository.NewInMemoryUserRepository()
	gameRepo := repository.NewInMemoryGameRepository()

	server := network.NewTCPServer(":5000", userRepo, gameRepo)

	log.Println("Server started on :5000")
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
