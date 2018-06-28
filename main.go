//go:generate go-bindata -o embedded/embedded.go -pkg embedded assets/...
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/coreyog/rubikstimer/config"
	"github.com/coreyog/rubikstimer/scenes"
	"github.com/coreyog/rubikstimer/scenes/settingsscene"
	"github.com/coreyog/rubikstimer/scenes/testscene"
	"github.com/coreyog/rubikstimer/scenes/timerscene"
	"github.com/coreyog/rubikstimer/util"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var undecorated = false
var currentScene scenes.SceneType

func main() {
	rand.Seed(time.Now().Unix())
	config.LoadConfig()

	help := []string{"-H", "--HELP", "/?", "HELP", "H"}
	undec := []string{"-U", "--U", "--UNDECORATED"}
	showHelp := false
	for _, arg := range os.Args {
		upperArg := strings.ToUpper(arg)
		for _, h := range help {
			if upperArg == h {
				showHelp = true
				break
			}
		}
		if !undecorated {
			for _, u := range undec {
				if upperArg == u {
					undecorated = true
					break
				}
			}
		}
	}

	if showHelp {
		printHelp()
		return
	}

	currentScene = scenes.TimerScene

	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:       "Rubik's Timer",
		Bounds:      pixel.R(0, 0, 1000, 400),
		VSync:       true,
		Undecorated: undecorated,
		Resizable:   false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	canvas := pixelgl.NewCanvas(win.Bounds())
	canvas.SetSmooth(true)

	dt := util.NewDeltaTimer(0)
	var mousePos pixel.Vec
	mouseDown := false
	refreshRate := float64(60)
	if win.Monitor() != nil {
		refreshRate = win.Monitor().RefreshRate()
	} else if pixelgl.PrimaryMonitor() != nil {
		refreshRate = pixelgl.PrimaryMonitor().RefreshRate()
	}
	fps := time.NewTicker(time.Second / time.Duration(refreshRate))

	timerscene.Init(win)
	testscene.Init(win)
	settingsscene.Init(win)

	for !win.Closed() {
		<-fps.C
		dt.Tick()

		if undecorated {
			if mouseDown {
				diff := win.MousePosition().Sub(mousePos)
				if diff.X != 0 && diff.Y != 0 {
					winPos := win.GetPos()
					winPos.X += diff.X
					winPos.Y -= diff.Y
					win.SetPos(winPos)
				}
			}

			if win.JustReleased(pixelgl.MouseButton1) {
				mouseDown = false
			}

			if win.JustPressed(pixelgl.MouseButton1) {
				mousePos = win.MousePosition()
				mouseDown = true
			}
		}

		if win.Pressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}

		var changeScene *scenes.SceneType
		switch currentScene {
		case scenes.TimerScene:
			changeScene = timerscene.Draw(canvas, win, dt)
			break
		case scenes.SettingsScene:
			changeScene = settingsscene.Draw(canvas, win, dt)
		case scenes.TestScene:
			changeScene = testscene.Draw(canvas, win, dt)
			break
		}

		canvas.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		if changeScene != nil {
			currentScene = *changeScene
			switch currentScene {
			case scenes.TimerScene:
				timerscene.OnShow()
				break
			case scenes.SettingsScene:
				settingsscene.OnShow()
				break
			case scenes.TestScene:
				testscene.OnShow()
				break
			}
		}

		win.Update()
	}
}

func printHelp() {
	fmt.Println("RubiksTimer [-u|--undecorated]")
	fmt.Println()
	fmt.Println("Adding an undecorated flag will remove the border from the window.")
	fmt.Println("Use Escape to close the program.")
	fmt.Println()
	fmt.Println("Hold both control keys on your keyboard to arm the timer.")
	fmt.Println("The timer will start when you release either control key.")
	fmt.Println("Press both control keys at the same time again to stop the timer.")
	fmt.Println("Pressing R will restart the timer and wait for both controls to be pressed again.")
	fmt.Println("F12 will flip between a Black and Magenta background (for use as a chroma key).")
}
