package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	windowWidth  = 640
	windowHeight = 480
	fontPath     = "DejaVuSans.ttf"
	fontSize     = 12
	numRectangles = 4
)

var (
	randomColors []sdl.Color
)

func init() {
	rand.Seed(time.Now().UnixNano())

	randomColors = make([]sdl.Color, numRectangles)
	for i := range randomColors {
		randomColors[i] = sdl.Color{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 255,
		}
	}
}

func main() {
	runtime.LockOSThread()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("SDL initialization failed:", err)
		return
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		fmt.Println("TTF initialization failed:", err)
		return
	}
	defer ttf.Quit()

	window, err := sdl.CreateWindow("SDL2 Random Rectangles and Text", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println("Window creation failed:", err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("Renderer creation failed:", err)
		return
	}
	defer renderer.Destroy()

	font, err := ttf.OpenFont(fontPath, fontSize)
	if err != nil {
		fmt.Println("Font loading failed:", err)
		return
	}
	defer font.Close()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if t.Keysym.Sym == 27 {
					running = false
				}
			}
		}

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		for i, color := range randomColors {
			// Draw a random color rectangle
			rectWidth := rand.Intn(100) + 50
			rectHeight := rand.Intn(100) + 50
			rectX := rand.Intn(windowWidth - rectWidth)
			rectY := rand.Intn(windowHeight - rectHeight)

			rect := sdl.Rect{X: int32(rectX), Y: int32(rectY), W: int32(rectWidth), H: int32(rectHeight)}
			renderer.SetDrawColor(color.R, color.G, color.B, color.A)
			renderer.FillRect(&rect)

			// Draw the text with rectangle dimensions
			text := fmt.Sprintf("[%d] %dx%d", i+1, rectWidth, rectHeight)
			textSurface, err := font.RenderUTF8Solid(text, sdl.Color{R: 255, G: 255, B: 255, A: 255})
			if err != nil {
				fmt.Println("Text rendering failed:", err)
				return
			}
			defer textSurface.Free()

			textTexture, err := renderer.CreateTextureFromSurface(textSurface)
			if err != nil {
				fmt.Println("Texture creation failed:", err)
				return
			}
			defer textTexture.Destroy()

			renderer.Copy(textTexture, nil, &sdl.Rect{X: rect.X + rect.W/2 - 50, Y: rect.Y + 5, W: 100, H: 20})
		}

		renderer.Present()

		sdl.Delay(2000) // Delay for 2000 milliseconds (2 seconds) between frames
	}
}
