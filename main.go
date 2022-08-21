package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"math"
	"os"
	testPlayer "raylib/playground/assets/player"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/resolv"
)

func test() testPlayer.Test {
	test := testPlayer.Test{Test: "123"}
	return test
}

type Enemy struct {
	sprite      Sprite
	obj         *resolv.Object
	health      int
	maxHealth   int
	hurtFrames  int
	deathFrames int
	dead        bool
}

type Sprite struct {
	src        rl.Rectangle
	dest       rl.Rectangle
	flipped    bool
	frameCount int
	frame      int
	rotation   float32
}

type Weapon struct {
	sprite              Sprite
	projectileSpriteSrc Sprite
	obj                 *resolv.Object
	handle              point
	reach               int
	attackSpeed         int
	cooldown            int
	tintColor           rl.Color

	idleRotation  float32
	attackRotator func(w Weapon, p Player) float32

	projectileCount         int
	projectilelength        int
	projectileSpreadDegrees int
	projectileTTLFrames     int
	projectileVelocity      int
}

// start code for player logic
type Player struct {
	sprite         Sprite
	obj            *resolv.Object
	weapon         *Weapon
	hand           point
	moving         bool
	attacking      bool
	attackFrame    int
	attackCooldown int
}

type Projectile struct {
	start rl.Vector2
	end   rl.Vector2
	ttl   int

	// for something like an arrow perhaps
	velocity   int
	trajectory float64 //degrees
	sprite     Sprite

	// sender     *interface{} at somepoint this would be good to have
}

const (
	screenWidth  = 1000
	screenHeight = 640
)

var (
	running  = true
	bkgColor = rl.NewColor(147, 211, 196, 255)

	texture             rl.Texture2D
	worldCollisionSpace *resolv.Space
	enemies             []*Enemy

	projectiles []Projectile

	player Player

	// TODO find a better way
	playerSword               *Weapon
	playerBow                 *Weapon
	playerBowThatShootSwords  *Weapon
	playerSwordThatShootsBows *Weapon

	playerSpeed   float32 = 3
	playerUp      bool
	playerDown    bool
	playerRight   bool
	playerLeft    bool
	playerFlipped bool

	enemy Enemy

	frameCount int

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
	music       rl.Music

	// TODO this should probably have some module to handle sound effects
	swordSound rl.Sound

	cam rl.Camera2D
)

func degreesToRadians(d float64) float64 {
	return d * (math.Pi / 180)
}

func radiansToDegrees(r float64) float64 {
	return r * (180 / math.Pi)
}

func (p *Projectile) draw() {
	w := p.sprite.dest.Width
	h := p.sprite.dest.Height
	dest := rl.NewRectangle(p.start.X, p.start.Y, w, h)
	rl.DrawTexturePro(texture, p.sprite.src, dest,
		rl.NewVector2(dest.Width/2, dest.Height), float32(180-p.trajectory), rl.White)

}

func (e *Enemy) draw() {
	tintColor := rl.White

	if frameCount%8 == 1 && !e.dead {
		e.sprite.frame++
	}
	if e.sprite.frame > 3 {
		e.sprite.frame = 0
	}

	if e.hurtFrames > 0 {
		tintColor = rl.Red
		e.hurtFrames--
	} else {
		tintColor = rl.White
	}
	if e.deathFrames > 0 {
		if e.sprite.rotation < 90 {
			e.sprite.rotation = float32(math.Min(90, float64(e.sprite.rotation)+8))
		}
		e.deathFrames--
	}

	e.sprite.src.X = 368                                                                       // pixel where rest idle starts
	e.sprite.src.X += float32(e.sprite.frame) * float32(math.Abs(float64(e.sprite.src.Width))) // rolling the animation

	rl.DrawTexturePro(texture, e.sprite.src, e.sprite.dest, rl.NewVector2(e.sprite.dest.Width, e.sprite.dest.Height), e.sprite.rotation, tintColor)

	if e.health != e.maxHealth && !e.dead {
		rl.DrawRectangle(int32(e.obj.X), int32(e.obj.Y-10), int32(e.obj.W), 4, rl.Red)
		rl.DrawRectangle(int32(e.obj.X), int32(e.obj.Y-10), int32(int(e.obj.W)*e.health/e.maxHealth), 4, rl.Green)
	}
}

func (e *Enemy) hurt() {
	e.hurtFrames = 16
	e.health -= 1
	if e.health <= 0 {
		e.die()
	}
}

