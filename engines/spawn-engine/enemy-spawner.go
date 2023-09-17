package spawn_engine

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	util "raylib/playground/game/utils"
	"raylib/playground/model"
	"raylib/playground/model/draw2d"
)

func NewEnemy() *model.Enemy {
	// TODO in the future we can handle enemy types like weapon types and load specific ones from model/bestiary/...

	// Test enemy orc_warrior_idle_anim 368 204 16 20 4
	enemySprite := model.Sprite{
		Src:        rl.NewRectangle(368, 204, 16, 24),
		Dest:       rl.NewRectangle(250, 250, 32, 48),
		Texture:    draw2d.Texture,
		FrameCount: 4,
	}

	enemyObj := util.ObjFromRect(enemySprite.Dest)
	enemyObj.Y += enemyObj.H * .2
	enemyObj.H *= .7
	enemyObj.X += enemyObj.W * .2
	enemyObj.W *= .8
	enemyObj.AddTags("enemy")

	return &model.Enemy{
		Sprite:    enemySprite,
		Obj:       enemyObj,
		Health:    12,
		MaxHealth: 12,
	}
}
