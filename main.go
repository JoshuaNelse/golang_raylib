package main

import (
	"math"
	mapdirector "raylib/playground/directors/map-director"
	audioengine "raylib/playground/engines/audio-engine"
	collisionengine "raylib/playground/engines/collision-engine"
	drawworldengine "raylib/playground/engines/draw-world-engine"
	projectileengine "raylib/playground/engines/projectile-engine"
	"raylib/playground/game"
	"raylib/playground/game/structs"
	"raylib/playground/game/structs/armory/bows"
	"raylib/playground/game/structs/armory/cannon"
	"raylib/playground/game/structs/armory/staves"
	"raylib/playground/game/structs/armory/swords"
	"raylib/playground/game/structs/draw2d"
	util "raylib/playground/game/utils"
	pointmodel "raylib/playground/models/point-model"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	bkgColor = rl.NewColor(147, 211, 196, 255)

	texture rl.Texture2D // eventually this should be remove and draw2d used everywhere
	enemies []*structs.Enemy

	firstPlayer structs.Player

	playerSpeed float32 = 3
	playerUp    bool
	playerDown  bool
	playerRight bool
	playerLeft  bool

	firstEnemy structs.Enemy

	tileDest rl.Rectangle
	tileSrc  rl.Rectangle

	collisionMapDebug []rl.Rectangle
	mapW              int
	mapH              int

	musicPaused bool = true
	cam         rl.Camera2D
	mapFile     = "resources/maps/second.map"
)

type lineDrawParam struct {
	x1, y1 int // start
	x2, y2 int // end
	color  rl.Color
}

func getCameraTarget() rl.Vector2 {
	playerCenterX := float32(firstPlayer.Obj.X + firstPlayer.Obj.W/2)
	playerCenterY := float32(firstPlayer.Obj.Y + firstPlayer.Obj.H/2)
	return rl.NewVector2(playerCenterX, playerCenterY)
}

