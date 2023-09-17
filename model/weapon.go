package model

import (
	pointmodel "raylib/playground/director-models/point-model"
	util "raylib/playground/game/utils"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/resolv"
)

type Weapon struct {
	Sprite              Sprite
	SpriteFlipped       bool
	ProjectileSpriteSrc Sprite
	Obj                 *resolv.Object
	Handle              pointmodel.Point
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
	w.Sprite.Dest.X = util.RectFromObj(w.Obj).X
	w.Sprite.Dest.Y = util.RectFromObj(w.Obj).Y
}

func (w *Weapon) AnchoredMove(x, y float64) {
	w.Sprite.Dest.X = float32(x)
	w.Sprite.Dest.Y = float32(y)
	w.Obj = util.ObjFromRect(w.Sprite.Dest)
	w.Obj.Update()
}

func (w *Weapon) Draw(frame int, next_frame bool, offset float32) {
	rotation := w.IdleRotation
	if w.AttackFrame >= 0 && w.AttackRotator != nil {
		rotation = w.AttackRotator(*w)
		w.AttackFrame++

		if w.AttackFrame >= w.AttackSpeed {
			w.AttackFrame = -1 // setting to -1 to symbolize attack is finished animating
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
		util.FlipRight(&w.Sprite.Src)
	}
	if w.SpriteFlipped {
		util.FlipLeft(&w.Sprite.Src)
		rotation *= -1
	}

	origin := rl.NewVector2(w.Handle.X, w.Handle.Y)
	dest := w.Sprite.Dest
	dest.Y += offset

	rl.DrawTexturePro(w.Sprite.Texture, w.Sprite.Src, dest,
		origin, rotation, w.TintColor)
}
