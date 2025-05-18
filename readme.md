# Snake Game in Go

A multiplayer-ready terminal-based Snake game implementation using the tcell library, featuring a client-server architecture.

## Features

- Client-server architecture for multiplayer support
- Classic snake gameplay with terminal graphics
- Arrow keys and WASD controls
- Real-time movement and collision detection
- Score tracking
- Terminal-based UI with Unicode characters
- Sound effects for game events
- Pause and restart functionality

## Prerequisites

- Go 1.24 or higher
- Terminal with Unicode support
- Audio support for sound effects

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/go_snake.git

# Change to project directory
cd go_snake

# Install server dependencies
cd snake_server
go mod tidy

# Install client dependencies
cd ../snake_client
go mod tidy
```

## Game Controls

- **Arrow Keys** or **WASD**: Control snake direction
- **ESC/Ctrl+C/q**: Quit game
- **P**: Pause/Unpause game
- **R**: Restart game
- **C**: Clear screen

## Game Elements

| Character | Description |
|-----------|-------------|
| `@` | Snake head |
| `O` | Snake body |
| `●` | Food |
| `║`, `═`, `╔`, `╗`, `╚`, `╝` | Game borders |

## Technical Details

### Game Configuration
- Board Size: 70x20
- Default Speed: 200ms per frame
- Initial Snake Length: 5 segments
- Default Server Port: 3000

### Game Mechanics
- Snake grows when eating food
- Game ends on wall or self-collision
- Cannot reverse direction directly
- Score increases with each food item collected
- Pause/Resume functionality
- Game restart capability

### Display Features
- Real-time score display
- Game status information
- Game pause status indicator
- Unicode-based game elements
- Sound effects for game events

### Network Features
- TCP-based client-server communication
- Binary protocol for game commands
- Real-time input synchronization
- Connection management for multiple clients

## Project Structure

```
go_snake/
├── snake_server/
│   ├── main.go               // Server implementation and game logic
│   ├── types.go              // Type definitions and constants
│   ├── connection_manager.go // Network connection handling
│   ├── player.go            // Sound player implementation
│   ├── command.go           // Game commands definitions
│   └── assets/              // Game sound effects
│       ├── game_over_beep.mp3
│       └── game_start_beep.mp3
├── snake_client/
│   └── main.go              // Client implementation
└── README.md                // Documentation
```

## Development

The game uses the following key components:

- `tcell` for terminal graphics and input handling
- `beep` library for sound effects playback
- Goroutines for concurrent game state management
- TCP sockets for network communication
- Mutex for thread-safe state changes
- Channels for game termination and event handling

## Future Improvements

- [ ] Multiple difficulty levels
- [ ] Multiplayer scoreboard
- [ ] Add a countdown timer for game start
- [ ] Enhanced game over screen
- [ ] Network game state synchronization
- [ ] Player nicknames and profiles
- [ ] Spectator mode
- [ ] Game replay functionality

## Building and Running

First, make sure you're in the project root directory:

```bash
# Build and run the server
cd snake_server
go build
./snake_server

# In a separate terminal, build and run the client
cd snake_client
go build
./snake_client
```

The server will start listening on port 3000 by default. You can run multiple clients to test the multiplayer functionality, though currently only one player can control the snake at a time.

### Development Setup

For development, you can use the following commands:

```bash
# Run server in development mode
cd snake_server
go run .

# Run client in development mode
cd snake_client
go run .
```

## License
MIT License requires a more detailed notice. Here's the proper license section:

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Note: You should also include a separate LICENSE file in your repository containing the full MIT License text.
[MIT]