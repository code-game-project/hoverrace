package hoverrace

import "github.com/code-game-project/go-server/cg"

type GameConfig struct {
	// The speed at which the thrust level reacts to user input. default = 1
	ThrottleSpeed float64 `json:"throttle_speed"`
	// The speed at which hovercrafts turn. default = 220
	TurnSpeed float64 `json:"turn_speed"`
	// The maximum acceleration of hovercrafts. default = 5
	MaxAcceleration float64 `json:"max_acceleration"`
	// The maximum velocity of hovercrafts. default = 20
	MaxVelocity float64 `json:"max_velocity"`
	// The amount of checkpoints per game. default = 10
	CheckpointCount int `json:"checkpoint_count"`
	// The time in milliseconds that a game is allowed to last. default = infinite
	Timeout int `json:"timeout"`
}

// Send the `ready` command to the server when you think the game should begin.
const ReadyCmd cg.CommandName = "ready"

type ReadyCmdData struct {
}

// The `control` command allows you to change your throttle level and direction.
// **NOTE:** These values are targets. The hovercraft needs some time to reach the desired values.
const ControlCmd cg.CommandName = "control"

type ControlCmdData struct {
	// Thrust level between -1 and 1.
	Thrust float64 `json:"thrust"`
	// The angle in degrees the hovercraft should be facing (right = 0°).
	Angle float64 `json:"angle"`
}

// The `start` event is sent when the race begins.
// The game begins once at least 2 players have joined and all players have sent the `ready` event.
const StartEvent cg.EventName = "start"

type StartEventData struct {
}

// The `in_progress` event is sent to sockets which connect to the game while it's running.
const InProgressEvent cg.EventName = "in_progress"

type InProgressEventData struct {
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
	Checkpoints []Vec `json:"checkpoints"`
	// The position of the finish line.
	FinishLine Vec `json:"finish_line"`
}

// The `ready_players` event contains a list of all players which are ready.
const ReadyPlayersEvent cg.EventName = "ready_players"

type ReadyPlayersEventData struct {
	// A list of all ready players.
	Players []string `json:"players"`
	// True if all players in the game are ready.
	Everyone bool `json:"everyone"`
}

// The `hovercraft` event tells all clients where all the hovercrafts are and how they are moving.
const HovercraftsEvent cg.EventName = "hovercrafts"

type HovercraftsEventData struct {
	// All hovercrafts mapped to their respective player IDs.
	Hovercrafts map[string]Hovercraft `json:"hovercrafts"`
	// The time in UNIX milliseconds when this event occured.
	Time int64 `json:"time"`
}

// The `finished_players` event contains a list of players that have finished the race.
const FinishedPlayersEvent cg.EventName = "finished_players"

type FinishedPlayersEventData struct {
	// A list of players that have finished the race sorted by their placement.
	Players []FinishedPlayer `json:"players"`
}

// The `game_over` event is sent when all players finished the game or the time runs out.
const GameOverEvent cg.EventName = "game_over"

type GameOverEventData struct {
	// The players that have finished the game before the time ran out.
	FinishedPlayers []FinishedPlayer `json:"finished_players"`
	// The IDs of the players that have not finished the game before the time ran out.
	UnfinishedPlayers []string `json:"unfinished_players"`
}

// A hovercraft is a circle with a diameter of 1 unit.
type Hovercraft struct {
	// The position of the center of the hovercraft.
	Pos Vec `json:"pos"`
	// The current velocity of the hovercraft.
	Velocity Vec `json:"velocity"`
	// The current thrust level of the hovercraft.
	Thrust float64 `json:"thrust"`
	// The angle in degrees the hovercraft is facing (up = 0°).
	Angle float64 `json:"angle"`
	// The amount of reached checkpoints.
	Checkpoints int `json:"checkpoints"`
}

// `finished_player` represents an entry in the final ranking.
type FinishedPlayer struct {
	// The ID of the player.
	Id string `json:"id"`
	// The place, the player has reached.
	Place int `json:"place"`
	// The amount of time in milliseconds the player needed to finish the race.
	Duration int64 `json:"duration"`
}

// One unit equals the width of the hovercrafts and checkpoints.
type Vec struct {
	// left to right
	X float64 `json:"x"`
	// bottom to top
	Y float64 `json:"y"`
}
