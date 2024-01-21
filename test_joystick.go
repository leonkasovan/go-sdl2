// Test input joystick SDL2
// author: Dhani Novan
// 21 Januari 2024 Jakarta Cempaka Putih

package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/gfx"
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
	// var points []sdl.Point
	// var rect sdl.Rect
	// var rects []sdl.Rect

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

		// renderer.SetDrawColor(0, 0, 255, 255)
		// renderer.DrawLine(0, 0, 200, 200)

		// points = []sdl.Point{{50, 50}, {100, 300}, {200, 50}, {50, 50}}
		// renderer.SetDrawColor(255, 255, 0, 255)
		// renderer.DrawLines(points)

		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.DrawRect(&sdl.Rect{0, 0, 640, 480})

		// rects = []sdl.Rect{{120, 230, 160, 240}, {80, 270, 200, 310}}
		// renderer.SetDrawColor(200, 200, 200, 255)
		// renderer.DrawRects(rects)

		// rect = sdl.Rect{250, 250, 200, 200}
		// renderer.SetDrawColor(0, 255, 0, 255)
		// renderer.FillRect(&rect)

		// Draw Direction Pad
		renderer.SetDrawColor(200, 200, 200, 255)
		renderer.DrawRects([]sdl.Rect{{120, 230, 40, 120}, {80, 270, 120, 40}})
		renderer.SetDrawColor(0, 255, 0, 255)
		if hatValue == 1 {
			renderer.FillRect(&sdl.Rect{121, 231, 39, 39})
		}
		if hatValue == 2 {
			renderer.FillRect(&sdl.Rect{160, 271, 39, 39})
		}
		if hatValue == 4 {
			renderer.FillRect(&sdl.Rect{121, 310, 39, 39})
		}
		if hatValue == 8 {
			renderer.FillRect(&sdl.Rect{81, 271, 39, 39})
		}

		// Draw Select + Start Button
		renderer.SetDrawColor(200, 200, 200, 255)
		renderer.DrawRects([]sdl.Rect{{245, 330, 60, 20}, {320, 330, 60, 20}})
		if buttonState[6] == 1 {
			renderer.SetDrawColor(0, 255, 0, 255)
			renderer.FillRect(&sdl.Rect{246, 331, 58, 18})
			gfx.StringRGBA(renderer, 251, 336, "SELECT", 0, 0, 0, 255)
		}else{
			gfx.StringRGBA(renderer, 251, 336, "SELECT", 255, 255, 255, 255)
		}

		if buttonState[7] == 1 {
			renderer.SetDrawColor(0, 255, 0, 255)
			renderer.FillRect(&sdl.Rect{321, 331, 58, 18})
			gfx.StringRGBA(renderer, 329, 336, "START", 0, 0, 0, 255)
		}else{
			gfx.StringRGBA(renderer, 329, 336, "START", 255, 255, 255, 255)
		}

		// Draw X Buttons
		gfx.CircleRGBA(renderer, 450, 290, 20, 200, 200, 200, 255)
		if buttonState[2] == 1 {
			gfx.FilledCircleRGBA(renderer, 450, 290, 18, 0, 255, 0, 255)
			gfx.StringRGBA(renderer, 447, 288, "X", 0, 0, 0, 255)
		}else{
			gfx.StringRGBA(renderer, 447, 288, "X", 255, 255, 255, 255)
		}
		
		// Draw B Buttons
		gfx.CircleRGBA(renderer, 540, 290, 20, 200, 200, 200, 255)
		if buttonState[1] == 1 {
			gfx.FilledCircleRGBA(renderer, 540, 290, 18, 0, 255, 0, 255)
			gfx.StringRGBA(renderer, 537, 288, "B", 0, 0, 0, 255)
		}else{
			gfx.StringRGBA(renderer, 537, 288, "B", 255, 255, 255, 255)
		}
		
		// Draw Y Buttons
		gfx.CircleRGBA(renderer, 495, 250, 20, 200, 200, 200, 255)
		if buttonState[3] == 1 {
			gfx.FilledCircleRGBA(renderer, 495, 250, 18, 0, 255, 0, 255)
			gfx.StringRGBA(renderer, 492, 248, "Y", 0, 0, 0, 255)
		}else{
			gfx.StringRGBA(renderer, 492, 248, "Y", 255, 255, 255, 255)
		}
		
		// Draw A Buttons
		gfx.CircleRGBA(renderer, 495, 330, 20, 200, 200, 200, 255)
		if joysticks[0].Button(0) == 1 {
		// if buttonState[0] == 1 {
			gfx.FilledCircleRGBA(renderer, 495, 330, 18, 0, 255, 0, 255)
			gfx.StringRGBA(renderer, 492, 328, "A", 255, 0, 255, 255)
		}else{
			gfx.StringRGBA(renderer, 492, 328, "A", 255, 255, 255, 255)
		}
		
		gfx.StringRGBA(renderer, 100, 10, "Test Input Joystick in SDL2 (EXIT: SELECT+START)", 255, 255, 255, 255)
		if msgJoystickInfo[0] != "" {
			gfx.StringRGBA(renderer, 50, 30, msgJoystickInfo[0], 0, 255, 0, 255)
			gfx.StringRGBA(renderer, 50, 30 + 16*1, msgJoystickInfo[1], 0, 255, 0, 255)
			gfx.StringRGBA(renderer, 50, 30 + 16*2, msgJoystickInfo[2], 0, 255, 0, 255)
			gfx.StringRGBA(renderer, 50, 30 + 16*3, msgJoystickInfo[3], 0, 255, 0, 255)
			gfx.StringRGBA(renderer, 50, 30 + 16*4, msgJoystickInfo[4], 0, 255, 0, 255)
		}
		gfx.StringRGBA(renderer, 50, 400, msgKeyboardEvent, 0, 255, 0, 255)

		if msgJoystickEvent[0] != "" {
			gfx.StringRGBA(renderer, 50, 416, msgJoystickEvent[0], 0, 255, 0, 255)
		}
		if msgJoystickEvent[1] != "" {
			gfx.StringRGBA(renderer, 50, 416 + 16*1, msgJoystickEvent[1], 0, 255, 0, 255)
		}
		if msgJoystickEvent[2] != "" {
			gfx.StringRGBA(renderer, 50, 416 + 16*2, msgJoystickEvent[2], 0, 255, 0, 255)
		}
		if msgJoystickEvent[3] != "" {
			gfx.StringRGBA(renderer, 50, 416 + 16*3, msgJoystickEvent[3], 0, 255, 0, 255)
		}
		renderer.Present()
		sdl.Delay(16)
	}

	return 0
}

func main() {
	os.Exit(run())
}
