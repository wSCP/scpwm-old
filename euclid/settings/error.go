package settings

import "fmt"

type settingsError struct {
	err  string
	vals []interface{}
}

func (m *settingsError) Error() string {
	return fmt.Sprintf(" %s", fmt.Sprintf(m.err, m.vals...))
}

func (m *settingsError) Out(vals ...interface{}) *settingsError {
	m.vals = vals
	return m
}

func Xrror(err string) *settingsError {
	return &settingsError{err: err}
}
