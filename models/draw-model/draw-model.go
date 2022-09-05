package drawmodel

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type DrawParams struct {
	Texture  rl.Texture2D
	SrcRec   rl.Rectangle
	DestRec  rl.Rectangle
	Origin   rl.Vector2
	Rotation float32
	Tint     color.RGBA
}
