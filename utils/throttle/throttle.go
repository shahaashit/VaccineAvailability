package throttle

import (
	"sync"
)

type Throttle struct {
	once sync.Once
	wg   sync.WaitGroup
	ch   chan struct{}
}

func (t *Throttle) Do() {
	t.ch <- struct{}{}
	t.wg.Add(1)
}

func (t *Throttle) Done() {
	select {
	case <-t.ch:
	default:
		panic("Throttle Do Done mismatch")
	}
	t.wg.Done()
}

func (t *Throttle) Finish() {
	t.once.Do(func() {
		t.wg.Wait()
		close(t.ch)
	})
}

func NewThrottle(max int) *Throttle {
	return &Throttle{
		ch: make(chan struct{}, max),
	}
}
