package app

import "github.com/go-gl/mathgl/mgl32"

type Ray struct {
	origin    mgl32.Vec3
	direction mgl32.Vec3
}

func (r Ray) IsLookingAt(b Box) (bool, BoxFace, mgl32.Vec3) {
	bmin := b.min
	bmax := b.max

	tMin := (bmin.X() - r.origin.X()) / r.direction.X()
	tMax := (bmax.X() - r.origin.X()) / r.direction.X()

	tmin := min(tMin, tMax)
	tmax := max(tMin, tMax)

	if tMin > tMax {
		tMin, tMax = tMax, tMin
	}

	tyMin := (bmin.Y() - r.origin.Y()) / r.direction.Y()
	tyMax := (bmax.Y() - r.origin.Y()) / r.direction.Y()

	tmin = max(tmin, min(tyMin, tyMax))
	tmax = min(tmax, max(tyMin, tyMax))

	if tyMin > tyMax {
		tyMin, tyMax = tyMax, tyMin
	}

	if tMin > tyMax || tyMin > tMax {
		return false, none, mgl32.Vec3{}
	}

	if tyMin > tMin {
		tMin = tyMin
	}
	if tyMax < tMax {
		tMax = tyMax
	}

	tzMin := (bmin.Z() - r.origin.Z()) / r.direction.Z()
	tzMax := (bmax.Z() - r.origin.Z()) / r.direction.Z()

	tmin = max(tmin, min(tzMin, tzMax))
	tmax = min(tmax, max(tzMin, tzMax))

	if tzMin > tzMax {
		tzMin, tzMax = tzMax, tzMin
	}

	if tMin > tzMax || tzMin > tMax {
		return false, none, mgl32.Vec3{}
	}

	// TODO: change ray
	hitPos := r.origin.Add(r.direction.Mul(tmin))

	var face BoxFace
	switch {
	case hitPos.X() == bmin.X():
		face = left
	case hitPos.X() == bmax.X():
		face = right
	case hitPos.Y() == bmin.Y():
		face = bottom
	case hitPos.Y() == bmax.Y():
		face = top
	case hitPos.Z() == bmin.Z():
		face = back
	case hitPos.Z() == bmax.Z():
		face = front
	default:
		face = none
	}

	return true, face, hitPos
}
