package swords

import (
	pointmodel "raylib/playground/director-models/point-model"
	util "raylib/playground/game/utils"
	data2 "raylib/playground/model"
	"raylib/playground/model/draw2d"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func RegularSword() *data2.Weapon {
	s := data2.Sprite{
		// TODO make some file that handles this mapping to abstract this out.
		Src: rl.NewRectangle(339, 26, 10, 21), // weapon_red_gem_sword

		Dest: rl.NewRectangle(
			0,
			0,
			10*1.35,
			21*1.35,
		),
		Texture: draw2d.Texture,
	}

	return &data2.Weapon{
		Sprite: s,
		Obj:    util.ObjFromRect(s.Dest),
		// handle is the origin offset for the sprite
		Handle:       pointmodel.Point{X: s.Dest.Width * .5, Y: s.Dest.Height * .9},
		AttackSpeed:  8,
		Cooldown:     24,
		IdleRotation: -30,
		AttackRotator: func(w data2.Weapon) float32 {
			return w.IdleRotation * -3 / float32(w.AttackSpeed) * float32(w.AttackFrame)
		},
		ProjectileCount:         3,
		ProjectileSpreadDegrees: 45,
		Projectilelength:        32,
		ProjectileTTLFrames:     10,
		TintColor:               rl.White,

		// stops the weapon from attack at creation
		AttackFrame: -1,
	}
}
