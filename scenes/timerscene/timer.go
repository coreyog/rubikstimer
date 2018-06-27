package timerscene

import (
	"fmt"
	"time"

	"bitbucket.org/coreyog/rubikstimer/scenes"
	"bitbucket.org/coreyog/rubikstimer/util"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

var sevenSegBigAtlas *text.Atlas
var sevenSegSmallAtlas *text.Atlas

var yellowIndicator *imdraw.IMDraw
var greenIndicator *imdraw.IMDraw

var bigSeven *text.Text
var smallSeven *text.Text

var streamerMode bool

type state int

var elapsed float64
var startTime time.Time

const (
	stateWaitingForHold state = iota
	stateWaitingForRelease
	stateRunning
	stateDone
)

var currentState = stateWaitingForHold
var gear *pixel.Sprite

// Init creates the resources for the Timer scene
func Init(win util.LimitedWindow) {
	sevenSegBigFont, err := util.LoadTTF("assets/DSEG7Modern-Bold.ttf", 200)
	if err != nil {
		fmt.Println()
		panic(err)
	}
	sevenSegSmallFont, err := util.LoadTTF("assets/DSEG7Modern-Bold.ttf", 100)
	if err != nil {
		fmt.Println()
		panic(err)
	}
	sevenSegBigAtlas = text.NewAtlas(sevenSegBigFont, text.ASCII)
	sevenSegSmallAtlas = text.NewAtlas(sevenSegSmallFont, text.ASCII)

	bigSeven = text.New(pixel.V(0, 0), sevenSegBigAtlas)
	smallSeven = text.New(pixel.V(0, 0), sevenSegSmallAtlas)

	yellowIndicator = imdraw.New(nil)
	yellowIndicator.Color = colornames.Yellow
	immediatePill(yellowIndicator, win)

	greenIndicator = imdraw.New(nil)
	greenIndicator.Color = colornames.Lime
	immediatePill(greenIndicator, win)

	pic, err := util.LoadPicture("assets/gear.png")
	if err != nil {
		panic(err)
	}

	gear = pixel.NewSprite(pic, pic.Bounds())

	streamerMode = false
	elapsed = 0
}

// Draw updates and renders the Timer scene
func Draw(canvas *pixelgl.Canvas, win util.LimitedWindow, dt *util.DeltaTimer) (change *scenes.SceneType) {
	if streamerMode {
		canvas.Clear(colornames.Magenta)
	} else {
		canvas.Clear(colornames.Black)
	}

	if win.JustPressed(pixelgl.KeyF12) {
		streamerMode = !streamerMode
	}

	switch currentState {
	case stateWaitingForHold:
		yellowIndicator.Draw(canvas)
		if win.Pressed(pixelgl.KeyLeftControl) && win.Pressed(pixelgl.KeyRightControl) {
			currentState = stateWaitingForRelease
		}
		break
	case stateWaitingForRelease:
		if blink(dt) {
			yellowIndicator.Draw(canvas)
		}
		if !win.Pressed(pixelgl.KeyLeftControl) || !win.Pressed(pixelgl.KeyRightControl) {
			currentState = stateRunning
			startTime = time.Now()
		}
		break
	case stateRunning:
		elapsed = time.Since(startTime).Seconds()
		greenIndicator.Draw(canvas)
		if win.Pressed(pixelgl.KeyLeftControl) && win.Pressed(pixelgl.KeyRightControl) {
			currentState = stateDone
		}
	case stateDone:
		if blink(dt) {
			greenIndicator.Draw(canvas)
		}
		if win.Pressed(pixelgl.KeyR) {
			reset()
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
	bigSeven.Draw(canvas, mat)
	mat = mat.Moved(pixel.V(bigSeven.Bounds().W()+30, 0))
	smallSeven.Draw(canvas, mat)

	if !streamerMode && (currentState == stateDone || currentState == stateWaitingForHold) {
		upperRight := pixel.V(canvas.Bounds().W()-25, canvas.Bounds().H()-25)
		mat = pixel.IM.Moved(upperRight)
		gear.Draw(canvas, mat)
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			mat = mat.Moved(gear.Frame().Center().Scaled(-1))
			pt := mat.Unproject(win.MousePosition())
			if gear.Frame().Contains(pt) {
				reset()
				change = new(scenes.SceneType)
				*change = scenes.SettingsScene
			}
		}
	}
	return change
}

func reset() {
	elapsed = 0
	currentState = stateWaitingForHold
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

func immediatePill(imd *imdraw.IMDraw, win util.LimitedWindow) {
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
