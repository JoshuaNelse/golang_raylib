package main

import (
	"raylib/playground/game"
)

func main() {

	game.Initialize(true)

	// Each Frame
	for game.Running {
		game.ReadPlayerInputs()
		game.Update()
		game.Render()
	}
	game.Quit()
}
