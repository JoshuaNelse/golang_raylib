package swords

import (
	pointmodel "raylib/playground/director-models/point-model"
	util "raylib/playground/game/utils"
	data2 "raylib/playground/model"
	"raylib/playground/model/draw2d"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func BowShooter() *data2.Weapon {
	s := data2.Sprite{
		Src: rl.NewRectangle(322, 81, 12, 30), // weapon_anime_sword 322 81 12 30

		Dest: rl.NewRectangle(
			0,
			0,
			12*1.1,
			30*1.1,
		),
		Texture: draw2d.Texture,
	}

	ps := data2.Sprite{
		Src:     rl.NewRectangle(325, 180, 7, 25),     // weapon_bow 325 180 7 25
		Dest:    rl.NewRectangle(0, 0, 7*1.5, 25*1.5), // only using h, w for scaling
		Texture: draw2d.Texture,
	}

	return &data2.Weapon{
		Sprite:              s,
		ProjectileSpriteSrc: ps,

		Obj: util.ObjFromRect(s.Dest),
		// handle is the origin offset or the sprite
		Handle:       pointmodel.Point{X: s.Dest.Width * .5, Y: s.Dest.Height * .99},
		AttackSpeed:  8,
		Cooldown:     24,
		IdleRotation: -30,
		AttackRotator: func(w data2.Weapon) float32 {
			return w.IdleRotation * -3 / float32(w.AttackSpeed) * float32(w.AttackFrame)
		},
		ProjectileCount:         3,
		ProjectileSpreadDegrees: 35,
		Projectilelength:        21,
		ProjectileTTLFrames:     32,
		ProjectileVelocity:      4,
		TintColor:               rl.White,

		// stops the weapon from attack at creation
		AttackFrame: -1,
	}
}
