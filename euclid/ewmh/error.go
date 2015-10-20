package ewmh

import "fmt"

type ewmhError struct {
	err  string
	vals []interface{}
}

func (e *ewmhError) Error() string {
	return fmt.Sprintf("[ewmh] %s", fmt.Sprintf(e.err, e.vals...))
}

func (e *ewmhError) Out(vals ...interface{}) *ewmhError {
	e.vals = vals
	return e
}

func Xrror(err string) *ewmhError {
	return &ewmhError{err: err}
}
