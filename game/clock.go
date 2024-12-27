package game

import "github.com/go-gl/glfw/v3.3/glfw"

// Provides a time delta on each tick.
type Clock struct {
	last float64
}

func newClock() *Clock {
	c := &Clock{}
	return c
}

// Starts the clock.
func (c *Clock) Start() float64 {
	c.last = glfw.GetTime()
	return c.last
}

// Returns the delta since the last delta call.
func (c *Clock) Delta() float64 {
	now := glfw.GetTime()
	delta := now - c.last
	c.last = now
	return delta
}
