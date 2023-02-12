package util

import (
	"github.com/faiface/pixel"
)

// IsClicked determines if mouse was in the bounds
func IsClicked(mat pixel.Matrix, bounds pixel.Rect, mouse pixel.Vec) bool {
	pt := mat.Unproject(mouse)
	return bounds.Contains(pt)
}

// RectIsClicked is a simpler IsClicked
func RectIsClicked(r pixel.Rect, pt pixel.Vec) bool {
	return pt.X >= r.Min.X && pt.X <= r.Max.X && pt.Y >= r.Min.Y && pt.Y <= r.Max.Y
}

// BuildMapper builds a matrix to map a buffer to a window
func BuildMapper(buffer pixel.Rect, win pixel.Rect) pixel.Matrix {
	return pixel.IM.Moved(win.Center()).ScaledXY(win.Center(), pixel.V(win.W()/buffer.W(), win.H()/buffer.H()))
}
