package main

import (
	"fmt"
	"os"

	"github.com/coreyog/rubikstimer/config"
	"github.com/coreyog/rubikstimer/scenes"
	"github.com/coreyog/rubikstimer/scenes/settingsscene"
	"github.com/coreyog/rubikstimer/scenes/timerscene"
	"github.com/coreyog/rubikstimer/util"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/jessevdk/go-flags"
)

type Arguments struct {
	Undecorated bool `short:"u" long:"undecorated" description:"Run the program without a window border"`
}

var undecorated = false
var currentScene scenes.SceneType
var args *Arguments

func main() {
	config.LoadConfig()

	args = &Arguments{}

	_, err := flags.Parse(args)
	if err != nil && !flags.WroteHelp(err) {
		fmt.Println(err)
		os.Exit(1)
	}

	currentScene = scenes.TimerScene

	pixelgl.Run(run)
}

func run() {
	cfg := config.GlobalConfig()

	pixelCfg := pixelgl.WindowConfig{
		Title:       "Rubik's Timer",
		Bounds:      pixel.R(0, 0, float64(cfg.WindowWidth), float64(cfg.WindowHeight)),
		VSync:       true,
		Undecorated: args.Undecorated,
		Resizable:   true,
	}

	win, err := pixelgl.NewWindow(pixelCfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(true)

	canvas := pixelgl.NewCanvas(win.Bounds())
	canvas.SetSmooth(true)

	dt := util.NewDeltaTimer(0)

	var mousePos pixel.Vec

	mouseDown := false

	game := util.NewGameWindowWithCanvas(win, canvas)

	timerscene.Init(game)
	settingsscene.Init(game)

	for !win.Closed() {
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
			changeScene = timerscene.Draw(canvas, game, dt)
		case scenes.SettingsScene:
			changeScene = settingsscene.Draw(canvas, game, dt)
		}

		// game.DrawMouse() // useful

		game.Draw()
		if changeScene != nil {
			currentScene = *changeScene
			switch currentScene {
			case scenes.TimerScene:
				timerscene.OnShow()
			case scenes.SettingsScene:
				settingsscene.OnShow()
			}
		}

		win.Update()
	}
}
