package util

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type GameWindow struct {
	window       *pixelgl.Window
	winBounds    *pixel.Rect
	canvas       *pixelgl.Canvas
	canvasBounds *pixel.Rect
	mapper       pixel.Matrix
}

func NewGameWindowWithCanvas(win *pixelgl.Window, buffer *pixelgl.Canvas) *GameWindow {
	bufBounds := buffer.Bounds()
	return &GameWindow{
		window:       win,
		winBounds:    Ptr(win.Bounds()),
		canvas:       buffer,
		canvasBounds: &bufBounds,
		mapper:       BuildMapper(bufBounds, win.Bounds()),
	}
}

func (win *GameWindow) ResizeAdjust() {
	wBounds := win.window.Bounds()
	if *win.winBounds != wBounds {
		// has the window been resized? adjust the canvas<=>window mapper
		// as well as current winBounds
		win.winBounds = &wBounds
		win.mapper = BuildMapper(*win.canvasBounds, wBounds)
	}
}

func (win *GameWindow) JustPressed(button pixelgl.Button) bool {
	return win.window.JustPressed(button)
}

func (win *GameWindow) Pressed(button pixelgl.Button) bool {
	return win.window.Pressed(button)
}

func (win *GameWindow) Bounds() pixel.Rect {
	return *win.canvasBounds
}

func (win *GameWindow) MousePosition() pixel.Vec {
	wMouse := win.window.MousePosition()
	wMouse = wMouse.Add(win.winBounds.Center())
	return win.mapper.Unproject(wMouse)
}

func (win *GameWindow) Draw() {
	win.canvas.Draw(win.window, win.mapper)
}

// DrawMouse draws a 10x10 golden rod rectangle on mapped mouse position
func (game *GameWindow) DrawMouse() {
	mousePos := game.MousePosition()

	mdebug := imdraw.New(nil)
	mdebug.Color = colornames.Goldenrod
	mdebug.Push(mousePos.Sub(pixel.V(5, 5)), mousePos.Add(pixel.V(5, 5)))
	mdebug.Rectangle(0)
	mdebug.Draw(game.canvas)
}
