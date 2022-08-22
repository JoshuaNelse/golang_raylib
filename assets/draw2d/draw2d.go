package draw2d

import rl "github.com/gen2brain/raylib-go/raylib"

type Sprite struct {
	src        rl.Rectangle
	dest       rl.Rectangle
	flipped    bool
	frameCount int
	frame      int
	rotation   float32
}
