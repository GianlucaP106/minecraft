package app

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Generates world terrain and content.
type WorldGenerator struct {
	// generates noise map for terrain generation
	noise *NoiseMapGenerator
}

func newWorldGenerator(seed int64) *WorldGenerator {
	w := &WorldGenerator{}
	w.noise = newNoiseMapGenerator()
	w.noise.Seed(seed)
	return w
}

func (w *WorldGenerator) Terrain(pos mgl32.Vec3) BlockTypes {
	pos2d := mgl32.Vec2{pos.X(), pos.Z()}
	tempurature := w.Tempurature(pos2d)

	// TODO: generalize to be able to fetch height easily
	biome := w.Biome(pos2d)
	flatHeights := w.FlatHeights(pos2d, 40)
	mountHeights := w.MountainHeights(pos2d, 170)
	out := newBlockTypes()

	// set terrain
	for x := 0; x < chunkWidth; x++ {
		for y := 0; y < chunkHeight; y++ {
			for z := 0; z < chunkWidth; z++ {
				curHeight := pos.Y() + float32(y)
				control := biome[x][z]
				targetHeight := exerp(control, flatHeights[x][z], mountHeights[x][z], 1.9)
				if curHeight > targetHeight {
					// if this block is higher than the the
					// height at this point in the terrain
					continue
				}

				if curHeight > 70 && tempurature[x][z] < 0.5 {
					out[x][y][z] = "dirt-snow"
					continue
				}

				if control > 0.5 {
					out[x][y][z] = "dirt-grass"
				} else {
					out[x][y][z] = "dirt-wet-grass"
				}
			}
		}
	}

	return out
}

// Generates tempuratures for a chunk at a given position.
func (w *WorldGenerator) Tempurature(pos mgl32.Vec2) [][]float32 {
	config2D := NoiseConfig2D{
		scale:     0.005,
		normalize: true,
		width:     float32(chunkWidth),
		height:    float32(chunkWidth),
		position:  pos,
		octaves:   1,
	}
	return w.noise.Generate2D(config2D)
}

// Generates tempuratures for a chunk at a given position.
func (w *WorldGenerator) Biome(pos mgl32.Vec2) [][]float32 {
	config2D := NoiseConfig2D{
		scale:     0.001,
		normalize: true,
		width:     float32(chunkWidth),
		height:    float32(chunkWidth),
		position:  pos,
		octaves:   1,
	}
	return w.noise.Generate2D(config2D)
}

// Generates mountain height ranges for a chunk at a given position.
func (w *WorldGenerator) MountainHeights(pos mgl32.Vec2, height float32) [][]float32 {
	config2D := NoiseConfig2D{
		scale:       0.015,
		normalize:   true,
		width:       float32(chunkWidth),
		height:      float32(chunkWidth),
		position:    pos,
		octaves:     4,
		persistence: 0.5,
		lacunarity:  2,
		f: func(noise float32, i, j int) float32 {
			return noise * height
		},
	}
	return w.noise.Generate2D(config2D)
}

// Generates flat height ranges for a chunk at a given position.
func (w *WorldGenerator) FlatHeights(pos mgl32.Vec2, height float32) [][]float32 {
	config2D := NoiseConfig2D{
		scale:     0.001,
		normalize: true,
		width:     float32(chunkWidth),
		height:    float32(chunkWidth),
		position:  pos,
		octaves:   1,
		f: func(noise float32, i, j int) float32 {
			return noise * height
		},
	}
	return w.noise.Generate2D(config2D)
}

func (w *WorldGenerator) TreeDistribution(pos mgl32.Vec2) [][]float32 {
	biome := w.Biome(pos)
	config2D := NoiseConfig2D{
		scale:     0.5,
		normalize: true,
		width:     float32(chunkWidth),
		height:    float32(chunkWidth),
		position:  pos,
		octaves:   1,
		f: func(noise float32, i, j int) float32 {
			control := biome[i][j]
			// amplify trees when biome is high
			return (1 - control) * noise
		},
	}
	return w.noise.Generate2D(config2D)
}

func (w *WorldGenerator) TreeFallout(width, height, depth float32) [][][]float32 {
	center := mgl32.Vec3{2.5, 0, 2.5}
	config := NoiseConfig3D{
		width:     width,
		depth:     depth,
		height:    height,
		normalize: true,
		octaves:   1,
		scale:     0.55,
		f: func(noise float32, x, y, z int) float32 {
			cur := mgl32.Vec3{float32(x), float32(y), float32(z)}
			dist := cur.Sub(center).Len()
			factor := 1 - dist/2.5
			if factor < 0 {
				factor = 0
			}
			return noise * factor
		},
	}

	return w.noise.Generate3D(config)
}
