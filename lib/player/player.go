package player

import (
	"raylib/playground/lib/draw2d"
	"raylib/playground/lib/weapon"

	"github.com/solarlune/resolv"
)

type Test struct {
	Test   string
	Sprite draw2d.Sprite
}

// start code for player logic
type Player struct {
	sprite         draw2d.Sprite
	obj            *resolv.Object
	weapon         *weapon.Weapon
	hand           draw2d.Point
	moving         bool
	attacking      bool
	attackFrame    int
	attackCooldown int
}
