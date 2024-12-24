package app

import (
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
)

// Generates perlin noise used for terrain generation.
type NoiseMapGenerator struct {
	perm []int
	seed int64
}

// Configurers a 3D generation.
type NoiseConfig3D struct {
	f           func(noise float32, i, j, k int) float32
	octaves     int
	position    mgl32.Vec3
	scale       float32
	persistence float32
	lacunarity  float32
	width       float32
	height      float32
	depth       float32
	normalize   bool
}

// Configurers a 2D generation.
type NoiseConfig2D struct {
	f           func(noise float32, i, j int) float32
	octaves     int
	position    mgl32.Vec2
	scale       float32
	persistence float32
	lacunarity  float32
	width       float32
	height      float32
	normalize   bool
}

func newNoiseMapGenerator() *NoiseMapGenerator {
	return &NoiseMapGenerator{}
}

// Seeds the generator by creating a permutation table.
func (n *NoiseMapGenerator) Seed(seed int64) {
	n.seed = seed
	n.perm = n.generatePermutation(seed)
}

// Generates a 3D noise map with the given configuration.
func (n *NoiseMapGenerator) Generate3D(config NoiseConfig3D) [][][]float32 {
	o := make([][][]float32, 0)
	idx := 0
	for i := int(floor(config.position.X())); i < int(floor(config.width+config.position.X())); i++ {
		o = append(o, make([][]float32, 0))
		jdx := 0
		for j := int(floor(config.position.Y())); j < int(floor(config.height+config.position.Y())); j++ {
			o[idx] = append(o[idx], make([]float32, 0))
			kdx := 0
			for k := int(floor(config.position.Z())); k < int(floor(config.depth+config.position.Z())); k++ {
				noise := n.octaveNoise3D(
					float32(i),
					float32(j),
					float32(k),
					config.scale,
					config.persistence,
					config.lacunarity,
					config.octaves,
				)

				if config.normalize {
					noise += 1
					noise /= 2
				}
				if config.f != nil {
					noise = config.f(noise, idx, jdx, kdx)
				}
				o[idx][jdx] = append(o[idx][jdx], noise)
				kdx++
			}
			jdx++
		}
		idx++
	}
	return o
}

// Generates a 2D noise map with the given configuration.
func (n *NoiseMapGenerator) Generate2D(config NoiseConfig2D) [][]float32 {
	idx := 0
	jdx := 0
	o := make([][]float32, 0)
	for i := int(floor(config.position.X())); i < int(floor(config.width+config.position.X())); i++ {
		o = append(o, make([]float32, 0))
		jdx = 0
		for j := int(floor(config.position.Y())); j < int(floor(config.height+config.position.Y())); j++ {
			noise := n.octaveNoise2D(
				float32(j),
				float32(i),
				config.scale,
				config.persistence,
				config.lacunarity,
				config.octaves,
			)

			if config.normalize {
				noise += 1
				noise /= 2
			}
			if config.f != nil {
				noise = config.f(noise, idx, jdx)
			}
			o[idx] = append(o[idx], noise)
			jdx++
		}
		idx++
	}
	return o
}

// Wrapper over perlinNoise3D to compute noise from octaves, persistence and lacunarity.
func (n *NoiseMapGenerator) octaveNoise3D(x, y, z, scale, persistence, lacunarity float32, octaves int) float32 {
	total := float32(0)
	frequency := float32(1)
	amplitude := float32(1)
	maxValue := float32(0)

	for i := 0; i < octaves; i++ {
		total += n.perlinNoise3D(x*scale*frequency, y*scale*frequency, z*scale*frequency, n.perm) * amplitude

		maxValue += amplitude
		amplitude *= persistence
		frequency *= lacunarity
	}

	return total / maxValue
}

// Wrapper over perlinNoise2D to compute noise from octaves, persistence and lacunarity.
func (n *NoiseMapGenerator) octaveNoise2D(x, y, scale, persistence, lacunarity float32, octaves int) float32 {
	total := float32(0)
	frequency := float32(1)
	amplitude := float32(1)
	maxValue := float32(0)

	for i := 0; i < octaves; i++ {
		total += n.perlinNoise2D(x*scale*frequency, y*scale*frequency, n.perm) * amplitude
		maxValue += amplitude
		amplitude *= persistence
		frequency *= lacunarity
	}

	return total / maxValue
}

