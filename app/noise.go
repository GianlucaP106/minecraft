package app

import (
	"math/rand"
)

type NoiseMapGenerator struct{}

func newNoiseMapGenerator() *NoiseMapGenerator {
	n := &NoiseMapGenerator{}
	return n
}

func (n *NoiseMapGenerator) Generate(width, height, depth int, scale float32, seed int64) [][]float32 {
	perm := n.generatePermutation(seed)
	o := make([][]float32, 0)
	for y := 0; y < height; y++ {
		o = append(o, make([]float32, 0))
		for x := 0; x < width; x++ {
			noise := n.perlinNoise(float32(x)*scale, float32(y)*scale, perm)
			noise += 1
			noise /= 2
			noise *= float32(depth)
			o[y] = append(o[y], noise)
		}
	}
	return o
}

// Fade function for smoothing
func (n *NoiseMapGenerator) fade(t float32) float32 {
	return t * t * t * (t*(t*6-15) + 10)
}

// Linear interpolation
func (n *NoiseMapGenerator) lerp(t, a, b float32) float32 {
	return a + t*(b-a)
}

// Gradient function
func (n *NoiseMapGenerator) grad(hash int, x, y float32) float32 {
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

// Perlin noise function
func (n *NoiseMapGenerator) perlinNoise(x, y float32, perm []int) float32 {
	x0 := floor(x)
	y0 := floor(y)
	x1 := x0 + 1
	y1 := y0 + 1

	// Compute relative coordinates
	relX := x - x0
	relY := y - y0

	// Wrap coordinates to permutation table
	x0i := int(x0) & 255
	y0i := int(y0) & 255
	x1i := int(x1) & 255
	y1i := int(y1) & 255

	// Calculate hash values
	h00 := perm[perm[x0i]+y0i]
	h10 := perm[perm[x1i]+y0i]
	h01 := perm[perm[x0i]+y1i]
	h11 := perm[perm[x1i]+y1i]

	// Compute gradients
	g00 := n.grad(h00, relX, relY)
	g10 := n.grad(h10, relX-1, relY)
	g01 := n.grad(h01, relX, relY-1)
	g11 := n.grad(h11, relX-1, relY-1)

	// Compute fade values
	u := n.fade(relX)
	v := n.fade(relY)

	// Interpolate
	lx0 := n.lerp(u, g00, g10)
	lx1 := n.lerp(u, g01, g11)
	return n.lerp(v, lx0, lx1)
}

// Generate a permutation table
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
