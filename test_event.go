// author: Jacky Boen

package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

var winTitle string = "Go-SDL2 Events"
var winWidth, winHeight int32 = 320, 200
var joysticks [16]*sdl.Joystick

func run() int {
	var window *sdl.Window
	var event sdl.Event
	var running bool
	var err error

	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer window.Destroy()

	sdl.JoystickEventState(sdl.ENABLE)

	running = true
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseMotionEvent:
				fmt.Printf("[%d ms] MouseMotion type:%d id:%d x:%d y:%d xrel:%d yrel:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
			case *sdl.MouseButtonEvent:
				// t.State == sdl.PRESSED
				fmt.Printf("[%d ms] MouseButton type:%d id:%d x:%d y:%d button:%d state:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
			case *sdl.MouseWheelEvent:
				fmt.Printf("[%d ms] MouseWheel type:%d id:%d x:%d y:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y)
			case *sdl.KeyboardEvent:
				fmt.Printf("[%d ms] Keyboard type:%d sym:%d scancode:%d modifiers:%d state:%d repeat:%d\n",
					t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Scancode, t.Keysym.Mod, t.State, t.Repeat)
			case *sdl.JoyAxisEvent:
				fmt.Printf("[%d ms] JoyAxis type:%d which:%c axis:%d value:%d\n",
					t.Timestamp, t.Type, t.Which, t.Axis, t.Value)
			case *sdl.JoyBallEvent:
				fmt.Printf("[%d ms] JoyBall type:%d which:%d ball:%d xrel:%d yrel:%d\n",
					t.Timestamp, t.Type, t.Which, t.Ball, t.XRel, t.YRel)
			case *sdl.JoyButtonEvent:
				fmt.Printf("[%d ms] JoyButton type:%d which:%d button:%d state:%d\n",
					t.Timestamp, t.Type, t.Which, t.Button, t.State)
			case *sdl.JoyHatEvent:
				// t.Value = sdl.HAT_RIGHT sdl.HAT_UP sdl.HAT_DOWN sdl.HAT_LEFT sdl.HAT_CENTERED sdl.HAT_LEFTUP sdl.HAT_RIGHTUP sdl.HAT_RIGHTDOWN sdl.HAT_LEFTDOWN
				fmt.Printf("[%d ms] JoyHat type:%d which:%d hat:%d value:%d\n",
					t.Timestamp, t.Type, t.Which, t.Hat, t.Value)
			case *sdl.JoyDeviceAddedEvent:
				// Open joystick for use
				joysticks[int(t.Which)] = sdl.JoystickOpen(int(t.Which))
				if joysticks[int(t.Which)] != nil {
					fmt.Printf("Joystick %v connected (%v)\n", t.Which, joysticks[int(t.Which)].Name())
				}
			case *sdl.JoyDeviceRemovedEvent:
				if joystick := joysticks[int(t.Which)]; joystick != nil {
					joystick.Close()
				}
				fmt.Println("Joystick", t.Which, "disconnected")
			default:
				fmt.Printf("Some event\n")
			}
		}

		sdl.Delay(16)
	}

	return 0
}

func main() {
	os.Exit(run())
}
