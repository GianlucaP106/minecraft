package game

import "github.com/go-gl/mathgl/mgl32"

// Computed vertex for a block (pre-world).
type TexturedVertex struct {
	Vertex
	tex mgl32.Vec2
}

type Vertex struct {
	pos  mgl32.Vec3
	norm mgl32.Vec3
	face Direction
}

func computeVertices() [chunkWidth][chunkHeight][chunkWidth][6][6]Vertex {
	// vertices for one block (6 quads)
	blockVertices := [6][6]mgl32.Vec3{}
	for face := 0; face < 6; face++ {
		face := Direction(face)
		quad := quad(face)
		blockVertices[face] = quad
	}

	// vertices for chunk
	chunkVertices := [chunkWidth][chunkHeight][chunkWidth][6][6]Vertex{}
	for i := 0; i < chunkWidth; i++ {
		for j := 0; j < chunkHeight; j++ {
			for k := 0; k < chunkWidth; k++ {
				translate := translateBlock(i, j, k)
				finalBlockVertices := [6][6]Vertex{}
				for f := 0; f < 6; f++ {
					face := Direction(f)
					norm := directions[face]
					for v := 0; v < 6; v++ {
						vertex := blockVertices[f][v]
						finalBlockVertices[f][v] = Vertex{
							pos:  translate.Mul4x1(vertex.Vec4(1)).Vec3(),
							face: face,
							norm: norm,
						}
					}
				}
				chunkVertices[i][j][k] = finalBlockVertices
			}
		}
	}
	return chunkVertices
}

func translateBlock(i, j, k int) mgl32.Mat4 {
	// half because block is size 2
	half := float32(blockSize / 2.0)
	return mgl32.Translate3D(
		float32(i)+half,
		float32(j)+half,
		float32(k)+half,
	)
}
