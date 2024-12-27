package game

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

// Generates the terrain for a chunk.
// Returns the block types that can be used to initialize the chunk.
func (w *WorldGenerator) Terrain(pos mgl32.Vec3) BlockTypes {
	pos2d := mgl32.Vec2{pos.X(), pos.Z()}
	biome := w.Biome(pos2d)
	heights := w.Heights(pos2d)
	out := newBlockTypes()
	caves := w.Caves(pos)

	// set terrain
	for x := 0; x < chunkWidth; x++ {
		for y := chunkHeight - 1; y >= 0; y-- {
			for z := 0; z < chunkWidth; z++ {
				curHeight := pos.Y() + float32(y)
				if curHeight > heights[x][z] {
					// if this block is higher than the the
					// height at this point in the terrain
					continue
				}

				if caves[x][y][z] > 0.725 {
					continue
				}

				if curHeight < 0.25 {
					out[x][y][z] = "stone"
					continue
				}

				if curHeight > 70 && biome >= 0.6 {
					out[x][y][z] = "dirt-snow"
					continue
				}

				// check if block is a stone
				if y < chunkHeight-4 {
					isStone := true
					for i := y + 1; i < y+4; i++ {
						if out[x][i][z] == "" {
							isStone = false
							break
						}
					}
					if isStone {
						out[x][y][z] = "stone"
						continue
					}
				}

				// define terrain based on biome
				switch {
				case biome <= 0.4:
					out[x][y][z] = "sand"
				case biome > 0.4 && biome < 0.7:
					if y == chunkHeight || out[x][y+1][z] == "" {
						out[x][y][z] = "dirt-grass"
					} else {
						out[x][y][z] = "dirt"
					}
				case biome >= 0.7:
					out[x][y][z] = "dirt-wet-grass"
				}
			}
		}
	}

	return out
}

func (w *WorldGenerator) Biome(pos mgl32.Vec2) float32 {
	biome := w.noise.OctaveNoise2D(pos.X(), pos.Y(), 0.0005, 0.0, 0, 1, true)
	return normsigmoid(biome)
}

func (w *WorldGenerator) Heights(pos mgl32.Vec2) [][]float32 {
	biome := w.Biome(pos)
	flatHeights := w.FlatHeights(pos, 170)
	mountHeights := w.MountainHeights(pos, 170)

	// set the new heights in place in mountHeights
	for i := 0; i < len(mountHeights); i++ {
		for j := 0; j < len(mountHeights[0]); j++ {
			targetHeight := exerp(normsigmoid(biome), flatHeights[i][j], mountHeights[i][j], 1.0)
			mountHeights[i][j] = targetHeight
		}
	}

	return mountHeights
}

func (w *WorldGenerator) MountainHeights(pos mgl32.Vec2, height float32) [][]float32 {
	config2D := NoiseConfig2D{
		scale:       0.01,
		normalize:   true,
		width:       float32(chunkWidth),
		height:      float32(chunkWidth),
		position:    pos,
		octaves:     8,
		persistence: 1,
		lacunarity:  1,
		f: func(noise float32, i, j int) float32 {
			return noise * height
		},
	}
	return w.noise.Generate2D(config2D)
}

func (w *WorldGenerator) FlatHeights(pos mgl32.Vec2, height float32) [][]float32 {
	config2D := NoiseConfig2D{
		scale:       0.001,
		normalize:   true,
		width:       float32(chunkWidth),
		height:      float32(chunkWidth),
		position:    pos,
		octaves:     4,
		persistence: 0.7,
		lacunarity:  2,
		f: func(noise float32, i, j int) float32 {
			return noise * height
		},
	}
	return w.noise.Generate2D(config2D)
}

func (w *WorldGenerator) TreeDistribution(pos mgl32.Vec2) [][]float32 {
	// biome := w.Biome(pos)
	config2D := NoiseConfig2D{
		scale:     0.5,
		normalize: true,
		width:     float32(chunkWidth),
		height:    float32(chunkWidth),
		position:  pos,
		octaves:   1,
		f: func(noise float32, i, j int) float32 {
			return noise
			// control := biome[i][j]
			// // amplify trees when biome is high
			// return (1 - control) * noise
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

func (w *WorldGenerator) Caves(pos mgl32.Vec3) [][][]float32 {
	config := NoiseConfig3D{
		scale:       0.1,
		normalize:   true,
		width:       float32(chunkWidth),
		height:      float32(chunkHeight),
		depth:       chunkWidth,
		position:    pos,
		octaves:     5,
		persistence: 0.7,
		lacunarity:  1,
	}

	return w.noise.Generate3D(config)
}
