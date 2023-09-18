package physics_engine

import (
	"fmt"
	util "raylib/playground/game/utils"
	"raylib/playground/model"
	"strings"
)

func CalculatePlayerMovement(player *model.Player) (bool, string) {
	if player.IsMoving() {
		// used for collision check
		dx := 0.0
		dy := 0.0

		if player.Moving.Up {
			dy -= float64(player.Speed)
		}
		if player.Moving.Down {
			dy += float64(player.Speed)
		}
		if player.Moving.Right {
			util.FlipRight(&player.Sprite.Src)
			dx += float64(player.Speed)
			player.SpriteFlipped = false
		}
		if player.Moving.Left {
			util.FlipLeft(&player.Sprite.Src)
			dx -= float64(player.Speed)
			player.SpriteFlipped = true
		}
		// check for collisions
		if collision := player.Obj.Check(0, dy, "env"); collision != nil {
			//fmt.Println("Y axis collision happened: ", collision)
			// heuristically stop movement on collision because the other way is buggy
			dy = 0
		}
		if collision := player.Obj.Check(dx, 0, "env"); collision != nil {
			//fmt.Println("X axis collision happened: ", collision)
			// heuristically stop movement on collision because the other way is buggy
			dx = 0
		}
		if collision := player.Obj.Check(dx, dy, "enemy"); collision != nil {
			dx /= 4
			dy /= 4
		}
		player.Move(dx, dy)
		return checkForCollisionOnNavigationTile(player, dx, dy)
	}
	return false, ""
}

func checkForCollisionOnNavigationTile(player *model.Player, dx, dy float64) (bool, string) {
	if collision := player.Obj.Check(dx, dy, "nav"); collision != nil {
		fmt.Println("We hit a door!")
		for _, tag := range collision.Objects[0].Tags() {
			if strings.HasPrefix(tag, "doorId") {
				return true, tag
			}
		}
	}
	return false, ""
}
