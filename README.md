# **🧩 Terminal Tic-Tac-Toe Game (Multiplayer over TCP)**

A networked Tic-Tac-Toe game written in Go, playable via a terminal client like** `nc` **or** `telnet`**. Supports both Two-Player and AI modes, with a Leaderboard to track player scores. The game runs as a TCP server, allowing multiple clients to connect, play in real-time, and compete.

## **Features**

- **Two-Player Mode:** Compete against another player over the network.
- **AI Mode:** Play against a strategic computer opponent.
- **Leaderboard:** Tracks wins (2 points for multiplayer, 1 point for AI, bonus for streaks).
- **Real-Time Updates:** Live board and turn updates.
- **Unique Usernames:** Ensures no username conflicts.
- **Graceful Exit:** Players can leave mid-game; opponents are notified.
- **Thread Safety:** Concurrency-safe using mutex locks.

## **Prerequisites**

- **Go:** v1.16 or higher (Download Go)
- **Terminal Client:** `netcat` (`nc`) or `telnet`
- **OS:** macOS, Linux, or Windows (terminal compatible)

## **Getting Started**

### **1. Clone the Repository**

```bash
git clone <repository-url>
cd tic-tac-toe
```

Alternatively, manually copy the project files into a local directory.

### **2. Ensure Project Structure**

```
tic-tac-toe/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── application/
│   │   ├── game_service.go
│   │   ├── leaderboard_service.go
│   │   └── matchmaking_service.go
│   ├── domain/
│   │   ├── ai/
│   │   │   └── ai.go
│   │   ├── game/
│   │   │   ├── game.go
│   │   │   └── repository.go
│   │   └── user/
│   │       ├── user.go
│   │       └── repository.go
│   ├── handler/
│   │   ├── command_handler.go
│   │   └── player.go
│   ├── infrastructure/
│   │   ├── network/
│   │   │   └── tcp_server.go
│   │   └── repository/
│   │       ├── in_memory_game_repository.go
│   │       └── in_memory_user_repository.go
│   └── types/
│       └── player.go
├── go.mod
└── README.md
```

### **3. Build and Run the Server**

Navigate to the project root and run:

```bash
go run cmd/server/main.go
```

The server starts on port `5000` (default) and logs `Server started on :5000`.

### **4. Connect to the Server**

Open a terminal and connect using `netcat`:

```bash
nc localhost 5000
```

Alternatively, use `telnet`:

```bash
telnet localhost 5000
```

## **Gameplay Instructions**

### **Start a Client**

1. Run `nc localhost 5000` in a terminal.
2. Enter a unique username when prompted (e.g., `abc`).
3. You’ll see a welcome message and available commands:

```
Welcome to Tic Tac Toe!
Enter username: abc
Welcome, abc
Commands: join <two-player|ai>, move <1-9>, leaderboard, exit
```

### **Join a Game**

- **Two-Player Mode:**

  - Type: `join two-player`
  - If no opponent is available, you’ll see: `Waiting for an opponent...`
  - When another player joins, the game starts, showing the board and whose turn it is: `Game started.`

- **AI Mode:**

  - Type: `join ai`
  - The game starts immediately against the AI: `Game started. Your turn.`

**Make Moves**

When it’s your turn, enter `move <position>` where `<position>` is a number from 1 to 9, corresponding to the grid:

```
1 | 2 | 3
---------
4 | 5 | 6
---------
7 | 8 | 9
```

### **View Leaderboard**

Type: `leaderboard`

Displays player scores and win streaks.

### **Exit the Game**

Type: `exit`

- You’ll see: `Goodbye!`, and the connection closes.
- If in a two-player game, the other player is notified: `<username> has left the game. You can start a new game.`
