package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	width  = 320
	height = 200

	vertexShaderSource = `
		#version 330 core

		layout(location = 0) in vec2 position;

		uniform mat4 model;

		void main() {
			gl_Position = model * vec4(position, 0.0, 1.0);
		}
	`

	fragmentShaderSource = `
		#version 330 core

		out vec4 FragColor;

		void main() {
			FragColor = vec4(gl_FragCoord.xy / 320.0, 0.5, 1.0);
		}
	`
)

type Particle struct {
	x, y, vx, vy float32
}

func init() {
	runtime.LockOSThread()
	rand.Seed(time.Now().UnixNano())
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatal(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Particle System in OpenGL 3.x API", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_OPENGL)
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()

	glContext, err := window.GLCreateContext()
	if err != nil {
		log.Fatal(err)
	}
	defer sdl.GLDeleteContext(glContext)

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Version: %v\nRenderer: %v\n", gl.GoStr(gl.GetString(gl.VERSION)), gl.GoStr((gl.GetString(gl.RENDERER))))
	fmt.Println("Press ESC to quit")

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// Create and compile shaders
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		log.Fatal(err)
	}
	defer gl.DeleteShader(vertexShader)

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		log.Fatal(err)
	}
	defer gl.DeleteShader(fragmentShader)

	// Create shader program
	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		logInfo := make([]byte, logLength)
		gl.GetProgramInfoLog(program, logLength, nil, &logInfo[0])

		log.Fatal("Error linking shader program:", string(logInfo))
	}
	gl.UseProgram(program)

	// Set up vertex data and buffers for particles
	particleCount := 1000
	particles := make([]Particle, particleCount)
	vertices := make([]float32, particleCount*2)

	for i := 0; i < particleCount; i++ {
		particles[i] = Particle{
			x: 0.0,
			y: 0.0,
			vx: randFloat(-0.01, 0.01),
			vy: randFloat(-0.01, 0.01),
		}
		vertices[i*2] = particles[i].x
		vertices[i*2+1] = particles[i].y
	}

	var vao, vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.DYNAMIC_DRAW)

	// position attribute
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BindVertexArray(0)

	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))

	// Main loop
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.KeyboardEvent:
				if event.Keysym.Sym == sdl.K_ESCAPE && event.State == sdl.PRESSED {
					return
				}
			}
		}

		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Update and render particles
		for i := 0; i < particleCount; i++ {
			particles[i].x += particles[i].vx
			particles[i].y += particles[i].vy

			// Reset particles when they go out of bounds
			if particles[i].x > 1.0 || particles[i].x < -1.0 || particles[i].y > 1.0 || particles[i].y < -1.0 {
				particles[i].x = 0.0
				particles[i].y = 0.0
			}

			vertices[i*2] = particles[i].x
			vertices[i*2+1] = particles[i].y
		}

		gl.BindVertexArray(vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices))
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)

		// Set the model matrix uniform (identity matrix for particles)
		gl.UniformMatrix4fv(modelUniform, 1, false, &identityMatrix[0])

		// Draw particles
		gl.DrawArrays(gl.POINTS, 0, int32(particleCount))

		gl.BindVertexArray(0)

		window.GLSwap()
		sdl.Delay(16) // Cap the frame rate to approximately 60 FPS
	}
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csource, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		logInfo := make([]byte, logLength)
		gl.GetShaderInfoLog(shader, logLength, nil, &logInfo[0])

		return 0, fmt.Errorf("compile shader error: %v", string(logInfo))
	}

	return shader, nil
}

func randFloat(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

var identityMatrix = [16]float32{
	1, 0, 0, 0,
	0, 1, 0, 0,
	0, 0, 1, 0,
	0, 0, 0, 1,
}
