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
