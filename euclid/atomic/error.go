package atomic

import "fmt"

type atomicError struct {
	err  string
	vals []interface{}
}

func (a *atomicError) Error() string {
	return fmt.Sprintf("[atomic] %s", fmt.Sprintf(a.err, a.vals...))
}

func (a *atomicError) Out(vals ...interface{}) *atomicError {
	a.vals = vals
	return a
}

func Xrror(err string) *atomicError {
	return &atomicError{err: err}
}
