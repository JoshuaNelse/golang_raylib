package physics_engine

import (
	"fmt"
	"raylib/playground/director-models/map-model"
	"raylib/playground/model/draw2d/texture-maps"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/resolv"
)

var WorldCollisionSpace *resolv.Space

func SetWorldSpaceCollidables(currentMap *map_model.MapModel) []rl.Rectangle {

	tileDest := rl.Rectangle{
		Height: currentMap.DestTileDimension.Height,
		Width:  currentMap.DestTileDimension.Width,
	}

	// load collision space
	objects := []*resolv.Object{}
	var collisionMapDebug []rl.Rectangle
	for i, col := range currentMap.CollisionMap {
		// probably more perfomant to just skip these "." collisionless tiles.
		if col == "." {
			continue
		}

		x := tileDest.Width * float32(i%currentMap.Width)  // 6 % 5 means x column 1
		y := tileDest.Height * float32(i/currentMap.Width) // 6 % 5 means y row of 1
		oX, oY, oW, oH := texturemaps.CollisionTileOffsetMap[col].GetTileCollisionOffset(x, y, tileDest.Width, tileDest.Height)
		collisionMapDebug = append(collisionMapDebug, rl.NewRectangle(oX, oY, oW, oH))
		// object offset is different from the regular sprite draw
		newObj := resolv.NewObject(float64(oX-oW), float64(oY-oH), float64(oW), float64(oH))
		newObj.AddTags("env")
		objects = append(objects, newObj)

		if col == "@" {
			fmt.Println("We actually created a \"nav\" collision")
			newObj.RemoveTags("env")
			newObj.AddTags("nav")
			newObj.AddTags("doorId-2")
		}

	}

	WorldCollisionSpace.Add(objects...)
	return collisionMapDebug
}
