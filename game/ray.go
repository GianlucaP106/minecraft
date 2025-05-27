package game

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Ray is a vector with an origin and projected length.
// Can be used as a line of sight.
type Ray struct {
	origin    mgl32.Vec3
	direction mgl32.Vec3
	length    float32
}

// The result of a ray march.
type March struct {
	hit      bool
	face     Direction
	blockPos mgl32.Vec3
	hitPos   mgl32.Vec3
	box      Box
}

// March marches in the direction of the ray, detection the first the block in sight,
// where the callback is used to determine if a block is present.
func (r Ray) March(find func(p mgl32.Vec3) *Box) March {
	// helper to find the smallest `t` such that `s + (ds * t)` is an integer
	// i.e finds the next block point
	intbound := func(s, ds mgl32.Vec3) mgl32.Vec3 {
		c := func(s, ds float32) float32 {
			if ds > 0 {
				return (float32(math.Ceil(float64(s))) - s) / float32(math.Abs(float64(ds)))
			} else if ds < 0 {
				return (s - float32(math.Floor(float64(s)))) / float32(math.Abs(float64(ds)))
			}
			return float32(math.Inf(1))
		}

		return mgl32.Vec3{
			c(s.X(), ds.X()),
			c(s.Y(), ds.Y()),
			c(s.Z(), ds.Z()),
		}
	}

	p := mgl32.Vec3{
		floor(r.origin.X()),
		floor(r.origin.Y()),
		floor(r.origin.Z()),
	}
	step := mgl32.Vec3{
		sign(r.direction.X()),
		sign(r.direction.Y()),
		sign(r.direction.Z()),
	}

	tmax := intbound(r.origin, r.direction)
	tdelta := mgl32.Vec3{
		step.X() / r.direction.X(),
		step.Y() / r.direction.Y(),
		step.Z() / r.direction.Z(),
	}
	radius := r.length / r.direction.Len()

	var face Direction
	var hit float32
	for {
		b := find(p)
		if b != nil {
			hitPos := r.origin.Add(r.direction.Mul(hit))
			return March{
				hit:      true,
				blockPos: p,
				hitPos:   hitPos,
				face:     face,
				box:      *b,
			}
		}

		if tmax.X() < tmax.Y() {
			if tmax.X() < tmax.Z() {
				hit = tmax.X()
				if tmax.X() > radius {
					break
				}

				p[0] += step.X()
				tmax[0] += tdelta.X()

				face = newDirection(mgl32.Vec3{-step.X(), 0, 0})
			} else {
				hit = tmax.Z()
				if tmax.Z() > radius {
					break
				}

				p[2] += step.Z()
				tmax[2] += tdelta.Z()
				face = newDirection(mgl32.Vec3{0, 0, -step.Z()})
			}
		} else {
			if tmax.Y() < tmax.Z() {
				hit = tmax.Y()
				if tmax.Y() > radius {
					break
				}

				p[1] += step.Y()
				tmax[1] += tdelta.Y()
				face = newDirection(mgl32.Vec3{0, -step.Y(), 0})
			} else {
				hit = tmax.Z()
				if tmax.Z() > radius {
					break
				}

				p[2] += step.Z()
				tmax[2] += tdelta.Z()
				face = newDirection(mgl32.Vec3{0, 0, -step.Z()})
			}
		}
	}
	return March{
		hit: false,
	}
}
