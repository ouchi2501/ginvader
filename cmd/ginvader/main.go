package main

import (
	"flag"
	"fmt"
	"os"

	"ginvader/internal/game"
)

const usage = `
Space Invaders CLI Game

Usage:
  ginvader [options]

Options:
  -h, --help     Show this help message

Controls:
  ←, →          Move left/right
  Space         Shoot
  ESC           Quit game

Game Rules:
  - Destroy all invaders to advance
  - Each invader destroyed gives 100 points
  - Game over if invaders reach your position
  - Game over if you get hit by invader's bullet
  - Press Space to restart after game over
`

func main() {
	help := flag.Bool("help", false, "Show usage information")
	flag.BoolVar(help, "h", false, "Show usage information")
	
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	g, err := game.NewGame()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing game: %v\n", err)
		os.Exit(1)
	}
	defer g.Close()

	g.Run()
}
