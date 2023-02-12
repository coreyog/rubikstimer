package util

import (
	"time"
)

// DeltaTimer manages TotalTime and LastFrame time
type DeltaTimer struct {
	fixedRate uint
	tickCount uint
	startTime time.Time
	lastFrame time.Time

	totalTime      float64
	sinceLastFrame float64
}

// NewDeltaTimer starts a new DeltaTimer
func NewDeltaTimer(fixedRate uint) (dt *DeltaTimer) {
	dt = &DeltaTimer{
		fixedRate: fixedRate,
	}

	dt.start()

	return dt
}

func (dt *DeltaTimer) start() {
	dt.startTime = time.Now()
	dt.lastFrame = dt.startTime
}

// Tick is called at the first of every time to make timing measurements
func (dt *DeltaTimer) Tick() {
	dt.tickCount++
	if dt.fixedRate <= 0 {
		now := time.Now()

		dt.sinceLastFrame = now.Sub(dt.lastFrame).Seconds()
		dt.totalTime = now.Sub(dt.startTime).Seconds()
		dt.lastFrame = now
	} else {
		dt.sinceLastFrame = 1 / float64(dt.fixedRate)
		dt.totalTime = float64(dt.tickCount) / float64(dt.fixedRate)
	}
}

// TotalTime returns the time since the DeltaTimer was created
func (dt *DeltaTimer) TotalTime() (seconds float64) {
	return dt.totalTime
}

// SinceLastFrame returns the time since the last frame
func (dt *DeltaTimer) SinceLastFrame() (seconds float64) {
	return dt.sinceLastFrame
}

// IsFixed returns if ticks count time in realtime or at a fixed rate
func (dt *DeltaTimer) IsFixed() (fixed bool) {
	return dt.fixedRate != 0
}

// TickCount returns how many ticks have passed
func (dt *DeltaTimer) TickCount() (count uint) {
	return dt.tickCount
}