func (n *NoiseMapGenerator) perlinNoise3D(x, y, z float32, perm []int) float32 {
	// get grid point
	x0 := floor(x)
	y0 := floor(y)
	z0 := floor(z)

	// offset in grid cell
	relX := x - x0
	relY := y - y0
	relZ := z - z0

	// fade curves for each coordinate
	u := fade(relX)
	v := fade(relY)
	w := fade(relZ)

	// wrap index to allow looking up into permutation
	x0i := int(x0) & 255
	y0i := int(y0) & 255
	z0i := int(z0) & 255

	// hash the cube corners (8 corners)
	h000 := perm[perm[perm[x0i]+y0i]+z0i]
	h001 := perm[perm[perm[x0i]+y0i]+z0i+1]
	h010 := perm[perm[perm[x0i]+y0i+1]+z0i]
	h011 := perm[perm[perm[x0i]+y0i+1]+z0i+1]
	h100 := perm[perm[perm[x0i+1]+y0i]+z0i]
	h101 := perm[perm[perm[x0i+1]+y0i]+z0i+1]
	h110 := perm[perm[perm[x0i+1]+y0i+1]+z0i]
	h111 := perm[perm[perm[x0i+1]+y0i+1]+z0i+1]

	// get gradients for corners
	g000 := n.grad3D(h000, relX, relY, relZ)
	g001 := n.grad3D(h001, relX, relY, relZ-1)
	g010 := n.grad3D(h010, relX, relY-1, relZ)
	g011 := n.grad3D(h011, relX, relY-1, relZ-1)
	g100 := n.grad3D(h100, relX-1, relY, relZ)
	g101 := n.grad3D(h101, relX-1, relY, relZ-1)
	g110 := n.grad3D(h110, relX-1, relY-1, relZ)
	g111 := n.grad3D(h111, relX-1, relY-1, relZ-1)

	// lerp along x, then y, then z
	lx00 := lerp(u, g000, g100)
	lx01 := lerp(u, g001, g101)
	lx10 := lerp(u, g010, g110)
	lx11 := lerp(u, g011, g111)

	ly0 := lerp(v, lx00, lx10)
	ly1 := lerp(v, lx01, lx11)

	return lerp(w, ly0, ly1)
}

func (n *NoiseMapGenerator) perlinNoise2D(x, y float32, perm []int) float32 {
	// get grid corners
	x0 := floor(x)
	y0 := floor(y)
	x1 := x0 + 1
	y1 := y0 + 1

	// get offset
	relX := x - x0
	relY := y - y0

	// compute fade values
	u := fade(relX)
	v := fade(relY)

	// wrap index to allow looking up into permutation
	x0i := int(x0) & 255
	y0i := int(y0) & 255
	x1i := int(x1) & 255
	y1i := int(y1) & 255

	// hash grid corners
	h00 := perm[perm[x0i]+y0i]
	h10 := perm[perm[x1i]+y0i]
	h01 := perm[perm[x0i]+y1i]
	h11 := perm[perm[x1i]+y1i]

	// compute gradients
	g00 := n.grad2D(h00, relX, relY)
	g10 := n.grad2D(h10, relX-1, relY)
	g01 := n.grad2D(h01, relX, relY-1)
	g11 := n.grad2D(h11, relX-1, relY-1)

	// lerp
	lx0 := lerp(u, g00, g10)
	lx1 := lerp(u, g01, g11)
	return lerp(v, lx0, lx1)
}

// Generates permutation table from a given seed.
func (n *NoiseMapGenerator) generatePermutation(seed int64) []int {
	perm := make([]int, 256)
	for i := range perm {
		perm[i] = i
	}

	gen := rand.New(rand.NewSource(seed))
	gen.Shuffle(len(perm), func(i, j int) {
		perm[i], perm[j] = perm[j], perm[i]
	})
	// Duplicate for easier wrapping
	return append(perm, perm...)
}

// Returns a gradient based on the input hash in a pseudo-random but repeatable way.
func (n *NoiseMapGenerator) grad3D(hash int, x, y, z float32) float32 {
	switch hash & 0xF {
	case 0x0:
		return x + y
	case 0x1:
		return -x + y
	case 0x2:
		return x - y
	case 0x3:
		return -x - y
	case 0x4:
		return x + z
	case 0x5:
		return -x + z
	case 0x6:
		return x - z
	case 0x7:
		return -x - z
	case 0x8:
		return y + z
	case 0x9:
		return -y + z
	case 0xA:
		return y - z
	case 0xB:
		return -y - z
	case 0xC:
		return y + x
	case 0xD:
		return -y + z
	case 0xE:
		return y - x
	case 0xF:
		return -y - z
	default:
		return 0
	}
}

// Returns a gradient based on the input hash in a pseudo-random but repeatable way.
func (n *NoiseMapGenerator) grad2D(hash int, x, y float32) float32 {
	h := hash & 3
	u := x
	if h&1 == 0 {
		u = -x
	}
	v := y
	if h&2 == 0 {
		v = -y
	}
	return u + v
}
