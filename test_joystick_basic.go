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
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return -1
	}
	defer sdl.Quit()

	running := true
	sdl.JoystickEventState(sdl.ENABLE)

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				fmt.Printf("Keyboard [%c] type:%d sym:%d scancode:%v modi:%d state:%d\n",
					t.Keysym.Sym, t.Type, t.Keysym.Sym, t.Keysym.Scancode, t.Keysym.Mod, t.State)
			case *sdl.JoyAxisEvent:
				fmt.Printf("JoyAxis type:%d which:%d axis:%d value:%d\n",
					t.Type, t.Which, t.Axis, t.Value)
			case *sdl.JoyBallEvent:
				fmt.Printf("JoyBall type:%d which:%d ball:%d xrel:%d yrel:%d\n",
					t.Type, t.Which, t.Ball, t.XRel, t.YRel)
			case *sdl.JoyButtonEvent:
				fmt.Printf("JoyButton type:%d which:%d button:%d state:%d\n",
					t.Type, t.Which, t.Button, t.State)
				buttonState[t.Button] = t.State
			case *sdl.JoyHatEvent:
				fmt.Printf("JoyHat type:%d which:%d hat:%d value:%d\n",
					t.Type, t.Which, t.Hat, t.Value)
				hatValue = t.Value
			case *sdl.JoyDeviceAddedEvent:
				// Open joystick for use
				joysticks[int(t.Which)] = sdl.JoystickOpen(int(t.Which))
				if joysticks[int(t.Which)] != nil {
					fmt.Printf("Joystick id=%v connected (%v)\n", t.Which, joysticks[int(t.Which)].Name())
					fmt.Printf("Joystick Name: %s\n", joysticks[int(t.Which)].Name())
					fmt.Printf("  - Number of Axes: %d\n", joysticks[int(t.Which)].NumAxes())
					fmt.Printf("  - Number of Buttons: %d\n", joysticks[int(t.Which)].NumButtons())
					fmt.Printf("  - Number of Balls: %d\n", joysticks[int(t.Which)].NumBalls())
					fmt.Printf("  - Number of Hats: %d\n", joysticks[int(t.Which)].NumHats())
				}
			}
		}

		if (buttonState[6] & buttonState[7]) == 1 {
			running = false
		}

		sdl.Delay(16)
	}

	return 0
}

func main() {
	os.Exit(run())
}
