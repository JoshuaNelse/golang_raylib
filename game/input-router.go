package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"raylib/playground/engines/audio-engine"
	"raylib/playground/model/armory/bows"
	"raylib/playground/model/armory/cannon"
	"raylib/playground/model/armory/staves"
	"raylib/playground/model/armory/swords"
)

func ReadPlayerInputs() {
	resetPlayer()

	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		MainPlayer.Moving.Up = true
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		MainPlayer.Moving.Down = true
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		MainPlayer.Moving.Left = true
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		MainPlayer.Moving.Right = true
	}
	if rl.IsKeyPressed(rl.KeyQ) {
		audio_engine.ToggleMusic()
	}
	if rl.IsKeyPressed(rl.KeyBackSlash) {
		DebugMode = !DebugMode
	}
	if rl.GetMouseWheelMove() != 0 {
		mouseMove := rl.GetMouseWheelMove()
		if mouseMove > 0 && Camera.Zoom < 2.0 {
			Camera.Zoom = float32(math.Min(2.0, float64(Camera.Zoom+float32(mouseMove)/15)))
		} else if mouseMove < 0 && Camera.Zoom > .75 {
			Camera.Zoom = float32(math.Max(.75, float64(Camera.Zoom+float32(mouseMove)/15)))
		}
	}
	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		MainPlayer.Attacking = true
	}
	if rl.IsKeyPressed(rl.KeyOne) {
		MainPlayer.EquipWeapon(staves.Keytar())
	} else if rl.IsKeyPressed(rl.KeyTwo) {
		MainPlayer.EquipWeapon(bows.RegularBow())
	} else if rl.IsKeyPressed(rl.KeyThree) {
		MainPlayer.EquipWeapon(bows.SwordShooter())
	} else if rl.IsKeyPressed(rl.KeyFour) {
		MainPlayer.EquipWeapon(swords.BowShooter())
	} else if rl.IsKeyPressed(rl.KeyFive) {
		MainPlayer.EquipWeapon(cannon.PeopleShooter())
	} else if rl.IsKeyPressed(rl.KeySix) {
		MainPlayer.EquipWeapon(staves.PizzaShooter())
	}
}

func resetPlayer() {
	MainPlayer.Moving.Up = false
	MainPlayer.Moving.Down = false
	MainPlayer.Moving.Left = false
	MainPlayer.Moving.Right = false
}
