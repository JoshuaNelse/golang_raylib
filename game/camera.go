package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var Camera rl.Camera2D

func LoadPlayerCamera() {
	Camera = rl.NewCamera2D(rl.NewVector2(ScreenWidth/2, ScreenHeight/2), getPlayerLocationAsCameraTarget(), 0.0, 1.25)
}

func LoadMapEditCamera() {
	Camera = rl.NewCamera2D(rl.NewVector2(ScreenWidth/2, ScreenHeight/2), rl.NewVector2(ScreenWidth/2, ScreenHeight/2), 0.0, 1.0)
}

func getPlayerLocationAsCameraTarget() rl.Vector2 {
	playerCenterX := float32(MainPlayer.Obj.X + MainPlayer.Obj.W/2)
	playerCenterY := float32(MainPlayer.Obj.Y + MainPlayer.Obj.H/2)
	return rl.NewVector2(playerCenterX, playerCenterY)
}

func UpdateCameraTargetToPlayerLocation() {
	Camera.Target = getPlayerLocationAsCameraTarget()
}
