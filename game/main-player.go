package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	pointModel "raylib/playground/director-models/point-model"
	util "raylib/playground/game/utils"
	"raylib/playground/model"
	"raylib/playground/model/armory/swords"
	"raylib/playground/model/draw2d"
)

var MainPlayer model.Player
var PlayerUp = false
var PlayerDown = false
var PlayerLeft = false
var PlayerRight = false

func LoadMainPlayer() {
	playerSprite := model.Sprite{
		Src:     rl.NewRectangle(128, 100, 16, 28),
		Dest:    rl.NewRectangle(285, 200, 32, 56),
		Texture: draw2d.Texture,
	}

	playerObj := util.ObjFromRect(playerSprite.Dest)
	playerObj.H *= .6

	MainPlayer = model.Player{
		Sprite: playerSprite,
		Obj:    playerObj,
		Hand:   pointModel.Point{X: float32(playerObj.W) * .5, Y: float32(playerObj.H) * .94},
	}

	MainPlayer.Sprite.Dest.X = util.RectFromObj(MainPlayer.Obj).X
	MainPlayer.Sprite.Dest.Y = util.RectFromObj(MainPlayer.Obj).Y
	MainPlayer.Obj.AddTags("Player")
	MainPlayer.EquipWeapon(swords.RegularSword())
}
