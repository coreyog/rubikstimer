package settingsscene

import (
	"bitbucket.org/coreyog/rubikstimer/scenes"
	"bitbucket.org/coreyog/rubikstimer/util"

	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// Init creates the resources for the Timer scene
func Init(win util.LimitedWindow) {

}

// Draw updates and renders the Timer scene
func Draw(canvas *pixelgl.Canvas, win util.LimitedWindow, dt *util.DeltaTimer) (change *scenes.SceneType) {
	canvas.Clear(colornames.Black)

	return change
}
