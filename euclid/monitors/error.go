package monitors

import "fmt"

type monitorsError struct {
	err  string
	vals []interface{}
}

func (m *monitorsError) Error() string {
	return fmt.Sprintf(" %s", fmt.Sprintf(m.err, m.vals...))
}

func (m *monitorsError) Out(vals ...interface{}) *monitorsError {
	m.vals = vals
	return m
}

func Xrror(err string) *monitorsError {
	return &monitorsError{err: err}
}
