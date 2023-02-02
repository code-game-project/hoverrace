package main

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/code-game-project/go-server/cg"
	"github.com/spf13/pflag"

	"github.com/code-game-project/hoverrace/hoverrace"
)

func main() {
	var port int
	pflag.IntVarP(&port, "port", "p", 0, "The network port of the game server.")
	pflag.Parse()

	if port == 0 {
		portStr, ok := os.LookupEnv("CG_PORT")
		if ok {
			port, _ = strconv.Atoi(portStr)
		}
	}

	if port == 0 {
		port = 8080
	}

	server := cg.NewServer("hoverrace", cg.ServerConfig{
		DisplayName:             "Hover Race",
		Description:             "Race against other hovercrafts from checkpoint to checkpoint.",
		Version:                 "0.4",
		RepositoryURL:           "https://github.com/code-game-project/hoverrace",
		WebsocketTimeout:        1 * time.Minute,
		MaxPlayersPerGame:       10,
		Port:                    port,
		CGEFilepath:             "events.cge",
		DeleteInactiveGameDelay: 30 * time.Minute,
		KickInactivePlayerDelay: 30 * time.Minute,
	})

	server.Run(func(cgGame *cg.Game, config json.RawMessage) {
		var gameConfig hoverrace.GameConfig
		err := json.Unmarshal(config, &gameConfig)
		if err != nil {
			cgGame.Log.Error("Failed to unmarshal game config: %s", err)
		}
		if gameConfig.ThrottleSpeed <= 0 {
			gameConfig.ThrottleSpeed = 1
		}
		if gameConfig.TurnSpeed <= 0 {
			gameConfig.TurnSpeed = 220
		}
		if gameConfig.MaxAcceleration <= 0 {
			gameConfig.MaxAcceleration = 5
		}
		if gameConfig.MaxVelocity <= 0 {
			gameConfig.MaxVelocity = 20
		}
		if gameConfig.CheckpointCount <= 0 {
			gameConfig.CheckpointCount = 10
		}
		cgGame.SetConfig(gameConfig)

		hoverrace.NewGame(cgGame, gameConfig).Run()
	})
}
