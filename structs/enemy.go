package structs

import (
	"math"
	collisionengine "raylib/playground/engines/collision-engine"
	"raylib/playground/game"
	"raylib/playground/structs/draw2d"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/resolv"
)

type Enemy struct {
	Sprite      Sprite
	Obj         *resolv.Object
	Health      int
	MaxHealth   int
	HurtFrames  int
	DeathFrames int
	Dead        bool
	TintColor   rl.Color
}

func (e *Enemy) Draw() {

	if game.FrameCount%8 == 1 && !e.Dead {
		e.Sprite.Frame++
	}
	if e.Sprite.Frame > 3 {
		e.Sprite.Frame = 0
	}

	if e.HurtFrames > 0 {
		e.TintColor = rl.Red
		e.HurtFrames--
	} else {
		e.TintColor = rl.White
	}
	if e.DeathFrames > 0 {
		if e.Sprite.Rotation < 90 {
			e.Sprite.Rotation = float32(math.Min(90, float64(e.Sprite.Rotation)+8))
		}
		e.DeathFrames--
	}

	e.Sprite.Src.X = 368                                                                       // pixel where rest idle starts
	e.Sprite.Src.X += float32(e.Sprite.Frame) * float32(math.Abs(float64(e.Sprite.Src.Width))) // rolling the animation

	rl.DrawTexturePro(draw2d.Texture, e.Sprite.Src, e.Sprite.Dest, rl.NewVector2(e.Sprite.Dest.Width, e.Sprite.Dest.Height), e.Sprite.Rotation, e.TintColor)

	if e.Health != e.MaxHealth && !e.Dead {
		rl.DrawRectangle(int32(e.Obj.X), int32(e.Obj.Y-10), int32(e.Obj.W), 4, rl.Red)
		rl.DrawRectangle(int32(e.Obj.X), int32(e.Obj.Y-10), int32(int(e.Obj.W)*e.Health/e.MaxHealth), 4, rl.Green)
	}
}

func (e *Enemy) Hurt() {
	e.HurtFrames = 16
	e.Health -= 1
	if e.Health <= 0 {
		e.Die()
	}
}

func (e *Enemy) Die() {
	e.DeathFrames = 32
	e.Dead = true
	// Really don't like calling an engine from struct. Need to find a better way
	collisionengine.WorldCollisionSpace.Remove(e.Obj)
}
