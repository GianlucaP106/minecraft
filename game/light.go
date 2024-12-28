package game

import (
	"log"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Maintains the light position and level.
type Light struct {
	time      *time.Ticker
	ascending bool
	pos       mgl32.Vec3
	level     float32
}

func newLight() *Light {
	l := &Light{}
	l.pos = mgl32.Vec3{0, 200, 0}
	l.level = 0.4
	return l
}

// Starts a timer to change the the light level at each interval.
func (l *Light) StartDay(interval time.Duration) {
	l.time = time.NewTicker(interval)
}

// Polls the ticker. Should be called at each iteration of game loop.
func (l *Light) HandleChange() {
	if l.time == nil {
		return
	}

	select {
	case <-l.time.C:
		var newLvl float32
		if l.ascending {
			newLvl = l.level + 0.1
		} else {
			newLvl = l.level - 0.1
		}
		if newLvl < 0 {
			newLvl = 0.1
			l.ascending = true
		} else if newLvl > 1.0 {
			newLvl = 1.0
			l.ascending = false
		}
		l.SetLevel(newLvl)
	default:
	}
}

// Sets the light level. (0.0-1.0)
func (l *Light) SetLevel(lvl float32) {
	if lvl < 0.0 || lvl > 1.0 {
		log.Panic("invalid light level ", lvl)
	}
	l.level = lvl
	gl.ClearColor(1.0*lvl, 1.0*lvl, 1.0*lvl, 1.0*lvl)
}
