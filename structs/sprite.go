package structs

import rl "github.com/gen2brain/raylib-go/raylib"

type Sprite struct {
	Src        rl.Rectangle
	Dest       rl.Rectangle
	Flipped    bool
	FrameCount int
	Frame      int
	Rotation   float32
}
