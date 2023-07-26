package timer

import (
	"time"
)

type timer struct {
	fireTime time.Time
	interval time.Duration
	callback TimerCallbackFunc
	addseq   uint
}

func (t *timer) Cancel() {
	t.callback = nil
}

func (t *timer) IsActive() bool {
	return t.callback != nil
}

type _TimerHeap struct {
	timers []*timer
}

func (h *_TimerHeap) Len() int {
	return len(h.timers)
}

func (h *_TimerHeap) Less(i, j int) bool {
	//log.Println(h.timers[i].fireTime, h.timers[j].fireTime)
	t1, t2 := h.timers[i].fireTime, h.timers[j].fireTime
	if t1.Before(t2) {
		return true
	}

	if t1.After(t2) {
		return false
	}
	// t1 == t2, making sure Timer with same deadline is fired according to their add order
	return h.timers[i].addseq < h.timers[j].addseq
}

func (h *_TimerHeap) Swap(i, j int) {
	tmp := h.timers[i]
	h.timers[i] = h.timers[j]
	h.timers[j] = tmp
}

func (h *_TimerHeap) Push(x any) {
	h.timers = append(h.timers, x.(*timer))
}

func (h *_TimerHeap) Pop() (ret any) {
	l := len(h.timers)
	h.timers, ret = h.timers[:l-1], h.timers[l-1]
	return
}
