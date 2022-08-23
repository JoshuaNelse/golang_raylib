package draw2d

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Texture rl.Texture2D

// TODO load textures but store reference in list to unload at program close
func InitTexture() rl.Texture2D {
	nilTexture := rl.Texture2D{}

	// Lazy load
	if Texture == nilTexture {
		fmt.Println("Am loading texture!")

		Texture = rl.LoadTexture("resources/sprites/0x72_DungeonTilesetII_v1.4.png")
	} else {
		fmt.Println("Not loading texture!")
	}
	return Texture
}

func UnloadTexture() {
	rl.UnloadTexture(Texture)
}
