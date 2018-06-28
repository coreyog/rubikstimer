package testscene

import (
	"fmt"

	"github.com/coreyog/rubikstimer/scenes"
	"github.com/coreyog/rubikstimer/util"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var gear *pixel.Sprite

// Init initializes resources for the scene
func Init(win util.LimitedWindow) {
	pic, err := util.LoadPicture("assets/gear.png")
	if err != nil {
		panic(err)
	}

	gear = pixel.NewSprite(pic, pic.Bounds())
}

// OnShow has some last minute prep for showing a scene
func OnShow() {}

// Draw updates and renders the Test scene
func Draw(canvas *pixelgl.Canvas, win util.LimitedWindow, dt *util.DeltaTimer) (change *scenes.SceneType) {
	canvas.Clear(colornames.Black)
	mat := pixel.IM.Moved(canvas.Bounds().Center())
	gear.Draw(canvas, mat)
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		mat = mat.Moved(gear.Frame().Center().Scaled(-1))
		pt := mat.Unproject(win.MousePosition())
		if gear.Frame().Contains(pt) {
			fmt.Println("CLICKED")
		}
	}
	return change
}
