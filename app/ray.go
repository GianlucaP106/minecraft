package app

import "github.com/go-gl/mathgl/mgl32"

type Ray struct {
	origin    mgl32.Vec3
	direction mgl32.Vec3
}

func (r Ray) IsLookingAt(b Box) (bool, string) {
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
		return false, ""
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
		return false, ""
	}

	// distanceTo := b.Distance(c.eye)
	// raySize := r.direction.Len()

	// TODO: change ray
	hitPos := r.origin.Add(r.direction.Mul(tmin))
	// fmt.Println(hitPos)

	if hitPos.X() == bmin.X() {
		return true, "left"
	}

	if hitPos.X() == bmax.X() {
		return true, "right"
	}

	if hitPos.Y() == bmin.Y() {
		return true, "bottom"
	}

	if hitPos.Y() == bmax.Y() {
		return true, "top"
	}

	if hitPos.Z() == bmin.Z() {
		return true, "back"
	}

	if hitPos.Z() == bmax.Z() {
		return true, "front"
	}

	return true, "none"

	// smallThresh := float32(0.000001)
	// bigThresh := float32(0.199990)
	// if mgl32.Abs(hitPos.X()-bmin.X()) < smallThresh {
	// 	return true, "left"
	// }
	//
	// if mgl32.Abs(hitPos.X()-bmin.X()) > bigThresh {
	// 	return true, "right"
	// }
	//
	// if mgl32.Abs(hitPos.Y()-bmin.Y()) < smallThresh {
	// 	return true, "bottom"
	// }
	//
	// if mgl32.Abs(hitPos.Y()-bmin.Y()) > bigThresh {
	// 	return true, "top"
	// }
	//
	// if mgl32.Abs(hitPos.Z()-bmin.Z()) < smallThresh {
	// 	return true, "front"
	// }
	//
	// if mgl32.Abs(hitPos.Z()-bmin.Z()) > bigThresh {
	// 	return true, "back"
	// }
}
