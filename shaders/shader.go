package shaders

import (
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"os"
	"strings"
)

func loadShaderSource(filePath string) (string, error) {
	source, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read shader file %s: %v", filePath, err)
	}
	return string(source) + "\x00", nil // Add null terminator
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	sources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, sources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile shader: %v", log)
	}

	return shader, nil
}

func NewProgramFromFiles(vertexShaderPath, fragmentShaderPath string) (uint32, error) {
	vertexShaderSource, err := loadShaderSource(vertexShaderPath)
	if err != nil {
		return 0, err
	}

	fragmentShaderSource, err := loadShaderSource(fragmentShaderPath)
	if err != nil {
		return 0, err
	}

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

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}
