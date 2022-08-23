package draw2d

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/resolv"
)

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

/*
Sprite image utility - if we don't have assest that face both
direction we can flip them programmatically
example: b -> d
*/
func FlipLeft(src *rl.Rectangle) {
	if !(src.Width < 0) {
		src.Width *= -1
	}
}

/*
Sprite image utility - if we don't have assest that face both
direction we can flip them programmatically
example: d -> b
*/
func FlipRight(src *rl.Rectangle) {
	if !(src.Width > 0) {
		src.Width *= -1
	}

}

type Point struct {
	X float32
	Y float32
}
