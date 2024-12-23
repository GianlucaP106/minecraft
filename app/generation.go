package app

import (
	"log"

	"github.com/go-gl/mathgl/mgl32"
)

// Generates world terrain and content.
type WorldGenerator struct {
	world *World

	// generates noise map for terrain generation
	noise *NoiseMapGenerator
}

func newWorldGenerator(world *World) *WorldGenerator {
	w := &WorldGenerator{}
	w.world = world
	w.noise = newNoiseMapGenerator()
	return w
}

// Generates terrain.
// TODO: clean up
func (w *WorldGenerator) Generate() {
	log.Println("Generating terrain...")
	w.noise.Seed(919)
	size := 500
	height := 150
	biomes := w.noise.Generate2D(size, size, 0.01, true, nil)
	w.noise.SetOctaves(4, 0.5, 2)
	missingChunks := w.noise.Generate3D(size, size, height, 0.25, true, nil)

	mountainHeights := w.noise.Generate2D(size, size, 0.015, true, func(noise float32, i, j int) float32 {
		return noise * float32(height)
	})

	flatHeights := w.noise.Generate2D(size, size, 0.001, true, func(noise float32, i, j int) float32 {
		return noise * float32(height)
	})
	w.noise.SetOctaves(1, 1, 1)

	trees := w.noise.Generate2D(size, size, 0.5, true, nil)
	temperature := w.noise.Generate2D(size, size, 0.005, true, nil)

	log.Println("Generating mountains...")
	for x, heights := range mountainHeights {
		for z, height := range heights {
			// biome seperation
			control := biomes[x][z]
			height = height*control + (1-control)*flatHeights[x][z]
			for y := 0; y < int(height); y++ {
				b := w.world.Block(mgl32.Vec3{float32(x), float32(y), float32(z)})
				b.active = true

				temp := temperature[x][z]
				if y > 70 && temp > 0.5 {
					b.blockType = "dirt-snow"
				} else if temp < 0.4 {
					b.blockType = "dirt-wet-grass"
				} else {
					b.blockType = "dirt-grass"
				}
			}
		}
	}

	// remove some spots
	for i, layer := range missingChunks {
		for j, row := range layer {
			for k, val := range row {
				if val < 0.3 {
					b := w.world.Block(mgl32.Vec3{float32(i), float32(j), float32(k)})
					b.active = false
				}
			}
		}
	}

	log.Println("Generating trees...")
	poses := []mgl32.Vec3{}
	for x, treeProbs := range trees {
		for z, prob := range treeProbs {
			if prob > 0.85 {
				groundBlock := w.world.Ground(float32(x), float32(z))
				if groundBlock != nil {
					base := groundBlock.WorldPos().Add(mgl32.Vec3{0, 1, 0})
					poses = append(poses, base)

				}
			}
		}
	}
	for _, p := range poses {
		w.Tree(p, 5.0, 7.0, 4.0, 1.75, 1.5)
	}
}

// Generates tree at the provided location with provided params.
func (w *WorldGenerator) Tree(base mgl32.Vec3, radius, trunkHeight, leafHeight, trunkFallout, leafFallout float32) {
	center := mgl32.Vec2{radius, radius}
	makeRadialFallout := func(maxDist float32) func(noise float32, x, y, z int) float32 {
		return func(noise float32, x, y, z int) float32 {
			cur := mgl32.Vec2{float32(x), float32(z)}
			dist := cur.Sub(center).Len()
			factor := 1 - dist/(maxDist)
			if factor < 0 {
				factor = 0
			}
			return noise * factor
		}
	}
	corner := base.Sub(mgl32.Vec3{radius, 0, radius})
	trunkNoise := w.noise.Generate3D(int(radius)*2, int(radius)*2, int(trunkHeight), 0.15, true, makeRadialFallout(radius*trunkFallout))
	leafNoise := w.noise.Generate3D(int(radius)*2, int(radius)*2, int(leafHeight), 0.25, true, makeRadialFallout(radius*leafFallout))

	for x, layers := range trunkNoise {
		for y, noises := range layers {
			for z, noise := range noises {
				if noise > 0.5 {
					pos := corner.Add(mgl32.Vec3{float32(x), float32(y), float32(z)})
					b := w.world.Block(pos)
					b.active = true
					b.blockType = "wood"
				}
			}
		}
	}

	corner = corner.Add(mgl32.Vec3{0, trunkHeight - 2, 0})
	for x, layers := range leafNoise {
		for y, noises := range layers {
			for z, noise := range noises {
				if noise > 0.3 {
					pos := corner.Add(mgl32.Vec3{float32(x), float32(y), float32(z)})
					b := w.world.Block(pos)
					b.active = true
					if noise < 0.35 {
						b.blockType = "leaves-flower"
					} else {
						b.blockType = "leaves"
					}

				}
			}
		}
	}
}
