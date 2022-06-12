# hoverrace
![CodeGame Version](https://img.shields.io/badge/CodeGame-v0.6-orange)
![CodeGame GameServer Version](https://img.shields.io/badge/GameServer-v0.1-yellow)
![CGE Version](https://img.shields.io/badge/CGE-v0.3-green)

Race against other hovercrafts from checkpoint to checkpoint.

## Known instances

- `games.code-game.org/hoverrace`

## Usage

```sh
# Run on default port 80
hoverrace

# Specify a custom port
hoverrace --port=8080

## Specify a custom port through an environment variable
CG_PORT=8080 hoverrace
```

### Running with Docker

Prerequisites:
- [Docker](https://docker.com/)

```sh
# Download image
docker pull codegameproject/hoverrace:0.1

# Run container
docker run -d -p <port-on-host-machine>:8080 --name hoverrace codegameproject/hoverrace:0.1
```

## Event Flow

1. Send a `ready` event to the server when you think the game should begin.
2. A `ready_players` event updates all players on the readiness of all players.
3. A `start` event is sent to all players when the race begins. It includes the positions of all checkpoints.
4. A `next_checkpoint` event is sent to every player to tell them where their next target is located.
5. Send a `throttle` event to begin moving.
6. A `hovercrafts` event is sent repeatedly to all players to update them on the state of all hovercrafts.
7. A `finished` event is sent to all players when a player crosses the finish line. The game keeps going until all players have finished.
8. Send a `ready` event if you want to play again.

## Building

### Prerequisites

- [Go](https://go.dev) 1.18+

```sh
git clone https://github.com/code-game-project/hoverrace.git
cd hoverrace
go build .
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
