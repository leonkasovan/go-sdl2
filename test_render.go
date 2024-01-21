// Test input joystick SDL2
// author: Dhani Novan
// 21 Januari 2024 Jakarta Cempaka Putih

package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var winTitle string = "Go-SDL2"
var winWidth, winHeight int32 = 640, 480
var msgKeyboardEvent string = ""
var msgJoystickEvent [4]string = [4]string{"", "", "", ""}
var joysticks [8]*sdl.Joystick
var msgJoystickInfo [5]string = [5]string{"", "", "", "", ""}
var buttonState [16]uint8
var hatValue uint8

func run() int {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var points []sdl.Point
	var rect sdl.Rect
	var rects []sdl.Rect

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return -1
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}
	defer renderer.Destroy()

	running := true
	sdl.JoystickEventState(sdl.ENABLE)

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				msgKeyboardEvent = fmt.Sprintf("Keyboard [%c] type:%d sym:%d scancode:%v modi:%d state:%d",
				t.Keysym.Sym, t.Type, t.Keysym.Sym, t.Keysym.Scancode, t.Keysym.Mod, t.State)
				if t.Keysym.Sym == 27 {
					running = false
				}
			case *sdl.JoyAxisEvent:
				msgJoystickEvent[2] = fmt.Sprintf("JoyAxis type:%d which:%d axis:%d value:%d",
					 t.Type, t.Which, t.Axis, t.Value)
			case *sdl.JoyBallEvent:
				msgJoystickEvent[3] = fmt.Sprintf("JoyBall type:%d which:%d ball:%d xrel:%d yrel:%d",
					 t.Type, t.Which, t.Ball, t.XRel, t.YRel)
			case *sdl.JoyButtonEvent:
				msgJoystickEvent[0] = fmt.Sprintf("JoyButton type:%d which:%d button:%d state:%d",
					 t.Type, t.Which, t.Button, t.State)
				buttonState[t.Button] = t.State
			case *sdl.JoyHatEvent:
				msgJoystickEvent[1] = fmt.Sprintf("JoyHat type:%d which:%d hat:%d value:%d",
					 t.Type, t.Which, t.Hat, t.Value)
				hatValue = t.Value
			case *sdl.JoyDeviceAddedEvent:
					// Open joystick for use
					joysticks[int(t.Which)] = sdl.JoystickOpen(int(t.Which))
					if joysticks[int(t.Which)] != nil {
						msgJoystickEvent[0] = fmt.Sprintf("Joystick id=%v connected (%v)", t.Which, joysticks[int(t.Which)].Name())
						msgJoystickInfo[0] = fmt.Sprintf("Joystick Name: %s", joysticks[int(t.Which)].Name())
						msgJoystickInfo[1] = fmt.Sprintf("  - Number of Axes: %d", joysticks[int(t.Which)].NumAxes())
						msgJoystickInfo[2] = fmt.Sprintf("  - Number of Buttons: %d", joysticks[int(t.Which)].NumButtons())
						msgJoystickInfo[3] = fmt.Sprintf("  - Number of Balls: %d", joysticks[int(t.Which)].NumBalls())
						msgJoystickInfo[4] = fmt.Sprintf("  - Number of Hats: %d", joysticks[int(t.Which)].NumHats())
					}
			}
		}

		if (buttonState[6] & buttonState[7]) == 1 {
			running = false
		} 

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		renderer.SetDrawColor(0, 0, 255, 255)
		renderer.DrawLine(0, 0, 200, 200)

		points = []sdl.Point{{50, 50}, {100, 300}, {200, 50}, {50, 50}}
		renderer.SetDrawColor(255, 255, 0, 255)
		renderer.DrawLines(points)
		renderer.SetDrawColor(255, 255, 0, 255)
		renderer.DrawLines([]sdl.Point{{10, 10}, {80, 30}, {60, 50}, {10, 10}})

		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.DrawRect(&sdl.Rect{0, 0, 640, 480})
		
		rect = sdl.Rect{250, 250, 200, 200}
		renderer.SetDrawColor(0, 255, 0, 255)
		renderer.FillRect(&rect)

		rects = []sdl.Rect{{120, 230, 160, 240}, {80, 270, 200, 310}}
		renderer.SetDrawColor(200, 200, 200, 255)
		renderer.DrawRects(rects)

		renderer.SetDrawColor(200, 200, 200, 255)
		renderer.DrawRects([]sdl.Rect{{120, 230, 40, 120}, {80, 270, 120, 40}})
		renderer.SetDrawColor(0, 255, 0, 255)

		renderer.SetDrawColor(200, 200, 200, 255)
		renderer.DrawRects([]sdl.Rect{{245, 330, 60, 20}, {320, 330, 60, 20}})
		renderer.SetDrawColor(0, 255, 0, 255)
		renderer.FillRect(&sdl.Rect{246, 331, 58, 18})

		renderer.Present()
		sdl.Delay(16)
	}

	return 0
}

func main() {
	os.Exit(run())
}
