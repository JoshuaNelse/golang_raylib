package staves

import (
	pointmodel "raylib/playground/director-models/point-model"
	util "raylib/playground/game/utils"
	data2 "raylib/playground/model"
	"raylib/playground/model/draw2d"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Keytar() *data2.Weapon {
	// TODO fix it as this thing is broke... it doesn't look right
	s := data2.Sprite{
		Src: rl.NewRectangle(0, 0, 32, 32),

		Dest: rl.NewRectangle(
			0,
			0,
			32*1.5,
			32*1.5,
		),
		Texture: draw2d.Keytar,
	}

	ps := data2.Sprite{
		Src:     rl.NewRectangle(0, 0, 32, 32),
		Dest:    rl.NewRectangle(0, 0, 32, 32),
		Texture: draw2d.MusicNote,
	}

	return &data2.Weapon{
		Sprite:              s,
		ProjectileSpriteSrc: ps,

		Obj: util.ObjFromRect(s.Dest),
		// handle is the origin offset for the sprite
		Handle:       pointmodel.Point{X: s.Dest.Width * .5, Y: s.Dest.Height * .55},
		AttackSpeed:  8,
		Cooldown:     24,
		IdleRotation: 0,
		AttackRotator: func(w data2.Weapon) float32 {
			// return 360 / float32(w.AttackFrame)

			return -32
		},
		ProjectileCount:         1,
		ProjectileSpreadDegrees: 0,
		Projectilelength:        21,
		ProjectileTTLFrames:     32,
		ProjectileVelocity:      8,
		TintColor:               rl.White,

		// stops the weapon from attack at creation
		AttackFrame: -1,
	}
}
