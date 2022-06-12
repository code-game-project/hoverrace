package hoverrace

import (
	"fmt"
	"time"

	"github.com/code-game-project/go-server/cg"
)

type Game struct {
	cg *cg.Game
}

const targetFrameTime time.Duration = 1 * time.Second / 30

func NewGame(cgGame *cg.Game) *Game {
	game := &Game{
		cg: cgGame,
	}
	game.cg.OnPlayerJoined = game.onPlayerJoined
	game.cg.OnPlayerLeft = game.onPlayerLeft
	game.cg.OnPlayerSocketConnected = game.onPlayerSocketConnected
	return game
}

func (g *Game) Run() {
	for g.cg.Running() {
		frameStart := time.Now()
		for {
			event, ok := g.cg.NextEvent()
			if !ok {
				break
			}
			g.handleEvent(event.Player, event.Event)
		}
		time.Sleep(targetFrameTime - time.Now().Sub(frameStart))
	}
}

func (g *Game) onPlayerJoined(player *cg.Player) {

}

func (g *Game) onPlayerLeft(player *cg.Player) {

}

func (g *Game) onPlayerSocketConnected(player *cg.Player, socket *cg.Socket) {

}

func (g *Game) handleEvent(player *cg.Player, event cg.Event) {
	switch event.Name {
	default:
		player.Send(player.Id, cg.ErrorEvent, cg.ErrorEventData{
			Message: fmt.Sprintf("unexpected event: %s", event.Name),
		})
	}
}
