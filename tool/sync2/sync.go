package sync2

import (
	"sync"
	"sync/atomic"
)

// OnceSucc is an object that will perform exactly one successful action
type OnceSucc struct {
	done uint32
	m    sync.Mutex
}

// Do func
func (o *OnceSucc) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) == 0 {
		// Outlined slow-path to allow inlining of the fast-path.
		return o.doSlow(f)
	}
	return nil
}

func (o *OnceSucc) doSlow(f func() error) error {
	var err error
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer func() {
			if err == nil {
				atomic.StoreUint32(&o.done, 1)
			}
		}()
		err = f()
		return err
	}
	return nil
}
