# Ginvader

A terminal-based Space Invaders clone written in Go.

## Description

Ginvader is a CLI game that brings the classic Space Invaders experience to your terminal. Fight against waves of alien invaders, dodge their bullets, and try to achieve the highest score possible!

## Features

- Terminal-based gameplay
- Colorful ASCII graphics
- Score tracking
- Enemy AI with shooting mechanics
- Game over screen with restart option
- Centered game display with border
- Cross-platform compatibility

## Requirements

- Go 1.21 or higher (only for building from source)
- Terminal with color support

## Installation

### Using Homebrew (recommended for macOS)

```bash
brew tap ouchi2501/tap
brew install ginvader
```

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/ouchi2501/ginvader.git
cd ginvader
```

2. Install dependencies:
```bash
go mod tidy
```

## How to Play

1. If installed via Homebrew:
```bash
ginvader
```

2. If building from source:
```bash
go run cmd/ginvader/main.go
```

3. Show help and game instructions:
```bash
ginvader -h
```

### Controls

- `←`, `→`: Move left/right
- `Space`: Shoot
- `ESC`: Quit game
- `Space`: Restart (when game over)

### Game Rules

- Destroy aliens to score points
- Each alien destroyed gives 100 points
- Avoid alien bullets
- Game over if:
  - An alien reaches your position
  - You get hit by an alien bullet

## Project Structure

```
ginvader/
├── cmd/
│   └── ginvader/
│       └── main.go       # Entry point
├── internal/
│   └── game/
│       ├── game.go       # Core game logic
│       ├── player.go     # Player implementation
│       └── enemy.go      # Enemy and bullet logic
├── go.mod               # Go module file
└── README.md           # This file
```

## Contributing

Feel free to submit issues, fork the repository, and create pull requests for any improvements.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by the classic Space Invaders arcade game
- Built with [tcell](https://github.com/gdamore/tcell) terminal library
