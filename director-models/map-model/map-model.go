package map_model

import rl "github.com/gen2brain/raylib-go/raylib"

type TileDimension struct {
	Height float32
	Width  float32
}

type MapModel struct {
	Width             int
	Height            int
	TileMap           []int
	SrcMap            []string
	CollisionMap      []string
	SrcTileDimension  TileDimension
	DestTileDimension TileDimension
	Texture           rl.Texture2D
}
