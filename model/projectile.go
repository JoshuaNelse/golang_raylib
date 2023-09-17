package model

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Projectile struct {
	Start rl.Vector2
	End   rl.Vector2
	Ttl   int

	// for something like an arrow perhaps
	Velocity   int
	Trajectory float64 //degrees
	Sprite     Sprite

	// sender     *interface{} at somepoint this would be good to have
}

func (p *Projectile) Draw() {
	w := p.Sprite.Dest.Width
	h := p.Sprite.Dest.Height
	dest := rl.NewRectangle(p.Start.X, p.Start.Y, w, h)
	rl.DrawTexturePro(p.Sprite.Texture, p.Sprite.Src, dest,
		rl.NewVector2(dest.Width/2, dest.Height), float32(180-p.Trajectory), rl.White)

}
