package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"raylib/playground/dev-tools/map-maker/gui"
	draw_model "raylib/playground/director-models/draw-model"
	map_director "raylib/playground/directors/map-director"
	audio_engine "raylib/playground/engines/audio-engine"
	draw_world_engine "raylib/playground/engines/draw-world-engine"
	"raylib/playground/game"
	"raylib/playground/model"
	"raylib/playground/model/draw2d"
)

func NewTile() *Tile {
	return &Tile{
		existingDrawParams: map[draw_model.DrawParams]bool{},
		orderedDrawParams:  []draw_model.DrawParams{},
	}
}

type Tile struct {
	// for constant time lookup if drawParams already exists
	existingDrawParams map[draw_model.DrawParams]bool
	orderedDrawParams  []draw_model.DrawParams
}

func (t *Tile) getDrawParams() []draw_model.DrawParams {
	return t.orderedDrawParams
}

func (t *Tile) draw() {
	for _, param := range t.orderedDrawParams {
		rl.DrawTexturePro(param.Texture, param.SrcRec, param.DestRec, param.Origin, param.Rotation, param.Tint)
	}
}

func (t *Tile) addDrawParam(params draw_model.DrawParams) {
	_, present := t.existingDrawParams[params]
	t.existingDrawParams[params] = true
	if !present {
		t.orderedDrawParams = append(t.orderedDrawParams, params)
	} else {
		//
		existingIndex := -1
		for i, p := range t.orderedDrawParams {
			if p == params {
				existingIndex = i
			}
		}
		if existingIndex >= 0 {
			t.orderedDrawParams = append(t.orderedDrawParams[:existingIndex], t.orderedDrawParams[existingIndex+1:]...)
			t.orderedDrawParams = append(t.orderedDrawParams, params)
		}
	}
}

type GameMap struct {
	// key is str(x) + str(y) of tile for quick lookup
	Tiles map[string]*Tile
}

func (g *GameMap) tileExistsAtDest(dest rl.Rectangle) bool {
	key := fmt.Sprintf("%d|%d", int(dest.X), int(dest.Y))
	_, exists := g.Tiles[key]
	return exists
}

func (g *GameMap) getTileAtDest(dest rl.Rectangle) *Tile {
	key := fmt.Sprintf("%d|%d", int(dest.X), int(dest.Y))
	return g.Tiles[key]
}

func (g *GameMap) addTile(dest rl.Rectangle, tile Tile) {
	key := fmt.Sprintf("%d|%d", int(dest.X), int(dest.Y))
	g.Tiles[key] = &tile
}

func (g *GameMap) deleteTile(dest rl.Rectangle) {
	key := fmt.Sprintf("%d|%d", int(dest.X), int(dest.Y))
	delete(g.Tiles, key)
}

func NewGameMap() *GameMap {
	return &GameMap{
		Tiles: map[string]*Tile{},
	}
}

var mousePosition rl.Vector2
var editorMode string

var gameMap = NewGameMap()
var debugRectangles []*rl.Rectangle
var selectedRectangle *rl.Rectangle

func Initialize(debugMode bool) {
	game.Running = true
	game.DebugMode = debugMode

	mapFile := "resources/maps/first.map"
	rl.InitWindow(game.ScreenWidth, game.ScreenHeight, "Raylib Playground :)")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	game.LoadMainPlayer()
	draw_world_engine.SetPlayer(&game.MainPlayer)
	game.LoadMapEditCamera()
	game.Enemies = []*model.Enemy{}
	draw_world_engine.SetEnemies(&game.Enemies)

	draw2d.InitTexture()
	audio_engine.InitializeAudio()

	map_director.LoadMap(mapFile, draw2d.Texture)

	for i := range [100]int{} {
		for j := range [100]int{} {
			rectangle := rl.NewRectangle(float32(j*32), float32(i*32), 32, 32)
			debugRectangles = append(debugRectangles, &rectangle)
		}
	}
}

func Render() {
	rl.BeginDrawing()
	rl.ClearBackground(game.BackgroundColor)
	rl.BeginMode2D(game.Camera)
	draw_world_engine.DrawScene(game.DebugMode)
	for _, tile := range gameMap.Tiles {
		tile.draw()
	}

	for _, rectangle := range debugRectangles {
		if rectangle == selectedRectangle {
			rl.DrawRectangleLines(int32(rectangle.X), int32(rectangle.Y), int32(rectangle.Width), int32(rectangle.Height), rl.Green)
		} else {
			rl.DrawRectangleLines(int32(rectangle.X), int32(rectangle.Y), int32(rectangle.Width), int32(rectangle.Height), rl.White)
		}
	}

	rl.EndMode2D()
	draw_world_engine.DrawUI(game.DebugMode)
	if !game.DebugMode {
		gui.DrawMapEditGUI()
	}
	rl.EndDrawing()
}

func Update() {
	mousePosition = rl.GetMousePosition()
	for _, rectangle := range debugRectangles {
		offsetMouse := rl.NewVector2(mousePosition.X, mousePosition.Y)
		if rl.CheckCollisionPointRec(offsetMouse, *rectangle) {
			selectedRectangle = rectangle
		}
	}
	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		for _, button := range gui.ClickableEditModeButtons {
			if rl.CheckCollisionPointRec(mousePosition, button.Rectangle) {
				for _, b := range gui.ClickableEditModeButtons {
					b.Selected = false
				}
				button.Selected = true
				editorMode = button.Name
				return
			}
		}
		for _, asset := range gui.ClickableAssets {
			if rl.CheckCollisionPointRec(mousePosition, asset.Rectangle) {
				gui.SelectedAsset = asset
				return
			}
		}

		if gui.SelectedAsset != nil {
			switch editorMode {
			case "Pencil":
				pencil()
			case "Erase":
				erase()
			default:
				fmt.Println("not implemented")

			}

		}
	}

	if rl.IsKeyPressed(rl.KeyBackSlash) {
		game.DebugMode = !game.DebugMode
	}

	//if rl.GetMouseWheelMove() != 0 {
	//	mouseMove := rl.GetMouseWheelMove()
	//	if mouseMove > 0 && game.Camera.Zoom < 2.0 {
	//		game.Camera.Zoom = float32(math.Min(2.0, float64(game.Camera.Zoom+float32(mouseMove)/15)))
	//	} else if mouseMove < 0 && game.Camera.Zoom > .75 {
	//		game.Camera.Zoom = float32(math.Max(.75, float64(game.Camera.Zoom+float32(mouseMove)/15)))
	//	}
	//}
}

func pencil() {
	newDrawParam := draw_model.DrawParams{
		Texture:  draw2d.Texture,
		SrcRec:   gui.SelectedAsset.Source,
		DestRec:  *selectedRectangle,
		Origin:   rl.NewVector2(0, 0),
		Rotation: 0,
		Tint:     rl.White,
	}
	if gameMap.tileExistsAtDest(*selectedRectangle) {
		existingTile := gameMap.getTileAtDest(*selectedRectangle)
		existingTile.addDrawParam(newDrawParam)
	} else {
		tile := NewTile()
		tile.addDrawParam(newDrawParam)
		gameMap.addTile(*selectedRectangle, *tile)
	}
}

func erase() {
	gameMap.deleteTile(*selectedRectangle)
}

func main() {
	Initialize(true)

	for game.Running {
		Update()
		Render()
	}
}
