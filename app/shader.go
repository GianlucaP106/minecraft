package app

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// ShaderManager manages references to shader programs.
type ShaderManager struct {
	shaders  map[string]uint32
	rootPath string
}

func newShaderManager(root string) *ShaderManager {
	s := &ShaderManager{}
	s.rootPath = root
	s.shaders = make(map[string]uint32)
	s.init()
	return s
}

// Initializes the shaders found in the rootPath.
func (s *ShaderManager) init() {
	dirEntries := func() []fs.FileInfo {
		dir, err := os.Open(s.rootPath)
		if err != nil {
			log.Panicln(err)
		}
		defer dir.Close()

		dirEntries, err := dir.Readdir(-1)
		if err != nil {
			log.Panicln(err)
		}
		return dirEntries
	}()

	for _, d := range dirEntries {
		if !d.IsDir() {
			continue
		}

		name := d.Name()
		vshader := filepath.Join(s.rootPath, name, "vert.glsl")
		fshader := filepath.Join(s.rootPath, name, "frag.glsl")
		vb, err := os.ReadFile(vshader)
		if err != nil {
			panic(err)
		}

		fb, err := os.ReadFile(fshader)
		if err != nil {
			panic(err)
		}

		vsrc := string(vb) + "\x00"
		fsrc := string(fb) + "\x00"
		program := s.createProgram(vsrc, fsrc)
		s.shaders[name] = program
	}
}

// Returns a stored reference to a program.
func (s *ShaderManager) Program(name string) uint32 {
	e, w := s.shaders[name]
	if !w {
		log.Panic("invalid shader: ", name)
	}
	return e
}

// Creates shader program from sources.
func (s *ShaderManager) createProgram(vertexShaderSource, fragmentShaderSource string) uint32 {
	vertexShader, err := s.compile(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := s.compile(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
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

		panic(fmt.Errorf("failed to link program: %v", log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program
}

// Compiles a shader program.
func (s *ShaderManager) compile(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
