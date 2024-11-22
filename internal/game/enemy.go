package game

import (
	"github.com/gdamore/tcell/v2"
)

type Enemy struct {
	x, y    int
	dx      int
	alive   bool
}

func NewEnemy(x, y int) *Enemy {
	return &Enemy{
		x:      x,
		y:      y,
		dx:     1,
		alive:  true,
	}
}

func (e *Enemy) update() {
	if !e.alive {
		return
	}
	e.x += e.dx
}

func (e *Enemy) reverseDirection() {
	e.dx = -e.dx
	e.y++
}

func (e *Enemy) draw(g *Game, screen tcell.Screen) {
	if !e.alive {
		return
	}
	style := tcell.StyleDefault.Foreground(tcell.ColorRed)
	g.setContent(e.x, e.y, 'M', nil, style)
}

type Bullet struct {
	x, y           int
	isPlayerBullet bool
}

func (b *Bullet) update() {
	if b.isPlayerBullet {
		b.y--
	} else {
		b.y++
	}
}

func (b *Bullet) draw(g *Game, screen tcell.Screen) {
	style := tcell.StyleDefault.Foreground(tcell.ColorYellow)
	if b.isPlayerBullet {
		g.setContent(b.x, b.y, '|', nil, style)
	} else {
		g.setContent(b.x, b.y, '*', nil, style)
	}
}
