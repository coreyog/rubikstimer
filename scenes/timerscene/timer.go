package timerscene

import (
	"fmt"
	"strings"
	"time"

	"github.com/coreyog/rubikstimer/config"
	"github.com/coreyog/rubikstimer/scenes"
	"github.com/coreyog/rubikstimer/util"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

var yellowIndicator *imdraw.IMDraw
var greenIndicator *imdraw.IMDraw

var bigSeven *text.Text
var smallSeven *text.Text
var galderLine1 *text.Text
var galderLine2 *text.Text

var streamerMode bool

type state int

var elapsed float64
var startTime time.Time
var scramble string
var lastStateChange time.Time

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
	bigSevenAtlas := util.LoadTTF("assets/DSEG7Modern-Bold.ttf", 200)
	smallSevenAtlas := util.LoadTTF("assets/DSEG7Modern-Bold.ttf", 100)
	galderAtlas := util.LoadTTF("assets/galderglynn titling rg.ttf", 30)

	bigSeven = text.New(pixel.V(0, 0), bigSevenAtlas)
	smallSeven = text.New(pixel.V(0, 0), smallSevenAtlas)
	galderLine1 = text.New(pixel.V(0, 0), galderAtlas)
	galderLine2 = text.New(pixel.V(0, 0), galderAtlas)

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

	doScramble()

	streamerMode = false
	elapsed = 0
	lastStateChange = time.Now()
}

// OnShow has some last minute prep for showing a scene
func OnShow() {
	elapsed = 0
	currentState = stateWaitingForHold

	if len(strings.Split(scramble, " ")) != config.GlobalConfig().ScrambleLength {
		doScramble()
	}
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

	// half second delay added after state change for situations like
	// using "any" to start the timer but pressing R to restart will count
	switch currentState {
	case stateWaitingForHold:
		yellowIndicator.Draw(canvas)
		if time.Since(lastStateChange).Seconds() > 0.5 && checkTriggerDown(win, config.GlobalConfig().TimerStartTrigger) {
			currentState = stateWaitingForRelease
			lastStateChange = time.Now()
		}

		if config.GlobalConfig().TimerStartTrigger != string(config.TriggerAny) && win.JustPressed(pixelgl.KeyR) {
			doScramble()
		}
	case stateWaitingForRelease:
		if blink(dt) {
			yellowIndicator.Draw(canvas)
		}

		if checkTriggerUp(win, config.GlobalConfig().TimerStartTrigger) {
			currentState = stateRunning
			startTime = time.Now()
			lastStateChange = time.Now()
		}
	case stateRunning:
		elapsed = time.Since(startTime).Seconds()
		greenIndicator.Draw(canvas)
		if time.Since(lastStateChange).Seconds() > 0.5 && checkTriggerDown(win, config.GlobalConfig().TimerEndTrigger) {
			currentState = stateDone
			lastStateChange = time.Now()
		}
	case stateDone:
		if blink(dt) {
			greenIndicator.Draw(canvas)
		}

		if time.Since(lastStateChange).Seconds() > 0.5 && win.Pressed(pixelgl.KeyR) {
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

	mat = pixel.IM.Moved(galderLine1.Bounds().Center().Scaled(-1)).Moved(pixel.V(win.Bounds().W()/2, win.Bounds().H()-45))
	if config.GlobalConfig().ScrambleLength > 15 {
		mat = mat.Moved(pixel.V(0, 20))
		galderLine1.Draw(canvas, mat)

		mat = pixel.IM.Moved(galderLine2.Bounds().Center().Scaled(-1)).Moved(pixel.V(win.Bounds().W()/2, win.Bounds().H()-45)).Moved(pixel.V(0, -18))
		galderLine2.Draw(canvas, mat)
	} else {
		galderLine1.Draw(canvas, mat)
	}

	if !streamerMode && (currentState == stateDone || currentState == stateWaitingForHold) {
		mat = pixel.IM.Moved(pixel.V(canvas.Bounds().W(), 0)).Moved(pixel.V(-25, 25))
		gear.Draw(canvas, mat)

		halfX := gear.Frame().W() / 2
		halfY := gear.Frame().H() / 2
		if win.JustPressed(pixelgl.MouseButtonLeft) && util.IsClicked(mat, pixel.R(-halfX, -halfY, halfX, halfY), win.MousePosition()) {
			change = util.Ptr(scenes.SettingsScene)
		}
	}

	return change
}

func checkTriggerDown(win util.LimitedWindow, t string) (fired bool) {
	switch config.Trigger(t) {
	case config.TriggerModifiers:
		return isLeftModifierPressed(win) && isRightModifierPressed(win)
	case config.TriggerSpacebar:
		return win.Pressed(pixelgl.KeySpace)
	case config.TriggerAny:
		for k := int(pixelgl.KeySpace); k < int(pixelgl.KeyLast); k++ {
			if win.Pressed(pixelgl.Button(k)) {
				return true
			}
		}
	}
	return false
}

func isLeftModifierPressed(win util.LimitedWindow) bool {
	return win.Pressed(pixelgl.KeyLeftShift) || win.Pressed(pixelgl.KeyLeftControl) || win.Pressed(pixelgl.KeyLeftSuper) || win.Pressed(pixelgl.KeyLeftAlt)
}

func isRightModifierPressed(win util.LimitedWindow) bool {
	return win.Pressed(pixelgl.KeyRightShift) || win.Pressed(pixelgl.KeyRightControl) || win.Pressed(pixelgl.KeyRightSuper) || win.Pressed(pixelgl.KeyRightAlt)
}

func checkTriggerUp(win util.LimitedWindow, t string) (fired bool) {
	switch config.Trigger(t) {
	case config.TriggerModifiers:
		return !isLeftModifierPressed(win) || !isRightModifierPressed(win)
	case config.TriggerSpacebar:
		return !win.Pressed(pixelgl.KeySpace)
	case config.TriggerAny:
		for k := int(pixelgl.KeySpace); k < int(pixelgl.KeyLast); k++ {
			if win.Pressed(pixelgl.Button(k)) {
				return false
			}
		}

		return true
	}

	return false
}

func reset() {
	elapsed = 0
	doScramble()
	currentState = stateWaitingForHold
	lastStateChange = time.Now()
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

func doScramble() {
	scramble = util.Scramble()

	galderLine1.Clear()
	galderLine2.Clear()

	if config.GlobalConfig().ScrambleLength > 15 {
		line1, line2 := splitScramble(scramble)
		fmt.Fprint(galderLine1, line1)
		fmt.Fprint(galderLine2, line2)
	} else {
		fmt.Fprint(galderLine1, scramble)
	}
}

func splitScramble(scramble string) (line1 string, line2 string) {
	split := -1
	count := config.GlobalConfig().ScrambleLength / 2
	if config.GlobalConfig().ScrambleLength%2 == 0 {
		count--
	}

	for i, r := range scramble {
		if r == ' ' {
			if count == 0 {
				split = i
				break
			}

			count--
		}
	}

	return scramble[:split], scramble[split+1:]
}
