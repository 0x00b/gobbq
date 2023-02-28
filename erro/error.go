package erro

import (
	"github.com/pkg/errors"
)

type Error interface {
	Code() int64
	Msg() string
}

type timeoutError interface {
	Timeout() bool // Is it a timeout error
}

type temperaryError interface {
	Temporary() bool
}

// IsTimeout checks if the error is a timeout error
func IsTimeout(err error) bool {
	if err == nil {
		return false
	}

	err = errors.Cause(err)
	ne, ok := err.(timeoutError)
	return ok && ne.Timeout()
}

// IsTimeout checks if the error is a timeout error
func IsTemporary(err error) bool {
	if err == nil {
		return false
	}

	err = errors.Cause(err)
	ne, ok := err.(temperaryError)
	return ok && ne.Temporary()
}
