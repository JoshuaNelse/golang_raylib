package cannon

import (
	pointmodel "raylib/playground/director-models/point-model"
	util "raylib/playground/game/utils"
	data2 "raylib/playground/model"
	"raylib/playground/model/draw2d"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func PeopleShooter() *data2.Weapon {
	s := data2.Sprite{
		Src: rl.NewRectangle(0, 0, 32, 32), // cannon texture

		Dest: rl.NewRectangle(
			0,
			0,
			32*1.35,
			32*1.35,
		),
		Texture: draw2d.CannonTexture,
	}

	ps := data2.Sprite{
		Src:     rl.NewRectangle(128, 4, 16, 28),       // elf_f_idle_anim 128 4 16 28 4
		Dest:    rl.NewRectangle(0, 0, 16*1.5, 28*1.5), // only using h, w for scaling
		Texture: draw2d.Texture,
	}

	return &data2.Weapon{
		Sprite:              s,
		ProjectileSpriteSrc: ps,

		Obj: util.ObjFromRect(s.Dest),
		// handle is the origin offset or the sprite
		Handle:       pointmodel.Point{X: s.Dest.Width * .55, Y: s.Dest.Height * .7},
		AttackSpeed:  8,
		Cooldown:     24,
		IdleRotation: 0,
		AttackRotator: func(w data2.Weapon) float32 {
			return w.IdleRotation * -3 / float32(w.AttackSpeed) * float32(w.AttackFrame)
		},
		ProjectileCount:         1,
		ProjectileSpreadDegrees: 0,
		Projectilelength:        21,
		ProjectileTTLFrames:     32,
		ProjectileVelocity:      4,
		TintColor:               rl.White,

		// stops the weapon from attack at creation
		AttackFrame: -1,
	}
}
