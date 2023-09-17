package main

import (
	"fmt"
	"math"
	"raylib/playground/directors/map-director"
	"raylib/playground/engines/audio-engine"
	"raylib/playground/engines/collision-engine"
	"raylib/playground/engines/projectile-engine"
	"raylib/playground/game"
	util "raylib/playground/game/utils"
	"raylib/playground/model"
	"raylib/playground/model/draw2d"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	playerSpeed float32 = 3
)

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

func update() {
	game.Running = !rl.WindowShouldClose()

	if game.MainPlayer.Moving {
		// used for collision check
		dx := 0.0
		dy := 0.0

		if game.PlayerUp {
			dy -= float64(playerSpeed)
		}
		if game.PlayerDown {
			dy += float64(playerSpeed)
		}
		if game.PlayerRight {
			util.FlipRight(&game.MainPlayer.Sprite.Src)
			dx += float64(playerSpeed)
			game.MainPlayer.SpriteFlipped = false
		}
		if game.PlayerLeft {
			util.FlipLeft(&game.MainPlayer.Sprite.Src)
			dx -= float64(playerSpeed)
			game.MainPlayer.SpriteFlipped = true
		}
		// check for collisions
		if collision := game.MainPlayer.Obj.Check(dx, dy, "nav"); collision != nil {

			fmt.Println("We hit a door!")

			for _, tag := range collision.Objects[0].Tags() {
				if strings.HasPrefix(tag, "doorId") {
					fmt.Println(tag)
				}
			}

			map_director.LoadMap("resources/maps/second.map", draw2d.Texture)
			game.MainPlayer.Obj.Space.Remove(game.MainPlayer.Obj)
			collision_engine.WorldCollisionSpace.Add(game.MainPlayer.Obj)

		}
		if collision := game.MainPlayer.Obj.Check(0, dy, "env"); collision != nil {
			//fmt.Println("Y axis collision happened: ", collision)
			// heuristically stop movement on collision because the other way is buggy
			dy = 0
		}
		if collision := game.MainPlayer.Obj.Check(dx, 0, "env"); collision != nil {
			//fmt.Println("X axis collision happened: ", collision)
			// heuristically stop movement on collision because the other way is buggy
			dx = 0
		}
		if collision := game.MainPlayer.Obj.Check(dx, dy, "enemy"); collision != nil {
			dx /= 4
			dy /= 4
		}
		game.MainPlayer.Move(dx, dy)
	}

	if game.MainPlayer.AttackCooldown > 0 {
		game.MainPlayer.AttackCooldown--
		game.MainPlayer.Attacking = false
	}
	if game.MainPlayer.Attacking {
		projectile_engine.Projectiles = append(projectile_engine.Projectiles, game.MainPlayer.Attack()...)
	}

	var nextFrameProjectiles []model.Projectile

	for _, p := range projectile_engine.Projectiles {
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
			for _, e := range game.Enemies {
				for _, line := range linesFromRect(util.RectFromObj(e.Obj)) {
					var collisionPoint *rl.Vector2
					collision = rl.CheckCollisionLines(p.Start, p.End, line.start, line.end, collisionPoint)
					if collision {
						break
					}
				}
				if collision {
					e.Hurt()
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
	projectile_engine.Projectiles = nextFrameProjectiles

	var survivingEnemies []*model.Enemy
	for _, e := range game.Enemies {
		if !e.Dead || e.DeathFrames > 0 {
			survivingEnemies = append(survivingEnemies, e)
		}
	}
	game.Enemies = survivingEnemies

	audio_engine.UpdateMusicStream()
	game.UpdateCameraTargetToPlayerLocation()

	game.PlayerUp, game.PlayerDown, game.PlayerRight, game.PlayerLeft = false, false, false, false
}

func main() {

	game.Initialize(true)

	// Each Frame
	for game.Running {
		game.ReadPlayerInputs()
		update()
		game.Render()
	}
	game.Quit()
}
