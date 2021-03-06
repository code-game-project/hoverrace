package hoverrace

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/code-game-project/go-server/cg"
)

type Game struct {
	cg              *cg.Game
	players         map[string]*Player
	hovercrafts     map[string]Hovercraft
	checkpoints     []Vec
	finishLine      Vec
	finishedPlayers []FinishedPlayer
	running         bool
	startTime       time.Time
}

const targetFrameTime time.Duration = 1 * time.Second / 30

func NewGame(cgGame *cg.Game) *Game {
	game := &Game{
		cg:              cgGame,
		players:         make(map[string]*Player),
		hovercrafts:     make(map[string]Hovercraft),
		finishedPlayers: make([]FinishedPlayer, 0),
		checkpoints:     make([]Vec, 0, 10),
	}
	game.cg.OnPlayerJoined = game.onPlayerJoined
	game.cg.OnPlayerLeft = game.onPlayerLeft
	game.cg.OnPlayerSocketConnected = game.onPlayerSocketConnected
	game.cg.OnSpectatorConnected = game.onSpectatorConnected
	return game
}

func (g *Game) Run() {
	deltaTime := targetFrameTime
	for g.cg.Running() {
		frameStart := time.Now()
		for {
			event, ok := g.cg.NextEvent()
			if !ok {
				break
			}
			g.handleEvent(event.Player, event.Event)
		}
		g.update(deltaTime)
		time.Sleep(targetFrameTime - time.Now().Sub(frameStart))

		deltaTime = time.Now().Sub(frameStart)
		// waited for countdown
		if deltaTime >= 5*time.Second {
			deltaTime = targetFrameTime
		}
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
		finished: true,
	}

	g.hovercrafts[cgPlayer.Id] = Hovercraft{}

	g.positionHovercrafts()

	readyPlayers := make([]string, 0, len(g.players))
	for _, p := range g.players {
		if p.ready {
			readyPlayers = append(readyPlayers, p.id)
		}
	}
	if len(readyPlayers) > 0 {
		cgPlayer.Send("server", ReadyPlayersEvent, ReadyPlayersEventData{
			Players:  readyPlayers,
			Everyone: len(readyPlayers) == len(g.players),
		})
	}
}

func (g *Game) onPlayerLeft(player *cg.Player) {
	_, ok := g.players[player.Id]
	if !ok {
		return
	}

	if !g.running {
		g.positionHovercrafts()
		for _, p := range g.players {
			if !p.ready {
				return
			}
		}
		g.start()
	} else {
		if len(g.finishedPlayers) == len(g.players) {
			g.finish()
		}
	}

}

func (g *Game) onPlayerSocketConnected(player *cg.Player, socket *cg.Socket) {
	if len(g.checkpoints) > 0 {
		socket.Send("server", HovercraftsEvent, HovercraftsEventData{
			Hovercrafts: g.hovercrafts,
			Time:        time.Now().UnixMilli(),
		})

		p := g.players[player.Id]
		socket.Send("server", CheckpointsEvent, CheckpointsEventData{
			Checkpoints: p.checkpoints,
			FinishLine:  g.finishLine,
		})
	}

	readyPlayers := make([]string, 0, len(g.players))
	for _, p := range g.players {
		if p.ready {
			readyPlayers = append(readyPlayers, p.id)
		}
	}
	if len(readyPlayers) > 0 {
		socket.Send("server", ReadyPlayersEvent, ReadyPlayersEventData{
			Players:  readyPlayers,
			Everyone: len(readyPlayers) == len(g.players),
		})
	}

	if !g.running {
		return
	}

	socket.Send("server", InProgressEvent, InProgressEventData{})

	if len(g.finishedPlayers) > 0 {
		socket.Send("server", FinishedPlayersEvent, FinishedPlayersEventData{
			Players: g.finishedPlayers,
		})
	}
}

