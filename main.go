//go:generate go-bindata -o embedded/embedded.go -pkg embedded assets/...
package main

import (
	"fmt"
	"time"

	"bitbucket.org/coreyog/rubixtimer/embedded"
	"bitbucket.org/coreyog/rubixtimer/util"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

var sevenSegBigAtlas *text.Atlas
var sevenSegSmallAtlas *text.Atlas

type state int

const (
	stateWaitingForHold state = iota
	stateWaitingForRelease
	stateRunning
	stateDone
)

var currentState = stateWaitingForHold

func main() {
	sevenSegBigFont, err := loadTTF("assets/DSEG7Modern-Bold.ttf", 200)
	if err != nil {
		fmt.Println()
		panic(err)
	}
	sevenSegSmallFont, err := loadTTF("assets/DSEG7Modern-Bold.ttf", 100)
	if err != nil {
		fmt.Println()
		panic(err)
	}
	sevenSegBigAtlas = text.NewAtlas(sevenSegBigFont, text.ASCII)
	sevenSegSmallAtlas = text.NewAtlas(sevenSegSmallFont, text.ASCII)
	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Rubix Timer",
		Bounds: pixel.R(0, 0, 1000, 400),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	bigSeven := text.New(pixel.V(0, 0), sevenSegBigAtlas)
	smallSeven := text.New(pixel.V(0, 0), sevenSegSmallAtlas)

	yellowIndicator := imdraw.New(nil)
	yellowIndicator.Color = colornames.Yellow
	immediatePill(yellowIndicator, win)

	greenIndicator := imdraw.New(nil)
	greenIndicator.Color = colornames.Lime
	immediatePill(greenIndicator, win)

	dt := util.NewDeltaTimer(0)
	backgroundFlip := true
	flipPressed := false
	var startTime time.Time

	for !win.Closed() {
		dt.Tick()
		if backgroundFlip {
			win.Clear(colornames.Black)
		} else {
			win.Clear(colornames.Magenta)
		}

		if win.Pressed(pixelgl.KeyF12) {
			if !flipPressed {
				backgroundFlip = !backgroundFlip
				flipPressed = true
			}
		} else {
			flipPressed = false
		}

		if win.Pressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}

		elapsed := float64(0)

		switch currentState {
		case stateWaitingForHold:
			yellowIndicator.Draw(win)
			if win.Pressed(pixelgl.KeyLeftControl) && win.Pressed(pixelgl.KeyRightControl) {
				currentState = stateWaitingForRelease
			}
			break
		case stateWaitingForRelease:
			if blink(dt) {
				yellowIndicator.Draw(win)
			}
			if !win.Pressed(pixelgl.KeyLeftControl) || !win.Pressed(pixelgl.KeyRightControl) {
				currentState = stateRunning
				startTime = time.Now()
			}
			break
		case stateRunning:
			elapsed = time.Since(startTime).Seconds()
			greenIndicator.Draw(win)
			if win.Pressed(pixelgl.KeyLeftControl) && win.Pressed(pixelgl.KeyRightControl) {
				currentState = stateDone
			}
		case stateDone:
			if blink(dt) {
				greenIndicator.Draw(win)
			}
			if win.Pressed(pixelgl.KeyR) {
				elapsed = 0
				currentState = stateWaitingForHold
			}
		}

		minsec, milli := buildTimer(elapsed)
		bigSeven.Clear()
		smallSeven.Clear()

		fmt.Fprint(bigSeven, minsec)
		fmt.Fprint(smallSeven, milli)

		offCenter := bigSeven.Bounds().Center().Scaled(-1)
		offCenter.X -= 130
		mat := pixel.IM.Moved(win.Bounds().Center()).Moved(offCenter)
		bigSeven.Draw(win, mat)
		mat = mat.Moved(pixel.V(bigSeven.Bounds().W()+30, 0))
		smallSeven.Draw(win, mat)
		win.Update()
	}
}

func buildTimer(t float64) (minsec string, milli string) {
	// return "01234:56789"
	minutes := int(t) / 60
	seconds := int(t) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds), fmt.Sprintf(".%03d", int(t*1000)%1000)
}

func blink(dt *util.DeltaTimer) bool {
	return (int(dt.TotalTime()*1000)/550)%2 == 1
}

func immediatePill(imd *imdraw.IMDraw, win *pixelgl.Window) {
	pt := win.Bounds().Center()
	pt.Y -= 150
	pt.X -= 125

	// lower left
	pt.X -= 100
	pt.Y -= 25
	imd.Push(pt)
	// lower right
	pt.X += 200
	imd.Push(pt)
	// upper right
	pt.Y += 50
	imd.Push(pt)
	// upper left
	pt.X -= 200
	imd.Push(pt)
	imd.Polygon(0)

	// left cap
	pt.Y -= 25
	imd.Push(pt)
	imd.Circle(25, 0)
	// right cap
	pt.X += 200
	imd.Push(pt)
	imd.Circle(25, 0)
}

func loadTTF(path string, size float64) (font.Face, error) {
	rawFont, err := embedded.Asset(path)
	if err != nil {
		return nil, err
	}
	ttfont, err := truetype.Parse(rawFont)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfont, &truetype.Options{
		Size:    size,
		Hinting: font.HintingFull,
	}), nil
}