func (e *Enemy) die() {
	e.deathFrames = 32
	e.dead = true
	worldCollisionSpace.Remove(e.obj)
}

func (p *Player) move(dx, dy float64) {
	p.obj.X += dx
	p.obj.Y += dy
	p.obj.Update()
	p.sprite.dest.X = rectFromObj(player.obj).X
	p.sprite.dest.Y = rectFromObj(player.obj).Y
	p.weapon.move(dx, dy)
}

func (p *Player) draw() {
	if frameCount%8 == 1 {
		p.sprite.frame++
	}
	if p.sprite.frame > 3 {
		p.sprite.frame = 0
	}
	var weaponOffset float32 = 0
	if p.moving {
		p.sprite.src.X = 192                                                                       // pixel where run animation starts
		p.sprite.src.X += float32(p.sprite.frame) * float32(math.Abs(float64(p.sprite.src.Width))) // rolling the animation
		weaponOffset = -4
	} else {
		if rl.GetScreenWidth()/2 <= int(rl.GetMouseX()) {
			flipRight(&p.sprite.src)
			playerFlipped = false
		} else {
			flipLeft(&p.sprite.src)
			playerFlipped = true
		}
		p.sprite.src.X = 128                                                                       // pixel where rest idle starts
		p.sprite.src.X += float32(p.sprite.frame) * float32(math.Abs(float64(p.sprite.src.Width))) // rolling the animation
	}
	player.moving = false
	rl.DrawTexturePro(texture, p.sprite.src, p.sprite.dest, rl.NewVector2(p.sprite.dest.Width, p.sprite.dest.Height), 0, rl.White)
	updateFrame := frameCount%8 == 0
	p.weapon.draw(p.sprite.frame, updateFrame, weaponOffset)
}

type lineDrawParam struct {
	x1, y1 int // start
	x2, y2 int // end
	color  rl.Color
}

func (p *Player) attack() {
	rl.PlaySound(swordSound)
	p.attackFrame = 0 // find a better way to trigger animation than this.
	player.attackCooldown = p.weapon.cooldown

	playerCenter := point{
		x: float32(player.obj.X + player.obj.W/2),
		y: float32(player.obj.Y + player.obj.H/2),
	}
	rl.DrawCircleLines(int32(playerCenter.x), int32(playerCenter.y), 32, rl.Green)
	angle := getPlayerToMouseAngleDegress()

	// TODO use weapon attributes in the future to determine this logic
	projectileCount := p.weapon.projectileCount
	projectileReach := p.weapon.projectilelength
	projectileSpread := p.weapon.projectileSpreadDegrees
	projectileTTL := p.weapon.projectileTTLFrames
	projectileVelocity := p.weapon.projectileVelocity
	projectileSpreadItter := int(float64(angle) - math.Floor(float64(projectileCount)/2)*float64(projectileSpread))

	for i := 0; i < projectileCount; i++ {
		x2 := int(float64(projectileReach) * math.Sin(degreesToRadians(float64(projectileSpreadItter))))
		y2 := int(float64(projectileReach) * math.Cos(degreesToRadians(float64(projectileSpreadItter))))
		var projectileTrajectory float64
		if projectileVelocity > 0 {
			projectileTrajectory = float64(projectileSpreadItter)
		}
		projectiles = append(projectiles,
			Projectile{
				start:      rl.NewVector2(playerCenter.x, playerCenter.y),
				end:        rl.NewVector2(playerCenter.x+float32(x2), playerCenter.y+float32(y2)),
				ttl:        projectileTTL,
				velocity:   projectileVelocity,
				trajectory: projectileTrajectory,
				sprite: Sprite{
					src:  p.weapon.projectileSpriteSrc.src,
					dest: p.weapon.projectileSpriteSrc.dest,
				},
			})
		projectileSpreadItter += projectileSpread
	}
	p.attacking = false
}

func (w *Weapon) move(dx, dy float64) {
	w.obj.X += dx
	w.obj.Y += dy
	w.obj.Update()
	w.sprite.dest.X = rectFromObj(w.obj).X
	w.sprite.dest.Y = rectFromObj(w.obj).Y
}

