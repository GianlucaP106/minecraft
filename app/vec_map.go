package app

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

type VecMap[T any] struct {
	m map[string]*T
}

func newVecMap[T any]() VecMap[T] {
	v := VecMap[T]{
		m: make(map[string]*T),
	}
	return v
}

func (v *VecMap[T]) Get(p mgl32.Vec3) *T {
	key := v.serialize(p)
	t, e := v.m[key]
	if !e {
		return nil
	}
	return t
}

func (v *VecMap[T]) Set(p mgl32.Vec3, t *T) {
	key := v.serialize(p)
	v.m[key] = t
}

func (v *VecMap[T]) All() []*T {
	out := make([]*T, 0)
	for _, v := range v.m {
		out = append(out, v)
	}
	return out
}

func (v *VecMap[T]) serialize(p mgl32.Vec3) string {
	return fmt.Sprintf("%f_%f_%f", p.X(), p.Y(), p.Z())
}
