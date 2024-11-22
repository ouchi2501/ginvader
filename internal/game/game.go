package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

type GameState int

const (
	StateRunning GameState = iota
	StateGameOver
)

const (
	MaxGameWidth  = 80
	MaxGameHeight = 24
)

type Game struct {
	screen tcell.Screen
	player *Player
	enemies []*Enemy
	bullets []*Bullet
	score   int
	isRunning bool
	state GameState
	lastEnemyShot time.Time
	offsetX, offsetY int // Screen centering offsets
	gameWidth, gameHeight int
}

func NewGame() (*Game, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	screenWidth, screenHeight := screen.Size()
	gameWidth := min(screenWidth, MaxGameWidth)
	gameHeight := min(screenHeight, MaxGameHeight)
	
	// Calculate offsets to center the game
	offsetX := (screenWidth - gameWidth) / 2
	offsetY := (screenHeight - gameHeight) / 2

	player := NewPlayer(gameWidth/2, gameHeight-2)
	
	game := &Game{
		screen: screen,
		player: player,
		enemies: make([]*Enemy, 0),
		bullets: make([]*Bullet, 0),
		score: 0,
		isRunning: true,
		state: StateRunning,
		lastEnemyShot: time.Now(),
		offsetX: offsetX,
		offsetY: offsetY,
		gameWidth: gameWidth,
		gameHeight: gameHeight,
	}

	// Initialize enemies
	for i := 0; i < 8; i++ {
		for j := 0; j < 3; j++ {
			enemy := NewEnemy(10+i*5, 3+j*2)
			game.enemies = append(game.enemies, enemy)
		}
	}
	
	return game, nil
}

func (g *Game) translateX(x int) int {
	return x + g.offsetX
}

func (g *Game) translateY(y int) int {
	return y + g.offsetY
}

func (g *Game) setContent(x, y int, mainc rune, combc []rune, style tcell.Style) {
	if x >= 0 && x < g.gameWidth && y >= 0 && y < g.gameHeight {
		g.screen.SetContent(g.translateX(x), g.translateY(y), mainc, combc, style)
	}
}

func (g *Game) drawBorder() {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	
	// Draw horizontal borders
	for x := -1; x <= g.gameWidth; x++ {
		g.screen.SetContent(g.translateX(x), g.translateY(-1), '-', nil, style)
		g.screen.SetContent(g.translateX(x), g.translateY(g.gameHeight), '-', nil, style)
	}
	
	// Draw vertical borders
	for y := -1; y <= g.gameHeight; y++ {
		g.screen.SetContent(g.translateX(-1), g.translateY(y), '|', nil, style)
		g.screen.SetContent(g.translateX(g.gameWidth), g.translateY(y), '|', nil, style)
	}
	
	// Draw corners
	g.screen.SetContent(g.translateX(-1), g.translateY(-1), '+', nil, style)
	g.screen.SetContent(g.translateX(g.gameWidth), g.translateY(-1), '+', nil, style)
	g.screen.SetContent(g.translateX(-1), g.translateY(g.gameHeight), '+', nil, style)
	g.screen.SetContent(g.translateX(g.gameWidth), g.translateY(g.gameHeight), '+', nil, style)
}

func (g *Game) Run() {
	go g.handleInput()

	ticker := time.NewTicker(50 * time.Millisecond)
	for g.isRunning {
		select {
		case <-ticker.C:
			if g.state == StateRunning {
				g.update()
			}
			g.draw()
		}
	}
}

func (g *Game) handleInput() {
	for g.isRunning {
		switch ev := g.screen.PollEvent().(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				g.isRunning = false
				return
			case tcell.KeyLeft:
				if g.state == StateRunning {
					g.player.moveLeft()
				}
			case tcell.KeyRight:
				if g.state == StateRunning {
					g.player.moveRight()
				}
			case tcell.KeyRune:
				if ev.Rune() == ' ' {
					if g.state == StateRunning {
						g.player.shoot(g)
					} else if g.state == StateGameOver {
						g.resetGame()
					}
				}
			}
		}
	}
}

