package hoverrace

import (
	"math"
	"time"

	"github.com/code-game-project/go-server/cg"
)

const throttleSpeed = 1
const turnSpeed = 220
const maxAcceleration = 5
const maxVelocity = 20

type Player struct {
	id       string
	username string

	ready bool

	finished bool

	checkpoints []Vec

	throttle       float64
	targetThrottle float64

	angle       float64
	targetAngle float64

	acc Vec
	vel Vec
	pos Vec

	cg   *cg.Player
	game *Game
}

func (p *Player) update(delta time.Duration) {
	if p.cg.SocketCount() > 0 {
		diff := AngleDifference(p.angle, p.targetAngle)
		p.targetAngle = p.angle + diff

		p.move(delta)
		p.checkCollisions()

		p.game.hovercrafts[p.id] = Hovercraft{
			Pos:         p.pos,
			Checkpoints: len(p.checkpoints),
			Velocity:    p.vel,
			Throttle:    p.throttle,
			Angle:       p.angle,
		}
	}
}

func (p *Player) move(delta time.Duration) {
	if !p.finished && p.game.running {
		if p.targetThrottle > p.throttle {
			p.throttle += throttleSpeed * delta.Seconds()
			if p.throttle > p.targetThrottle {
				p.throttle = p.targetThrottle
			}
		} else if p.targetThrottle < p.throttle {
			p.throttle -= throttleSpeed * delta.Seconds()
			if p.throttle < p.targetThrottle {
				p.throttle = p.targetThrottle
			}
		}

		if p.targetAngle > p.angle {
			p.angle += turnSpeed * delta.Seconds()
			if p.angle > p.targetAngle {
				p.angle = p.targetAngle
			}
		} else if p.targetAngle < p.angle {
			p.angle -= turnSpeed * delta.Seconds()
			if p.angle < p.targetAngle {
				p.angle = p.targetAngle
			}
		}

		p.acc = VecFromAngle(p.angle).Mul(maxAcceleration * p.throttle)
	} else {
		p.throttle = 0

		if p.vel.MagnitudeSquared() != 0 {
			velMag := p.vel.Magnitude()

			p.acc = Vec{
				X: -p.vel.X / velMag,
				Y: -p.vel.Y / velMag,
			}.Mul(math.Min(maxAcceleration, math.Abs(velMag)))
		}

		if p.vel.MagnitudeSquared() < 0.01 {
			p.acc = p.acc.Mul(0)
			p.vel = p.vel.Mul(0)
		}
	}

	p.vel = p.vel.Add(p.acc.Mul(delta.Seconds()))

	p.pos = p.pos.Add(p.vel.Mul(delta.Seconds()))
}

func (p *Player) checkCollisions() {
	if !p.game.running || p.finished {
		return
	}

	done := false
outer:
	for !done {
		for i, point := range p.checkpoints {
			if math.Abs(point.Sub(p.pos).MagnitudeSquared()) > 1 {
				continue
			}

			p.checkpoints[i] = p.checkpoints[len(p.checkpoints)-1]
			p.checkpoints = p.checkpoints[:len(p.checkpoints)-1]

			p.cg.Send("server", CheckpointsEvent, CheckpointsEventData{
				Checkpoints: p.checkpoints,
				FinishLine:  p.game.finishLine,
			})
			continue outer
		}
		done = true
	}

	if len(p.checkpoints) == 0 && math.Abs(p.game.finishLine.Sub(p.pos).MagnitudeSquared()) <= 1 {
		p.finished = true
		p.game.finishedPlayers = append(p.game.finishedPlayers, FinishedPlayer{
			Id:       p.id,
			Place:    len(p.game.finishedPlayers) + 1,
			Duration: time.Now().Sub(p.game.startTime).Milliseconds(),
		})
		p.game.cg.Send(p.id, FinishedPlayersEvent, FinishedPlayersEventData{
			Players: p.game.finishedPlayers,
		})

		if len(p.game.finishedPlayers) == len(p.game.players) {
			p.game.finish()
		}
	}
}
