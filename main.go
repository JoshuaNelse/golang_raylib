package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"math"
	"os"
	audioengine "raylib/playground/engines/audio-engine"
	collisionengine "raylib/playground/engines/collision-engine"
	projectileengine "raylib/playground/engines/projectile-engine"
	"raylib/playground/game"
	util "raylib/playground/game/utils"
	"raylib/playground/structs"
	"raylib/playground/structs/draw2d"
	"strconv"

	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/resolv"
)

const (
	screenWidth  = 1000
	screenHeight = 640
)

var (
	running  = true
	bkgColor = rl.NewColor(147, 211, 196, 255)

	texture rl.Texture2D // eventually this should be remove and draw2d used everywhere
	enemies []*structs.Enemy

	firstPlayer structs.Player

	// TODO find a better way
	playerSword               *structs.Weapon
	playerBow                 *structs.Weapon
	playerBowThatShootSwords  *structs.Weapon
	playerSwordThatShootsBows *structs.Weapon

	playerSpeed   float32 = 3
	playerUp      bool
	playerDown    bool
	playerRight   bool
	playerLeft    bool
	playerFlipped bool

	firstEnemy structs.Enemy

	tileDest          rl.Rectangle
	tileSrc           rl.Rectangle
	tileMap           []int
	srcMap            []string
	collisionMap      []string
	collisionMapDebug []rl.Rectangle
	debugMode         bool = true
	mapW              int
	mapH              int

	musicPaused bool = true
	cam         rl.Camera2D
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

func loadMap(mapFile string) {
	fmt.Println("Attempting to load map:", mapFile)
	file, err := ioutil.ReadFile(mapFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	remNewLines := strings.Replace(string(file), "\n", " ", -1)
	sliced := strings.Split(remNewLines, " ")

	// map dimensions
	mapW, err = strconv.Atoi(sliced[0])
	mapH, err = strconv.Atoi(sliced[1])

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// floor tile dimensions
	tileDest.Height = 32
	tileDest.Width = 32
	tileSrc.Height = 16
	tileSrc.Width = 16

	//pixel level collision
	spaceWidth := mapW * int(2*tileDest.Width)
	spaceHeight := mapH * int(2*tileDest.Height)
	spaceCellWidth := 1
	spaceCellHeight := 1
	collisionengine.WorldCollisionSpace = resolv.NewSpace(spaceWidth, spaceHeight, spaceCellWidth, spaceCellHeight)

	for i, val := range sliced[2:] {
		if i < mapW*mapH {
			if m, err := strconv.Atoi(val); err != nil {
				fmt.Println(err)
			} else {
				tileMap = append(tileMap, m)
			}
		} else if i < mapW*mapH*2 {
			srcMap = append(srcMap, val)
		} else {
			collisionMap = append(collisionMap, val)
		}
	}

	// load collision space
	objects := []*resolv.Object{}
	for i, col := range collisionMap {
		// probably more perfomant to just skip these "." collisionless tiles.
		if col == "." {
			continue
		}
		x := tileDest.Width * float32(i%mapW)  // 6 % 5 means x column 1
		y := tileDest.Height * float32(i/mapW) // 6 % 5 means y row of 1
		oX, oY, oW, oH := collisionTileOffsetMap[col].getTileCollisionOffset(x, y, tileDest.Width, tileDest.Height)
		collisionMapDebug = append(collisionMapDebug, rl.NewRectangle(oX, oY, oW, oH))
		// object offset is different than the regular sprite draw
		newObj := resolv.NewObject(float64(oX-oW), float64(oY-oH), float64(oW), float64(oH))
		newObj.AddTags("env")
		objects = append(objects, newObj)

	}

	collisionengine.WorldCollisionSpace.Add(objects...)
}

/*
Right and Left offset have to consider each other
Top and Bottom offset have to consider each other
*/
func (c collisionOffset) getTileCollisionOffset(x, y, w, h float32) (float32, float32, float32, float32) {
	offsetW := w*c.R - (w * (1 - c.L))
	offsetH := h*c.B - (h * (1 - c.T))
	offsetX := x - w*(1-c.R)
	offsetY := y - h*(1-c.B)
	return offsetX, offsetY, offsetW, offsetH
}

type collisionOffset struct {
	// percentage offset for collision (left, Right, Top, Bottom)
	L float32
	R float32
	T float32
	B float32
}

var (
	tileMapIndex = map[string]map[int]structs.Point{
		"_": emptyTileMap,
		"f": floorTileMap,
		"w": wallTileMap,
		"d": decorTileMap,
	}
	collisionTileOffsetMap = map[string]collisionOffset{
		"+": {1, 1, 1, 1},
		".": {0, 0, 0, 0}, // probably more efficient to just skip writting a collision here
		"d": {L: .8, R: .8, T: .3, B: 1},
	}
	emptyTileMap = map[int]structs.Point{
		0: {X: 0, Y: 0},
	}
	floorTileMap = map[int]structs.Point{
		1: {X: 16, Y: 64}, // plain floor
		2: {X: 32, Y: 64}, // really cracked floor
		3: {X: 48, Y: 64}, // kinda cracked floor
		4: {X: 16, Y: 80}, // big hole
		5: {X: 32, Y: 80}, // somewhat cracked floor
		6: {X: 48, Y: 80}, // bottomr right hole
		7: {X: 16, Y: 96}, // top right hole
		8: {X: 32, Y: 96}, // top left whole
		9: {X: 48, Y: 96}, // ladder
	}
	wallTileMap = map[int]structs.Point{
		1:  {X: 16, Y: 0},   // top left
		2:  {X: 32, Y: 0},   // top mid
		3:  {X: 48, Y: 0},   // top right
		4:  {X: 16, Y: 16},  // left
		5:  {X: 32, Y: 16},  // mid
		6:  {X: 48, Y: 16},  // right
		7:  {X: 0, Y: 112},  // wall_side_top_left 0 112 16 16
		8:  {X: 16, Y: 112}, // wall_side_top_right 16 112 16 16
		9:  {X: 0, Y: 128},  // wall_side_mid_left 0 128 16 16
		10: {X: 16, Y: 128}, // wall_side_mid_right 16 128 16 16
		11: {X: 0, Y: 144},  // wall_side_front_left 0 144 16 16
		12: {X: 16, Y: 144}, // wall_side_front_right 16 144 16 16
		13: {X: 32, Y: 112}, // wall_corner_top_left 32 112 16 16
		14: {X: 48, Y: 112}, // wall_corner_top_right 48 112 16 16
		15: {X: 32, Y: 128}, // wall_corner_left 32 128 16 16
		16: {X: 48, Y: 128}, // wall_corner_right 48 128 16 16
		17: {X: 32, Y: 144}, // wall_corner_bottom_left 32 144 16 16
		18: {X: 48, Y: 144}, // wall_corner_bottom_right 48 144 16 16
		19: {X: 32, Y: 160}, // wall_corner_front_left 32 160 16 16
		20: {X: 48, Y: 160}, // wall_corner_front_right 48 160 16 16
		21: {X: 80, Y: 128}, // wall_inner_corner_l_top_left 80 128 16 16
		22: {X: 64, Y: 128}, // wall_inner_corner_l_top_rigth 64 128 16 16
		23: {X: 80, Y: 144}, // wall_inner_corner_mid_left 80 144 16 16
		24: {X: 64, Y: 144}, // wall_inner_corner_mid_rigth 64 144 16 16
		25: {X: 80, Y: 160}, // wall_inner_corner_t_top_left 80 160 16 16
		26: {X: 64, Y: 160}, // wall_inner_corner_t_top_rigth 64 160 16 16
	}
	decorTileMap = map[int]structs.Point{
		1:  {X: 64, Y: 0},   // wall_fountain_top 64 0 16 16
		2:  {X: 64, Y: 16},  // wall_fountain_mid_red_anim 64 16 16 16 3
		3:  {X: 64, Y: 32},  // wall_fountain_basin_red_anim 64 32 16 16 3
		4:  {X: 64, Y: 48},  // wall_fountain_mid_blue_anim 64 48 16 16 3
		5:  {X: 64, Y: 64},  // wall_fountain_basin_blue_anim 64 64 16 16 3
		6:  {X: 16, Y: 32},  // wall_banner_red 16 32 16 16
		7:  {X: 32, Y: 32},  // wall_banner_blue 32 32 16 16
		8:  {X: 16, Y: 48},  // wall_banner_green 16 48 16 16
		9:  {X: 32, Y: 48},  // wall_banner_yellow 32 48 16 16
		10: {X: 96, Y: 80},  // wall_column_top 96 80 16 16
		11: {X: 96, Y: 96},  // wall_column_mid 96 96 16 16
		12: {X: 96, Y: 112}, // wall_coulmn_base 96 112 16 16
		13: {X: 80, Y: 80},  // column_top 80 80 16 16
		14: {X: 80, Y: 96},  // column_mid 80 96 16 16
		15: {X: 80, Y: 112}, // coulmn_base 80 112 16 16
	}
)

type drawParams struct {
	texture  rl.Texture2D
	srcRec   rl.Rectangle
	destRec  rl.Rectangle
	origin   rl.Vector2
	rotation float32
	tint     color.RGBA
}

func drawMapBackground() []drawParams {

	foreGroundDrawParams := []drawParams{}
	for i, tileInt := range tileMap {
		if tileInt == 0 {
			continue
		}
		tileDest.X = tileDest.Width * float32(i%mapW) // 6 % 5 means x column 1
		tileDest.Y = tileDest.Width * float32(i/mapW) // 6 % 5 means y row of 1
		tileMap := tileMapIndex[strings.ToLower(srcMap[i])]
		tileSrc.X = tileMap[tileInt].X
		tileSrc.Y = tileMap[tileInt].Y

		if strings.ToUpper(srcMap[i]) == srcMap[i] {
			// TODO make this fill square more sofistacted - maybe random or something
			fillTile := tileSrc
			fillTile.X = 16
			fillTile.Y = 64
			rl.DrawTexturePro(texture, fillTile, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)

			// draw behind player if Y is "behind" player, but skip this with Walls
			if tileDest.Y > firstPlayer.Sprite.Dest.Y || strings.ToLower(srcMap[i]) == "w" {
				foreGroundDrawParams = append(
					foreGroundDrawParams,
					drawParams{
						texture:  texture,
						srcRec:   tileSrc,
						destRec:  tileDest,
						origin:   rl.NewVector2(tileDest.Width, tileDest.Height),
						rotation: 0,
						tint:     rl.White,
					})
			} else {
				rl.DrawTexturePro(texture, tileSrc, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)
			}
		} else {
			rl.DrawTexturePro(texture, tileSrc, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)
		}
	}

	return foreGroundDrawParams
}

func drawScene() {
	foreGround := drawMapBackground()

	/*
		Thinking to make a method drawEnv
			players, monsters, and projectiles should all be part of an env array
			Here I could have an interface for SomeEnvObject.Draw()
			This draw method in implemetation could handle drawing sub components / nested images
			for example:
				The Player.Draw() could make sure that the wielded weapon
				is drawn as well as the player itself.
	*/
	firstPlayer.Draw()
	for _, e := range enemies {
		e.Draw()
	}
	for _, p := range projectileengine.Projectiles {
		p.Draw()
	}
	// draw foreground after player
	for _, draw := range foreGround {
		rl.DrawTexturePro(draw.texture, draw.srcRec, draw.destRec, draw.origin, draw.rotation, draw.tint)
	}

	// draw debug collision objects
	if debugMode {
		for _, o := range collisionMapDebug {
			mo := util.ObjFromRect(o)
			rl.DrawRectangleLines(int32(mo.X), int32(mo.Y), int32(mo.W), int32(mo.H), rl.White)
		}

		for _, p := range projectileengine.Projectiles {
			rl.DrawLine(int32(p.Start.X), int32(p.Start.Y), int32(p.End.X), int32(p.End.Y), rl.Pink)
		}

		// debug player collision box
		po := firstPlayer.Obj
		rl.DrawRectangleLines(int32(po.X), int32(po.Y), int32(po.W), int32(po.H), rl.Orange)

		for _, e := range enemies {
			rl.DrawRectangleLines(int32(e.Obj.X), int32(e.Obj.Y), int32(e.Obj.W), int32(e.Obj.H), rl.White)

		}

		playerCenter := structs.Point{
			X: float32(firstPlayer.Obj.X + firstPlayer.Obj.W/2),
			Y: float32(firstPlayer.Obj.Y + firstPlayer.Obj.H/2),
		}
		rl.DrawCircleLines(int32(playerCenter.X), int32(playerCenter.Y), 32, rl.Green)
		angle := util.GetPlayerToMouseAngleDegress()
		rl.DrawCircleSectorLines(rl.NewVector2(playerCenter.X, playerCenter.Y), 32, angle, angle-45, 5, rl.White)
		rl.DrawCircleSectorLines(rl.NewVector2(playerCenter.X, playerCenter.Y), 32, angle, angle+45, 5, rl.White)
	}

}

func drawUI() {
	if debugMode {
		rl.DrawRectangleRounded(rl.NewRectangle(3, 3, 500, 90), .1, 10, rl.DarkGray)
		rl.DrawRectangleRoundedLines(rl.NewRectangle(3, 3, 500, 90), .1, 10, 3, rl.White)
		rl.DrawText(fmt.Sprintf("FPS: %v", rl.GetFPS()), 10, 10, 16, rl.White)
		rl.DrawText(fmt.Sprintf("player {X: %v, Y:%v}", firstPlayer.Obj.X, firstPlayer.Obj.Y), 10, 30, 16, rl.White)
		rl.DrawText(fmt.Sprintf("mouse  {X: %v, Y:%v}", rl.GetMouseX(), rl.GetMouseY()), 10, 50, 16, rl.White)

		// wierd thing where rise/run are opposite directions (think it has to do with x/y being negative flipped)
		rise := float64(rl.GetMouseX()) - float64(rl.GetScreenWidth()/2)
		run := float64(rl.GetMouseY()) - float64(rl.GetScreenHeight())/2

		angle := util.GetPlayerToMouseAngleDegress()
		rl.DrawText(fmt.Sprintf("mouse->player  {X: %v, Y:%v}", rise, run), 10, 70, 16, rl.White)
		rl.DrawText(fmt.Sprintf("Atan(%v/%v) = %v degrees", rise, run, int(angle)), 250, 10, 16, rl.White)
		rl.DrawText(fmt.Sprintf("Live Projectiles: %v", len(projectileengine.Projectiles)), 250, 30, 16, rl.White)
	}
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
		debugMode = !debugMode
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
		firstPlayer.EquipWeapon(playerSword)
	} else if rl.IsKeyPressed(rl.KeyTwo) {
		firstPlayer.EquipWeapon(playerBow)
	} else if rl.IsKeyPressed(rl.KeyThree) {
		firstPlayer.EquipWeapon(playerBowThatShootSwords)
	} else if rl.IsKeyPressed(rl.KeyFour) {
		firstPlayer.EquipWeapon(playerSwordThatShootsBows)
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
	running = !rl.WindowShouldClose()

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
			playerFlipped = false
		}
		if playerLeft {
			util.FlipLeft(&firstPlayer.Sprite.Src)
			dx -= float64(playerSpeed)
			playerFlipped = true
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
	drawScene()

	rl.EndMode2D()
	drawUI()
	rl.EndDrawing()
}

func initialize() {
	rl.InitWindow(screenWidth, screenHeight, "Raylib Playground :)")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	draw2d.InitTexture()
	texture = draw2d.Texture
	audioengine.InitializeAudio()

	playerSprite := structs.Sprite{
		Src:  rl.NewRectangle(128, 100, 16, 28),
		Dest: rl.NewRectangle(285, 200, 32, 56),
	}

	playerObj := util.ObjFromRect(playerSprite.Dest)
	playerObj.H *= .6

	firstPlayer = structs.Player{
		Sprite: playerSprite,
		Obj:    playerObj,
		Hand:   structs.Point{X: float32(playerObj.W) * .5, Y: float32(playerObj.H) * .94},
	}

	firstPlayer.Sprite.Dest.X = util.RectFromObj(firstPlayer.Obj).X
	firstPlayer.Sprite.Dest.Y = util.RectFromObj(firstPlayer.Obj).Y
	firstPlayer.Obj.AddTags("Player")

	swordSprite := structs.Sprite{
		// src: rl.NewRectangle(307, 26, 10, 21), // rusty sword
		// src: rl.NewRectangle(339, 114, 10, 29), // weapon_knight_sword
		// src: rl.NewRectangle(310, 124, 8, 19), // cleaver
		// src: rl.NewRectangle(325, 113, 9, 30), // weapon_duel_sword
		// src: rl.NewRectangle(322, 81, 12, 30), // weapon_anime_sword
		Src: rl.NewRectangle(339, 26, 10, 21), // weapon_red_gem_sword

		Dest: rl.NewRectangle(
			firstPlayer.Hand.X+float32(firstPlayer.Obj.X),
			firstPlayer.Hand.Y+float32(firstPlayer.Obj.Y),
			10*1.35,
			21*1.35,
		),
	}

	playerSword = &structs.Weapon{
		Sprite: swordSprite,
		Obj:    util.ObjFromRect(swordSprite.Dest),
		// handle is the origin offset for the sprite
		Handle:       structs.Point{X: swordSprite.Dest.Width * .5, Y: swordSprite.Dest.Height * .9},
		AttackSpeed:  8,
		Cooldown:     24,
		IdleRotation: -30,
		AttackRotator: func(w structs.Weapon) float32 {
			return w.IdleRotation * -3 / float32(w.AttackSpeed) * float32(w.AttackFrame)
		},
		ProjectileCount:         3,
		ProjectileSpreadDegrees: 45,
		Projectilelength:        32,
		ProjectileTTLFrames:     10,
		TintColor:               rl.White,
	}

	bowSprite := structs.Sprite{
		Src: rl.NewRectangle(325, 180, 7, 25), // weapon_bow 325 180 7 25

		Dest: rl.NewRectangle(
			firstPlayer.Hand.X+float32(firstPlayer.Obj.X),
			firstPlayer.Hand.Y+float32(firstPlayer.Obj.Y),
			7*1.1,
			25*1.1,
		),
	}

	arrowSpriteSource := structs.Sprite{
		Src:  rl.NewRectangle(308, 186, 7, 21),     // weapon_arrow 308 186 7 21
		Dest: rl.NewRectangle(0, 0, 7*1.5, 21*1.5), // only using h, w for scaling
	}

	playerBow = &structs.Weapon{
		Sprite:              bowSprite,
		ProjectileSpriteSrc: arrowSpriteSource,

		Obj: util.ObjFromRect(bowSprite.Dest),
		// handle is the origin offset for the sprite
		Handle:       structs.Point{X: bowSprite.Dest.Width * .5, Y: bowSprite.Dest.Height * .75},
		AttackSpeed:  8,
		Cooldown:     24,
		IdleRotation: 20,
		AttackRotator: func(w structs.Weapon) float32 {
			// TODO make it follow mouse -- for now make it 0
			return 0
		},
		ProjectileCount:         1,
		ProjectileSpreadDegrees: 0,
		Projectilelength:        21,
		ProjectileTTLFrames:     32,
		ProjectileVelocity:      8,
		TintColor:               rl.White,
	}

	bowThatShootsSwordsSprite := structs.Sprite{
		Src: rl.NewRectangle(325, 180, 7, 25), // weapon_bow 325 180 7 25

		Dest: rl.NewRectangle(
			firstPlayer.Hand.X+float32(firstPlayer.Obj.X),
			firstPlayer.Hand.Y+float32(firstPlayer.Obj.Y),
			7*1.1,
			25*1.1,
		),
	}

	swordArrowSpriteSource := structs.Sprite{
		Src:  rl.NewRectangle(339, 114, 10, 29),     // weapon_knight_sword 339 114 10 29
		Dest: rl.NewRectangle(0, 0, 10*1.5, 29*1.5), // only using h, w for scaling
	}

	playerBowThatShootSwords = &structs.Weapon{
		Sprite:              bowThatShootsSwordsSprite,
		ProjectileSpriteSrc: swordArrowSpriteSource,

		Obj: util.ObjFromRect(bowSprite.Dest),
		// handle is the origin offset for the sprite
		Handle:       structs.Point{X: bowSprite.Dest.Width * .5, Y: bowSprite.Dest.Height * .75},
		AttackSpeed:  8,
		Cooldown:     24,
		IdleRotation: 20,
		AttackRotator: func(w structs.Weapon) float32 {
			// TODO make it follow mouse -- for now make it 0
			return 0
		},
		ProjectileCount:         1,
		ProjectileSpreadDegrees: 0,
		Projectilelength:        21,
		ProjectileTTLFrames:     32,
		ProjectileVelocity:      8,
		TintColor:               rl.Blue,
	}

	swordThatShootsBowsSprite := structs.Sprite{
		Src: rl.NewRectangle(322, 81, 12, 30), // weapon_anime_sword 322 81 12 30

		Dest: rl.NewRectangle(
			firstPlayer.Hand.X+float32(firstPlayer.Obj.X),
			firstPlayer.Hand.Y+float32(firstPlayer.Obj.Y),
			12*1.1,
			30*1.1,
		),
	}

	bowArrowSpriteSource := structs.Sprite{
		Src:  rl.NewRectangle(325, 180, 7, 25),     // weapon_bow 325 180 7 25
		Dest: rl.NewRectangle(0, 0, 7*1.5, 25*1.5), // only using h, w for scaling
	}

	playerSwordThatShootsBows = &structs.Weapon{
		Sprite:              swordThatShootsBowsSprite,
		ProjectileSpriteSrc: bowArrowSpriteSource,

		Obj: util.ObjFromRect(swordSprite.Dest),
		// handle is the origin offset or the sprite
		Handle:       structs.Point{X: swordSprite.Dest.Width * .5, Y: swordSprite.Dest.Height * .99},
		AttackSpeed:  8,
		Cooldown:     24,
		IdleRotation: -30,
		AttackRotator: func(w structs.Weapon) float32 {
			return w.IdleRotation * -3 / float32(w.AttackSpeed) * float32(w.AttackFrame)
		},
		ProjectileCount:         3,
		ProjectileSpreadDegrees: 35,
		Projectilelength:        21,
		ProjectileTTLFrames:     32,
		ProjectileVelocity:      4,
		TintColor:               rl.White,
	}

	firstPlayer.Weapon = playerSword

	// Test enemy orc_warrior_idle_anim 368 204 16 20 4
	enemySprite := structs.Sprite{
		Src:        rl.NewRectangle(368, 204, 16, 24),
		Dest:       rl.NewRectangle(250, 250, 32, 48),
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

	cam = rl.NewCamera2D(rl.NewVector2(screenWidth/2, screenHeight/2), getCameraTarget(), 0.0, 1.25)
	loadMap("resources/maps/first.map")
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
	for running {
		input()
		update()
		render()
	}
	quit()
}
