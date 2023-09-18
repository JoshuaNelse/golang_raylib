package physics_engine

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	util "raylib/playground/game/utils"
	"raylib/playground/model"
)

var Projectiles []model.Projectile

type Line struct {
	start rl.Vector2
	end   rl.Vector2
}

func linesFromRect(rect rl.Rectangle) []Line {
	x := rect.X
	y := rect.Y
	w := rect.Width
	h := rect.Height
	t := Line{rl.NewVector2(x, y), rl.NewVector2(x-w, y)}
	l := Line{rl.NewVector2(x, y), rl.NewVector2(x, y-h)}
	r := Line{rl.NewVector2(x-w, y-h), rl.NewVector2(x, y-h)}
	b := Line{rl.NewVector2(x-w, y-h), rl.NewVector2(x-w, y)}
	return []Line{t, l, r, b}
}

func CalculatePlayerProjectileOutcome(enemies *[]*model.Enemy) {

	var nextFrameProjectiles []model.Projectile

	for _, p := range Projectiles {
		if p.Ttl > 0 {
			var collision bool
			p.Ttl--
			if p.Velocity > 0 {
				// find delta x and y
				var dx float32
				var dy float32
				switch p.Trajectory {
				case -90:
					dx, dy = float32(-p.Velocity), 0 // -x
				case 0:
					dx, dy = 0, float32(p.Velocity) // +y
				case 90:
					dx, dy = float32(p.Velocity), 0 // +x
				case 180:
					dx, dy = 0, float32(-p.Velocity) // -y
				default:
					dx = float32(p.Velocity) * float32(math.Sin(util.DegreesToRadians(p.Trajectory)))
					dy = float32(p.Velocity) * float32(math.Cos(util.DegreesToRadians(p.Trajectory)))
				}
				// add delta x and y to both x1,y1 and x2,y2
				p.Start.X += dx
				p.Start.Y += dy
				p.End.X += dx
				p.End.Y += dy
			}
			for _, e := range *enemies {
				for _, line := range linesFromRect(util.RectFromObj(e.Obj)) {
					var collisionPoint *rl.Vector2
					collision = rl.CheckCollisionLines(p.Start, p.End, line.start, line.end, collisionPoint)
					if collision {
						break
					}
				}
				if collision {
					e.Hurt()
					if e.Dead {
						WorldCollisionSpace.Remove(e.Obj)
					}
					break
				}
			}
			if !collision {
				nextFrameProjectiles = append(nextFrameProjectiles, p)
			}
		} else {
			// If TTL not > 0 then let this projectile "fade away"
		}
	}
	Projectiles = nextFrameProjectiles

	var survivingEnemies []*model.Enemy
	for _, e := range *enemies {
		if !e.Dead || e.DeathFrames > 0 {
			survivingEnemies = append(survivingEnemies, e)
		}
	}
	enemies = &survivingEnemies
}

func FireProjects(projectiles []model.Projectile) {
	Projectiles = append(Projectiles, projectiles...)
}
