package app

import "github.com/go-gl/glfw/v3.3/glfw"

// Provides a timedelta for game simulation.
type Clock struct {
	last float64
}

func newClock() *Clock {
	c := &Clock{}
	return c
}

func (c *Clock) Start() float64 {
	c.last = glfw.GetTime()
	return c.last
}

func (c *Clock) Delta() float64 {
	now := glfw.GetTime()
	delta := now - c.last
	c.last = now
	return delta
}
