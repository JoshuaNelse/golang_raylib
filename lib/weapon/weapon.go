package weapon

import (
	"raylib/playground/lib/draw2d"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/resolv"
)

type Weapon struct {
	Sprite              draw2d.Sprite
	SpriteFlipped       bool
	ProjectileSpriteSrc draw2d.Sprite
	Obj                 *resolv.Object
	Handle              draw2d.Point
	Reach               int
	AttackSpeed         int
	Cooldown            int
	TintColor           rl.Color
	AttackFrame         int

	IdleRotation  float32
	AttackRotator func(w Weapon) float32

	ProjectileCount         int
	Projectilelength        int
	ProjectileSpreadDegrees int
	ProjectileTTLFrames     int
	ProjectileVelocity      int
}

func (w *Weapon) Move(dx, dy float64) {
	w.Obj.X += dx
	w.Obj.Y += dy
	w.Obj.Update()
	w.Sprite.Dest.X = draw2d.RectFromObj(w.Obj).X
	w.Sprite.Dest.Y = draw2d.RectFromObj(w.Obj).Y
}

func (w *Weapon) Draw(frame int, next_frame bool, offset float32) {
	rotation := w.IdleRotation
	if w.AttackFrame >= 0 && w.AttackRotator != nil {
		rotation = w.AttackRotator(*w)
		w.AttackFrame++
		if w.AttackFrame >= w.AttackSpeed {
			w.AttackFrame = -1 // need to find a better way to manage attack animations
			w.Move(0, 0)       // recenter weapon after attack animation
		}

	} else if next_frame {

		if frame == 0 || frame == 1 {
			w.Sprite.Dest.Y += 1
		} else {
			w.Sprite.Dest.Y -= 1
		}
	}

	if !w.SpriteFlipped {
		draw2d.FlipRight(&w.Sprite.Src)
	}
	if w.SpriteFlipped {
		draw2d.FlipLeft(&w.Sprite.Src)
		rotation *= -1
	}

	origin := rl.NewVector2(w.Handle.X, w.Handle.Y)
	dest := w.Sprite.Dest
	dest.Y += offset

	rl.DrawTexturePro(draw2d.Texture, w.Sprite.Src, dest,
		origin, rotation, w.TintColor)
}
