package gui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"raylib/playground/game"
	"raylib/playground/model/draw2d"
	texturemaps "raylib/playground/model/draw2d/texture-maps"
)

type Button struct {
	Name      string
	Rectangle rl.Rectangle
	Selected  bool
}

type Asset struct {
	Source    rl.Rectangle
	Rectangle rl.Rectangle
}

// Editor mode buttons
var editModeButtonNames = []string{"Pencil", "Erase", "Area"}
var ClickableEditModeButtons []*Button

// Map Layer Selection Buttons
var layerButtonNames = []string{"Background", "MiddleGround", "Foreground", "All"}
var clickableLayerButtons []*Button

// Collision editor button
var collisionEditorButtonName = "Collision"
var clickableCollisionEditorButtons []*Button

// Show Grid Button
var showGridButtonName = "Grid"
var clickableShowGridButtons []*Button

// Save Button
var saveButtonName = "Save"

var ClickableAssets []*Asset
var SelectedAsset *Asset
var initialized = false

func initializeMapEditGUI() {
	for i, buttonName := range editModeButtonNames {
		newButton := Button{
			Name:      buttonName,
			Rectangle: rl.NewRectangle(8, float32(i*64+8), 46, 46),
		}
		ClickableEditModeButtons = append(ClickableEditModeButtons, &newButton)
	}
	ClickableEditModeButtons[0].Selected = true

	for i, sourceTile := range texturemaps.FloorTileMap {
		baseXOffset := 64 + 4
		xOffset := float32(baseXOffset + 64*(i-1))
		newAsset := Asset{
			Source:    rl.NewRectangle(sourceTile.X, sourceTile.Y, 16, 16),
			Rectangle: rl.NewRectangle(xOffset, game.ScreenHeight-game.ScreenHeight/5+4+game.ScreenHeight/40, 48, 48),
		}
		ClickableAssets = append(ClickableAssets, &newAsset)
	}
}

func DrawMapEditGUI() {
	if !initialized {
		initializeMapEditGUI()
		initialized = true
	}
	// draw sidebar
	rl.DrawRectangleRec(rl.NewRectangle(0, 0, 64, game.ScreenHeight), rl.LightGray)

	// sidebar buttons
	for _, button := range ClickableEditModeButtons {
		drawRoundedButton(*button, 8)
	}

	// draw asset set selection window
	assetSelectionWindowRectangle := rl.NewRectangle(64, game.ScreenHeight-game.ScreenHeight/5, game.ScreenWidth, game.ScreenHeight/40)
	rl.DrawRectangleRec(assetSelectionWindowRectangle, rl.RayWhite)
	rl.DrawRectangleLinesEx(assetSelectionWindowRectangle, 3, rl.LightGray)
	floorTileButton := Button{
		Name:      "FloorTiles",
		Rectangle: rl.NewRectangle(64, game.ScreenHeight-game.ScreenHeight/5, 128, game.ScreenHeight/40),
	}
	anotherTileButton := Button{
		Name:      "AnotherTile",
		Rectangle: rl.NewRectangle(192, game.ScreenHeight-game.ScreenHeight/5, 128, game.ScreenHeight/40),
	}
	drawButton(floorTileButton, 8)
	drawButton(anotherTileButton, 8)

	// draw asset window
	rl.DrawRectangleRec(rl.NewRectangle(64, game.ScreenHeight-game.ScreenHeight/5+game.ScreenHeight/40, game.ScreenWidth, game.ScreenHeight/5), rl.RayWhite)
	for _, asset := range ClickableAssets {
		drawAsset(*asset)
	}

}

func drawRoundedButton(button Button, fontSize int32) {
	textWidth := float32(rl.MeasureText(button.Name, fontSize))
	textX := button.Rectangle.X + (button.Rectangle.Width-textWidth)/2
	textY := button.Rectangle.Y + (button.Rectangle.Height-float32(fontSize))/2
	textColor := rl.Gray
	lineThickness := 2

	if button.Selected {
		textColor = rl.DarkGray
		lineThickness = 3
	}

	rl.DrawRectangleRoundedLines(button.Rectangle, .2, 20, float32(lineThickness), textColor)
	rl.DrawText(button.Name, int32(textX), int32(textY), fontSize, textColor)
}
func drawButton(button Button, fontSize int32) {
	textWidth := float32(rl.MeasureText(button.Name, fontSize))
	textX := button.Rectangle.X + (button.Rectangle.Width-textWidth)/2
	textY := button.Rectangle.Y + (button.Rectangle.Height-float32(fontSize))/2
	textColor := rl.Gray
	lineThickness := 2.0

	if button.Selected {
		textColor = rl.DarkGray
		lineThickness = 3.0
	}

	rl.DrawRectangleLinesEx(button.Rectangle, float32(lineThickness), textColor)
	rl.DrawText(button.Name, int32(textX), int32(textY), fontSize, textColor)
}

func drawAsset(asset Asset) {
	color := rl.White
	if SelectedAsset != nil && asset == *SelectedAsset {
		outlineRectangle := rl.NewRectangle(asset.Rectangle.X-2, asset.Rectangle.Y-2, asset.Rectangle.Width+4, asset.Rectangle.Height+4)
		rl.DrawRectangleLinesEx(outlineRectangle, 2, rl.Blue)
	}
	rl.DrawTexturePro(draw2d.Texture, asset.Source, asset.Rectangle, rl.NewVector2(0, 0), 0, color)
}