func (g *Game) onSpectatorConnected(socket *cg.Socket) {
	if len(g.checkpoints) > 0 {
		socket.Send("server", HovercraftsEvent, HovercraftsEventData{
			Hovercrafts: g.hovercrafts,
			Time:        time.Now().UnixMilli(),
		})

		socket.Send("server", CheckpointsEvent, CheckpointsEventData{
			Checkpoints: g.checkpoints,
			FinishLine:  g.finishLine,
		})
	}

	readyPlayers := make([]string, 0, len(g.players))
	for _, p := range g.players {
		if p.ready {
			readyPlayers = append(readyPlayers, p.id)
		}
	}
	if len(readyPlayers) > 0 {
		socket.Send("server", ReadyPlayersEvent, ReadyPlayersEventData{
			Players:  readyPlayers,
			Everyone: len(readyPlayers) == len(g.players),
		})
	}

	if !g.running {
		return
	}

	socket.Send("server", InProgressEvent, InProgressEventData{})

	if len(g.finishedPlayers) > 0 {
		socket.Send("server", FinishedPlayersEvent, FinishedPlayersEventData{
			Players: g.finishedPlayers,
		})
	}
}

func (g *Game) update(delta time.Duration) {
	for _, player := range g.players {
		player.update(delta)
	}

	g.cg.Send("server", HovercraftsEvent, HovercraftsEventData{
		Hovercrafts: g.hovercrafts,
		Time:        time.Now().UnixMilli(),
	})
}

func (g *Game) handleEvent(player *cg.Player, event cg.Event) {
	switch event.Name {
	case ReadyEvent:
		g.handleReady(player.Id)
	case ThrottleEvent:
		g.handleThrottle(player.Id, event)
	default:
		player.Send(player.Id, cg.ErrorEvent, cg.ErrorEventData{
			Message: fmt.Sprintf("unexpected event: %s", event.Name),
		})
	}
}

func (g *Game) handleReady(playerId string) {
	player := g.players[playerId]
	if !player.finished {
		player.cg.Send(playerId, cg.ErrorEvent, cg.ErrorEventData{
			Message: "the game has already begun",
		})
		return
	}

	player.ready = true

	readyPlayers := make([]string, 0, len(g.players))
	for _, p := range g.players {
		if p.ready {
			readyPlayers = append(readyPlayers, p.id)
		}
	}

	g.cg.Send(playerId, ReadyPlayersEvent, ReadyPlayersEventData{
		Players:  readyPlayers,
		Everyone: len(readyPlayers) == len(g.players),
	})

	if len(readyPlayers) == len(g.players) {
		g.start()
	}
}

func (g *Game) handleThrottle(playerId string, event cg.Event) {
	if !g.running {
		return
	}

	var data ThrottleEventData
	event.UnmarshalData(&data)

	if data.Level > 1 {
		data.Level = 1
	} else if data.Level < -1 {
		data.Level = -1
	}

	data.Angle = NormalizeAngle(data.Angle)

	player := g.players[playerId]
	player.targetThrottle = data.Level
	player.targetAngle = data.Angle
}

func (g *Game) start() {
	g.finishedPlayers = g.finishedPlayers[:0]
	g.createCheckpoints()

	for _, player := range g.players {
		player.ready = false
		player.vel = Vec{}
		player.acc = Vec{}
		player.angle = 0
		player.targetAngle = 0
		player.throttle = 0
		player.targetThrottle = 0

		player.finished = false
		player.checkpoints = make([]Vec, len(g.checkpoints))
		copy(player.checkpoints, g.checkpoints)
	}
	g.positionHovercrafts()

	g.update(0)

	g.cg.Send("server", CheckpointsEvent, CheckpointsEventData{
		Checkpoints: g.checkpoints,
		FinishLine:  g.finishLine,
	})

	for countdown := 5; countdown > 0; countdown-- {
		g.cg.Send("server", CountdownEvent, CountdownEventData{
			Value: countdown,
		})
		time.Sleep(1 * time.Second)
	}

	g.running = true
	g.startTime = time.Now()

	g.cg.Send("server", StartEvent, StartEventData{})
}

func (g *Game) positionHovercrafts() {
	x := 0.0
	i := 0
	for _, player := range g.players {
		player.pos = Vec{
			X: x,
			Y: 0,
		}
		if i%2 != 0 {
			player.pos.X = -x
		} else {
			x += 1.5
		}
		i++
	}
}

func (g *Game) finish() {
	g.running = false
	g.checkpoints = g.checkpoints[:0]
}

func (g *Game) createCheckpoints() {
	g.checkpoints = g.checkpoints[:cap(g.checkpoints)]
	for i := range g.checkpoints {
		g.checkpoints[i] = Vec{
			X: rand.Float64()*50*2 - 50,
			Y: rand.Float64()*50*2 - 50,
		}
	}

	g.finishLine = Vec{
		X: rand.Float64()*50*2 - 50,
		Y: rand.Float64()*50*2 - 50,
	}
}
