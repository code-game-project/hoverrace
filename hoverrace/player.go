package hoverrace

import (
	"math"
	"time"

	"github.com/code-game-project/go-server/cg"
)

type Player struct {
	id       string
	username string

	ready bool

	finished bool

	checkpoints []Vec

	thrust       float64
	targetThrust float64

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
			Thrust:      p.thrust,
			Angle:       p.angle,
		}
	}
}

func (p *Player) move(delta time.Duration) {
	if !p.finished && p.game.running {
		if p.targetThrust > p.thrust {
			p.thrust += p.game.config.ThrottleSpeed * delta.Seconds()
			if p.thrust > p.targetThrust {
				p.thrust = p.targetThrust
			}
		} else if p.targetThrust < p.thrust {
			p.thrust -= p.game.config.ThrottleSpeed * delta.Seconds()
			if p.thrust < p.targetThrust {
				p.thrust = p.targetThrust
			}
		}

		if p.targetAngle > p.angle {
			p.angle += p.game.config.TurnSpeed * delta.Seconds()
			if p.angle > p.targetAngle {
				p.angle = p.targetAngle
			}
		} else if p.targetAngle < p.angle {
			p.angle -= p.game.config.TurnSpeed * delta.Seconds()
			if p.angle < p.targetAngle {
				p.angle = p.targetAngle
			}
		}
		p.angle = NormalizeAngle(p.angle)

		p.acc = VecFromAngle(p.angle).Mul(p.game.config.MaxAcceleration * p.thrust)
	} else {
		p.thrust = 0

		if p.vel.MagnitudeSquared() != 0 {
			velMag := p.vel.Magnitude()

			p.acc = Vec{
				X: -p.vel.X / velMag,
				Y: -p.vel.Y / velMag,
			}.Mul(math.Min(p.game.config.MaxAcceleration, math.Abs(velMag)))
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

			p.cg.Send(CheckpointsEvent, CheckpointsEventData{
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
			Duration: time.Since(p.game.startTime).Milliseconds(),
		})
		p.game.cg.Send(FinishedPlayersEvent, FinishedPlayersEventData{
			Players: p.game.finishedPlayers,
		})

		if len(p.game.finishedPlayers) == len(p.game.players) {
			p.game.finish()
		}
	}
}
