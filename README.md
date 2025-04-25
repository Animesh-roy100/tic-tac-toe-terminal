# Terminal Tic-Tac-Toe Game (Multiplayer over TCP)

A networked Tic-Tac-Toe game written in Go, playable via a terminal client (e.g., nc or telnet). Supports both two-player and AI modes, with a leaderboard to track player scores. The game runs as a TCP server, allowing multiple clients to connect, play games, and compete.
Features

## Features

- Two-Player Mode: Compete against another player over the network. Players take turns placing their symbols (X or O) on a 3x3 grid until one wins or the game ends in a draw.
- AI Mode: Play against a computer opponent that makes strategic moves to win or block your victory.
- Leaderboard: Tracks player scores based on wins (2 points for two-player, 1 point for AI wins, bonus points for win streaks).
- Real-Time Updates: Players receive immediate board updates and turn notifications during gameplay.
- Unique Usernames: Ensures each player has a unique username to prevent conflicts.
- Graceful Exit: Players can exit a game, notifying opponents and allowing them to start a new game.
- Thread Safety: Uses mutex locks to handle concurrent player actions safely.

## Prerequisites

Go: Version 1.16 or higher (download from golang.org).
Terminal Client: Tools like netcat (nc) or telnet for connecting to the server.
Operating System: Compatible with macOS, Linux, or Windows with a terminal.

## How to Run

Clone the Repository (if applicable):
git clone <repository-url>
cd tic-tac-toe

Or copy the project files to a local directory.

Ensure Project Structure:The project should have the following structure:
tic-tac-toe/
├── cmd/
│ └── server/
│ └── main.go
├── internal/
│ ├── application/
│ │ ├── game_service.go
│ │ ├── leaderboard_service.go
│ │ └── matchmaking_service.go
│ ├── domain/
│ │ ├── ai/
│ │ │ └── ai.go
│ │ ├── game/
│ │ │ ├── game.go
│ │ │ └── repository.go
│ │ └── user/
│ │ ├── user.go
│ │ └── repository.go
│ ├── handler/
│ │ ├── command_handler.go
│ │ └── player.go
│ ├── infrastructure/
│ │ ├── network/
│ │ │ └── tcp_server.go
│ │ └── repository/
│ │ ├── in_memory_game_repository.go
│ │ └── in_memory_user_repository.go
│ └── types/
│ └── player.go
├── go.mod
└── README.md

Build and Run the Server:Navigate to the project root and run:
go run cmd/server/main.go

The server starts on port 5000 (default) and logs: Server started on :5000.

Connect to the Server:Open a terminal and connect using netcat:
nc localhost 5000

Alternatively, use telnet:
telnet localhost 5000

## Gameplay Instructions

Start a Client:

Run nc localhost 5000 in a terminal.
Enter a unique username when prompted (e.g., abc).
You’ll see a welcome message and available commands:
Welcome to Tic Tac Toe!
Enter username: abc
Welcome, abc
Commands: join <two-player|ai>, move <1-9>, leaderboard, exit

Join a Game:

Two-Player Mode:
Type: join two-player
If no opponent is available, you’ll see: Waiting for an opponent...
When another player joins, the game starts, showing the board and whose turn it is:Game started.

AI Mode:
Type: join ai
The game starts immediately against the AI:Game started. Your turn.

Make Moves:

When it’s your turn, enter move <position> where <position> is a number from 1 to 9, corresponding to the grid:

View Leaderboard:

Type: leaderboard
Displays player scores and win streaks

Exit the Game:

Type: exit
You’ll see: Goodbye!, and the connection closes.
If in a two-player game, the other player is notified: <username> has left the game. You can start a new game.
You can reconnect and start a new game.
