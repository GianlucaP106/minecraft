package app

import (
	"github.com/go-gl/mathgl/mgl32"
)

type World struct {
	// player camera
	camera *Camera

	// shader program that draws the chunks
	chunkShader uint32
}

func newWorld(cam *Camera, chunkShader uint32) *World {
	w := &World{}
	w.camera = cam
	w.chunkShader = chunkShader
	return w
}

func (w *World) SpawnChunk(pos mgl32.Vec3) *Chunk {
	// init chunk, attribs and pointers
	chunk := newChunk(w.chunkShader, w.camera, pos)
	chunk.Init()
	chunk.Buffer()
	return chunk
}
