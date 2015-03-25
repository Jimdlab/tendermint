package common

import "time"
import "sync"

/*
RepeatTimer repeatedly sends a struct{}{} to .Ch after each "dur" period.
It's good for keeping connections alive.
*/
type RepeatTimer struct {
	Ch chan time.Time

	mtx    sync.Mutex
	name   string
	ticker *time.Ticker
	quit   chan struct{}
	dur    time.Duration
}

func NewRepeatTimer(name string, dur time.Duration) *RepeatTimer {
	var t = &RepeatTimer{
		Ch:     make(chan time.Time),
		ticker: time.NewTicker(dur),
		quit:   make(chan struct{}),
		name:   name,
		dur:    dur,
	}
	go t.fireRoutine(t.ticker)
	return t
}

func (t *RepeatTimer) fireRoutine(ticker *time.Ticker) {
	for {
		select {
		case t_ := <-ticker.C:
			t.Ch <- t_
		case <-t.quit:
			return
		}
	}
}

// Wait the duration again before firing.
func (t *RepeatTimer) Reset() {
	t.mtx.Lock() // Lock
	defer t.mtx.Unlock()

	if t.ticker != nil {
		t.ticker.Stop()
	}
	t.ticker = time.NewTicker(t.dur)
	go t.fireRoutine(t.ticker)
}

func (t *RepeatTimer) Stop() bool {
	t.mtx.Lock() // Lock
	defer t.mtx.Unlock()

	exists := t.ticker != nil
	if exists {
		t.ticker.Stop()
		t.ticker = nil
	}
	return exists
}
