# Snake Game in Go

A terminal-based Snake game implementation using the tcell library.

## Features

- Classic snake gameplay with terminal graphics
- Arrow keys and WASD controls
- Real-time movement and collision detection
- Score tracking
- Terminal-based UI with Unicode characters

## Prerequisites

- Go 1.16 or higher
- Terminal with Unicode support

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/go_snake.git

# Install dependencies
go get github.com/gdamore/tcell/v2
```

## Game Controls

- **Arrow Keys** or **WASD**: Control snake direction
- **ESC/Ctrl+C/q**: Quit game

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
- Initial Snake Length: 3 segments

### Game Mechanics
- Snake grows when eating food
- Game ends on wall or self-collision
- Cannot reverse direction directly
- Score increases with each food item collected

### Display Features
- Real-time score display
- Game status information
- Frame rate counter
- Position tracking for debugging

## Project Structure

```
go_snake/
├── main.go     // Game logic and rendering
├── types.go    // Type definitions and constants
└── README.md   // Documentation
```

## Development

The game uses the following key components:

- `tcell` for terminal graphics and input handling
- Goroutines for concurrent input processing
- Mutex for thread-safe direction changes
- Channels for game termination

## Future Improvements

- [x] Configurable board size and game speed
- [x] Configurable snake size
- [ ] Multiple difficulty levels
- [x] Pause functionality
- [x] Sound effects
- [ ] Add a count down timer for when the game is starting
- [ ] Game over screen enhancement

## Building and Running

```bash
# Build the game
go build

# Run the game
./go_snake
```

## License
MIT License requires a more detailed notice. Here's the proper license section:

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Note: You should also include a separate LICENSE file in your repository containing the full MIT License text.
[MIT]