func (w *Weapon) draw(frame int, next_frame bool, offset float32) {
	rotation := w.idleRotation
	if player.attackFrame >= 0 && w.attackRotator != nil {
		rotation = w.attackRotator(*w, player)
		player.attackFrame++
		if player.attackFrame >= w.attackSpeed {
			player.attackFrame = -1 // need to find a better way to manage attack animations
			w.move(0, 0)            // recenter weapon after attack animation
		}

	} else if next_frame {

		if frame == 0 || frame == 1 {
			w.sprite.dest.Y += 1
		} else {
			w.sprite.dest.Y -= 1
		}
	}

	if !playerFlipped {
		flipRight(&w.sprite.src)
	}
	if playerFlipped {
		flipLeft(&w.sprite.src)
		rotation *= -1
	}

	origin := rl.NewVector2(w.handle.x, w.handle.y)
	dest := w.sprite.dest
	dest.Y += offset

	rl.DrawTexturePro(texture, w.sprite.src, dest,
		origin, rotation, w.tintColor)
}

func (p *Player) equipWeapon(w *Weapon) {
	// create new object from updated dest X/Y
	w.sprite.dest.X = p.hand.x + float32(p.obj.X)
	w.sprite.dest.Y = p.hand.y + float32(p.obj.Y)
	w.obj = objFromRect(w.sprite.dest)

	// update player weapon
	p.weapon = w
}

func objFromRect(rect rl.Rectangle) *resolv.Object {
	x, y, w, h := float64(rect.X), float64(rect.Y), float64(rect.Width), float64(rect.Height)
	// Janky fix: x = x-w and y = y-h to resolve difference in rl & resolv packages
	return resolv.NewObject(x-w, y-h, w, h)
}

func rectFromObj(obj *resolv.Object) rl.Rectangle {
	x, y, w, h := float32(obj.X), float32(obj.Y), float32(obj.W), float32(obj.H)
	// Janky fix: x = x+w and y = y+h to resolve difference in rl & resolv packages
	return rl.NewRectangle(x+w, y+h, w, h)
}

func getCameraTarget() rl.Vector2 {
	playerCenterX := float32(player.obj.X + player.obj.W/2)
	playerCenterY := float32(player.obj.Y + player.obj.H/2)
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
	worldCollisionSpace = resolv.NewSpace(spaceWidth, spaceHeight, spaceCellWidth, spaceCellHeight)

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

	worldCollisionSpace.Add(objects...)
}

