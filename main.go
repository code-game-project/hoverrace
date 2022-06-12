package main

import (
	"os"
	"strconv"
	"time"

	"github.com/code-game-project/go-server/cg"
	"github.com/code-game-project/hoverrace/hoverrace"
	"github.com/spf13/pflag"
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
		port = 80
	}

	server := cg.NewServer("hoverrace", cg.ServerConfig{
		DisplayName:             "Hover Race",
		Description:             "Race against other hovercrafts from checkpoint to checkpoint.",
		Version:                 "0.1",
		RepositoryURL:           "https://github.com/code-game-project/hoverrace",
		WebsocketTimeout:        1 * time.Minute,
		MaxPlayersPerGame:       10,
		Port:                    port,
		CGEFilepath:             "events.cge",
		DeleteInactiveGameDelay: 30 * time.Minute,
		KickInactivePlayerDelay: 1 * time.Minute,
	})

	server.Run(func(cgGame *cg.Game) {
		hoverrace.NewGame(cgGame).Run()
	})
}
