package commander

import "fmt"

type commanderError struct {
	err  string
	vals []interface{}
}

func (x *commanderError) Error() string {
	return fmt.Sprintf("[commander] %s", fmt.Sprintf(x.err, x.vals...))
}

func (x *commanderError) Out(vals ...interface{}) *commanderError {
	x.vals = vals
	return x
}

func Xrror(err string) *commanderError {
	return &commanderError{err: err}
}
