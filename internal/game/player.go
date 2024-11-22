package game

import (
	"github.com/gdamore/tcell/v2"
)

type Player struct {
	x, y int
	width int
}

func NewPlayer(x, y int) *Player {
	return &Player{
		x: x,
		y: y,
		width: 3,
	}
}

func (p *Player) moveLeft() {
	if p.x > 0 {
		p.x--
	}
}

func (p *Player) moveRight() {
	// Screen width check will be done in the game loop
	p.x++
}

func (p *Player) shoot(g *Game) {
	bullet := &Bullet{
		x: p.x + p.width/2,
		y: p.y - 1,
		isPlayerBullet: true,
	}
	g.bullets = append(g.bullets, bullet)
}

func (p *Player) draw(g *Game, screen tcell.Screen) {
	style := tcell.StyleDefault.Foreground(tcell.ColorGreen)
	g.setContent(p.x, p.y, '^', nil, style)
	g.setContent(p.x-1, p.y, '-', nil, style)
	g.setContent(p.x+1, p.y, '-', nil, style)
}
