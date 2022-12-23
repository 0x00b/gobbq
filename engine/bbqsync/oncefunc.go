package bbqsync

import (
	"sync"
)

// OnceFunc returns a function wrapping f which ensures f is only executed
// once even if the returned function is executed multiple times.
func OnceFunc(f func()) func() {
	var once sync.Once
	return func() {
		once.Do(f)
	}
}
