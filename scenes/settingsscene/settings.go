package settingsscene

import (
	"fmt"
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

var leftStartArrow *imdraw.IMDraw
var rightStartArrow *imdraw.IMDraw
var leftEndArrow *imdraw.IMDraw
var rightEndArrow *imdraw.IMDraw

var leftStartArrowRect pixel.Rect
var rightStartArrowRect pixel.Rect
var leftEndArrowRect pixel.Rect
var rightEndArrowRect pixel.Rect

var triggerOrder = []string{string(config.TriggerControls), string(config.TriggerSpacebar), string(config.TriggerAny)}
var labelOrder = []*text.Text{}
var startTriggerIndex int
var endTriggerIndex int

// Init creates the resources for the Timer scene
func Init(win util.LimitedWindow) {
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

	labelOrder = append(labelOrder, controlsText, spacebarText, anykeyText)

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

	pt = win.Bounds().Max.ScaledXY(pixel.V(3.0/4, 3.0/4)).Add(pixel.V(-endTrigger.Bounds().Center().X-15, -endTrigger.Bounds().Center().Y-25))
	buildArrow(pt, leftEndArrow, true)
	leftEndArrowRect = pixel.R(pt.X-arrowSize, pt.Y-arrowSize, pt.X, pt.Y+arrowSize)

	pt = win.Bounds().Max.ScaledXY(pixel.V(3.0/4, 3.0/4)).Add(pixel.V(endTrigger.Bounds().Center().X+15, -endTrigger.Bounds().Center().Y-25))
	buildArrow(pt, rightEndArrow, false)
	rightEndArrowRect = pixel.R(pt.X, pt.Y-arrowSize, pt.X+arrowSize, pt.Y+arrowSize)

	startTriggerIndex = util.IndexOfString(triggerOrder, config.GlobalConfig().TimerStartTrigger)
	endTriggerIndex = util.IndexOfString(triggerOrder, config.GlobalConfig().TimerEndTrigger)
}

// Draw updates and renders the Timer scene
func Draw(canvas *pixelgl.Canvas, win util.LimitedWindow, dt *util.DeltaTimer) (change *scenes.SceneType) {
	canvas.Clear(colornames.Black)

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