func input() {
	/*
		Thinking about making an inputDiretor module that checks for all input and returns
		A map of input managers that need to be invoked and their corresponding inputs to manage
	*/
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		firstPlayer.Moving = true
		playerUp = true
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		firstPlayer.Moving = true
		playerDown = true
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		firstPlayer.Moving = true
		playerLeft = true
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		firstPlayer.Moving = true
		playerRight = true
	}
	if rl.IsKeyPressed(rl.KeyQ) {
		musicPaused = !musicPaused
	}
	if rl.IsKeyPressed(rl.KeyBackSlash) {
		game.DebugMode = !game.DebugMode
	}
	if rl.GetMouseWheelMove() != 0 {
		mouseMove := rl.GetMouseWheelMove()
		if mouseMove > 0 && cam.Zoom < 2.0 {
			cam.Zoom = float32(math.Min(2.0, float64(cam.Zoom+float32(mouseMove)/15)))
		} else if mouseMove < 0 && cam.Zoom > .75 {
			cam.Zoom = float32(math.Max(.75, float64(cam.Zoom+float32(mouseMove)/15)))
		}
	}
	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		firstPlayer.Attacking = true
	}
	if rl.IsKeyPressed(rl.KeyOne) {
		firstPlayer.EquipWeapon(swords.Key())
	} else if rl.IsKeyPressed(rl.KeyTwo) {
		firstPlayer.EquipWeapon(bows.RegularBow())
	} else if rl.IsKeyPressed(rl.KeyThree) {
		firstPlayer.EquipWeapon(bows.SwordShooter())
	} else if rl.IsKeyPressed(rl.KeyFour) {
		firstPlayer.EquipWeapon(swords.BowShooter())
	} else if rl.IsKeyPressed(rl.KeyFive) {
		firstPlayer.EquipWeapon(cannon.PeopleShooter())
	} else if rl.IsKeyPressed(rl.KeySix) {
		firstPlayer.EquipWeapon(staves.PizzaShooter())
	}

}

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

	if firstPlayer.Moving {
		// used for collision check
		dx := 0.0
		dy := 0.0

		if playerUp {
			dy -= float64(playerSpeed)
		}
		if playerDown {
			dy += float64(playerSpeed)
		}
		if playerRight {
			util.FlipRight(&firstPlayer.Sprite.Src)
			dx += float64(playerSpeed)
			firstPlayer.SpriteFlipped = false
		}
		if playerLeft {
			util.FlipLeft(&firstPlayer.Sprite.Src)
			dx -= float64(playerSpeed)
			firstPlayer.SpriteFlipped = true
		}
		// check for collisions
		if collision := firstPlayer.Obj.Check(0, dy, "env"); collision != nil {
			//fmt.Println("Y axis collision happened: ", collision)
			// hueristically stop movement on collision because the other way is buggy
			dy = 0
		}
		if collision := firstPlayer.Obj.Check(dx, 0, "env"); collision != nil {
			//fmt.Println("X axis collision happened: ", collision)
			// hueristically stop movement on collision because the other way is buggy
			dx = 0
		}
		if collision := firstPlayer.Obj.Check(dx, dy, "enemy"); collision != nil {
			dx /= 4
			dy /= 4
		}
		firstPlayer.Move(dx, dy)
	}

	if firstPlayer.AttackCooldown > 0 {
		firstPlayer.AttackCooldown--
		firstPlayer.Attacking = false
	}
	if firstPlayer.Attacking {
		projectileengine.Projectiles = append(projectileengine.Projectiles, firstPlayer.Attack()...)
	}

	var nextFrameProjectiles []structs.Projectile

	for _, p := range projectileengine.Projectiles {
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
			for _, e := range enemies {
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
	projectileengine.Projectiles = nextFrameProjectiles

	var survivingEnemies []*structs.Enemy
	for _, e := range enemies {
		if !e.Dead || e.DeathFrames > 0 {
			survivingEnemies = append(survivingEnemies, e)
		}
	}
	enemies = survivingEnemies

	game.FrameCount++

	audioengine.UpdateMusicStream()
	if musicPaused {
		audioengine.PauseMusicStream()
	} else {
		audioengine.ResumeMusicStream()
	}

	cam.Target = getCameraTarget()

	playerUp, playerDown, playerRight, playerLeft = false, false, false, false
}

func render() {
	rl.BeginDrawing()
	rl.ClearBackground(bkgColor)
	rl.BeginMode2D(cam)
	drawworldengine.DrawScene()

	rl.EndMode2D()
	drawworldengine.DrawUI()
	rl.EndDrawing()
}

func initialize() {
	game.Initialize(true)

	draw2d.InitTexture()
	texture = draw2d.Texture
	audioengine.InitializeAudio()

	playerSprite := structs.Sprite{
		Src:     rl.NewRectangle(128, 100, 16, 28),
		Dest:    rl.NewRectangle(285, 200, 32, 56),
		Texture: draw2d.Texture,
	}

	playerObj := util.ObjFromRect(playerSprite.Dest)
	playerObj.H *= .6

	firstPlayer = structs.Player{
		Sprite: playerSprite,
		Obj:    playerObj,
		Hand:   pointmodel.Point{X: float32(playerObj.W) * .5, Y: float32(playerObj.H) * .94},
	}

	firstPlayer.Sprite.Dest.X = util.RectFromObj(firstPlayer.Obj).X
	firstPlayer.Sprite.Dest.Y = util.RectFromObj(firstPlayer.Obj).Y
	firstPlayer.Obj.AddTags("Player")
	firstPlayer.EquipWeapon(swords.RegularSword())

	drawworldengine.SetPlayer(&firstPlayer)

	// Test enemy orc_warrior_idle_anim 368 204 16 20 4
	enemySprite := structs.Sprite{
		Src:        rl.NewRectangle(368, 204, 16, 24),
		Dest:       rl.NewRectangle(250, 250, 32, 48),
		Texture:    draw2d.Texture,
		FrameCount: 4,
	}

	enemyObj := util.ObjFromRect(enemySprite.Dest)
	enemyObj.Y += enemyObj.H * .2
	enemyObj.H *= .7
	enemyObj.X += enemyObj.W * .2
	enemyObj.W *= .8
	enemyObj.AddTags("enemy")
	firstEnemy = structs.Enemy{
		Sprite:    enemySprite,
		Obj:       enemyObj,
		Health:    12,
		MaxHealth: 12,
	}
	enemies = append(enemies, &firstEnemy)
	drawworldengine.SetEnemies(&enemies)

	cam = rl.NewCamera2D(rl.NewVector2(game.ScreenWidth/2, game.ScreenHeight/2), getCameraTarget(), 0.0, 1.25)

	mapdirector.LoadMap(mapFile, draw2d.Texture)
	collisionengine.WorldCollisionSpace.Add(firstPlayer.Obj, enemyObj)
}

func quit() {
	draw2d.UnloadTexture()
	audioengine.UnloadAudioComponents()
	rl.CloseWindow()
}

func main() {
	initialize()

	// Each Frame
	for game.Running {
		input()
		update()
		render()
	}
	quit()
}
