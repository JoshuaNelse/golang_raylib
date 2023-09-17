package draw_world_engine

import (
	"fmt"
	"raylib/playground/director-models/draw-model"
	"raylib/playground/director-models/map-model"
	"raylib/playground/director-models/point-model"
	"raylib/playground/engines/projectile-engine"
	util "raylib/playground/game/utils"
	"raylib/playground/model"
	"raylib/playground/model/draw2d/texture-maps"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	player            *model.Player
	enemies           *[]*model.Enemy
	currentMap        *map_model.MapModel
	collisionMapDebug []rl.Rectangle
	frameCount        int32
)

func SetCurrentMap(_currentMap *map_model.MapModel) {
	currentMap = _currentMap
}

func SetCollisionMapDebug(_collisionMapDebug []rl.Rectangle) {
	collisionMapDebug = _collisionMapDebug
}

func SetPlayer(_player *model.Player) {
	player = _player
}

func SetEnemies(_enemies *[]*model.Enemy) {
	enemies = _enemies
}

func DrawMapBackground() []draw_model.DrawParams {

	tileSrc := rl.Rectangle{
		Height: currentMap.SrcTileDimension.Height,
		Width:  currentMap.SrcTileDimension.Width,
	}
	tileDest := rl.Rectangle{
		Height: currentMap.DestTileDimension.Height,
		Width:  currentMap.DestTileDimension.Width,
	}

	var foreGroundDrawParams []draw_model.DrawParams
	for i, tileInt := range currentMap.TileMap {
		if tileInt == 0 {
			continue
		}
		tileDest.X = tileDest.Width * float32(i%currentMap.Width) // 6 % 5 means x column 1
		tileDest.Y = tileDest.Width * float32(i/currentMap.Width) // 6 % 5 means y row of 1
		tileMap := texturemaps.TileMapIndex[strings.ToLower(currentMap.SrcMap[i])]
		tileSrc.X = tileMap[tileInt].X
		tileSrc.Y = tileMap[tileInt].Y

		if strings.ToUpper(currentMap.SrcMap[i]) == currentMap.SrcMap[i] {
			// TODO make this fill square more sophisticated - maybe random or something
			fillTile := tileSrc
			fillTile.X = 16
			fillTile.Y = 64
			rl.DrawTexturePro(currentMap.Texture, fillTile, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)

			// draw behind player if Y is "behind" player, but skip this with Walls
			if tileDest.Y > player.Sprite.Dest.Y || strings.ToLower(currentMap.SrcMap[i]) == "w" {
				foreGroundDrawParams = append(
					foreGroundDrawParams,
					draw_model.DrawParams{
						Texture:  currentMap.Texture,
						SrcRec:   tileSrc,
						DestRec:  tileDest,
						Origin:   rl.NewVector2(tileDest.Width, tileDest.Height),
						Rotation: 0,
						Tint:     rl.White,
					})
			} else {
				rl.DrawTexturePro(currentMap.Texture, tileSrc, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)
			}
		} else {
			rl.DrawTexturePro(currentMap.Texture, tileSrc, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)
		}
	}

	return foreGroundDrawParams
}

func DrawScene(debugMode bool) {
	foreGround := DrawMapBackground()

	/*
		Thinking to make a method drawEnv
			players, monsters, and projectiles should all be part of an env array
			Here I could have an interface for SomeEnvObject.Draw()
			This draw method in implemetation could handle drawing sub components / nested images
			for example:
				The Player.Draw() could make sure that the wielded weapon
				is drawn as well as the player itself.
	*/
	player.Draw(frameCount)
	for _, e := range *enemies {
		e.Draw(frameCount)
	}
	for _, p := range projectile_engine.Projectiles {
		p.Draw()
	}
	// drawing foreground after player so it appears "in-front"
	for _, draw := range foreGround {
		rl.DrawTexturePro(draw.Texture, draw.SrcRec, draw.DestRec, draw.Origin, draw.Rotation, draw.Tint)
	}

	// draw debug collision objects
	if debugMode {
		for _, o := range collisionMapDebug {
			mo := util.ObjFromRect(o)
			rl.DrawRectangleLines(int32(mo.X), int32(mo.Y), int32(mo.W), int32(mo.H), rl.White)
		}

		for _, p := range projectile_engine.Projectiles {
			rl.DrawLine(int32(p.Start.X), int32(p.Start.Y), int32(p.End.X), int32(p.End.Y), rl.Pink)
		}

		// debug player collision box
		po := player.Obj
		rl.DrawRectangleLines(int32(po.X), int32(po.Y), int32(po.W), int32(po.H), rl.Orange)

		for _, e := range *enemies {
			rl.DrawRectangleLines(int32(e.Obj.X), int32(e.Obj.Y), int32(e.Obj.W), int32(e.Obj.H), rl.White)

		}

		playerCenter := point_model.Point{
			X: float32(player.Obj.X + player.Obj.W/2),
			Y: float32(player.Obj.Y + player.Obj.H/2),
		}
		rl.DrawCircleLines(int32(playerCenter.X), int32(playerCenter.Y), 32, rl.Green)
		angle := util.GetPlayerToMouseAngleDegrees()
		rl.DrawCircleSectorLines(rl.NewVector2(playerCenter.X, playerCenter.Y), 32, angle, angle-45, 5, rl.White)
		rl.DrawCircleSectorLines(rl.NewVector2(playerCenter.X, playerCenter.Y), 32, angle, angle+45, 5, rl.White)
	}

	frameCount = (frameCount + 1) % 256
}

func DrawUI(debugMode bool) {
	if debugMode {
		rl.DrawRectangleRounded(rl.NewRectangle(3, 3, 500, 90), .1, 10, rl.DarkGray)
		rl.DrawRectangleRoundedLines(rl.NewRectangle(3, 3, 500, 90), .1, 10, 3, rl.White)
		rl.DrawText(fmt.Sprintf("FPS: %v", rl.GetFPS()), 10, 10, 16, rl.White)
		rl.DrawText(fmt.Sprintf("player {X: %v, Y:%v}", player.Obj.X, player.Obj.Y), 10, 30, 16, rl.White)
		rl.DrawText(fmt.Sprintf("mouse  {X: %v, Y:%v}", rl.GetMouseX(), rl.GetMouseY()), 10, 50, 16, rl.White)

		// wierd thing where rise/run are opposite directions (think it has to do with x/y being negative flipped)
		rise := float64(rl.GetMouseX()) - float64(rl.GetScreenWidth()/2)
		run := float64(rl.GetMouseY()) - float64(rl.GetScreenHeight())/2

		angle := util.GetPlayerToMouseAngleDegrees()
		rl.DrawText(fmt.Sprintf("mouse->player  {X: %v, Y:%v}", rise, run), 10, 70, 16, rl.White)
		rl.DrawText(fmt.Sprintf("Atan(%v/%v) = %v degrees", rise, run, int(angle)), 250, 10, 16, rl.White)
		rl.DrawText(fmt.Sprintf("Live Projectiles: %v", len(projectile_engine.Projectiles)), 250, 30, 16, rl.White)
	}
}