func (g *Game) update() {
	// Update enemies
	moveDown := false
	for _, enemy := range g.enemies {
		enemy.update()
		if enemy.alive {
			// Check if enemy hits the walls
			if enemy.x <= 0 || enemy.x >= g.gameWidth-1 {
				moveDown = true
			}
		}
	}

	if moveDown {
		for _, enemy := range g.enemies {
			enemy.reverseDirection()
		}
	}

	// Enemy shooting
	if time.Since(g.lastEnemyShot) > time.Second {
		g.enemyShoot()
		g.lastEnemyShot = time.Now()
	}

	// Update bullets
	for i := len(g.bullets) - 1; i >= 0; i-- {
		bullet := g.bullets[i]
		bullet.update()
		
		// Remove bullets that are out of screen
		if bullet.y < 0 || bullet.y > g.gameHeight {
			g.bullets = append(g.bullets[:i], g.bullets[i+1:]...)
		}
	}

	g.checkCollisions()
	g.checkGameOver()
}

func (g *Game) enemyShoot() {
	var aliveEnemies []*Enemy
	for _, enemy := range g.enemies {
		if enemy.alive {
			aliveEnemies = append(aliveEnemies, enemy)
		}
	}

	if len(aliveEnemies) > 0 {
		shooter := aliveEnemies[rand.Intn(len(aliveEnemies))]
		bullet := &Bullet{
			x: shooter.x,
			y: shooter.y + 1,
			isPlayerBullet: false,
		}
		g.bullets = append(g.bullets, bullet)
	}
}

func (g *Game) draw() {
	g.screen.Clear()
	g.drawBorder()

	if g.state == StateRunning {
		// Draw player
		g.player.draw(g, g.screen)

		// Draw enemies
		for _, enemy := range g.enemies {
			enemy.draw(g, g.screen)
		}

		// Draw bullets
		for _, bullet := range g.bullets {
			bullet.draw(g, g.screen)
		}

		// Draw score
		scoreStr := fmt.Sprintf("Score: %d", g.score)
		for i, r := range scoreStr {
			g.setContent(i, 0, r, nil, tcell.StyleDefault)
		}
	} else if g.state == StateGameOver {
		gameOver := "GAME OVER"
		score := fmt.Sprintf("Final Score: %d", g.score)
		restart := "Press SPACE to restart"
		
		style := tcell.StyleDefault.Foreground(tcell.ColorRed)
		
		// Draw game over text
		for i, r := range gameOver {
			g.setContent(g.gameWidth/2-len(gameOver)/2+i, g.gameHeight/2-1, r, nil, style)
		}
		
		// Draw score
		for i, r := range score {
			g.setContent(g.gameWidth/2-len(score)/2+i, g.gameHeight/2+1, r, nil, style)
		}
		
		// Draw restart instruction
		for i, r := range restart {
			g.setContent(g.gameWidth/2-len(restart)/2+i, g.gameHeight/2+3, r, nil, style)
		}
	}

	g.screen.Show()
}

func (g *Game) checkCollisions() {
	// Check player bullet collisions with enemies
	for i := len(g.bullets) - 1; i >= 0; i-- {
		bullet := g.bullets[i]
		if bullet.isPlayerBullet {
			for _, enemy := range g.enemies {
				if enemy.alive && bullet.x == enemy.x && bullet.y == enemy.y {
					enemy.alive = false
					g.bullets = append(g.bullets[:i], g.bullets[i+1:]...)
					g.score += 100
					break
				}
			}
		} else {
			// Check enemy bullet collision with player
			if bullet.x >= g.player.x-1 && bullet.x <= g.player.x+1 && bullet.y == g.player.y {
				g.state = StateGameOver
				break
			}
		}
	}
}

func (g *Game) checkGameOver() {
	// Check if any enemy reached the player's level
	for _, enemy := range g.enemies {
		if enemy.alive && enemy.y >= g.player.y {
			g.state = StateGameOver
			return
		}
	}

	// Check if all enemies are defeated
	allDefeated := true
	for _, enemy := range g.enemies {
		if enemy.alive {
			allDefeated = false
			break
		}
	}

	if allDefeated {
		g.resetGame()
	}
}

func (g *Game) resetGame() {
	g.player = NewPlayer(g.gameWidth/2, g.gameHeight-2)
	g.bullets = make([]*Bullet, 0)
	g.enemies = make([]*Enemy, 0)
	g.score = 0
	g.state = StateRunning

	// Initialize new enemies
	for i := 0; i < 8; i++ {
		for j := 0; j < 3; j++ {
			enemy := NewEnemy(10+i*5, 3+j*2)
			g.enemies = append(g.enemies, enemy)
		}
	}
}

func (g *Game) Close() {
	g.screen.Fini()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
