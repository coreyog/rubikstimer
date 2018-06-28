package settingsscene

import (
	"fmt"
	"math"
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

var arrowSize = float64(10)

var startTrigger *text.Text
var endTrigger *text.Text
var scrambleLength *text.Text
var controlsText *text.Text
var spacebarText *text.Text
var anykeyText *text.Text
var scrambleNum *text.Text
var saveBtn *text.Text
var cancelBtn *text.Text

var saveMatrix pixel.Matrix
var cancelMatrix pixel.Matrix
var saveRect pixel.Rect
var cancelRect pixel.Rect

var leftStartArrow *imdraw.IMDraw
var rightStartArrow *imdraw.IMDraw
var leftEndArrow *imdraw.IMDraw
var rightEndArrow *imdraw.IMDraw
var tickMarks *imdraw.IMDraw
var leftScrambleArrow *imdraw.IMDraw
var rightScrambleArrow *imdraw.IMDraw
var saveBox *imdraw.IMDraw
var cancelBox *imdraw.IMDraw

var leftStartArrowRect pixel.Rect
var rightStartArrowRect pixel.Rect
var leftEndArrowRect pixel.Rect
var rightEndArrowRect pixel.Rect
var leftScrambleArrowRect pixel.Rect
var rightScrambleArrowRect pixel.Rect

var leftScramblePt pixel.Vec
var lastScrambleAdjust time.Time
var indicatorGrab bool

var triggerOrder = []string{string(config.TriggerControls), string(config.TriggerSpacebar), string(config.TriggerAny)}
var labelOrder []*text.Text
var startTriggerIndex int
var endTriggerIndex int

var tempConfig config.Config

// Init creates the resources for the Timer scene
func Init(win util.LimitedWindow) {
	tempConfig = config.GlobalConfig()
	indicatorGrab = false

	galderHeaderAtlas := util.LoadTTF("assets/galderglynn titling rg.ttf", 32)
	galderAtlas := util.LoadTTF("assets/galderglynn titling rg.ttf", 18)
	startTrigger = text.New(pixel.ZV, galderHeaderAtlas)
	fmt.Fprint(startTrigger, "Start Trigger")

	endTrigger = text.New(pixel.ZV, galderHeaderAtlas)
	fmt.Fprint(endTrigger, "End Trigger")

	scrambleLength = text.New(pixel.ZV, galderHeaderAtlas)
	fmt.Fprint(scrambleLength, "Scramble Length")

	controlsText = text.New(pixel.ZV, galderAtlas)
	fmt.Fprint(controlsText, "Controls")

	spacebarText = text.New(pixel.ZV, galderAtlas)
	fmt.Fprint(spacebarText, "Spacebar")

	anykeyText = text.New(pixel.ZV, galderAtlas)
	fmt.Fprint(anykeyText, "Any Key")

	labelOrder = []*text.Text{controlsText, spacebarText, anykeyText}

	scrambleNum = text.New(pixel.ZV, galderAtlas)

	saveBtn = text.New(pixel.ZV, galderHeaderAtlas)
	fmt.Fprint(saveBtn, "Save")

	saveBox = imdraw.New(nil)
	saveMatrix = pixel.IM.Moved(pixel.V(win.Bounds().W(), 0)).Moved(saveBtn.Bounds().Center().ScaledXY(pixel.V(-1.5, -1))).Moved(pixel.V(-45, 45))
	saveBox.SetMatrix(saveMatrix)
	saveRect = boxText(saveBtn, saveBox)

	cancelBtn = text.New(pixel.ZV, galderHeaderAtlas)
	fmt.Fprint(cancelBtn, "Cancel")

	cancelBox = imdraw.New(nil)
	cancelMatrix = pixel.IM.Moved(pixel.V(win.Bounds().W()-saveBtn.Bounds().W()/2, 0)).Moved(saveBtn.Bounds().Center().Scaled(-1)).Moved(pixel.V(-25, 45)).Moved(pixel.V(-cancelBtn.Bounds().W()-30, 0))
	cancelBox.SetMatrix(cancelMatrix)
	cancelRect = boxText(cancelBtn, cancelBox)

	leftStartArrow = imdraw.New(nil)
	rightStartArrow = imdraw.New(nil)

	pt := win.Bounds().Max.ScaledXY(pixel.V(1.0/4, 3.0/4)).Add(pixel.V(-startTrigger.Bounds().Center().X-15, -startTrigger.Bounds().Center().Y-25))
	buildArrow(pt, leftStartArrow, true)
	leftStartArrowRect = pixel.R(pt.X-arrowSize, pt.Y-arrowSize, pt.X, pt.Y+arrowSize)

	pt = win.Bounds().Max.ScaledXY(pixel.V(1.0/4, 3.0/4)).Add(pixel.V(startTrigger.Bounds().Center().X+15, -startTrigger.Bounds().Center().Y-25))
	buildArrow(pt, rightStartArrow, false)
	rightStartArrowRect = pixel.R(pt.X, pt.Y-arrowSize, pt.X+arrowSize, pt.Y+arrowSize)

	leftEndArrow = imdraw.New(nil)
	rightEndArrow = imdraw.New(nil)
	leftScrambleArrow = imdraw.New(nil)
	rightScrambleArrow = imdraw.New(nil)

	pt = win.Bounds().Max.ScaledXY(pixel.V(3.0/4, 3.0/4)).Add(pixel.V(-endTrigger.Bounds().Center().X-15, -endTrigger.Bounds().Center().Y-25))
	buildArrow(pt, leftEndArrow, true)
	leftEndArrowRect = pixel.R(pt.X-arrowSize, pt.Y-arrowSize, pt.X, pt.Y+arrowSize)

	pt = win.Bounds().Max.ScaledXY(pixel.V(3.0/4, 3.0/4)).Add(pixel.V(endTrigger.Bounds().Center().X+15, -endTrigger.Bounds().Center().Y-25))
	buildArrow(pt, rightEndArrow, false)
	rightEndArrowRect = pixel.R(pt.X, pt.Y-arrowSize, pt.X+arrowSize, pt.Y+arrowSize)

	startTriggerIndex = util.IndexOfString(triggerOrder, tempConfig.TimerStartTrigger)
	endTriggerIndex = util.IndexOfString(triggerOrder, tempConfig.TimerEndTrigger)

	tickMarks = imdraw.New(nil)
	buildTickMarks(win.Bounds().Center().Add(pixel.V(0, -75)), tickMarks)

	pt = win.Bounds().Center().Add(pixel.V(-175, -75))
	buildArrow(pt, leftScrambleArrow, true)
	leftScrambleArrowRect = pixel.R(pt.X-arrowSize, pt.Y-arrowSize, pt.X, pt.Y+arrowSize)
	leftScramblePt = pt.Add(pixel.V(25, 0))

	pt = win.Bounds().Center().Add(pixel.V(175, -75))
	buildArrow(pt, rightScrambleArrow, false)
	rightScrambleArrowRect = pixel.R(pt.X, pt.Y-arrowSize, pt.X+arrowSize, pt.Y+arrowSize)

	lastScrambleAdjust = time.Now()
}

// OnShow has some last minute prep for showing a scene
func OnShow() {
	tempConfig = config.GlobalConfig()
	startTriggerIndex = util.IndexOfString(triggerOrder, tempConfig.TimerStartTrigger)
	endTriggerIndex = util.IndexOfString(triggerOrder, tempConfig.TimerEndTrigger)
	indicatorGrab = false
}

// Draw updates and renders the Timer scene
func Draw(canvas *pixelgl.Canvas, win util.LimitedWindow, dt *util.DeltaTimer) (change *scenes.SceneType) {
	canvas.Clear(colornames.Black)
	if !win.Pressed(pixelgl.MouseButtonLeft) {
		indicatorGrab = false
	}
	cB := canvas.Bounds()
	// StartTriggerLabel
	mat := pixel.IM.Moved(startTrigger.Bounds().Center().Scaled(-1)).Moved(pixel.V(cB.W()/4, cB.H()*3/4))
	startTrigger.Draw(canvas, mat)

	leftStartArrow.Draw(canvas)
	if win.JustPressed(pixelgl.MouseButtonLeft) && util.RectIsClicked(leftStartArrowRect, win.MousePosition()) {
		startTriggerIndex = (startTriggerIndex + len(triggerOrder) - 1) % len(triggerOrder)
	}
	rightStartArrow.Draw(canvas)
	if win.JustPressed(pixelgl.MouseButtonLeft) && util.RectIsClicked(rightStartArrowRect, win.MousePosition()) {
		startTriggerIndex = (startTriggerIndex + 1) % len(triggerOrder)
	}

	// StartTriggerValue
	mat = pixel.IM.Moved(labelOrder[startTriggerIndex].Bounds().Center().Scaled(-1)).Moved(pixel.V(cB.W()/4, cB.H()*3/4-36))
	labelOrder[startTriggerIndex].Draw(canvas, mat)

	// EndTriggerLabel
	mat = pixel.IM.Moved(endTrigger.Bounds().Center().Scaled(-1)).Moved(pixel.V(cB.W()*3/4, cB.H()*3/4))
	endTrigger.Draw(canvas, mat)

	leftEndArrow.Draw(canvas)
	if win.JustPressed(pixelgl.MouseButtonLeft) && util.RectIsClicked(leftEndArrowRect, win.MousePosition()) {
		endTriggerIndex = (endTriggerIndex + len(triggerOrder) - 1) % len(triggerOrder)
	}
	rightEndArrow.Draw(canvas)
	if win.JustPressed(pixelgl.MouseButtonLeft) && util.RectIsClicked(rightEndArrowRect, win.MousePosition()) {
		endTriggerIndex = (endTriggerIndex + 1) % len(triggerOrder)
	}

	// EndTriggerValue
	mat = pixel.IM.Moved(labelOrder[endTriggerIndex].Bounds().Center().Scaled(-1)).Moved(pixel.V(cB.W()*3/4, cB.H()*3/4-36))
	labelOrder[endTriggerIndex].Draw(canvas, mat)

	// ScrambleLengthLabel
	mat = pixel.IM.Moved(canvas.Bounds().Center()).Moved(scrambleLength.Bounds().Center().Scaled(-1))
	scrambleLength.Draw(canvas, mat)

	// Progress Bar (width 300)
	tickMarks.Draw(canvas)
	leftScrambleArrow.Draw(canvas)
	rightScrambleArrow.Draw(canvas)

	indicator := imdraw.New(nil)
	pt := leftScramblePt.Add(pixel.V(float64(tempConfig.ScrambleLength-10)*7.5, 0))
	if indicatorGrab && win.Pressed(pixelgl.MouseButtonLeft) {
		pt.X = math.Max(leftScramblePt.X, win.MousePosition().X)
		pt.X = math.Min(leftScramblePt.X+300, pt.X)
		offX := int((pt.X-leftScramblePt.X)/7.5 + 10)
		tempConfig.ScrambleLength = offX
	}
	indBox := buildIndicator(pt, indicator)
	indicator.Draw(canvas)

	if win.JustPressed(pixelgl.MouseButtonLeft) && util.IsClicked(pixel.IM.Moved(pt), indBox, win.MousePosition()) {
		indicatorGrab = true
	}

	scrambleNum.Clear()
	fmt.Fprint(scrambleNum, tempConfig.ScrambleLength)
	mat = pixel.IM.Moved(pt).Moved(pixel.V(0, -40)).Moved(scrambleNum.Bounds().Center().Scaled(-1))
	scrambleNum.Draw(canvas, mat)

	if !indicatorGrab && win.Pressed(pixelgl.MouseButtonLeft) && util.RectIsClicked(leftScrambleArrowRect, win.MousePosition()) && time.Since(lastScrambleAdjust).Seconds() > 0.15 {
		tempConfig.ScrambleLength = int(math.Max(10, float64(tempConfig.ScrambleLength)-1))
		lastScrambleAdjust = time.Now()
	}
	if !indicatorGrab && win.Pressed(pixelgl.MouseButtonLeft) && util.RectIsClicked(rightScrambleArrowRect, win.MousePosition()) && time.Since(lastScrambleAdjust).Seconds() > 0.15 {
		tempConfig.ScrambleLength = int(math.Min(50, float64(tempConfig.ScrambleLength)+1))
		lastScrambleAdjust = time.Now()
	}

	saveBtn.Draw(canvas, saveMatrix)
	saveBox.Draw(canvas)
	if win.JustPressed(pixelgl.MouseButtonLeft) && util.IsClicked(saveMatrix, saveRect, win.MousePosition()) {
		tempConfig.TimerStartTrigger = triggerOrder[startTriggerIndex]
		tempConfig.TimerEndTrigger = triggerOrder[endTriggerIndex]
		config.SaveConfig(tempConfig)
		change = new(scenes.SceneType)
		*change = scenes.TimerScene
	}

	cancelBtn.Draw(canvas, cancelMatrix)
	cancelBox.Draw(canvas)
	if win.JustPressed(pixelgl.MouseButtonLeft) && util.IsClicked(cancelMatrix, cancelRect, win.MousePosition()) {
		change = new(scenes.SceneType)
		*change = scenes.TimerScene
	}

	return change
}

func buildArrow(pt pixel.Vec, imd *imdraw.IMDraw, pointsLeft bool) {
	imd.Color = colornames.White

	imd.SetMatrix(pixel.IM.Moved(pt))
	imd.Push(pixel.V(0, arrowSize), pixel.V(0, -arrowSize))
	if pointsLeft {
		imd.Push(pixel.V(-arrowSize, 0))
	} else {
		imd.Push(pixel.V(arrowSize, 0))
	}
	imd.Polygon(0)
}

func buildTickMarks(pt pixel.Vec, imd *imdraw.IMDraw) {
	imd.Color = colornames.White

	imd.SetMatrix(pixel.IM.Moved(pt))
	for i := float64(-5); i <= 5; i++ {
		if int(i) == -5 || int(i) == 5 {
			imd.Push(pixel.V(30*i, -20), pixel.V(30*i, 20))
			imd.Line(3)
		} else if int(math.Abs(i))%2 == 1 {
			imd.Push(pixel.V(30*i, -15), pixel.V(30*i, 15))
			imd.Line(3)
		} else {
			imd.Push(pixel.V(30*i, -10), pixel.V(30*i, 10))
			imd.Line(1)
		}
	}
}

func buildIndicator(pt pixel.Vec, imd *imdraw.IMDraw) pixel.Rect {
	imd.Color = colornames.White

	imd.SetMatrix(pixel.IM.Moved(pt))

	imd.Push(pixel.V(-5, 25), pixel.V(5, -25))
	imd.Rectangle(0)
	return pixel.R(-5, -25, 5, 25)
}

func boxText(tx *text.Text, imd *imdraw.IMDraw) (box pixel.Rect) {
	bounds := tx.Bounds()

	imd.Color = colornames.White
	offset := pixel.V(8, 10)
	min := bounds.Min.Sub(offset)
	max := bounds.Max.Add(offset)
	imd.Push(min, max)
	imd.Rectangle(3)
	return pixel.R(min.X, min.Y, max.X, max.Y)
}
