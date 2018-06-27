package util

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// LimitedWindow has all the utility without any of the drawing ability
type LimitedWindow interface {
	// NOTE: As I need things that are on pixelgl.Window,
	// add them here, but avoid adding any Draw functions.
	// This is used by scenes so they can access input,
	// window dimensions, etc. but scenes should draw to
	// the provided canvas and not directly to the window
	JustPressed(pixelgl.Button) bool
	Pressed(pixelgl.Button) bool
	Bounds() pixel.Rect
	MousePosition() pixel.Vec
}
