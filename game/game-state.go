package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	map_director "raylib/playground/directors/map-director"
	"raylib/playground/engines/audio-engine"
	collision_engine "raylib/playground/engines/collision-engine"
	"raylib/playground/engines/draw-world-engine"
	spawn_engine "raylib/playground/engines/spawn-engine"
	"raylib/playground/model"
	"raylib/playground/model/draw2d"
)

var (
	DebugMode bool
	Running   bool

	Enemies []*model.Enemy

	mapFile = "resources/maps/first.map"
)

func Initialize(debugMode bool) {
	Running = true
	DebugMode = debugMode

	rl.InitWindow(ScreenWidth, ScreenHeight, "Raylib Playground :)")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	draw2d.InitTexture()
	audio_engine.InitializeAudio()

	LoadMainPlayer()

	draw_world_engine.SetPlayer(&MainPlayer)

	draw_world_engine.SetEnemies(&Enemies)
	LoadPlayerCamera()

	map_director.LoadMap(mapFile, draw2d.Texture)

	enemy := spawn_engine.NewEnemy()
	Enemies = append(Enemies, enemy)

	collision_engine.WorldCollisionSpace.Add(MainPlayer.Obj, enemy.Obj)
}

func Quit() {
	draw2d.UnloadTexture()
	audio_engine.UnloadAudioComponents()
	rl.CloseWindow()
}

func Render() {
	rl.BeginDrawing()
	rl.ClearBackground(BackgroundColor)
	rl.BeginMode2D(Camera)
	draw_world_engine.DrawScene(DebugMode)

	rl.EndMode2D()
	draw_world_engine.DrawUI(DebugMode)
	rl.EndDrawing()
}
