package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	width       = 800
	height      = 600
	imageFile   = "image.png"
	vertexShader = `
		#version 140
		in vec2 position;
		out vec2 TexCoords;
		void main()
		{
			gl_Position = vec4(position, 0.0, 1.0);
			TexCoords = (position + vec2(1.0, 1.0)) / 2.0;
		}
	` + "\x00"
	fragmentShader = `
		#version 140
		in vec2 TexCoords;
		out vec4 FragColor;
		uniform sampler2D image;
		void main()
		{
			vec4 texColor = texture(image, TexCoords);
			float scanline = 0.007; // Adjust scanline intensity
			float scanlinePos = mod(TexCoords.y, scanline * 2.0);
			if (scanlinePos > scanline) {
				texColor.rgb *= 0.5; // Adjust darkness of scanline
			}
			FragColor = texColor;
		}
	` + "\x00"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatal(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("SDL2 OpenGL Filter", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_OPENGL)
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()

	context, err := window.GLCreateContext()
	if err != nil {
		log.Fatal(err)
	}
	defer sdl.GLDeleteContext(context)

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Version: %v\nRenderer: %v\n", gl.GoStr(gl.GetString(gl.VERSION)), gl.GoStr((gl.GetString(gl.RENDERER))))
	fmt.Println("Press ESC to quit")

	// Load image
	textureID, err := loadImage(imageFile)
	if err != nil {
		log.Fatal(err)
	}

	// Compile and link shaders
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		log.Fatal(err)
	}
	gl.UseProgram(program)

	// Set up vertex data and attribute pointers
	vertices := []float32{
		// Positions
		-1.0, -1.0,
		1.0, -1.0,
		1.0, 1.0,
		-1.0, 1.0,
	}
	var vao, vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	for !shouldQuit(window) {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(program)
		gl.BindVertexArray(vao)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, textureID)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("image\x00")), 0)

		gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)

		window.GLSwap()
		sdl.Delay(16) // Cap the frame rate to approximately 60 FPS
	}

	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteTextures(1, &textureID)
}

func loadImage(filename string) (uint32, error) {
	surface, err := img.Load(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to load image: %v", err)
	}
	defer surface.Free()

	textureID := createTexture(surface)

	return textureID, nil
}

func createTexture(surface *sdl.Surface) uint32 {
	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	// Set the texture wrapping parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	// Set the texture filtering parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// Upload the image data
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(surface.W), int32(surface.H), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(surface.Pixels()))

	return textureID
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(vertexShader)

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(fragmentShader)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		logInfo := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(logInfo))

		return 0, fmt.Errorf("link program error: %v", logInfo)
	}

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

		logInfo := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(logInfo))

		return 0, fmt.Errorf("compile shader error: %v", logInfo)
	}

	return shader, nil
}

func shouldQuit(window *sdl.Window) bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event := event.(type) {
		case *sdl.QuitEvent:
			return true
		case *sdl.KeyboardEvent:
			if event.Keysym.Sym == sdl.K_ESCAPE && event.State == sdl.PRESSED {
				return true
			}
		}
	}
	return false
}
