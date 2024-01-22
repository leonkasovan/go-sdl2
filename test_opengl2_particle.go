package main

import (
	"fmt"
	"math/rand"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 640
	windowHeight = 480
	numParticles = 1000
)

var (
	particles []Particle
	window    *sdl.Window
)

// Particle represents a single particle in the fire effect
type Particle struct {
	x, y     float32
	speed    float32
	life     int
	maxLife  int
	r, g, b   float32
}

func main() {
	// Initialize SDL and OpenGL
	runtime.LockOSThread()
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("SDL initialization failed:", err)
		return
	}

	// Set OpenGL version and profile
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 2)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)

	// Create window
	window, err := sdl.CreateWindow("Particle Fire", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		windowWidth, windowHeight, sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Println("Window creation failed:", err)
		return
	}
	defer window.Destroy()

	// Create OpenGL context
	glContext, err := window.GLCreateContext()
	if err != nil {
		fmt.Println("OpenGL context creation failed:", err)
		return
	}
	defer sdl.GLDeleteContext(glContext)

	// Initialize OpenGL
	if err := gl.Init(); err != nil {
		fmt.Println("OpenGL initialization failed:", err)
		return
	}

	gl.ClearColor(0, 0, 0, 1)
	gl.Viewport(0, 0, windowWidth, windowHeight)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, windowWidth, 0, windowHeight, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	// Initialize particles
	particles = make([]Particle, numParticles)
	for i := range particles {
		particles[i] = NewParticle()
	}

	// Main loop
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.KeyboardEvent:
				if e.Keysym.Sym == sdl.K_ESCAPE && e.State == sdl.PRESSED {
					return
				}
			}
		}

		// Update particles
		for i := range particles {
			particles[i].Update()
		}

		// Render particles
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Begin(gl.POINTS)
		for _, p := range particles {
			gl.Color3f(p.r, p.g, p.b)
			gl.Vertex2f(p.x, p.y)
		}
		gl.End()

		// Swap buffers
		window.GLSwap()

		// Delay to control frame rate
		sdl.Delay(16)
	}
}

// NewParticle initializes a new particle with random properties
func NewParticle() Particle {
	return Particle{
		x:       float32(rand.Intn(windowWidth)),
		y:       float32(rand.Intn(windowHeight)),
		speed:   float32(rand.Intn(5) + 1),
		life:    rand.Intn(100) + 100,
		maxLife: rand.Intn(100) + 100,
		r:       rand.Float32(),
		g:       rand.Float32(),
		b:       rand.Float32(),
	}
}

// Update updates the particle's position and life
func (p *Particle) Update() {
	if p.life > 0 {
		p.y += p.speed
		p.life--
	} else {
		*p = NewParticle()
	}
}
