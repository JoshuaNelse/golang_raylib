package texturemaps

import (
	pointmodel "raylib/playground/models/point-model"
)

// TODO move this out - not sure where yet
type CollisionOffset struct {
	// percentage offset for collision (left, Right, Top, Bottom)
	L float32
	R float32
	T float32
	B float32
}

/*
Right and Left offset have to consider each other
Top and Bottom offset have to consider each other
*/
func (c CollisionOffset) GetTileCollisionOffset(x, y, w, h float32) (float32, float32, float32, float32) {
	offsetW := w*c.R - (w * (1 - c.L))
	offsetH := h*c.B - (h * (1 - c.T))
	offsetX := x - w*(1-c.R)
	offsetY := y - h*(1-c.B)
	return offsetX, offsetY, offsetW, offsetH
}

var (
	TileMapIndex = map[string]map[int]pointmodel.Point{
		"_": EmptyTileMap,
		"f": FloorTileMap,
		"w": WallTileMap,
		"d": DecorTileMap,
		"n": NavigationTileMap,
	}
	CollisionTileOffsetMap = map[string]CollisionOffset{
		"+": {1, 1, 1, 1},
		".": {0, 0, 0, 0}, // probably more efficient to just skip writting a collision here
		"d": {L: .8, R: .8, T: .3, B: 1},
		"@": {1, 1, 1, 1},
	}
	EmptyTileMap = map[int]pointmodel.Point{
		0: {X: 0, Y: 0},
	}
	FloorTileMap = map[int]pointmodel.Point{
		1: {X: 16, Y: 64}, // plain floor
		2: {X: 32, Y: 64}, // really cracked floor
		3: {X: 48, Y: 64}, // kinda cracked floor
		4: {X: 16, Y: 80}, // big hole
		5: {X: 32, Y: 80}, // somewhat cracked floor
		6: {X: 48, Y: 80}, // bottomr right hole
		7: {X: 16, Y: 96}, // top right hole
		8: {X: 32, Y: 96}, // top left whole
		9: {X: 48, Y: 96}, // ladder
	}
	WallTileMap = map[int]pointmodel.Point{
		1:  {X: 16, Y: 0},   // top left
		2:  {X: 32, Y: 0},   // top mid
		3:  {X: 48, Y: 0},   // top right
		4:  {X: 16, Y: 16},  // left
		5:  {X: 32, Y: 16},  // mid
		6:  {X: 48, Y: 16},  // right
		7:  {X: 0, Y: 112},  // wall_side_top_left 0 112 16 16
		8:  {X: 16, Y: 112}, // wall_side_top_right 16 112 16 16
		9:  {X: 0, Y: 128},  // wall_side_mid_left 0 128 16 16
		10: {X: 16, Y: 128}, // wall_side_mid_right 16 128 16 16
		11: {X: 0, Y: 144},  // wall_side_front_left 0 144 16 16
		12: {X: 16, Y: 144}, // wall_side_front_right 16 144 16 16
		13: {X: 32, Y: 112}, // wall_corner_top_left 32 112 16 16
		14: {X: 48, Y: 112}, // wall_corner_top_right 48 112 16 16
		15: {X: 32, Y: 128}, // wall_corner_left 32 128 16 16
		16: {X: 48, Y: 128}, // wall_corner_right 48 128 16 16
		17: {X: 32, Y: 144}, // wall_corner_bottom_left 32 144 16 16
		18: {X: 48, Y: 144}, // wall_corner_bottom_right 48 144 16 16
		19: {X: 32, Y: 160}, // wall_corner_front_left 32 160 16 16
		20: {X: 48, Y: 160}, // wall_corner_front_right 48 160 16 16
		21: {X: 80, Y: 128}, // wall_inner_corner_l_top_left 80 128 16 16
		22: {X: 64, Y: 128}, // wall_inner_corner_l_top_rigth 64 128 16 16
		23: {X: 80, Y: 144}, // wall_inner_corner_mid_left 80 144 16 16
		24: {X: 64, Y: 144}, // wall_inner_corner_mid_rigth 64 144 16 16
		25: {X: 80, Y: 160}, // wall_inner_corner_t_top_left 80 160 16 16
		26: {X: 64, Y: 160}, // wall_inner_corner_t_top_rigth 64 160 16 16
	}
	DecorTileMap = map[int]pointmodel.Point{
		1:  {X: 64, Y: 0},   // wall_fountain_top 64 0 16 16
		2:  {X: 64, Y: 16},  // wall_fountain_mid_red_anim 64 16 16 16 3
		3:  {X: 64, Y: 32},  // wall_fountain_basin_red_anim 64 32 16 16 3
		4:  {X: 64, Y: 48},  // wall_fountain_mid_blue_anim 64 48 16 16 3
		5:  {X: 64, Y: 64},  // wall_fountain_basin_blue_anim 64 64 16 16 3
		6:  {X: 16, Y: 32},  // wall_banner_red 16 32 16 16
		7:  {X: 32, Y: 32},  // wall_banner_blue 32 32 16 16
		8:  {X: 16, Y: 48},  // wall_banner_green 16 48 16 16
		9:  {X: 32, Y: 48},  // wall_banner_yellow 32 48 16 16
		10: {X: 96, Y: 80},  // wall_column_top 96 80 16 16
		11: {X: 96, Y: 96},  // wall_column_mid 96 96 16 16
		12: {X: 96, Y: 112}, // wall_coulmn_base 96 112 16 16
		13: {X: 80, Y: 80},  // column_top 80 80 16 16
		14: {X: 80, Y: 96},  // column_mid 80 96 16 16
		15: {X: 80, Y: 112}, // coulmn_base 80 112 16 16
	}
	NavigationTileMap = map[int]pointmodel.Point{
		1: {X: 304, Y: 288}, // chest_empty_open_anim 304 288 16 16 3
	}
)
