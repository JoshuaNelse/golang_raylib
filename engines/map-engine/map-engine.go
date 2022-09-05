package mapengine

import (
	"fmt"
	"io/ioutil"
	"os"
	collisionengine "raylib/playground/engines/collision-engine"
	mapmodel "raylib/playground/models/map-model"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/resolv"
)

func LoadMap(mapFile string, texture rl.Texture2D) *mapmodel.MapModel {
	fmt.Println("Attempting to load map:", mapFile)

	file, err := ioutil.ReadFile(mapFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	remNewLines := strings.Replace(string(file), "\n", " ", -1)
	sliced := strings.Split(remNewLines, " ")

	// map dimensions
	mapW, err := strconv.Atoi(sliced[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mapH, err := strconv.Atoi(sliced[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	srcTileDimension := mapmodel.TileDimension{Width: 16, Height: 16}
	destTileDimension := mapmodel.TileDimension{Width: 32, Height: 32}

	//pixel level collision
	spaceWidth := mapW * int(2*destTileDimension.Width)
	spaceHeight := mapH * int(2*destTileDimension.Height)
	spaceCellWidth := 1
	spaceCellHeight := 1
	collisionengine.WorldCollisionSpace = resolv.NewSpace(spaceWidth, spaceHeight, spaceCellWidth, spaceCellHeight)

	var tileMap []int
	var srcMap []string
	var collisionMap []string

	for i, val := range sliced[2:] {
		if i < mapW*mapH {
			if m, err := strconv.Atoi(val); err != nil {
				fmt.Println(err)
			} else {
				tileMap = append(tileMap, m)
			}
		} else if i < mapW*mapH*2 {
			srcMap = append(srcMap, val)
		} else {
			collisionMap = append(collisionMap, val)
		}
	}

	return &mapmodel.MapModel{
		Width:             mapW,
		Height:            mapH,
		SrcMap:            srcMap,
		TileMap:           tileMap,
		CollisionMap:      collisionMap,
		SrcTileDimension:  srcTileDimension,
		DestTileDimension: destTileDimension,
		Texture:           texture,
	}
}
