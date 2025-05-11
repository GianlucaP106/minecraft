package game

import (
	"time"
)

// Provides a time delta on each tick.
// Implements fixed time step.
type Clock struct {
	// last ticked
	last time.Time

	// 1 / fps
	tickInterval time.Duration

	// fixed time step accumulator
	frameTimer time.Duration
}

func newClock() *Clock {
	c := &Clock{
		tickInterval: time.Millisecond * 16,
		frameTimer:   0,
	}
	return c
}

// Starts the clock.
func (c *Clock) Start() {
	c.last = time.Now()
}

// Ticks simulation forward.
func (c *Clock) Tick() {
	now := time.Now()
	delta := now.Sub(c.last)
	c.last = now
	c.frameTimer += delta
}

// Returns true if we should consume a simulation step.
func (c *Clock) ShouldSimulate() bool {
	return c.frameTimer >= c.tickInterval
}

// Consume a simulation step in fixed time step.
func (c *Clock) ConsumeStep() {
	c.frameTimer -= c.tickInterval
}

// Returns the delta value in seconds to use for simulation.
func (c *Clock) SimulationDelta() float64 {
	return c.tickInterval.Seconds()
}
