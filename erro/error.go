package erro

import (
	"fmt"

	"github.com/pkg/errors"
)

var _ IError = &Error{}

type CodeError interface {
	error
	Code() int32
	Message() string
}
type IError interface {
	CodeError
	Is(error) bool
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

type Error struct {
	ErrCode    int32
	ErrMessage string
}

func NewError(code int32, msg string) *Error {
	return &Error{
		ErrCode:    code,
		ErrMessage: msg,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%d]%s", e.ErrCode, e.ErrMessage)
}

func (e *Error) Code() int32 {
	if e == nil {
		return 0
	}
	return e.ErrCode
}

func (e *Error) Is(err error) bool {
	if e == nil || err == nil {
		return e == err
	}

	if t, ok := err.(CodeError); ok {
		return t.Code() == e.Code()
	}

	return false
}

func (e *Error) Message() string {
	if e == nil {
		return "Error(nil)"
	}
	return e.ErrMessage
}

func (e *Error) WithMessage(msg string) IError {
	return &Error{
		ErrCode:    e.ErrCode,
		ErrMessage: msg,
	}
}
