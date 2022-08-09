# hoverrace
![CodeGame Version](https://img.shields.io/badge/CodeGame-v0.7-orange)
![CGE Version](https://img.shields.io/badge/CGE-v0.4-green)

Race against other hovercrafts from checkpoint to checkpoint.

## Known instances

- `games.code-game.org/hoverrace`

## Usage

```sh
# Run on default port 8080
hoverrace

# Specify a custom port
hoverrace --port=5000

## Specify a custom port through an environment variable
CG_PORT=5000 hoverrace
```

### Running with Docker

Prerequisites:
- [Docker](https://docker.com/)

```sh
# Download image
docker pull codegameproject/hoverrace:0.3

# Run container
docker run -d --restart on-failure -p <port-on-host-machine>:8080 --name hoverrace codegameproject/hoverrace:0.3
```

## Event Flow

1. Send a `ready` event to the server when you think the game should begin.
2. The `ready_players` event updates all players on the readiness of all players.
3. The `checkpoints` event contains all checkpoints and the finish line position.
4. The `countdown` event counts down 5 seconds.
5. The `start` event is sent to all players when the race begins.
6. Send a `throttle` event to begin moving.
7. The `hovercrafts` event is sent repeatedly to all players to update them on the state of all hovercrafts.
8. The `checkpoints` event is sent again when you cross a checkpoint.
9. The `finished_players` event is sent to all players when a player crosses the finish line. It contains all of the players that have finished the race. The game keeps going until all players have finished.
10. Send a `ready` event if you want to play again.

## Building

### Prerequisites

- [Go](https://go.dev) 1.18+

```sh
git clone https://github.com/code-game-project/hoverrace.git
cd hoverrace
codegame build
```
## License

Copyright (C) 2022 Julian Hofmann

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