type point struct {
	x float32
	y float32
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
	tileMapIndex = map[string]map[int]point{
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
	emptyTileMap = map[int]point{
		0: {0, 0},
	}
	floorTileMap = map[int]point{
		1: {16, 64}, // plain floor
		2: {32, 64}, // really cracked floor
		3: {48, 64}, // kinda cracked floor
		4: {16, 80}, // big hole
		5: {32, 80}, // somewhat cracked floor
		6: {48, 80}, // bottomr right hole
		7: {16, 96}, // top right hole
		8: {32, 96}, // top left whole
		9: {48, 96}, // ladder
	}
	wallTileMap = map[int]point{
		1:  {16, 0},   // top left
		2:  {32, 0},   // top mid
		3:  {48, 0},   // top right
		4:  {16, 16},  // left
		5:  {32, 16},  // mid
		6:  {48, 16},  // right
		7:  {0, 112},  // wall_side_top_left 0 112 16 16
		8:  {16, 112}, // wall_side_top_right 16 112 16 16
		9:  {0, 128},  // wall_side_mid_left 0 128 16 16
		10: {16, 128}, // wall_side_mid_right 16 128 16 16
		11: {0, 144},  // wall_side_front_left 0 144 16 16
		12: {16, 144}, // wall_side_front_right 16 144 16 16
		13: {32, 112}, // wall_corner_top_left 32 112 16 16
		14: {48, 112}, // wall_corner_top_right 48 112 16 16
		15: {32, 128}, // wall_corner_left 32 128 16 16
		16: {48, 128}, // wall_corner_right 48 128 16 16
		17: {32, 144}, // wall_corner_bottom_left 32 144 16 16
		18: {48, 144}, // wall_corner_bottom_right 48 144 16 16
		19: {32, 160}, // wall_corner_front_left 32 160 16 16
		20: {48, 160}, // wall_corner_front_right 48 160 16 16
		21: {80, 128}, // wall_inner_corner_l_top_left 80 128 16 16
		22: {64, 128}, // wall_inner_corner_l_top_rigth 64 128 16 16
		23: {80, 144}, // wall_inner_corner_mid_left 80 144 16 16
		24: {64, 144}, // wall_inner_corner_mid_rigth 64 144 16 16
		25: {80, 160}, // wall_inner_corner_t_top_left 80 160 16 16
		26: {64, 160}, // wall_inner_corner_t_top_rigth 64 160 16 16
	}
	decorTileMap = map[int]point{
		1:  {64, 0},   // wall_fountain_top 64 0 16 16
		2:  {64, 16},  // wall_fountain_mid_red_anim 64 16 16 16 3
		3:  {64, 32},  // wall_fountain_basin_red_anim 64 32 16 16 3
		4:  {64, 48},  // wall_fountain_mid_blue_anim 64 48 16 16 3
		5:  {64, 64},  // wall_fountain_basin_blue_anim 64 64 16 16 3
		6:  {16, 32},  // wall_banner_red 16 32 16 16
		7:  {32, 32},  // wall_banner_blue 32 32 16 16
		8:  {16, 48},  // wall_banner_green 16 48 16 16
		9:  {32, 48},  // wall_banner_yellow 32 48 16 16
		10: {96, 80},  // wall_column_top 96 80 16 16
		11: {96, 96},  // wall_column_mid 96 96 16 16
		12: {96, 112}, // wall_coulmn_base 96 112 16 16
		13: {80, 80},  // column_top 80 80 16 16
		14: {80, 96},  // column_mid 80 96 16 16
		15: {80, 112}, // coulmn_base 80 112 16 16
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
		tileSrc.X = tileMap[tileInt].x
		tileSrc.Y = tileMap[tileInt].y

		if strings.ToUpper(srcMap[i]) == srcMap[i] {
			// TODO make this fill square more sofistacted - maybe random or something
			fillTile := tileSrc
			fillTile.X = 16
			fillTile.Y = 64
			rl.DrawTexturePro(texture, fillTile, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)

			// draw behind player if Y is "behind" player, but skip this with Walls
			if tileDest.Y > player.sprite.dest.Y || strings.ToLower(srcMap[i]) == "w" {
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
	player.draw()
	for _, e := range enemies {
		e.draw()
	}
	for _, p := range projectiles {
		p.draw()
	}
	// draw foreground after player
	for _, draw := range foreGround {
		rl.DrawTexturePro(draw.texture, draw.srcRec, draw.destRec, draw.origin, draw.rotation, draw.tint)
	}

	// draw debug collision objects
	if debugMode {
		for _, o := range collisionMapDebug {
			mo := objFromRect(o)
			rl.DrawRectangleLines(int32(mo.X), int32(mo.Y), int32(mo.W), int32(mo.H), rl.White)
		}

		for _, p := range projectiles {
			rl.DrawLine(int32(p.start.X), int32(p.start.Y), int32(p.end.X), int32(p.end.Y), rl.Pink)
		}

		// debug player collision box
		po := player.obj
		rl.DrawRectangleLines(int32(po.X), int32(po.Y), int32(po.W), int32(po.H), rl.Orange)

		for _, e := range enemies {
			rl.DrawRectangleLines(int32(e.obj.X), int32(e.obj.Y), int32(e.obj.W), int32(e.obj.H), rl.White)

		}

		playerCenter := point{
			x: float32(player.obj.X + player.obj.W/2),
			y: float32(player.obj.Y + player.obj.H/2),
		}
		rl.DrawCircleLines(int32(playerCenter.x), int32(playerCenter.y), 32, rl.Green)
		angle := getPlayerToMouseAngleDegress()
		rl.DrawCircleSectorLines(rl.NewVector2(playerCenter.x, playerCenter.y), 32, angle, angle-45, 5, rl.White)
		rl.DrawCircleSectorLines(rl.NewVector2(playerCenter.x, playerCenter.y), 32, angle, angle+45, 5, rl.White)
	}

}

/*
returns degrees mouse is from player
rise/run seem to be flipped because x/y are 90 degrees off in game engines
*/
func getPlayerToMouseAngleDegress() float32 {
	rise := float64(rl.GetMouseX()) - float64(rl.GetScreenWidth()/2)
	run := float64(rl.GetMouseY()) - float64(rl.GetScreenHeight()/2)
	angle := float32(radiansToDegrees(math.Atan(rise / run)))
	if run < 0 {
		angle += 180
	}
	return angle
}

func drawUI() {
	if debugMode {
		rl.DrawRectangleRounded(rl.NewRectangle(3, 3, 500, 90), .1, 10, rl.DarkGray)
		rl.DrawRectangleRoundedLines(rl.NewRectangle(3, 3, 500, 90), .1, 10, 3, rl.White)
		rl.DrawText(fmt.Sprintf("FPS: %v", rl.GetFPS()), 10, 10, 16, rl.White)
		rl.DrawText(fmt.Sprintf("player {X: %v, Y:%v}", player.obj.X, player.obj.Y), 10, 30, 16, rl.White)
		rl.DrawText(fmt.Sprintf("mouse  {X: %v, Y:%v}", rl.GetMouseX(), rl.GetMouseY()), 10, 50, 16, rl.White)

		// wierd thing where rise/run are opposite directions (think it has to do with x/y being negative flipped)
		rise := float64(rl.GetMouseX()) - float64(rl.GetScreenWidth()/2)
		run := float64(rl.GetMouseY()) - float64(rl.GetScreenHeight())/2

		angle := getPlayerToMouseAngleDegress()
		rl.DrawText(fmt.Sprintf("mouse->player  {X: %v, Y:%v}", rise, run), 10, 70, 16, rl.White)
		rl.DrawText(fmt.Sprintf("Atan(%v/%v) = %v degrees", rise, run, int(angle)), 250, 10, 16, rl.White)
		rl.DrawText(fmt.Sprintf("Live Projectiles: %v", len(projectiles)), 250, 30, 16, rl.White)
	}
}

func input() {
	/*
		Thinking about making an inputDiretor module that checks for all input and returns
		A map of input managers that need to be invoked and their corresponding inputs to manage
	*/
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		player.moving = true
		playerUp = true
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		player.moving = true
		playerDown = true
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		player.moving = true
		playerLeft = true
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		player.moving = true
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
		player.attacking = true
	}
	if rl.IsKeyPressed(rl.KeyOne) {
		player.equipWeapon(playerSword)
	} else if rl.IsKeyPressed(rl.KeyTwo) {
		player.equipWeapon(playerBow)
	} else if rl.IsKeyPressed(rl.KeyThree) {
		player.equipWeapon(playerBowThatShootSwords)
	} else if rl.IsKeyPressed(rl.KeyFour) {
		player.equipWeapon(playerSwordThatShootsBows)
	}

}

/*
Sprite image utility - if we don't have assest that face both
direction we can flip them programmatically
example: b -> d
*/
func flipLeft(src *rl.Rectangle) {
	if !(src.Width < 0) {
		src.Width *= -1
	}
}

/*
Sprite image utility - if we don't have assest that face both
direction we can flip them programmatically
example: d -> b
*/
func flipRight(src *rl.Rectangle) {
	if !(src.Width > 0) {
		src.Width *= -1
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

	if player.moving {
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
			flipRight(&player.sprite.src)
			dx += float64(playerSpeed)
			playerFlipped = false
		}
		if playerLeft {
			flipLeft(&player.sprite.src)
			dx -= float64(playerSpeed)
			playerFlipped = true
		}
		// check for collisions
		if collision := player.obj.Check(0, dy, "env"); collision != nil {
			//fmt.Println("Y axis collision happened: ", collision)
			// hueristically stop movement on collision because the other way is buggy
			dy = 0
		}
		if collision := player.obj.Check(dx, 0, "env"); collision != nil {
			//fmt.Println("X axis collision happened: ", collision)
			// hueristically stop movement on collision because the other way is buggy
			dx = 0
		}
		if collision := player.obj.Check(dx, dy, "enemy"); collision != nil {
			dx /= 4
			dy /= 4
		}
		player.move(dx, dy)
	}

	if player.attackCooldown > 0 {
		player.attackCooldown--
		player.attacking = false
	}
	if player.attacking {
		player.attack()
	}

	var nextFrameProjectiles []Projectile

	for _, p := range projectiles {
		if p.ttl > 0 {
			var collision bool
			p.ttl--
			if p.velocity > 0 {
				// find delta x and y
				var dx float32
				var dy float32
				switch p.trajectory {
				case -90:
					dx, dy = float32(-p.velocity), 0 // -x
				case 0:
					dx, dy = 0, float32(p.velocity) // +y
				case 90:
					dx, dy = float32(p.velocity), 0 // +x
				case 180:
					dx, dy = 0, float32(-p.velocity) // -y
				default:
					dx = float32(p.velocity) * float32(math.Sin(degreesToRadians(p.trajectory)))
					dy = float32(p.velocity) * float32(math.Cos(degreesToRadians(p.trajectory)))
				}
				// add delta x and y to both x1,y1 and x2,y2
				p.start.X += dx
				p.start.Y += dy
				p.end.X += dx
				p.end.Y += dy
			}
			for _, e := range enemies {
				for _, line := range linesFromRect(rectFromObj(e.obj)) {
					var collisionPoint *rl.Vector2
					collision = rl.CheckCollisionLines(p.start, p.end, line.start, line.end, collisionPoint)
					if collision {
						break
					}
				}
				if collision {
					e.hurt()
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
	projectiles = nextFrameProjectiles

	var survivingEnemies []*Enemy
	for _, e := range enemies {
		if !e.dead || e.deathFrames > 0 {
			survivingEnemies = append(survivingEnemies, e)
		}
	}
	enemies = survivingEnemies

	frameCount++

	rl.UpdateMusicStream(music)

	if musicPaused {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
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

	texture = rl.LoadTexture("resources/sprites/0x72_DungeonTilesetII_v1.4.png")

	playerSprite := Sprite{
		src:  rl.NewRectangle(128, 100, 16, 28),
		dest: rl.NewRectangle(285, 200, 32, 56),
	}

	playerObj := objFromRect(playerSprite.dest)
	playerObj.H *= .6

	player = Player{
		sprite: playerSprite,
		obj:    playerObj,
		hand:   point{float32(playerObj.W) * .5, float32(playerObj.H) * .94},
	}

	player.sprite.dest.X = rectFromObj(player.obj).X
	player.sprite.dest.Y = rectFromObj(player.obj).Y
	player.obj.AddTags("Player")

	swordSprite := Sprite{
		// src: rl.NewRectangle(307, 26, 10, 21), // rusty sword
		// src: rl.NewRectangle(339, 114, 10, 29), // weapon_knight_sword
		// src: rl.NewRectangle(310, 124, 8, 19), // cleaver
		// src: rl.NewRectangle(325, 113, 9, 30), // weapon_duel_sword
		// src: rl.NewRectangle(322, 81, 12, 30), // weapon_anime_sword
		src: rl.NewRectangle(339, 26, 10, 21), // weapon_red_gem_sword

		dest: rl.NewRectangle(
			player.hand.x+float32(player.obj.X),
			player.hand.y+float32(player.obj.Y),
			10*1.35,
			21*1.35,
		),
	}

	playerSword = &Weapon{
		sprite: swordSprite,
		obj:    objFromRect(swordSprite.dest),
		// handle is the origin offset for the sprite
		handle:       point{swordSprite.dest.Width * .5, swordSprite.dest.Height * .9},
		attackSpeed:  8,
		cooldown:     24,
		idleRotation: -30,
		attackRotator: func(w Weapon, p Player) float32 {
			return w.idleRotation * -3 / float32(w.attackSpeed) * float32(player.attackFrame)
		},
		projectileCount:         3,
		projectileSpreadDegrees: 45,
		projectilelength:        32,
		projectileTTLFrames:     10,
		tintColor:               rl.White,
	}

	bowSprite := Sprite{
		src: rl.NewRectangle(325, 180, 7, 25), // weapon_bow 325 180 7 25

		dest: rl.NewRectangle(
			player.hand.x+float32(player.obj.X),
			player.hand.y+float32(player.obj.Y),
			7*1.1,
			25*1.1,
		),
	}

	arrowSpriteSource := Sprite{
		src:  rl.NewRectangle(308, 186, 7, 21),     // weapon_arrow 308 186 7 21
		dest: rl.NewRectangle(0, 0, 7*1.5, 21*1.5), // only using h, w for scaling
	}

	playerBow = &Weapon{
		sprite:              bowSprite,
		projectileSpriteSrc: arrowSpriteSource,

		obj: objFromRect(bowSprite.dest),
		// handle is the origin offset for the sprite
		handle:       point{bowSprite.dest.Width * .5, bowSprite.dest.Height * .75},
		attackSpeed:  8,
		cooldown:     24,
		idleRotation: 20,
		attackRotator: func(w Weapon, p Player) float32 {
			// TODO make it follow mouse -- for now make it 0
			return 0
		},
		projectileCount:         1,
		projectileSpreadDegrees: 0,
		projectilelength:        21,
		projectileTTLFrames:     32,
		projectileVelocity:      8,
		tintColor:               rl.White,
	}

	bowThatShootsSwordsSprite := Sprite{
		src: rl.NewRectangle(325, 180, 7, 25), // weapon_bow 325 180 7 25

		dest: rl.NewRectangle(
			player.hand.x+float32(player.obj.X),
			player.hand.y+float32(player.obj.Y),
			7*1.1,
			25*1.1,
		),
	}

	swordArrowSpriteSource := Sprite{
		src:  rl.NewRectangle(339, 114, 10, 29),     // weapon_knight_sword 339 114 10 29
		dest: rl.NewRectangle(0, 0, 10*1.5, 29*1.5), // only using h, w for scaling
	}

	playerBowThatShootSwords = &Weapon{
		sprite:              bowThatShootsSwordsSprite,
		projectileSpriteSrc: swordArrowSpriteSource,

		obj: objFromRect(bowSprite.dest),
		// handle is the origin offset for the sprite
		handle:       point{bowSprite.dest.Width * .5, bowSprite.dest.Height * .75},
		attackSpeed:  8,
		cooldown:     24,
		idleRotation: 20,
		attackRotator: func(w Weapon, p Player) float32 {
			// TODO make it follow mouse -- for now make it 0
			return 0
		},
		projectileCount:         1,
		projectileSpreadDegrees: 0,
		projectilelength:        21,
		projectileTTLFrames:     32,
		projectileVelocity:      8,
		tintColor:               rl.Blue,
	}

	swordThatShootsBowsSprite := Sprite{
		src: rl.NewRectangle(322, 81, 12, 30), // weapon_anime_sword 322 81 12 30

		dest: rl.NewRectangle(
			player.hand.x+float32(player.obj.X),
			player.hand.y+float32(player.obj.Y),
			12*1.1,
			30*1.1,
		),
	}

	bowArrowSpriteSource := Sprite{
		src:  rl.NewRectangle(325, 180, 7, 25),     // weapon_bow 325 180 7 25
		dest: rl.NewRectangle(0, 0, 7*1.5, 25*1.5), // only using h, w for scaling
	}

	playerSwordThatShootsBows = &Weapon{
		sprite:              swordThatShootsBowsSprite,
		projectileSpriteSrc: bowArrowSpriteSource,

		obj: objFromRect(swordSprite.dest),
		// handle is the origin offset or the sprite
		handle:       point{swordSprite.dest.Width * .5, swordSprite.dest.Height * .99},
		attackSpeed:  8,
		cooldown:     24,
		idleRotation: -30,
		attackRotator: func(w Weapon, p Player) float32 {
			return w.idleRotation * -3 / float32(w.attackSpeed) * float32(player.attackFrame)
		},
		projectileCount:         3,
		projectileSpreadDegrees: 35,
		projectilelength:        21,
		projectileTTLFrames:     32,
		projectileVelocity:      4,
		tintColor:               rl.White,
	}

	player.weapon = playerSword

	// Test enemy orc_warrior_idle_anim 368 204 16 20 4
	enemySprite := Sprite{
		src:        rl.NewRectangle(368, 204, 16, 24),
		dest:       rl.NewRectangle(250, 250, 32, 48),
		frameCount: 4,
	}

	enemyObj := objFromRect(enemySprite.dest)
	enemyObj.Y += enemyObj.H * .2
	enemyObj.H *= .7
	enemyObj.X += enemyObj.W * .2
	enemyObj.W *= .8
	enemyObj.AddTags("enemy")
	enemy = Enemy{
		sprite:    enemySprite,
		obj:       enemyObj,
		health:    12,
		maxHealth: 12,
	}
	enemies = append(enemies, &enemy)

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("resources/audio/tracks/ting ting.mp3")

	// TODO find a more sophisticated way to handle sound effect
	swordSound = rl.LoadSound("resources/audio/effects/swing-whoosh.mp3")

	// music = rl.LoadMusicStream("resources/audio/Underworld Coffee Shop.mp3")
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(rl.NewVector2(screenWidth/2, screenHeight/2), getCameraTarget(), 0.0, 1.25)
	loadMap("resources/maps/first.map")
	worldCollisionSpace.Add(player.obj, enemyObj)
}

func quit() {
	rl.UnloadTexture(texture)
	rl.UnloadMusicStream(music)
	rl.UnloadSound(swordSound)
	rl.CloseWindow()
}

func main() {
	initialize()
	fmt.Println(test())
	// Each Frame
	for running {
		input()
		update()
		render()
	}

	quit()

}
