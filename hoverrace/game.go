package hoverrace

import (
	"fmt"
	"time"

	"github.com/code-game-project/go-server/cg"
)

type Game struct {
	cg              *cg.Game
	players         map[string]*Player
	checkpoints     []Position
	finishedPlayers int
	running         bool
	startTime       time.Time
}

const targetFrameTime time.Duration = 1 * time.Second / 30

func NewGame(cgGame *cg.Game) *Game {
	game := &Game{
		cg:      cgGame,
		players: make(map[string]*Player),
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
			g.update()
		}
		time.Sleep(targetFrameTime - time.Now().Sub(frameStart))
	}
}

func (g *Game) onPlayerJoined(cgPlayer *cg.Player) {
	if g.running {
		cgPlayer.Send("server", cg.ErrorEvent, cg.ErrorEventData{
			Message: "the game has already begun",
		})
		cgPlayer.Leave()
		return
	}

	g.players[cgPlayer.Id] = &Player{
		id:       cgPlayer.Id,
		cg:       cgPlayer,
		username: cgPlayer.Username,
		game:     g,
	}
}

func (g *Game) onPlayerLeft(player *cg.Player) {
	_, ok := g.players[player.Id]
	if !ok {
		return
	}

	delete(g.players, player.Id)

	if !g.running {
		for _, p := range g.players {
			if !p.ready {
				return
			}
		}
		g.start()
	} else {
		if g.finishedPlayers == len(g.players) {
			g.finish()
		}
	}

}

func (g *Game) onPlayerSocketConnected(player *cg.Player, socket *cg.Socket) {
	if !g.running {
		return
	}

	p := g.players[player.Id]
	p.cg.Send("server", CheckpointsEvent, CheckpointsEventData{
		Checkpoints: p.checkpoints,
	})

	socket.Send("server", StartEvent, StartEventData{})

	for _, player := range g.players {
		if player.finished {
			socket.Send(player.id, FinishedEvent, FinishedEventData{
				Place:    player.place,
				Duration: player.duration,
			})
		}
	}
}

func (g *Game) update() {
	if !g.running {
		return
	}
}

func (g *Game) handleEvent(player *cg.Player, event cg.Event) {
	switch event.Name {
	case ReadyEvent:
		g.handleReady(player.Id)
	default:
		player.Send(player.Id, cg.ErrorEvent, cg.ErrorEventData{
			Message: fmt.Sprintf("unexpected event: %s", event.Name),
		})
	}
}

func (g *Game) handleReady(playerId string) {
	player := g.players[playerId]
	if g.running {
		player.cg.Send(playerId, cg.ErrorEvent, cg.ErrorEventData{
			Message: "the game has already begun",
		})
		return
	}

	player.ready = true

	if len(g.players) < 2 {
		return
	}

	readyPlayers := make([]string, 0, len(g.players))
	for _, p := range g.players {
		if p.ready {
			readyPlayers = append(readyPlayers, p.id)
		}
	}

	g.cg.Send("server", ReadyPlayersEvent, ReadyPlayersEventData{
		Players: readyPlayers,
	})

	if len(readyPlayers) == len(g.players) {
		g.start()
	}
}

func (g *Game) start() {
	g.finishedPlayers = 0
	g.createCheckpoints()

	for _, player := range g.players {
		player.checkpoints = g.checkpoints
	}

	g.running = true
	g.startTime = time.Now()

	g.cg.Send("server", CheckpointsEvent, CheckpointsEventData{
		Checkpoints: g.checkpoints,
	})

	for countdown := 5; countdown > 0; countdown-- {
		g.cg.Send("server", CountdownEvent, CountdownEventData{
			Value: countdown,
		})
		time.Sleep(1 * time.Second)
	}

	g.cg.Send("server", StartEvent, StartEventData{})
}

func (g *Game) finish() {
	for _, player := range g.players {
		player.reset()
	}
	g.running = false
}

func (g *Game) createCheckpoints() {
	g.checkpoints = make([]Position, 5)
}
