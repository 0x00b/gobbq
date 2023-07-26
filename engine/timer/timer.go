package timer

import (
	"container/heap"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/0x00b/gobbq/tool/secure"
)

const (
	MIN_TIMER_INTERVAL = 1 * time.Millisecond
)

// CallbackFunc Type of callback function
type CallbackFunc func()

// TimerCallbackFunc Type of callback function, return false represent stop this timer
type TimerCallbackFunc func() bool

type Timer struct {
	timerHeap     _TimerHeap
	timerHeapLock sync.Mutex
	nextAddSeq    uint
}

func (t *Timer) Init() {
	t.nextAddSeq = 1
	heap.Init(&t.timerHeap)
}

// Add a callback which will be called after specified duration
func (t *Timer) AddCallback(d time.Duration, callback CallbackFunc) *timer {
	timer := &timer{
		fireTime: time.Now().Add(d),
		interval: d,
		callback: func() bool {
			callback()
			return false
		},
	}

	t.timerHeapLock.Lock()
	defer t.timerHeapLock.Unlock()

	timer.addseq = t.nextAddSeq // set addseq when locked
	t.nextAddSeq += 1

	heap.Push(&t.timerHeap, timer)
	return timer
}

// Add a timer which calls callback periodly
func (t *Timer) AddTimer(d time.Duration, callback TimerCallbackFunc) *timer {
	if d < MIN_TIMER_INTERVAL {
		d = MIN_TIMER_INTERVAL
	}

	timer := &timer{
		fireTime: time.Now().Add(d),
		interval: d,
		callback: callback,
	}

	t.timerHeapLock.Lock()
	defer t.timerHeapLock.Unlock()

	timer.addseq = t.nextAddSeq // set addseq when locked
	t.nextAddSeq += 1

	heap.Push(&t.timerHeap, timer)
	return timer
}

// Tick once for timers
func (t *Timer) Tick() {
	now := time.Now()

	t.timerHeapLock.Lock()
	defer t.timerHeapLock.Unlock()

	for {
		if t.timerHeap.Len() <= 0 {
			break
		}

		nextFireTime := t.timerHeap.timers[0].fireTime
		//xlog.Tracef(">>> nextFireTime %s, now is %s\n", nextFireTime, now)
		if nextFireTime.After(now) {
			break
		}

		timer := heap.Pop(&t.timerHeap).(*timer)

		callback := timer.callback
		if callback == nil {
			continue
		}

		stop := true
		func() {
			// unlock the lock to run callback, because callback may add more callbacks / timers
			t.timerHeapLock.Unlock()
			defer t.timerHeapLock.Lock()
			stop = !t.runCallback(callback)
		}()

		if !stop {
			// add Timer back to heap
			timer.fireTime = timer.fireTime.Add(timer.interval)
			if !timer.fireTime.After(now) { // might happen when interval is very small
				timer.fireTime = now.Add(timer.interval)
			}
			timer.addseq = t.nextAddSeq
			t.nextAddSeq += 1
			heap.Push(&t.timerHeap, timer)
		}
	}
}

// Start the self-ticking routine, which ticks per tickInterval
func (t *Timer) StartTicks(tickInterval time.Duration) {
	secure.GO(func() {
		t.selfTickRoutine(tickInterval)
	})
}

func (t *Timer) selfTickRoutine(tickInterval time.Duration) {
	for {
		time.Sleep(tickInterval)
		t.Tick()
	}
}

func (t *Timer) runCallback(callback TimerCallbackFunc) bool {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Callback %v paniced: %v\n", callback, err)
			debug.PrintStack()
		}
	}()
	return callback()
}
