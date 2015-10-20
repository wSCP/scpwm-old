package manager

import "fmt"

type managerError struct {
	err  string
	vals []interface{}
}

func (m *managerError) Error() string {
	return fmt.Sprintf(" %s", fmt.Sprintf(m.err, m.vals...))
}

func (m *managerError) Out(vals ...interface{}) *managerError {
	m.vals = vals
	return m
}

func Xrror(err string) *managerError {
	return &managerError{err: err}
}
