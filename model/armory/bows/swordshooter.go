package bows

import (
	pointmodel "raylib/playground/director-models/point-model"
	util "raylib/playground/game/utils"
	data2 "raylib/playground/model"
	"raylib/playground/model/draw2d"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func SwordShooter() *data2.Weapon {
	s := data2.Sprite{
		Src: rl.NewRectangle(325, 180, 7, 25), // weapon_bow 325 180 7 25

		Dest: rl.NewRectangle(
			0,
			0,
			7*1.1,
			25*1.1,
		),
		Texture: draw2d.Texture,
	}

	ps := data2.Sprite{
		Src:     rl.NewRectangle(339, 114, 10, 29),     // weapon_knight_sword 339 114 10 29
		Dest:    rl.NewRectangle(0, 0, 10*1.5, 29*1.5), // only using h, w for scaling
		Texture: draw2d.Texture,
	}

	return &data2.Weapon{
		Sprite:              s,
		ProjectileSpriteSrc: ps,

		Obj: util.ObjFromRect(s.Dest),
		// handle is the origin offset for the sprite
		Handle:       pointmodel.Point{X: s.Dest.Width * .5, Y: s.Dest.Height * .75},
		AttackSpeed:  8,
		Cooldown:     24,
		IdleRotation: 20,
		AttackRotator: func(w data2.Weapon) float32 {
			// TODO make it follow mouse -- for now make it 0
			return 0
		},
		ProjectileCount:         1,
		ProjectileSpreadDegrees: 0,
		Projectilelength:        21,
		ProjectileTTLFrames:     32,
		ProjectileVelocity:      8,
		TintColor:               rl.Blue,

		// stops the weapon from attack at creation
		AttackFrame: -1,
	}

}
