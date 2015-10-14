package main

import "fmt"

type euclidError struct {
	err  string
	vals []interface{}
}

func (s *euclidError) Error() string {
	return fmt.Sprintf(" %s", fmt.Sprintf(s.err, s.vals...))
}

func (s *euclidError) Out(vals ...interface{}) *euclidError {
	s.vals = vals
	return s
}

func Xrror(err string) *euclidError {
	return &euclidError{err: err}
}
