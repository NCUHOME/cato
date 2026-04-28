package httpc

import (
	"sync/atomic"
	"time"

	"github.com/rfyiamcool/go-timewheel"
)

type timeoutWheel struct {
	wheel *timewheel.TimeWheel
}

type wheelTask struct {
	task    *timewheel.Task
	wheel   *timewheel.TimeWheel
	stopped atomic.Bool
}

func newTimeoutWheel(tick time.Duration, slots int) *timeoutWheel {
	if tick <= 0 {
		tick = timeoutWheelTick
	}
	if slots <= 0 {
		slots = timeoutWheelSlots
	}

	wheel, err := timewheel.NewTimeWheel(tick, slots)
	if err != nil {
		panic(err)
	}
	wheel.Start()

	return &timeoutWheel{wheel: wheel}
}

func (w *timeoutWheel) Schedule(after time.Duration, fn func()) *wheelTask {
	return &wheelTask{
		task:  w.wheel.Add(after, fn),
		wheel: w.wheel,
	}
}

func (t *wheelTask) Stop() bool {
	if t == nil || t.task == nil {
		return false
	}
	if !t.stopped.CompareAndSwap(false, true) {
		return false
	}
	_ = t.wheel.Remove(t.task)
	return true
}

func (t *wheelTask) Release() {}
