package hoverrace

import (
	"github.com/code-game-project/go-server/cg"
)

type Player struct {
	id       string
	username string

	ready bool

	finished bool
	place    int
	duration int64

	checkpoints []Position

	cg   *cg.Player
	game *Game
}

func (p *Player) reset() {
	p.ready = false
	p.finished = false
}
