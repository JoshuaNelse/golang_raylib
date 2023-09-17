package map_director

import (
	"raylib/playground/engines/collision-engine"
	"raylib/playground/engines/draw-world-engine"
	"raylib/playground/engines/map-engine"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadMap(mapFile string, texture rl.Texture2D) {

	mapModel := map_engine.LoadMap(mapFile, texture)
	collisionMapDebug := collision_engine.SetWorldSpaceCollideables(mapModel)

	draw_world_engine.SetCurrentMap(mapModel)
	draw_world_engine.SetCollisionMapDebug(collisionMapDebug)
}
