package hoverrace

import "github.com/code-game-project/go-server/cg"

// The `start` event is sent when the race begins.
// The game begins once at least 2 players have joined and all players have sent the `ready` event.
const StartEvent cg.EventName = "start"

type StartEventData struct {
	// The positions of all checkpoints.
	Checkpoints []Position `json:"checkpoints"`
}

// Send the `ready` event to the server when you think the game should begin.
const ReadyEvent cg.EventName = "ready"

type ReadyEventData struct {
}

// The `ready_players` event contains a list of all players which are ready.
const ReadyPlayersEvent cg.EventName = "ready_players"

type ReadyPlayersEventData struct {
	// A list of all ready players.
	Players []string `json:"players"`
}

// The `throttle` event allows you to change your throttle level and direction.
// **NOTE:** This values are targets. The hovercraft needs some time to reach the desired values.
const ThrottleEvent cg.EventName = "throttle"

type ThrottleEventData struct {
	// Throttle level between 0-1.
	Level float64 `json:"level"`
	// The angle in degrees the hovercraft should be facing (up = 0°).
	Angle float64 `json:"angle"`
}

// The `hovercraft` event tells all clients where all the hovercrafts are and how they are moving.
const HovercraftsEvent cg.EventName = "hovercrafts"

type HovercraftsEventData struct {
	// All hovercrafts mapped to their respective player IDs.
	Hovercrafts map[string]Hovercraft `json:"hovercrafts"`
	// The time in UNIX seconds when this event occured.
	Time int64 `json:"time"`
}

// The `next_checkpoint` tells you the position of the next checkpoint.
const NextCheckpointEvent cg.EventName = "next_checkpoint"

type NextCheckpointEventData struct {
	// The position of the next checkpoint.
	Pos Position `json:"pos"`
}

// The `finished` event is sent when a player crosses the finish line.
const FinishedEvent cg.EventName = "finished"

type FinishedEventData struct {
	// The ID of the player who reached the finish line.
	Player string `json:"player"`
	// The place which the player has reached.
	Place int `json:"place"`
}

type Hovercraft struct {
	// The position of the hovercraft.
	Pos Position `json:"pos"`
	// The current speed of the hovercraft.
	Speed float64 `json:"speed"`
	// The angle in degrees the hovercraft is facing (up = 0°).
	Angle float64 `json:"angle"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
