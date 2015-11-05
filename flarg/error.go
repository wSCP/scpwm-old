package flarg

import "fmt"

type flargError struct {
	err  string
	vals []interface{}
}

func (a *flargError) Error() string {
	return fmt.Sprintf("[flarg] %s", fmt.Sprintf(a.err, a.vals...))
}

func (a *flargError) Out(vals ...interface{}) *flargError {
	a.vals = vals
	return a
}

func Xrror(err string) *flargError {
	return &flargError{err: err}
}
