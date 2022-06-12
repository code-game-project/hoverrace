package hoverrace

import "github.com/code-game-project/go-server/cg"

// The `start` event is sent when the race begins.
// The game begins once at least 2 players have joined and all players have sent the `ready` event.
const StartEvent cg.EventName = "start"

type StartEventData struct {
}

// The `countdown` counts down from 5. When the value reaches 0 a `start` event will be sent instead of the `countdown` event.
const CountdownEvent cg.EventName = "countdown"

type CountdownEventData struct {
	// The current value of the countdown (5-1).
	Value int `json:"value"`
}

// The `checkpoints` event contains all of the remaining checkpoints.
const CheckpointsEvent cg.EventName = "checkpoints"

type CheckpointsEventData struct {
	// The positions of all the remaining checkpoints.
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
	// The time in UNIX milliseconds when this event occured.
	Time int64 `json:"time"`
}

// The `finished` event is sent when a player crosses the finish line.
const FinishedEvent cg.EventName = "finished"

type FinishedEventData struct {
	// The place which the player has reached.
	Place int `json:"place"`
	// The amount of time in milliseconds the player needed to finish the race.
	Duration int64 `json:"duration"`
}

// A hovercraft is a circle with a diameter of 1 unit.
type Hovercraft struct {
	// The position of the center of the hovercraft.
	Pos Position `json:"pos"`
	// The current speed of the hovercraft.
	Speed float64 `json:"speed"`
	// The angle in degrees the hovercraft is facing (up = 0°).
	Angle float64 `json:"angle"`
	// The amount of reached checkpoints.
	Checkpoints int `json:"checkpoints"`
}

// One unit equals the width of the hovercrafts and checkpoints.
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
