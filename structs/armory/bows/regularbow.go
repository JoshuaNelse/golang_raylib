package bows

import (
	util "raylib/playground/game/utils"
	"raylib/playground/structs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func RegularBow() *structs.Weapon {
	s := structs.Sprite{
		Src: rl.NewRectangle(325, 180, 7, 25), // weapon_bow 325 180 7 25

		Dest: rl.NewRectangle(
			0,
			0,
			7*1.1,
			25*1.1,
		),
	}

	ps := structs.Sprite{
		Src:  rl.NewRectangle(308, 186, 7, 21),     // weapon_arrow 308 186 7 21
		Dest: rl.NewRectangle(0, 0, 7*1.5, 21*1.5), // only using h, w for scaling
	}

	return &structs.Weapon{
		Sprite:              s,
		ProjectileSpriteSrc: ps,

		Obj: util.ObjFromRect(s.Dest),
		// handle is the origin offset for the sprite
		Handle:       structs.Point{X: s.Dest.Width * .5, Y: s.Dest.Height * .75},
		AttackSpeed:  8,
		Cooldown:     24,
		IdleRotation: 20,
		AttackRotator: func(w structs.Weapon) float32 {
			// TODO make it follow mouse -- for now make it 0
			return 0
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
