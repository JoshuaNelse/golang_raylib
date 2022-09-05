package mapdirector

import (
	collisionengine "raylib/playground/engines/collision-engine"
	drawworldengine "raylib/playground/engines/draw-world-engine"
	mapengine "raylib/playground/engines/map-engine"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadMap(mapFile string, texture rl.Texture2D) {

	mapModel := mapengine.LoadMap(mapFile, texture)
	collisionMapDebug := collisionengine.SetWorldSpaceCollideables(mapModel)

	drawworldengine.SetCurrentMap(mapModel)
	drawworldengine.SetCollisionMapDebug(collisionMapDebug)
}
