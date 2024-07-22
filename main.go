package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/oaorotti/poor-engine/camera"
	"github.com/oaorotti/poor-engine/shaders"
	"github.com/oaorotti/poor-engine/textures"
)

var (
	lastX, lastY float64 = 640, 360
	firstMouse   bool    = true
	deltaTime    float32
	lastFrame    float32
)

var myCamera = camera.NewCamera(mgl32.Vec3{0.0, 0.0, 3.0}, mgl32.Vec3{0.0, 1.0, 0.0}, -90.0, 0.0)

func mouseCallback(window *glfw.Window, xpos, ypos float64) {
	if firstMouse {
		lastX = xpos
		lastY = ypos
		firstMouse = false
	}

	xoffset := float32(xpos - lastX)
	yoffset := float32(lastY - ypos)
	lastX = xpos
	lastY = ypos

	myCamera.ProcessMouseMovement(xoffset, yoffset, true)
}

func scrollCallback(window *glfw.Window, xoffset, yoffset float64) {
	myCamera.ProcessMouseMovement(float32(xoffset), float32(yoffset), false)
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyW) == glfw.Press {
		myCamera.ProcessKeyboard("FORWARD", deltaTime)
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		myCamera.ProcessKeyboard("BACKWARD", deltaTime)
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		myCamera.ProcessKeyboard("LEFT", deltaTime)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		myCamera.ProcessKeyboard("RIGHT", deltaTime)
	}
}

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func ResizeWindow(window *glfw.Window, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func SetNewScroll(window *glfw.Window) {
	if window.GetInputMode(glfw.CursorMode) == glfw.CursorDisabled {
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		}
	} else if window.GetInputMode(glfw.CursorMode) == glfw.CursorNormal {
		if window.GetMouseButton(glfw.MouseButton1) == glfw.Press {
			window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
		}
	}
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(1280, 720, "Poor Engine", nil, nil)
	if err != nil {
		panic(err)
	}

	window.SetFramebufferSizeCallback(ResizeWindow)
	window.SetCursorPosCallback(mouseCallback)
	window.SetScrollCallback(scrollCallback)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}

	shader, err := shaders.NewProgramFromFiles("assets/vertex.glsl", "assets/fragment.glsl")
	if err != nil {
		log.Fatalln(err)
	}

	gl.UseProgram(shader)

	vertices := []float32{
		// Positions        // Texture Coords
		0.0, 0.5, 0.0, 0.5, 1.0,
		-0.5, -0.5, 0.0, 0.0, 0.0,
		0.5, -0.5, 0.0, 1.0, 0.0,
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	textureID, err := textures.NewTexture("assets/textures/manoel.jpeg")
	if err != nil {
		log.Fatalln(err)
	}

	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.Uniform1i(gl.GetUniformLocation(shader, gl.Str("ourTexture\x00")), 0)

	// Light and view positions
	lightPos := mgl32.Vec3{1.2, 1.0, 1.0}
	viewPos := mgl32.Vec3{0.0, 0.0, 3.0}
	lightColor := mgl32.Vec3{1.0, 1.0, 1.0}
	objectColor := mgl32.Vec3{1.0, 0.5, 0.31}

	gl.Uniform3fv(gl.GetUniformLocation(shader, gl.Str("lightPos\x00")), 1, &lightPos[0])
	gl.Uniform3fv(gl.GetUniformLocation(shader, gl.Str("viewPos\x00")), 1, &viewPos[0])
	gl.Uniform3fv(gl.GetUniformLocation(shader, gl.Str("lightColor\x00")), 1, &lightColor[0])
	gl.Uniform3fv(gl.GetUniformLocation(shader, gl.Str("objectColor\x00")), 1, &objectColor[0])

	modelLoc := gl.GetUniformLocation(shader, gl.Str("model\x00"))
	viewLoc := gl.GetUniformLocation(shader, gl.Str("view\x00"))
	projectionLoc := gl.GetUniformLocation(shader, gl.Str("projection\x00"))

	for !window.ShouldClose() {
		currentFrame := float32(glfw.GetTime())
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ClearColor(0.2, 0.2, 0.2, 1.0)

		processInput(window)
		SetNewScroll(window)

		gl.UseProgram(shader)

		model := mgl32.Ident4()
		view := myCamera.GetViewMatrix()
		projection := mgl32.Perspective(mgl32.DegToRad(myCamera.Zoom), float32(1280)/float32(720), 0.1, 100.0)

		gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
		gl.UniformMatrix4fv(projectionLoc, 1, false, &projection[0])

		gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
