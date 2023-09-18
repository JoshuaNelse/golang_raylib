package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	map_director "raylib/playground/directors/map-director"
	audio_engine "raylib/playground/engines/audio-engine"
	physics_engine "raylib/playground/engines/physics-engine"
	"raylib/playground/model/draw2d"
)

func Update() {
	Running = !rl.WindowShouldClose()

	collidedWithNavigationTile, _ := physics_engine.CalculatePlayerMovement(&MainPlayer)
	if collidedWithNavigationTile {
		// TODO use doorId from CalculatePlayerMovement to determine new Map
		map_director.LoadMap("resources/maps/second.map", draw2d.Texture)
		MainPlayer.Obj.Space.Remove(MainPlayer.Obj)
		physics_engine.WorldCollisionSpace.Add(MainPlayer.Obj)
	}

	if MainPlayer.AttackCooldown > 0 {
		MainPlayer.AttackCooldown--
		MainPlayer.Attacking = false
	}
	if MainPlayer.Attacking {
		physics_engine.FireProjects(MainPlayer.Attack())
	}

	physics_engine.CalculatePlayerProjectileOutcome(&Enemies)

	audio_engine.UpdateMusicStream()
	UpdateCameraTargetToPlayerLocation()
}
