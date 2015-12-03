package icccm

import "fmt"

type icccmError struct {
	err  string
	vals []interface{}
}

func (e *icccmError) Error() string {
	return fmt.Sprintf("[icccm] %s", fmt.Sprintf(e.err, e.vals...))
}

func (e *icccmError) Out(vals ...interface{}) *icccmError {
	e.vals = vals
	return e
}

func Xrror(err string) *icccmError {
	return &icccmError{err: err}
}
