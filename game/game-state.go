package game

import rl "github.com/gen2brain/raylib-go/raylib"

var (
	DebugMode bool
	Running   bool
)

func Initialize(debugMode bool) {
	Running = true
	DebugMode = debugMode

	rl.InitWindow(ScreenWidth, ScreenHeight, "Raylib Playground :)")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)
}
