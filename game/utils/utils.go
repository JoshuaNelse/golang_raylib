package util

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/resolv"
)

// GetPlayerToMouseAngleDegrees
/*
	returns degrees mouse is from player
	rise/run seem to be flipped because x/y are 90 degrees off in game engines
*/
func GetPlayerToMouseAngleDegrees() float32 {
	rise := float64(rl.GetMouseX()) - float64(rl.GetScreenWidth()/2)
	run := float64(rl.GetMouseY()) - float64(rl.GetScreenHeight()/2)
	angle := float32(RadiansToDegrees(math.Atan(rise / run)))
	if run < 0 {
		angle += 180
	}
	return angle
}

func DegreesToRadians(d float64) float64 {
	return d * (math.Pi / 180)
}

func RadiansToDegrees(r float64) float64 {
	return r * (180 / math.Pi)
}

func ObjFromRect(rect rl.Rectangle) *resolv.Object {
	x, y, w, h := float64(rect.X), float64(rect.Y), float64(rect.Width), float64(rect.Height)
	// Janky fix: x = x-w and y = y-h to resolve difference in rl & resolv packages
	return resolv.NewObject(x-w, y-h, w, h)
}

func RectFromObj(obj *resolv.Object) rl.Rectangle {
	x, y, w, h := float32(obj.X), float32(obj.Y), float32(obj.W), float32(obj.H)
	// Janky fix: x = x+w and y = y+h to resolve difference in rl & resolv packages
	return rl.NewRectangle(x+w, y+h, w, h)
}

// FlipLeft
/*
	Sprite image utility - if we don't have assests that face both
	direction we can flip them programmatically
	example: b -> d
*/
func FlipLeft(src *rl.Rectangle) {
	if !(src.Width < 0) {
		src.Width *= -1
	}
}

// FlipRight
/*
	Sprite image utility - if we don't have assests that face both
	direction we can flip them programmatically
	example: d -> b
*/
func FlipRight(src *rl.Rectangle) {
	if !(src.Width > 0) {
		src.Width *= -1
	}

}
