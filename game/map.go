package game

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

// Vector map providing lookup of objects by 3D coordinate.
type SpatialMap[T any] struct {
	m map[string]*T
}

func newVecMap[T any]() SpatialMap[T] {
	v := SpatialMap[T]{
		m: make(map[string]*T),
	}
	return v
}

// Returns the object stored at the coordinate.
func (v *SpatialMap[T]) Get(p mgl32.Vec3) *T {
	key := v.serialize(p)
	t, e := v.m[key]
	if !e {
		return nil
	}
	return t
}

// Sets the object at the coordinate.
func (v *SpatialMap[T]) Set(p mgl32.Vec3, t *T) {
	key := v.serialize(p)
	v.m[key] = t
}

// Deletes the object at the coordinate.
func (v *SpatialMap[T]) Delete(p mgl32.Vec3) {
	key := v.serialize(p)
	delete(v.m, key)
}

// Returns all the objects in a list.
func (v *SpatialMap[T]) All() []*T {
	out := make([]*T, 0)
	for _, v := range v.m {
		out = append(out, v)
	}
	return out
}

// Serializes the coordinate.
func (v *SpatialMap[T]) serialize(p mgl32.Vec3) string {
	return fmt.Sprintf("%f_%f_%f", p.X(), p.Y(), p.Z())
}
