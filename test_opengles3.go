/*
Compile: go build -tags=gles2 -x -v test_opengles3.1.go
*/
package main

import (
	"fmt"
	gl "github.com/leonkasovan/gl/v3.1/gles2"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"runtime"
)

const (
	vertexShaderSource = `
	#version 100

	attribute vec3 position;

	void main() {
		gl_Position = vec4(position, 1.0);
	}
	`

	fragmentShaderSource = `
	#version 100

	void main() {
		gl_FragColor = vec4(1.0, 0.5, 0.2, 1.0);
	}
	`
)

var (
	vertices = []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}
)

func main() {
	runtime.LockOSThread()

	// SDL2 initialization
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL: %s\n", err)
		os.Exit(1)
	}

	// Set OpenGL version
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_ES)
	// Create SDL window
	window, err := sdl.CreateWindow("Go SDL2 OpenGLES", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_OPENGL|sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		os.Exit(1)
	}
	defer window.Destroy()

	// Create SDL OpenGL context
	context, err := window.GLCreateContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create OpenGL context: %s\n", err)
		os.Exit(1)
	}
	defer sdl.GLDeleteContext(context)

	// Initialize GLES2
	if err := gl.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize OpenGL: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Version: %v\nRenderer: %v\n", gl.GoStr(gl.GetString(gl.VERSION)), gl.GoStr((gl.GetString(gl.RENDERER))))
	fmt.Println("Press ESC to quit")

	program, err := newProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		fmt.Println("Failed to create shader program:", err)
		return
	}
	gl.UseProgram(program)

	vao, vbo := makeVBO(vertices)

	// Main loop
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		// Clear the screen
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)/3))
		gl.BindVertexArray(0)

		// Swap the buffers
		window.GLSwap()

		sdl.Delay(16) // Cap frame rate
	}

	// Cleanup
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteProgram(program)
	sdl.Quit()
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := make([]byte, logLength)
		gl.GetProgramInfoLog(program, logLength, nil, &log[0])

		return 0, fmt.Errorf("linking program failed: %v", string(log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
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

		log := make([]byte, logLength)
		gl.GetShaderInfoLog(shader, logLength, nil, &log[0])

		return 0, fmt.Errorf("compiling shader failed: %v", string(log))
	}

	return shader, nil
}

func makeVBO(vertices []float32) (uint32, uint32) {
	var vbo, vao uint32

	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return vao, vbo
}