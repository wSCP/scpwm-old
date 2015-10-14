package main

import "fmt"

type keterError struct {
	err  string
	vals []interface{}
}

func (k *keterError) Error() string {
	return fmt.Sprintf("[keter] %s", fmt.Sprintf(k.err, k.vals...))
}

func (k *keterError) Out(vals ...interface{}) *keterError {
	k.vals = vals
	return k
}

func Krror(err string) *keterError {
	return &keterError{err: err}
}
